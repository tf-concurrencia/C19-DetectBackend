package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/tf-concurrencia/model-training/internal/store"
)

type TreeNode struct {
	ColumnNo int //column number
	Value    interface{}
	Left     *TreeNode
	Right    *TreeNode
	Labels   map[string]int
}

type Tree struct {
	Root *TreeNode
}
type Forest struct {
	Trees []*Tree
}
type ResponseDataset struct {
	Inputs  [][]interface{} `json:"inputs"`
	Targets []string        `json:"targets"`
}

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// DatasetService provides operations on strings.
type TrainService interface {
	TrainModel(n_tree int) (interface{}, float64, float64, error)
}

// datasetService is a concrete implementation of DatasetService
type trainService struct{}

func NewTrainModelService() TrainService {
	return &trainService{}
}

func (trainService) TrainModel(n_tree int) (interface{}, float64, float64, error) {
	start := time.Now()
	resp, err := http.Get("http://host.docker.internal:8081/load-dataset")
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		panic(err2.Error())
	}
	var result ResponseDataset
	if err3 := json.Unmarshal(body, &result); err3 != nil { // Parse []byte to the go struct pointer
		panic(err3.Error())
	}
	// declarar arreglos de entrenamiento
	train_inputs := make([][]interface{}, 0)
	train_targets := make([]string, 0)
	// declarar arreglos de testeo
	test_inputs := make([][]interface{}, 0)
	test_targets := make([]string, 0)
	// dividir entrenamiento y test, preguntando si la fila es divisible entre 3
	for i, x := range result.Inputs {
		if i%3 == 0 {
			test_inputs = append(test_inputs, x)
		} else {
			train_inputs = append(train_inputs, x)
		}
	}
	for i, y := range result.Targets {
		if i%3 == 0 {
			test_targets = append(test_targets, y)
		} else {
			train_targets = append(train_targets, y)
		}
	}
	// Tuneo del arbol: entregar la lista de inputs, targets, y numero de arboles
	forest := DefaultForest(train_inputs, train_targets, n_tree)
	DumpForest(forest, "rf.bin")
	forest = LoadForest("rf.bin")
	// calcular error
	err_count := 0.0
	for i := 0; i < len(test_inputs); i++ {
		output := forest.Predict(test_inputs[i])
		expect := test_targets[i]
		if output != expect {
			err_count += 1
		}
	}
	successs := 1.0 - err_count/float64(len(test_inputs))
	ti := time.Since(start)
	d := store.ForestModel{Forest: forest}
	store.Model = &d
	return forest, successs, ti.Seconds(), nil
}

const CAT = "cat"
const NUMERIC = "numeric"

func getRandomRange(N int, M int) []int {
	tmp := make([]int, N)
	for i := 0; i < N; i++ {
		tmp[i] = i
	}
	for i := 0; i < M; i++ {
		j := i + int(rand.Float64()*float64(N-i))
		tmp[i], tmp[j] = tmp[j], tmp[i]
	}

	return tmp[:M]
}

func getSamples(ary [][]interface{}, index []int) [][]interface{} {
	result := make([][]interface{}, len(index))
	for i := 0; i < len(index); i++ {
		result[i] = ary[index[i]]
	}
	return result
}

func getLabels(ary []string, index []int) []string {
	result := make([]string, len(index))
	for i := 0; i < len(index); i++ {
		result[i] = ary[index[i]]
	}
	return result
}
func (self *Forest) Predict(input []interface{}) string {
	counter := make(map[string]float64)
	for i := 0; i < len(self.Trees); i++ {
		tree_counter := PredictTree(self.Trees[i], input)
		total := 0.0
		for _, v := range tree_counter {
			total += float64(v)
		}
		for k, v := range tree_counter {
			counter[k] += float64(v) / total
		}
	}

	max_c := 0.0
	max_label := ""
	for k, v := range counter {
		if v >= max_c {
			max_c = v
			max_label = k
		}
	}
	return max_label
}
func PredictTree(tree *Tree, input []interface{}) map[string]int {
	return predict(tree.Root, input)
}
func getEntropy(ep_map map[string]float64, total int) float64 {

	for k, _ := range ep_map {
		ep_map[k] = ep_map[k] / float64(total) //normalize
	}

	entropy := 0.0
	for _, v := range ep_map {
		entropy += v * math.Log(1.0/v)
	}

	return entropy
}

func getBestGain(samples [][]interface{}, c int, samples_labels []string, column_type string, current_entropy float64) (float64, interface{}, int, int) {
	var best_value interface{}
	best_gain := 0.0
	best_total_r := 0
	best_total_l := 0

	uniq_values := make(map[interface{}]int)
	for i := 0; i < len(samples); i++ {
		uniq_values[samples[i][c]] = 1
	}

	for value, _ := range uniq_values {
		map_l := make(map[string]float64)
		map_r := make(map[string]float64)
		total_l := 0
		total_r := 0
		if column_type == CAT {
			for j := 0; j < len(samples); j++ {
				if samples[j][c] == value {
					total_l += 1
					map_l[samples_labels[j]] += 1.0
				} else {
					total_r += 1
					map_r[samples_labels[j]] += 1.0
				}
			}
		}
		if column_type == NUMERIC {
			for j := 0; j < len(samples); j++ {
				if samples[j][c].(float64) <= value.(float64) {
					total_l += 1
					map_l[samples_labels[j]] += 1.0
				} else {
					total_r += 1
					map_r[samples_labels[j]] += 1.0
				}
			}
		}

		p1 := float64(total_r) / float64(len(samples))
		p2 := float64(total_l) / float64(len(samples))

		new_entropy := p1*getEntropy(map_r, total_r) + p2*getEntropy(map_l, total_l)
		//fmt.Println(new_entropy,current_entropy)
		entropy_gain := current_entropy - new_entropy

		if entropy_gain >= best_gain {
			best_gain = entropy_gain
			best_value = value
			best_total_l = total_l
			best_total_r = total_r
		}
	}

	return best_gain, best_value, best_total_l, best_total_r
}

func splitSamples(samples [][]interface{}, column_type string, c int, value interface{}, part_l *[]int, part_r *[]int) {
	if column_type == CAT {
		for j := 0; j < len(samples); j++ {
			if samples[j][c] == value {
				*part_l = append(*part_l, j)
			} else {
				*part_r = append(*part_r, j)
			}
		}
	}
	if column_type == NUMERIC {
		for j := 0; j < len(samples); j++ {
			if samples[j][c].(float64) <= value.(float64) {
				*part_l = append(*part_l, j)
			} else {
				*part_r = append(*part_r, j)
			}
		}
	}
}

func buildTree(samples [][]interface{}, samples_labels []string, selected_feature_count int) *TreeNode {
	//fmt.Println(len(samples))
	//find a best splitter
	column_count := len(samples[0])
	//split_count := int(math.Log(float64(column_count)))
	split_count := selected_feature_count
	columns_choosen := getRandomRange(column_count, split_count)

	best_gain := 0.0
	var best_part_l []int = make([]int, 0, len(samples))
	var best_part_r []int = make([]int, 0, len(samples))
	var best_total_l int = 0
	var best_total_r int = 0
	var best_value interface{}
	var best_column int
	var best_column_type string

	current_entropy_map := make(map[string]float64)
	for i := 0; i < len(samples_labels); i++ {
		current_entropy_map[samples_labels[i]] += 1
	}

	current_entropy := getEntropy(current_entropy_map, len(samples_labels))

	for _, c := range columns_choosen {
		column_type := CAT
		if _, ok := samples[0][c].(float64); ok {
			column_type = NUMERIC
		}

		gain, value, total_l, total_r := getBestGain(samples, c, samples_labels, column_type, current_entropy)
		//fmt.Println("kkkkk",gain,part_l,part_r)
		if gain >= best_gain {
			best_gain = gain
			best_value = value
			best_column = c
			best_column_type = column_type
			best_total_l = total_l
			best_total_r = total_r
		}
	}

	if best_gain > 0 && best_total_l > 0 && best_total_r > 0 {
		node := &TreeNode{}
		node.Value = best_value
		node.ColumnNo = best_column
		splitSamples(samples, best_column_type, best_column, best_value, &best_part_l, &best_part_r)
		node.Left = buildTree(getSamples(samples, best_part_l), getLabels(samples_labels, best_part_l), selected_feature_count)
		node.Right = buildTree(getSamples(samples, best_part_r), getLabels(samples_labels, best_part_r), selected_feature_count)
		return node
	}

	return genLeafNode(samples_labels)

}

func genLeafNode(labels []string) *TreeNode {
	counter := make(map[string]int)
	for _, v := range labels {
		counter[v] += 1
	}

	node := &TreeNode{}
	node.Labels = counter
	//fmt.Println(node)
	return node
}

func predict(node *TreeNode, input []interface{}) map[string]int {
	if node.Labels != nil { //leaf node
		return node.Labels
	}

	c := node.ColumnNo
	value := input[c]

	switch value.(type) {
	case float64:
		if value.(float64) <= node.Value.(float64) && node.Left != nil {
			return predict(node.Left, input)
		} else if node.Right != nil {
			return predict(node.Right, input)
		}
	case string:
		if value == node.Value && node.Left != nil {
			return predict(node.Left, input)
		} else if node.Right != nil {
			return predict(node.Right, input)
		}
	}

	return nil
}

func BuildTree(inputs [][]interface{}, labels []string, samples_count, selected_feature_count int) *Tree {

	samples := make([][]interface{}, samples_count)
	samples_labels := make([]string, samples_count)
	for i := 0; i < samples_count; i++ {
		j := int(rand.Float64() * float64(len(inputs)))
		samples[i] = inputs[j]
		samples_labels[i] = labels[j]
	}

	tree := &Tree{}
	tree.Root = buildTree(samples, samples_labels, selected_feature_count)

	return tree
}

func BuildForest(inputs [][]interface{}, labels []string, treesAmount, samplesAmount, selectedFeatureAmount int) *Forest {
	rand.Seed(time.Now().UnixNano())
	forest := &Forest{}
	forest.Trees = make([]*Tree, treesAmount)
	done_flag := make(chan bool)
	prog_counter := 0
	mutex := &sync.Mutex{}
	for i := 0; i < treesAmount; i++ {
		go func(x int) {
			//1
			//fmt.Printf(">> %v building %vth tree...\n", time.Now(), x)
			forest.Trees[x] = BuildTree(inputs, labels, samplesAmount, selectedFeatureAmount)
			//fmt.Printf("<< %v the %vth tree is done.\n",time.Now(), x)
			mutex.Lock()
			prog_counter += 1
			//2
			//fmt.Printf("%v trainning progress %.0f%%\n", time.Now(), float64(prog_counter)/float64(treesAmount)*100)
			mutex.Unlock()
			done_flag <- true
		}(i)
	}

	for i := 1; i <= treesAmount; i++ {
		<-done_flag
	}
	//3
	//fmt.Println("all done.")
	return forest
}

func DefaultForest(inputs [][]interface{}, labels []string, treesAmount int) *Forest {
	m := int(math.Sqrt(float64(len(inputs[0]))))
	n := int(math.Sqrt(float64(len(inputs))))
	return BuildForest(inputs, labels, treesAmount, n, m)
}

func DumpForest(forest *Forest, fileName string) {
	out_f, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		panic("failed to create " + fileName)
	}
	defer out_f.Close()
	encoder := json.NewEncoder(out_f)
	encoder.Encode(forest)
}

func LoadForest(fileName string) *Forest {
	in_f, err := os.Open(fileName)
	if err != nil {
		panic("failed to open " + fileName)
	}
	defer in_f.Close()
	decoder := json.NewDecoder(in_f)
	forest := &Forest{}
	decoder.Decode(forest)
	return forest
}

package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tf-concurrencia/read-dataset-svc/internal/entity"
)

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

type ResponseModel struct {
	Forest *Forest `json:"forest"`
}
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

// DatasetService provides operations on strings.
type PredictService interface {
	PrediceCovid(entity.PredictCovidRequest) (string, error)
}

// datasetService is a concrete implementation of DatasetService
type predictService struct{}

func NewPredictService() PredictService {
	return &predictService{}
}

func (s *predictService) PrediceCovid(predict entity.PredictCovidRequest) (string, error) {
	resp, err := http.Get("http://host.docker.internal:8082/train-model")
	if err != nil {
		fmt.Println("No response from request")
	}
	defer resp.Body.Close()
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println("Fail read Body")
	}
	var result ResponseModel
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to the go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	lista := convert(predict.Inputs)
	rpta := result.Forest.Predicate(lista)
	return rpta, nil
}
func convert(tmp []float64) []interface{} {
	X := make([]interface{}, 0)
	for i := 0; i < len(tmp); i++ {
		X = append(X, tmp[i])
	}
	return X
}
func (self *Forest) Predicate(input []interface{}) string {
	counter := make(map[string]float64)
	for i := 0; i < len(self.Trees); i++ {
		tree_counter := PredicateTree(self.Trees[i], input)
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
func PredicateTree(tree *Tree, input []interface{}) map[string]int {
	return predicate(tree.Root, input)
}
func predicate(node *TreeNode, input []interface{}) map[string]int {
	if node.Labels != nil { //leaf node
		return node.Labels
	}

	c := node.ColumnNo
	value := input[c]

	switch value.(type) {
	case float64:
		if value.(float64) <= node.Value.(float64) && node.Left != nil {
			return predicate(node.Left, input)
		} else if node.Right != nil {
			return predicate(node.Right, input)
		}
	case string:
		if value == node.Value && node.Left != nil {
			return predicate(node.Left, input)
		} else if node.Right != nil {
			return predicate(node.Right, input)
		}
	}

	return nil
}

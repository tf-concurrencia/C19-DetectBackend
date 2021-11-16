package service

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/tf-concurrencia/read-dataset-svc/internal/store"
)

// ErrEmpty is returned when an input string is empty.
var ErrEmpty = errors.New("empty string")

// DatasetService provides operations on strings.
type DatasetService interface {
	LoadDataset(path string, n_rows int) ([][]interface{}, []string, error)
}

// datasetService is a concrete implementation of DatasetService
type datasetService struct{}

func NewDatasetService() DatasetService {
	return &datasetService{}
}

func (datasetService) LoadDataset(path string, n_rows int) ([][]interface{}, []string, error) {
	// Lectura del dataset
	f, _ := os.Open(path) //"TB_F00_SICOVID.csv"
	defer f.Close()
	// Leer todo lo demás del dataset
	content, _ := ioutil.ReadAll(f)
	s_content := string(content)
	lines := strings.Split(s_content, "\n")

	// Declarar inputs y target
	inputs := make([][]interface{}, 0)
	targets := make([]string, 0)
	for i, line := range lines {
		// primera linea
		line = strings.TrimRight(line, "\r\n")
		// si esta vacia, ignorar esta iteracion
		if len(line) == 0 {
			continue
		}
		// ignorar header
		if i == 0 {
			continue
		}
		// crear arreglo
		tup := strings.Split(line, ",")
		// arreglo con indices de 4-18 los cuales seran los inputs
		pattern := tup[4:18]
		// arreglo con el target que se encuentra en el indice 2
		target := tup[2]
		// arreglo X que obtendra los inputs convertidos de dato string a float64
		X := make([]interface{}, 0)
		for _, x := range pattern {
			f_x, _ := strconv.ParseFloat(x, 64)
			X = append(X, f_x)
		}
		// se alamcena en el arreglo inputs
		inputs = append(inputs, X)
		// se alamacena en el arreglo target
		targets = append(targets, target)
		// si llegamos al limite de filas, romper el bucle
		if i == n_rows {
			break
		}
	}

	d := store.Dataset{Inputs: inputs, Targets: targets}
	store.Data = &d
	return inputs, targets, nil
}

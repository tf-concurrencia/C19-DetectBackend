package entity

type LoadDatasetResponse struct {
	Inputs  [][]interface{} `json:"inputs"`
	Targets []string        `json:"targets"`
}

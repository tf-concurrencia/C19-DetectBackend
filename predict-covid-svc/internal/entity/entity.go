package entity

type PredictResponse struct {
	Rpta string `json:"rpta"`
}
type PredictCovidRequest struct {
	Inputs []float64 `json:"inputs"`
}

package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(trainModelEndpoint endpoint.Endpoint) {

	trainModelHandler := httptransport.NewServer(
		trainModelEndpoint,
		decodeTrainModelRequest,
		encodeResponse,
	)

	http.Handle("/train-model", trainModelHandler)
	// http.Handle("/get-dataset", loadDatasetHandler)
}
func decodeTrainModelRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request = struct{}{}
	return request, nil
}
func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

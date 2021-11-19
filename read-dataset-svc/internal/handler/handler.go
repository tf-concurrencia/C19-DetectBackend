package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewHttpHandler(loadDatasetEndpoint endpoint.Endpoint) {

	loadDatasetHandler := httptransport.NewServer(
		loadDatasetEndpoint,
		decodeLoadDatasetRequest,
		encodeResponse,
	)

	http.Handle("/load-dataset", loadDatasetHandler)
	// http.Handle("/get-dataset", loadDatasetHandler)
}

func decodeLoadDatasetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request = struct{}{}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	return json.NewEncoder(w).Encode(response)
}

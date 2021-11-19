package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/tf-concurrencia/read-dataset-svc/internal/entity"
)

func NewHttpHandler(predictCovidEndpoint endpoint.Endpoint) {

	predictCovidtHandler := httptransport.NewServer(
		predictCovidEndpoint,
		decodeLoadDatasetRequest,
		encodeResponse,
	)

	http.Handle("/predict-covid", predictCovidtHandler)
	// http.Handle("/get-dataset", loadDatasetHandler)
}

func decodeLoadDatasetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request entity.PredictCovidRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(response)
}

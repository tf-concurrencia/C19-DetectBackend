package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/tf-concurrencia/read-dataset-svc/internal/entity"
	"github.com/tf-concurrencia/read-dataset-svc/internal/service"
)

func MakePredictCovidEndpoint(svc service.PredictService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(entity.PredictCovidRequest)
		rpta, err := svc.PrediceCovid(req)
		if err != nil {
			return nil, errors.New("Error Predict")
		}
		return entity.PredictResponse{Rpta: rpta}, nil
	}
}

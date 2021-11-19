package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/tf-concurrencia/model-training/internal/entity"
	"github.com/tf-concurrencia/model-training/internal/service"
)

func MakeTrainingModelEndpoint(svc service.TrainService, n_tree int) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {

		forest, success, tiempo, err := svc.TrainModel(n_tree)
		if err != nil {
			return nil, errors.New("Error Train Model")
		}
		return entity.TrainModelResponse{Success: success, Tiempo: tiempo, Forest: forest}, nil
	}
}

package endpoint

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/tf-concurrencia/read-dataset-svc/internal/entity"
	"github.com/tf-concurrencia/read-dataset-svc/internal/service"
)

func MakeLoadDatasetEndpoint(svc service.DatasetService, path string, num_rows int) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		inputs, targets, err := svc.LoadDataset(path, num_rows)
		if err != nil {
			return nil, errors.New("Error loading dataset")
		}
		return entity.LoadDatasetResponse{Inputs: inputs, Targets: targets}, nil
	}
}

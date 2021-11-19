package main

import (
	"net/http"

	"github.com/tf-concurrencia/model-training/internal/endpoint"
	"github.com/tf-concurrencia/model-training/internal/handler"
	"github.com/tf-concurrencia/model-training/internal/service"
)

func main() {

	svc := service.NewTrainModelService()
	trainingModelEndpoint := endpoint.MakeTrainingModelEndpoint(svc, 100)
	handler.NewHttpHandler(trainingModelEndpoint)
	http.ListenAndServe(":8082", nil)

}

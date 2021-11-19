package main

import (
	"net/http"

	"github.com/tf-concurrencia/read-dataset-svc/internal/endpoint"
	"github.com/tf-concurrencia/read-dataset-svc/internal/handler"
	"github.com/tf-concurrencia/read-dataset-svc/internal/service"
)

func main() {

	svc := service.NewPredictService()
	predictCovidEndpoint := endpoint.MakePredictCovidEndpoint(svc)
	handler.NewHttpHandler(predictCovidEndpoint)
	http.ListenAndServe(":8083", nil)

}

package main

import (
	"net/http"

	"github.com/tf-concurrencia/read-dataset-svc/internal/endpoint"
	"github.com/tf-concurrencia/read-dataset-svc/internal/handler"
	"github.com/tf-concurrencia/read-dataset-svc/internal/service"
)

func main() {

	svc := service.NewDatasetService()
	loadDatasetEndpoint := endpoint.MakeLoadDatasetEndpoint(svc, "TB_F00_SICOVID.csv", 1000)
	handler.NewHttpHandler(loadDatasetEndpoint)
	http.ListenAndServe(":8081", nil)

}

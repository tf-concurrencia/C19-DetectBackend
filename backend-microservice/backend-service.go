package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Sintomas struct {
	Inputs []float64 `json:"inputs"`
}
type RepPredict struct {
	resp int `json:"rpta"`
}
type BackendResp struct{
  input []float64
  output int
}

func postSintomas(resp http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		if req.Header.Get("Content-Type") == "application/json" {
			log.Println("Ingresando datos de sintomas")
			log.Println(req.Body)
			cuerpoMsg, err := ioutil.ReadAll(req.Body)

			if err != nil {
				http.Error(resp, "Error interno al leer el body", http.StatusInternalServerError)
			}

			var oSintomas Sintomas
			json.Unmarshal(cuerpoMsg, &oSintomas)
			json_data, err := json.Marshal(oSintomas)
			if err != nil {
				log.Println("Fail conver to json")
			}

			predict_req, err := http.Post("http://20.64.73.44:8083/predict-covid", "application/json", bytes.NewBuffer(json_data))
			if err != nil {
				http.Error(resp, "Error con maquinas virtuales", http.StatusInternalServerError)
			}
			log.Println("Respuesta de las maquinas virtuales")
			log.Println(predict_req.Body)
			cuerpoResputa, err := ioutil.ReadAll(predict_req.Body)
			if err != nil {
				http.Error(resp, "Error al leer la respuesta de la prediccion", http.StatusInternalServerError)
			}
			var oResp RepPredict
			json.Unmarshal(cuerpoResputa, &oResp)
			log.Println(oResp)
      //Generando respuesta
      var oBackendResp BackendResp
      oBackendResp.input = oSintomas.Inputs
      oBackendResp.output = oResp.resp
      json_resp, err := json.Marshal(oBackendResp)
      if err != nil{
				http.Error(resp, "Error al leer la respuesta de la prediccion", http.StatusInternalServerError)
      }
			//Respuesta
			resp.Header().Set("Content-Type", "application-json")
      resp.Write(json_resp)
			//io.WriteString(resp,`
                        //{
                          //"inputs": ,
                          //"output":
                        //}
                        //`)
		} else {
			http.Error(resp, "Contenido invalido", http.StatusBadRequest)
		}
	} else {
		http.Error(resp, "Método invalido", http.StatusMethodNotAllowed)
	}
}

func testRequest(resp http.ResponseWriter, req *http.Request) {
	log.Println("Entre al test")
}
func manejadorSolicitudes() {
	mux := http.NewServeMux()
	//endpoint
	mux.HandleFunc("/predict", postSintomas)
	mux.HandleFunc("/test", testRequest)
	//puerto
	log.Fatal(http.ListenAndServe(":8084", mux))
}

func main() {
	log.Println("Comienza")
	manejadorSolicitudes()
}

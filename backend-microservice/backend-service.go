package main

import (
  "bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Sintomas struct {
  Inputs []float64 `json:"inputs"`
}
type Rep struct {
  resp int `json:"rpta"`
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
      json_data,err := json.Marshal(oSintomas)
      if err != nil {
        log.Println("Fail conver to json")
      }

      predict_req, err := http.Post("http://20.64.73.44:8083/predict-covid","application/json",bytes.NewBuffer(json_data))
      if err != nil{
        http.Error(resp, "Error con maquinas virtuales", http.StatusInternalServerError)
      }
      log.Println("Respuesta de las maquinas virtuales")
      log.Println(predict_req.Body)
      cuerpoResputa, err := ioutil.ReadAll(predict_req.Body)
      if err != nil {
        http.Error(resp, "Error al leer la respuesta de la prediccion", http.StatusInternalServerError)
      }
      var oResp Rep
      json.Unmarshal(cuerpoResputa, &oResp)
      log.Println(oResp)
      //Respuesta
			resp.Header().Set("Content-Type", "application-json")
			io.WriteString(resp, `
                        {
                                "respuesta":"Sintomas recibidos"
                        }
                        `)
		} else {
			http.Error(resp, "Contenido invalido", http.StatusBadRequest)
		}
	} else {
		http.Error(resp, "Método invalido", http.StatusMethodNotAllowed)
	}
}

func testRequest(resp http.ResponseWriter, req *http.Request){
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

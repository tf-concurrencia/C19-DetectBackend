package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Sintomas struct {
  Inputs []float64 `json:"inputs"`
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
      err = json.Unmarshal([]byte(cuerpoMsg), &oSintomas)
      log.Panicln(err)
      //Respuesta
			resp.Header().Set("Content-Type", "application-json")
			io.WriteString(resp, `
                        {
                                "respuesta":"Sintomas recibidos"
                        }
                        `)
			log.Println(oSintomas)

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

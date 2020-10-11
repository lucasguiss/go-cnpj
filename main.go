package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

//CnpjBody - Struct to receive in the request body
type CnpjBody struct {
	Cnpj string
}

// ReceitaWSResponse - Struct for the response of the receita API
type ReceitaWSResponse struct {
	SituationDate string `json:"data_situacao"`
	Type          string `json:"tipo"`
	Name          string `json:"nome"`
	Uf            string `json:"uf"`
	Phone         string `json:"telefone"`
	Situation     string `json:"situacao"`
	District      string `json:"bairro"`
	Street        string `json:"logradouro"`
	Number        string `json:"numero"`
	ZipCode       string `json:"cep"`
	City          string `json:"municipio"`
	Opening       string `json:"abertura"`
	FantasyName   string `json:"fantasia"`
	JuridicNature string `json:"natureza_juridica"`
}

// CreateCnpj - Receives the request body and transforms into CnpjBody
func CreateCnpj(w http.ResponseWriter, r *http.Request) {
	var cnpjBody CnpjBody
	err := json.NewDecoder(r.Body).Decode(&cnpjBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(RequestWSReceita(cnpjBody))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write([]byte(resp))
}

// RequestWSReceita - Retrieves data from receita API
func RequestWSReceita(cnpjBody CnpjBody) ReceitaWSResponse {
	response, err := http.Get("https://www.receitaws.com.br/v1/cnpj/" + cnpjBody.Cnpj)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(response.Body)
	var responseWS ReceitaWSResponse
	json.Unmarshal(bodyBytes, &responseWS)
	fmt.Printf("API Response as struct %+v\n", responseWS)
	return responseWS
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", CreateCnpj).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", router))
}

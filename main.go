package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type ResultadoBusca struct {
	Dados  interface{}
	Origem string
	Err    error
}

func main() {

	http.HandleFunc("/", BuscaCepHandler)
	http.ListenAndServe(":8080", nil)

}

func BuscaCepHandler(w http.ResponseWriter, r *http.Request) {
	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		http.Error(w, "Informe o parâmetro 'cep'", http.StatusBadRequest)
		return
	}

	resultChan := make(chan ResultadoBusca, 2)

	// Goroutines paralelas
	go func() {
		dados, err := BuscaViaCep(cepParam)
		resultChan <- ResultadoBusca{Dados: dados, Origem: "ViaCEP", Err: err}
	}()

	go func() {
		dados, err := BuscaBrasilApi(cepParam)
		resultChan <- ResultadoBusca{Dados: dados, Origem: "BrasilAPI", Err: err}
	}()

	select {
	case resultado := <-resultChan:
		if resultado.Err != nil {
			fmt.Println("Erro ao consultar", resultado.Origem, ":", resultado.Err)
			http.Error(w, "Erro ao consultar o CEP", http.StatusInternalServerError)
			return
		}
		fmt.Println("Resposta recebida de:", resultado.Origem)
		jsonBytes, _ := json.MarshalIndent(resultado.Dados, "", "  ")
		fmt.Println("Dados do endereço:\n", string(jsonBytes))

		// Envia ao navegador apenas como confirmação
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonBytes)

	case <-time.After(1 * time.Second):
		fmt.Println("Timeout: Nenhuma API respondeu em 1 segundo")
		http.Error(w, "Timeout: nenhuma API respondeu a tempo", http.StatusGatewayTimeout)
	}
}

func BuscaViaCep(cep string) (*ViaCEP, error) {

	resp, error := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)

	if error != nil {
		return nil, error
	}

	var c ViaCEP
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}

	return &c, nil

}

func BuscaBrasilApi(cep string) (*BrasilApi, error) {

	resp, error := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep + "")

	if error != nil {
		return nil, error
	}

	defer resp.Body.Close()

	body, error := io.ReadAll(resp.Body)

	if error != nil {
		return nil, error
	}

	var c BrasilApi
	error = json.Unmarshal(body, &c)
	if error != nil {
		return nil, error
	}

	return &c, nil

}

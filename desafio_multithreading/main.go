package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BrasilAPIResponse struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"cidade"`
	UF         string `json:"uf"`
}

type ViaCEPResponse struct {
	CEP        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Cidade     string `json:"localidade"`
	UF         string `json:"uf"`
}

func searchBrasilAPI(cep string, brasilAPI chan<- BrasilAPIResponse) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return
	}

	var data BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
		return
	}

	brasilAPI <- data
}

func searchViaCEP(cep string, viaCEP chan<- ViaCEPResponse) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(err)
		return
	}

	var data ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		fmt.Println(err)
		return
	}

	viaCEP <- data
}

func main() {
	cep := "80610090"
	timeout := 1 * time.Second

	brasilAPI := make(chan BrasilAPIResponse)
	viaCEP := make(chan ViaCEPResponse)

	go searchBrasilAPI(cep, brasilAPI)
	go searchViaCEP(cep, viaCEP)

	select {
	case address := <-brasilAPI:
		fmt.Println(address)
	case address := <-viaCEP:
		fmt.Println(address)
	case <-time.After(timeout):
		fmt.Println("timeout")
	}
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	url := "http://localhost:8080/quotation"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Panic(resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Quotations: %s\n", body)

	if err := ioutil.WriteFile("quotation.txt", body, 0644); err != nil {
		log.Panic(err)
	}
}

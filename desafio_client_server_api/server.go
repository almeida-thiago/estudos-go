package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type QuotationAPI struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

const dbFile = "./quotation.db"

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS quotations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		value TEXT,
		data TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func addQuotation(db *sql.DB, value string) error {
	stmt, err := db.Prepare("INSERT INTO quotations (value) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(value)
	return err
}

func getQuotation() (string, error) {
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var quotation QuotationAPI
	if err := json.NewDecoder(resp.Body).Decode(&quotation); err != nil {
		return "", err
	}

	return quotation.USDBRL.Bid, nil
}

func quotationHandler(w http.ResponseWriter, r *http.Request) {
	quotation, err := getQuotation()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	db, err := initDB()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	if err := addQuotation(db, quotation); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"quotations": "%s"}`, quotation)))
}

func main() {
	http.HandleFunc("/quotation", quotationHandler)
	log.Println("http://localhost:8080/quotation")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

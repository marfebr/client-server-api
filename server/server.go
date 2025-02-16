package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/marfebr/client-server-api/util"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./currency.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableIfNotExists()

	mux := http.NewServeMux()
	mux.HandleFunc("/", Home)
	mux.HandleFunc("/cotacao", UsdtoBrl)

	http.ListenAndServe(":8080", mux)

}

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func UsdtoBrl(w http.ResponseWriter, r *http.Request) {

	var err error
	var currency []util.Currency
	ctx, cancel := context.WithTimeout(r.Context(), time.Millisecond*200)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/USD-BRL", nil)
	if err != nil {
		http.Error(w, "Erro ao criar requisição", http.StatusInternalServerError)
		log.Println(err)

		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("tempo de execução insuficiente")
		http.Error(w, "Erro ao fazer requisição", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Erro ao ler corpo da resposta", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	// log.Println(body)
	err = json.Unmarshal(body, &currency)
	if err != nil {
		http.Error(w, "Erro ao converter JSON", http.StatusInternalServerError)
		log.Println(err)

		return
	}

	log.Println(currency)

	err = insertCurrencyData(r.Context(), currency)

	if err != nil {
		http.Error(w, "Erro ao inserir dados no banco", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func insertCurrencyData(ctx context.Context, currency []util.Currency) error {
	stmt, err := db.Prepare(`INSERT INTO currency (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*10)
	defer cancel()

	for _, c := range currency {
		if _, err := stmt.ExecContext(ctx, c.Code, c.Codein, c.Name, c.High, c.Low, c.VarBid, c.PctChange, c.Bid, c.Ask, c.Timestamp, c.CreateDate); err != nil {
			log.Printf("Erro ao inserir dados no banco: %v", err)
			return err

		}
	}

	return nil
}

func createTableIfNotExists() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS currency (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT,
		codein TEXT,
		name TEXT,
		high TEXT,
		low TEXT,
		varBid TEXT,
		pctChange TEXT,
		bid TEXT,
		ask TEXT,
		timestamp TEXT,
		create_date TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

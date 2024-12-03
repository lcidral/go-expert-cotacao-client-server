package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Quote struct {
	USDBRL struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type QuoteServerResponse struct {
	Bid string `json:"bid"`
}

func main() {
	db, err := sql.Open("sqlite3", "./quotes.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS quotes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        bid TEXT,
        created_at DATETIME
    )`)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Contexto para a chamada à API externa
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {
			log.Printf("Erro ao criar request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Erro ao fazer request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		var quote Quote
		if err := json.NewDecoder(resp.Body).Decode(&quote); err != nil {
			log.Printf("Erro ao decodificar resposta: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Contexto para salvar no banco de dados
		ctxDB, cancelDB := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancelDB()

		// Salvar no banco de dados com timeout de 10ms
		query := `INSERT INTO quotes (bid, created_at) VALUES (?, ?)`
		_, err = db.ExecContext(ctxDB, query, quote.USDBRL.Bid, time.Now())
		if err != nil {
			log.Printf("Erro ao salvar no banco: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := QuoteServerResponse{
			Bid: quote.USDBRL.Bid,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Servidor rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

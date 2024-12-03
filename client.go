package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type QuoteClientResponse struct {
	Bid string `json:"bid"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var quote QuoteClientResponse
	if err := json.Unmarshal(body, &quote); err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf("Dólar: %s", quote.Bid)
	err = ioutil.WriteFile("cotacao.txt", []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Cotação atual: %s\n", quote.Bid)
	fmt.Println("Cotação salva com sucesso no arquivo cotacao.txt")
}

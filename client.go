package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/marfebr/client-server-api/util"
)

func main() {
	// Code for Desafio1/client-server-api/client.go

	var currency []util.Currency
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Millisecond*300)
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {

		log.Fatal("Erro na requisição")
		// panic(err)
	}
	defer cancel()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {

		fmt.Println(err)
		log.Fatal("tempo de execução insuficiente")
		// panic(err)

	}
	defer resp.Body.Close()

	// log.Println(resp.Status)
	if resp.StatusCode != http.StatusOK {
		log.Fatal("Error ao acessar o servidor")
	}

	err = json.NewDecoder(resp.Body).Decode(&currency)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("A cotação do dolar hoje é: %s", currency[0].Bid)
	f, err := os.Create("cotacao.txt")

	if err != nil {
		fmt.Println(err)
		log.Fatal("Erroa ao criar o arquivo")
	}
	defer f.Close()

	size, err := f.Write([]byte(fmt.Sprintf("Dólar: %s ", currency[0].Bid)))
	if err != nil {
		fmt.Printf("Erro gravando o arquivo: %v", err)
	}
	fmt.Printf("Ultima cotação salva no arquivo: cotacao.txt tamanho %d bytes", size)

}

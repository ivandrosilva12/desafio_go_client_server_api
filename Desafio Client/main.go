package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

type CotacaoFinal struct {
	Valor float64 `json:"Dolar"`
}

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if ctx.Err() != nil {
		println("Excedeu o tempo limite para a requisicao")
		return
	}

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	f, err := os.Create("cotacao.txt")

	if err != nil {
		panic(err)
	}
	_, err = f.Write([]byte(body))

	if err != nil {
		panic(err)
	}

	f.Close()
}

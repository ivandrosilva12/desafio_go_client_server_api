package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Usdbrl struct {
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
	gorm.Model
}

type CotacaoFinal struct {
	Valor float64 `json:"Dolar"`
}

type Cotacao struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}

func main() {

	http.HandleFunc("/cotacao", PesquisarCotacao)
	http.ListenAndServe(":8080", nil)

}

func PesquisarCotacao(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ctx := r.Context()

	cotacao, error := buscarCotacao(ctx)

	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var c Cotacao
	err := json.Unmarshal(cotacao, &c)
	if err != nil {
		println(err)
	}
	error = saveToDatabase(ctx, &c)

	if error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	vf, _ := strconv.ParseFloat(c.Usdbrl.Bid, 64)
	var cf CotacaoFinal = CotacaoFinal{Valor: vf}

	json.NewEncoder(w).Encode(cf)
}

func buscarCotacao(ctx context.Context) ([]byte, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Nanosecond*200)

	defer cancel()

	req, err := http.NewRequest("GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		println(err)
	}

	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if ctx.Err() != nil {
		println("Timeout - API request time exceeded!!!")
	}

	if err != nil {
		println(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		println(err)
	}

	return body, err
}

func saveToDatabase(ctx context.Context, c *Cotacao) error {

	ctx, cancel := context.WithTimeout(ctx, time.Nanosecond*10)

	defer cancel()

	db, err := gorm.Open(sqlite.Open("cotacao.db?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})

	if err != nil {
		println("failed to connect database")
	}

	db.AutoMigrate(&Usdbrl{})

	db.Create(&c.Usdbrl)

	if ctx.Err() != nil {
		println("Timeout - Database request time exceeded!!!")
	}

	return err

}

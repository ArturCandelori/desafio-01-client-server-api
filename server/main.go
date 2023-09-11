package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type APIResponse struct {
	USDBRL ExchangeRate `json:"USDBRL"`
}

type ExchangeRate struct {
	ID         string `gorm:"primaryKey"`
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
}

type ResponseToClient struct {
	Cotacao string `json:"cotacao"`
}

func main() {
	// connect to db and create table
	db, err := gorm.Open(sqlite.Open("cotacoes.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&ExchangeRate{})

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		exchangeRate := getExchangeRate()

		// save response to db
		err = saveExchangeRate(db, exchangeRate)
		if err != nil {
			panic(err)
		}

		// send data to client
		responseBody, err := json.Marshal(ResponseToClient{Cotacao: exchangeRate.Bid})
		if err != nil {
			panic(err)
		}
		w.Write(responseBody)
	})
	http.ListenAndServe(":8080", nil)
}

func getExchangeRate() *ExchangeRate {
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	// send request
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// parse response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	var apiResponse APIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		panic(err)
	}

	return &apiResponse.USDBRL
}

func saveExchangeRate(db *gorm.DB, exchangeRate *ExchangeRate) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	return db.WithContext(ctx).Create(&exchangeRate).Error
}

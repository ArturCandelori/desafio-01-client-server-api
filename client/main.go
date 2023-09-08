package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	result := getResult()
	saveToFile(result)
}

func getResult() []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	result, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return result
}

func saveToFile(content []byte) {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		panic(err)
	}
	_, err = f.Write(content)
	if err != nil {
		panic(err)
	}
}

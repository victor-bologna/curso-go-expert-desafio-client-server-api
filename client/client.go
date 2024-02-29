package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	body := getCotacaoWithContext(ctx)
	writeCotacao(body)
}

func getCotacaoWithContext(ctx context.Context) []byte {
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func writeCotacao(body []byte) {
	file, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.Write([]byte(fmt.Sprint("Dolár:", string(body))))
}

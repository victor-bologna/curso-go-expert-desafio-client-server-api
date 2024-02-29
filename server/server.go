package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Cambio struct {
	USDBRL struct {
		Code       string `json:"code"`
		CodeIn     string `json:"codein"`
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

var db *sql.DB

func main() {
	db = initializeDB()
	defer db.Close()
	http.HandleFunc("/cotacao", Cotacao)
	http.HandleFunc("/verifyCotacoes", VerifyCotacoes)
	http.ListenAndServe(":8080", nil)
}

func Cotacao(w http.ResponseWriter, r *http.Request) {
	ctxHttp, cancelHttp := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancelHttp()

	body := getCambioHttpWithContext(ctxHttp, w)
	cambio := buildResponse(body)
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(cambio.USDBRL.Bid)

	ctxDB, cancelDB := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancelDB()
	persistDataWithContext(ctxDB, db, cambio)
}

func VerifyCotacoes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM cambio")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	defer rows.Close()

	var cambios []Cambio
	for rows.Next() {
		var cambio Cambio
		err := rows.Scan(
			&cambio.USDBRL.Code,
			&cambio.USDBRL.CodeIn,
			&cambio.USDBRL.Name,
			&cambio.USDBRL.High,
			&cambio.USDBRL.Low,
			&cambio.USDBRL.VarBid,
			&cambio.USDBRL.PctChange,
			&cambio.USDBRL.Bid,
			&cambio.USDBRL.Ask,
			&cambio.USDBRL.Timestamp,
			&cambio.USDBRL.CreateDate,
		)
		if err != nil {
			panic(err)
		}
		cambios = append(cambios, cambio)
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cambios)

}

func getCambioHttpWithContext(ctx context.Context, w http.ResponseWriter) []byte {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		panic(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return body
}

func buildResponse(body []byte) Cambio {
	var cambio Cambio
	err := json.Unmarshal(body, &cambio)
	if err != nil {
		panic(err)
	}
	return cambio
}

func persistDataWithContext(ctx context.Context, db *sql.DB, cambio Cambio) {
	stmt, err := db.PrepareContext(ctx, "INSERT INTO cambio (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx,
		cambio.USDBRL.Code,
		cambio.USDBRL.CodeIn,
		cambio.USDBRL.Name,
		cambio.USDBRL.High,
		cambio.USDBRL.Low,
		cambio.USDBRL.VarBid,
		cambio.USDBRL.PctChange,

		cambio.USDBRL.Bid,
		cambio.USDBRL.Ask,
		cambio.USDBRL.Timestamp,
		cambio.USDBRL.CreateDate,
	)
	if err != nil {
		panic(err)
	}
}

func initializeDB() *sql.DB {
	os.Remove("./cambio.db")
	db, err := sql.Open("sqlite3", "./cambio.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cambio (
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
	return db
}

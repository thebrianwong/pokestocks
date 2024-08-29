package main

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
)

const (
	numOfPokemon = 1025
)

type stockInfo struct {
	symbol string
	name   string
}

func readCsv(file string) [][]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal("Error opening csv file")
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Error reading csv file")
	}

	return records
}

func getTopStocks(stocks [][]string) [][]string {
	return stocks[1:(numOfPokemon + 1)]
}

func insertIntoDb(ctx context.Context, db *pgx.Conn, stocks [][]string) {
	options := pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		log.Fatalf("Error starting db transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	for _, stock := range stocks {
		data := stockInfo{symbol: stock[0], name: stock[1]}
		query := "INSERT INTO stocks (symbol, name) VALUES ($1, $2)"
		_, err := tx.Conn().Exec(ctx, query, data.symbol, data.name)
		if err != nil {
			log.Fatalf("Error inserting "+data.symbol+" into db: %v", err)
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Fatalf("Error committing db transaction: %v", err)
		return
	}
}

func main() {
	utils.LoadEnvVars("../../.env")
	conn := utils.ConnectToDb()

	records := readCsv("../../data/stock_data.csv")
	stocks := getTopStocks(records)

	insertIntoDb(context.Background(), conn, stocks)
}

package main

import (
	"context"
	"encoding/csv"
	"os"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		utils.LogFailureError("Error opening csv file", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		utils.LogFailureError("Error reading csv file", err)
	}

	return records
}

func getTopStocks(stocks [][]string) [][]string {
	return stocks[1:(numOfPokemon + 1)]
}

func insertIntoDb(ctx context.Context, db *pgxpool.Pool, stocks [][]string) {
	options := pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		utils.LogFailureError("Error starting db transaction", err)
	}

	defer tx.Rollback(ctx)

	batch := pgx.Batch{}
	for _, stock := range stocks {
		data := stockInfo{symbol: stock[0], name: stock[1]}
		query := `
			INSERT INTO stocks (symbol, name)
			SELECT symbol, name
			FROM (VALUES ($1, $2)) AS data(symbol, name)
			WHERE NOT EXISTS (
				SELECT 1
				FROM stocks
				WHERE stocks.symbol = data.symbol
			)
		`
		batch.Queue(query, data.symbol, data.name)
		if err != nil {
			utils.LogFailureError("Error inserting "+data.symbol+" into db", err)
		}
	}

	err = db.SendBatch(ctx, &batch).Close()
	if err != nil {
		utils.LogFailureError("Error sending batch inserts", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		utils.LogFailureError("Error committing db transaction", err)
	}
}

func main() {
	utils.LoadEnvVars("../../.env")
	conn := utils.ConnectToDb()

	records := readCsv("../../data/stock_data.csv")
	stocks := getTopStocks(records)

	insertIntoDb(context.Background(), conn, stocks)
}

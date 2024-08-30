package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type typeData struct {
	Name   string `json:"name"`
	Sprite string `json:"url"`
}

func readTypesJson(file string) []typeData {
	data := []typeData{}

	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Error opening json file: %v", err)
	}
	defer f.Close()

	caser := cases.Title(language.English, cases.NoLower)

	rawData := []typeData{}
	dec := json.NewDecoder(f)
	err = dec.Decode(&rawData)
	if err != nil {
		log.Fatalf("Error reading type data into memory: %v", err)
	}

	for _, dataObj := range rawData {
		formattedData := typeData{Name: caser.String(dataObj.Name), Sprite: dataObj.Sprite}
		data = append(data, formattedData)
	}

	return data
}

func insertIntoDb(ctx context.Context, db *pgx.Conn, tData []typeData) {
	options := pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		log.Fatalf("Error starting db transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	batch := pgx.Batch{}
	for _, data := range tData {
		query := `
			INSERT INTO pokemon_types (type, sprite_url)
			SELECT type, sprite_url
			FROM (VALUES ($1, $2)) AS data(type, sprite_url)
			WHERE NOT EXISTS (
				SELECT 1
				FROM pokemon_types
				WHERE pokemon_types.type = data.type
			)
		`
		batch.Queue(query, data.Name, data.Sprite)
	}

	err = db.SendBatch(ctx, &batch).Close()
	if err != nil {
		log.Fatalf("Error sending batch inserts: %v", err)
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

	types := readTypesJson("../../data/types.json")
	insertIntoDb(context.Background(), conn, types)
}

package main

import (
	"context"
	"math/rand"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	numOfPokemon = 1025
)

type pokemonStockPair struct {
	pokemonPokedexNumber int
	stockId              int
}

func randomPokedexNumbers() []int {
	random := []int{}

	for i := 1; i <= numOfPokemon; i++ {
		random = append(random, i)
	}

	for i := range random {
		j := rand.Intn(i + 1)
		random[i], random[j] = random[j], random[i]
	}

	return random
}

func getStockIds(ctx context.Context, db *pgxpool.Pool) []int {
	query := "SELECT id FROM stocks"
	rows, err := db.Query(ctx, query)
	if err != nil {
		utils.LogFailureError("Error querying stock ids", err)
	}
	defer rows.Close()

	var stockIds []int

	for rows.Next() {
		var stockId int

		err := rows.Scan(&stockId)
		if err != nil {
			utils.LogFailureError("Error reading stock id", err)
		}

		stockIds = append(stockIds, stockId)
	}

	if err = rows.Err(); err != nil {
		utils.LogFailureError("Error while finishing up reading rows", err)
	}

	return stockIds
}

func insertRandomPokemonStocksIntoDb(ctx context.Context, db *pgxpool.Pool, pokedexNumbers []int, stockIds []int, seasonName string) {
	if len(pokedexNumbers) != len(stockIds) {
		utils.LogFailure("The number of Pokemon do not match the number of stocks. Query the stocks table and check how many rows exist.")
	}

	options := pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		utils.LogFailureError("Error starting db transaction", err)
	}

	defer tx.Rollback(ctx)

	batch := pgx.Batch{}
	for i := 0; i < len(pokedexNumbers); i++ {
		pokemonStockPair := pokemonStockPair{
			pokemonPokedexNumber: pokedexNumbers[i],
			stockId:              stockIds[i],
		}
		query := `
			INSERT INTO pokemon_stock_pairs (pokemon_id, stock_id, season_id)
			SELECT pokemon_id, stock_id, season_id
			FROM (VALUES (
				(SELECT id FROM pokemon WHERE pokedex_number=$1),
				$2::BIGINT,
				(SELECT id FROM seasons WHERE name=$3)
			))
			AS data(pokemon_id, stock_id, season_id)
			WHERE NOT EXISTS (
				SELECT 1
				FROM pokemon_stock_pairs
				WHERE pokemon_stock_pairs.stock_id = data.stock_id
				AND pokemon_stock_pairs.season_id = data.season_id
			)
		`
		batch.Queue(query, pokemonStockPair.pokemonPokedexNumber, pokemonStockPair.stockId, seasonName)
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
	seasonName := utils.GetSeasonName()

	utils.LoadEnvVars("../../.env")
	conn := utils.ConnectToDb()

	pokedexNumbers := randomPokedexNumbers()
	stockIds := getStockIds(context.Background(), conn)
	insertRandomPokemonStocksIntoDb(context.Background(), conn, pokedexNumbers, stockIds, seasonName)
}

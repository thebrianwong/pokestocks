package main

import (
	"context"
	"log"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
)

func insertIntoDb(ctx context.Context, db *pgx.Conn, seasonName string) {
	options := pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		log.Fatalf("Error starting db transaction: %v", err)
	}

	defer tx.Rollback(ctx)

	query := `
		INSERT INTO seasons (name, active)
		SELECT name, active
		FROM (VALUES ($1, TRUE)) AS data(name, active)
		WHERE NOT EXISTS (
			SELECT 1
			FROM seasons
			WHERE seasons.name = data.name
		)
	`
	_, err = tx.Conn().Exec(ctx, query, seasonName)
	if err != nil {
		log.Fatalf("Error inserting season "+seasonName+" into db: %v", err)
		return
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

	seasonName := utils.GetSeasonName()

	insertIntoDb(context.Background(), conn, seasonName)
}

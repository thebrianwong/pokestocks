package main

import (
	"context"
	"pokestocks/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func insertIntoDb(ctx context.Context, db *pgxpool.Pool, seasonName string) {
	options := pgx.TxOptions{IsoLevel: pgx.RepeatableRead, AccessMode: pgx.ReadWrite, DeferrableMode: pgx.Deferrable}
	tx, err := db.BeginTx(ctx, options)
	if err != nil {
		utils.LogFailureError("Error starting db transaction", err)
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
		utils.LogFailureError("Error inserting season "+seasonName+" into db", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		utils.LogFailureError("Error committing db transaction", err)
	}
}

func main() {
	utils.LoadEnvVars("../../.env")
	conn := utils.ConnectToDb()

	seasonName := utils.GetSeasonName()

	insertIntoDb(context.Background(), conn, seasonName)
}

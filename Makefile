include .env

dbUser := ${PG_USER}
dbPassword := ${PG_PASSWORD}
dbHost := ${PG_HOST}
dbPort := ${PG_PORT}
dbName := ${PG_NAME}
dbUrl := "postgres://${dbUser}:${dbPassword}@${dbHost}:${dbPort}/${dbName}"

migrate_up:
	goose -dir "./migrations" postgres ${dbUrl} up

migrate_down_one:
	goose -dir "./migrations" postgres ${dbUrl} down

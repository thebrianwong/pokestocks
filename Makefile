include .env

dbUser := ${PG_USER}
dbPassword := ${PG_PASSWORD}
dbHost := ${PG_HOST}
dbPort := ${PG_PORT}
dbName := ${PG_NAME}
dbUrl := "postgres://${dbUser}:${dbPassword}@${dbHost}:${dbPort}/${dbName}"

seasonName := ${MAKEFILE_SEASON_NAME}

migrate_up:
	goose -dir "./migrations" postgres ${dbUrl} up

migrate_down_one:
	goose -dir "./migrations" postgres ${dbUrl} down

seed_db:
	cd scripts/add_types/; \
		go run main.go
	cd scripts/add_season/; \
		go run main.go ${seasonName}
	cd scripts/add_pokemon/; \
		go run main.go
	cd scripts/add_stocks/; \
		go run main.go
	cd scripts/random_mapping/; \
		go run main.go ${seasonName}

compose: migrate_up seed_db remake_psp_index seed_psp_index
	air

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/**/*.proto

remake_psp_index:
	cd scripts/delete_pokemon_stock_pairs_index/; \
		go run main.go
	cd scripts/create_pokemon_stock_pairs_index/; \
		go run main.go

seed_psp_index:
	cd scripts/index_pokemon_stock_pairs/; \
		go run main.go
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

.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/common/common.proto
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/pokemon_stock_pair/pokemon_stock_pair.proto

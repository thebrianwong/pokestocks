package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func LoadEnvVars(path string) {
	var err error
	if path == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(path)
	}
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func ConnectToDb() *pgx.Conn {
	dbUrl := os.Getenv("DATABASE_URL")
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatalln("Unable to connect to database:", err)
	}
	return conn
}

func GetSeasonName() string {
	if len(os.Args) != 2 {
		log.Fatalln("You must provide a season name.\nUsage: go run main.go [name]")
		fmt.Println("go run main.go $SEASON_NAME")
	}

	seasonName := os.Args[1]

	return seasonName
}

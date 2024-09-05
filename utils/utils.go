package utils

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
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

func ConnectToDb() *pgxpool.Pool {
	dbUser := os.Getenv("PG_USER")
	dbPassword := os.Getenv("PG_PASSWORD")
	dbHost := os.Getenv("PG_HOST")
	dbPort := os.Getenv("PG_PORT")
	dbName := os.Getenv("PG_NAME")
	dbUrl := "postgres://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
	conn, err := pgxpool.New(context.Background(), dbUrl)
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

func ConnectToElastic(certPath string) *elasticsearch.TypedClient {
	elasticUsername := os.Getenv("ELASTIC_USERNAME")
	elasticPassword := os.Getenv("ELASTIC_PASSWORD")
	// elasticApiKey := os.Getenv("ELASTIC_API_KEY")
	elasticEndpoint := os.Getenv(("ELASTIC_ENDPOINT"))
	// cert, err := os.ReadFile("../../http_ca.crt")
	cert, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatalf("Error reading Elasticsearch certificate: %v", err)
	}
	elasticClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		// APIKey:    elasticApiKey,
		Addresses: []string{elasticEndpoint},
		Username:  elasticUsername,
		Password:  elasticPassword,
		CACert:    cert,
	})

	if err != nil {
		log.Fatalf("Error connecting to Elasticsearch: %v", err)
	}

	return elasticClient
}

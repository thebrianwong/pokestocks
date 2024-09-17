package utils

import (
	"context"
	"log"
	"os"

	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Reset  = "\033[0m"
)

func LoadEnvVars(path string) {
	var err error
	if path == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(path)
	}
	if err != nil {
		LogFailureError("Error loading .env file", err)
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
		LogFailureError("Unable to connect to database", err)
	}
	return conn
}

func GetSeasonName() string {
	if len(os.Args) != 2 {
		LogFailure("You must provide a season name.\n" + Cyan + "Usage: go run main.go [name]")
	}

	seasonName := os.Args[1]

	return seasonName
}

func CreateTypedElasticClient(certPath string) *elasticsearch.TypedClient {
	elasticUsername := os.Getenv("ELASTIC_USERNAME")
	elasticPassword := os.Getenv("ELASTIC_PASSWORD")
	// elasticApiKey := os.Getenv("ELASTIC_API_KEY")
	elasticEndpoint := os.Getenv(("ELASTIC_ENDPOINT"))
	var cert []byte
	var err error
	if certPath == "" {
		cert, err = os.ReadFile("./http_ca.crt")
	} else {
		cert, err = os.ReadFile(certPath + "http_ca.crt")
	}
	if err != nil {
		LogFailureError("Error reading Elasticsearch certificate", err)
	}
	elasticClient, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		// APIKey:    elasticApiKey,
		Addresses: []string{elasticEndpoint},
		Username:  elasticUsername,
		Password:  elasticPassword,
		CACert:    cert,
	})

	if err != nil {
		LogFailureError("Error connecting to Elasticsearch", err)
	}

	return elasticClient
}

func CreateRegularElasticClient(certPath string) *elasticsearch.Client {
	elasticUsername := os.Getenv("ELASTIC_USERNAME")
	elasticPassword := os.Getenv("ELASTIC_PASSWORD")
	// elasticApiKey := os.Getenv("ELASTIC_API_KEY")
	elasticEndpoint := os.Getenv(("ELASTIC_ENDPOINT"))
	var cert []byte
	var err error
	if certPath == "" {
		cert, err = os.ReadFile("./http_ca.crt")
	} else {
		cert, err = os.ReadFile(certPath + "http_ca.crt")
	}
	if err != nil {
		LogFailureError("Error reading Elasticsearch certificate", err)
	}
	elasticClient, err := elasticsearch.NewClient(elasticsearch.Config{
		// APIKey:    elasticApiKey,
		Addresses: []string{elasticEndpoint},
		Username:  elasticUsername,
		Password:  elasticPassword,
		CACert:    cert,
	})

	if err != nil {
		LogFailureError("Error creating regular Elasticsearch client", err)
	}

	return elasticClient
}

func CreateAlpacaMarketDataClient() *marketdata.Client {
	key := os.Getenv("ALPACA_TRADING_KEY_ID")
	secret := os.Getenv("ALPACA_TRADING_SECRET_KEY")
	alpacaClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    key,
		APISecret: secret,
	})
	return alpacaClient
}

func CreateRedisClient() *redis.Client {
	redisAddress := os.Getenv("REDIS_ADDRESS")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       0,
	})

	res, err := redisClient.Ping(context.Background()).Result()

	if !(res == "PONG" && err == nil) {
		LogFailureError("Error creating Redis client", err)
	}

	return redisClient
}

func LogFailureError(message string, err error) {
	log.Fatalf(Red+message+": %v"+Reset, err)
}

func LogFailure(message string) {
	log.Fatalf(Red + message + Reset)
}

func LogSuccess(message string) {
	log.Println(Green + message + Reset)
}

func LogWarning(message string) {
	log.Println(Yellow + message)
}

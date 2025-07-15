package util

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/vuphan121/quotum/storage"
)

type Config struct {
	Rate      int
	Interval  time.Duration
	Store     storage.Storage
	Logging   bool
	Algorithm string
	APIKey    string
	RedisURL  string
}

var AppConfig Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	rateStr := os.Getenv("RATE")
	if rateStr == "" {
		log.Fatal("RATE not in .env")
	}
	rate, err := strconv.Atoi(rateStr)
	if err != nil {
		log.Fatalf("Invalid RATE: %v", err)
	}

	intervalStr := os.Getenv("INTERVAL")
	if intervalStr == "" {
		log.Fatal("INTERVAL not in .env")
	}
	interval, err := time.ParseDuration(intervalStr)
	if err != nil {
		log.Fatalf("Invalid INTERVAL: %v", err)
	}

	storeType := os.Getenv("STORE")
	if storeType == "" {
		log.Fatal("STORE not in .env ('memory' or 'redis')")
	}

	var store storage.Storage
	var redis_url string

	switch strings.ToLower(storeType) {
	case "memory":
		store = storage.SharedMemoryStorage
	case "redis":
		redis_url = os.Getenv("REDIS_URL")
	default:
		log.Fatalf("Invalid STORE: %s ('memory' or 'redis')", storeType)
	}

	algorithm := strings.ToLower(os.Getenv("ALGORITHM"))
	if algorithm == "" {
		algorithm = "fixed_window"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY not in .env")
	}

	logging := strings.ToLower(os.Getenv("LOGGING")) == "true"

	AppConfig = Config{
		Rate:      rate,
		Interval:  interval,
		Store:     store,
		Logging:   logging,
		Algorithm: algorithm,
		APIKey:    apiKey,
		RedisURL:  redis_url,
	}
}

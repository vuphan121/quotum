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
	Rate     int
	Interval time.Duration
	Store    storage.Storage
	Logging  bool
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
	switch strings.ToLower(storeType) {
	case "memory":
		store = storage.SharedMemoryStorage
	//case "redis":
	//	store = storage.NewRedisStorage() //
	default:
		log.Fatalf("Invalid STORE: %s ('memory' or 'redis')", storeType)
	}

	logging := strings.ToLower(os.Getenv("LOGGING")) == "true"

	AppConfig = Config{
		Rate:     rate,
		Interval: interval,
		Store:    store,
		Logging:  logging,
	}
}

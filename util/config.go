package util

import (
	"crypto/tls"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/vuphan121/quotum/storage"
)

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}
type Config struct {
	Rate      int
	Interval  time.Duration
	Store     storage.Storage
	Logging   bool
	Algorithm string
	APIKey    string
	Redis     *RedisConfig
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
	var redisCfg *RedisConfig

	switch strings.ToLower(storeType) {
	case "memory":
		store = storage.SharedMemoryStorage
	case "redis":
		addr := os.Getenv("REDIS_ADDR")
		username := os.Getenv("REDIS_USERNAME")
		password := os.Getenv("REDIS_PASSWORD")
		dbStr := os.Getenv("REDIS_DB")
		tlsEnabled := os.Getenv("REDIS_TLS") == "true"
		db := 0
		if dbStr != "" {
			db, err = strconv.Atoi(dbStr)
			if err != nil {
				log.Fatalf("Invalid REDIS_DB: %v", err)
			}
		}

		redisCfg = &RedisConfig{
			Addr:     addr,
			Username: username,
			Password: password,
			DB:       db,
		}

		redisOptions := &redis.Options{
			Addr:     redisCfg.Addr,
			Username: redisCfg.Username,
			Password: redisCfg.Password,
			DB:       redisCfg.DB,
		}

		if tlsEnabled {
			redisOptions.TLSConfig = &tls.Config{}
			fmt.Println("[CONFIG] TLS enabled for Redis")
		} else {
			fmt.Println("[CONFIG] TLS not enabled for Redis")
		}

		redis_client := redis.NewClient(redisOptions)
		store = storage.NewRedisStorage(redis_client)
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
		Redis:     redisCfg,
	}
}

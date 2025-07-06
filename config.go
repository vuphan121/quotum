package quotum

import (
	"quotum/storage"
	"time"
)

type LimiterConfig struct {
	Rate     int
	Interval time.Duration
	Key      string
	Storage  storage.Storage
}

func NewLimiter(config LimiterConfig) *Limiter {
	return &Limiter{
		rate:     config.Rate,
		interval: config.Interval,
		storage:  config.Storage,
		key:      config.Key,
	}
}

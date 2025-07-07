package limiter

import (
	"github.com/vuphan121/quotum/storage"
	"time"
)

type LimiterConfig struct {
	Key      string
	Rate     int
	Interval time.Duration
	Storage  storage.Storage
}

func NewLimiter(cfg LimiterConfig) *Limiter {
	return &Limiter{
		rate:     cfg.Rate,
		interval: cfg.Interval,
		key:      cfg.Key,
		storage:  cfg.Storage,
	}
}

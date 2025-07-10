package limiter

import (
	"github.com/vuphan121/quotum/storage"
	"time"
)

type Config struct {
	Key       string
	Rate      int
	Interval  time.Duration
	Storage   storage.Storage
	Algorithm string
}

func NewLimiter(cfg Config) *Limiter {
	return &Limiter{
		rate:      cfg.Rate,
		interval:  cfg.Interval,
		key:       cfg.Key,
		storage:   cfg.Storage,
		algorithm: cfg.Algorithm,
	}
}

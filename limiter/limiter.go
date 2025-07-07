package limiter

import (
	"github.com/vuphan121/quotum/algorithms"
	"github.com/vuphan121/quotum/storage"
	"time"
)

type Limiter struct {
	rate     int
	interval time.Duration
	storage  storage.Storage
	key      string
}

func (l *Limiter) Allow() bool {
	return algorithms.FixedWindow(l.key, l.rate, l.interval, l.storage)
}

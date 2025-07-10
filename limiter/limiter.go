package limiter

import (
	"github.com/vuphan121/quotum/algorithms"
	"github.com/vuphan121/quotum/storage"
	"time"
)

type Limiter struct {
	rate      int
	interval  time.Duration
	storage   storage.Storage
	key       string
	algorithm string
}

func (l *Limiter) Allow() (bool, *time.Time) {
	switch l.algorithm {
	case "fixed":
		return algorithms.FixedWindow(l.key, l.rate, l.interval, l.storage)
	// case "sliding":
	// 	return algorithms.SlidingWindow(...)
	default:
		return algorithms.FixedWindow(l.key, l.rate, l.interval, l.storage)
	}
}

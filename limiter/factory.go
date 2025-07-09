package limiter

import (
	"github.com/vuphan121/quotum/storage"
	"time"
)

func CreateUserLimiter(key string, rate int, interval time.Duration, storage storage.Storage) *Limiter {
	return NewLimiter(Config{
		Key:      key,
		Rate:     rate,
		Interval: interval,
		Storage:  storage,
	})
}

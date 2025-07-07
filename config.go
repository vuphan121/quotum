package quotum

import (
	"github.com/vuphan121/quotum/storage"
	"time"
)

// AppConfig contains shared limiter configuration
var AppConfig = struct {
	Rate     int
	Interval time.Duration
	Store    storage.Storage
}{
	Rate:     5,
	Interval: 1 * time.Second,
	Store:    storage.SharedMemoryStorage,
}

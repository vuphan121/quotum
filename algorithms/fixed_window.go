package algorithms

import (
	"github.com/vuphan121/quotum/storage"
	"time"
)

func FixedWindow(key string, rate int, interval time.Duration, store storage.Storage) bool {
	now := time.Now()

	state, err := store.GetState(key)
	if err != nil {
		state = storage.LimiterState{
			WindowStart:  now,
			RequestCount: 0,
		}
	}

	if now.Sub(state.WindowStart) >= interval {
		state.WindowStart = now
		state.RequestCount = 1
		store.SetState(key, state)
		return true
	}

	if state.RequestCount < rate {
		state.RequestCount++
		store.SetState(key, state)
		return true
	}

	return false
}

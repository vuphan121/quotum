package algorithms

import (
	"github.com/vuphan121/quotum/storage"
	"github.com/vuphan121/quotum/util"
	"time"
)

func FixedWindow(key string, rate int, interval time.Duration, store storage.Storage) bool {
	now := time.Now()
	windowStart := now.Truncate(interval)

	state, err := store.GetState(key)
	if err != nil || !state.WindowStart.Equal(windowStart) {
		newState := storage.LimiterState{
			WindowStart:  windowStart,
			RequestCount: 1,
		}
		_ = store.SetState(key, newState)

		util.Log("[RateLimiter] [%s] New window started at %v. Request allowed (1/%d).\n", key, windowStart, rate)
		return true
	}

	if state.RequestCount < rate {
		state.RequestCount++
		_ = store.SetState(key, state)

		util.Log("[RateLimiter] [%s] Request allowed (%d/%d).\n", key, state.RequestCount, rate)
		return true
	}

	util.Log("[RateLimiter] [%s] Request rejected. Window started at %v (%d/%d).\n", key, state.WindowStart, state.RequestCount, rate)
	return false
}

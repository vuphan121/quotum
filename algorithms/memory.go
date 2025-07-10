package algorithms

import (
	"github.com/vuphan121/quotum/storage"
	"time"
)

func FixedWindow(key string, rate int, interval time.Duration, store storage.Storage) (bool, *time.Time) {
	now := time.Now()
	windowStart := now.Truncate(interval)

	state, err := store.GetState(key)
	if err != nil {
		state = storage.LimiterState{
			WindowStart:  windowStart,
			RequestCount: 0,
			BannedUntil:  nil,
		}
	}

	// If user is banned, deny until ban expires
	if state.BannedUntil != nil && now.Before(*state.BannedUntil) {
		return false, state.BannedUntil
	}

	// If we're in a new window, reset count and window start
	if state.WindowStart.Before(windowStart) {
		state.WindowStart = windowStart
		state.RequestCount = 0
		state.BannedUntil = nil
	}

	// Allow if within rate limit
	if state.RequestCount < rate {
		state.RequestCount++
		state.BannedUntil = nil
		store.SetState(key, state)
		return true, nil
	}

	// Exceeded rate limit → ban until end of this window
	bannedUntil := windowStart.Add(interval)
	state.BannedUntil = &bannedUntil
	store.SetState(key, state)
	return false, &bannedUntil
}

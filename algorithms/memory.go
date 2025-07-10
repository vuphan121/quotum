package algorithms

import (
	"github.com/vuphan121/quotum/storage"
	"github.com/vuphan121/quotum/util"
	"time"
)

func FixedWindow(key string, rate int, interval time.Duration, store storage.Storage) (bool, *time.Time) {
	now := time.Now()
	windowStart := now.Truncate(interval)

	state, err := store.GetState(key)
	if err != nil {
		util.Log("[ERROR] Failed to get state for %s: %v\n", key, err)
		state = storage.LimiterState{
			WindowStart:  windowStart,
			RequestCount: 0,
			BannedUntil:  nil,
		}
	}

	util.Log("[DEBUG] %s: count=%d, start=%v, now=%v\n", key, state.RequestCount, state.WindowStart, now)

	if state.BannedUntil != nil && now.Before(*state.BannedUntil) {
		util.Log("[DEBUG] %s is banned until %v\n", key, *state.BannedUntil)
		return false, state.BannedUntil
	}

	if state.WindowStart.Before(windowStart) {
		util.Log("[DEBUG] New window. Resetting count for %s\n", key)
		state.WindowStart = windowStart
		state.RequestCount = 0
		state.BannedUntil = nil
	}

	if state.RequestCount < rate {
		util.Log("[DEBUG] Allowed. Incremented count to %d for %s\n", state.RequestCount, key)
		state.RequestCount++
		state.BannedUntil = nil
		store.SetState(key, state)
		return true, nil
	}

	bannedUntil := windowStart.Add(interval)
	state.BannedUntil = &bannedUntil
	store.SetState(key, state)
	util.Log("[DEBUG] Blocked. %s banned until %v\n", key, bannedUntil)
	return false, &bannedUntil
}

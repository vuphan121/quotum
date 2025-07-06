package quotum

import (
	"quotum/storage"
	"time"
)

type Limiter struct {
	rate     int
	interval time.Duration
	storage  storage.Storage
	key      string
}

func (l *Limiter) Allow() bool {
	state, _ := l.storage.GetState(l.key)

	now := time.Now()
	if now.Sub(state.WindowStart) >= l.interval {
		state.WindowStart = now
		state.RequestCount = 0
	}

	if state.RequestCount < l.rate {
		state.RequestCount++
		l.storage.SetState(l.key, state)
		return true
	}

	l.storage.SetState(l.key, state)
	return false
}

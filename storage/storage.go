package storage

import "time"

type LimiterState struct {
	RequestCount int
	WindowStart  time.Time
}

type Storage interface {
	GetState(key string) (LimiterState, error)
	SetState(key string, state LimiterState) error
}

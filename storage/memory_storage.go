package storage

import (
	"sync"
	"time"
)

type MemoryStorage struct {
	data map[string]LimiterState
	mu   sync.Mutex
}

var SharedMemoryStorage = NewMemoryStorage()

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[string]LimiterState),
	}
}

func (m *MemoryStorage) GetState(key string) (LimiterState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, exists := m.data[key]
	if !exists {
		state = LimiterState{
			RequestCount: 0,
			WindowStart:  time.Now(),
		}
		m.data[key] = state
	}
	return state, nil
}

func (m *MemoryStorage) SetState(key string, state LimiterState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = state
	return nil
}

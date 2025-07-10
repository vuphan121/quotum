package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(client *redis.Client) *RedisStorage {
	return &RedisStorage{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisStorage) GetState(key string) (LimiterState, error) {
	val, err := r.client.Get(r.ctx, key).Result()

	if err == redis.Nil {
		return LimiterState{
			RequestCount: 0,
			WindowStart:  time.Now(),
			BannedUntil:  nil,
		}, nil
	} else if err != nil {
		return LimiterState{}, err
	}

	var state LimiterState
	err = json.Unmarshal([]byte(val), &state)
	if err != nil {
		return LimiterState{}, err
	}
	return state, nil
}

func (r *RedisStorage) SetState(key string, state LimiterState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, key, data, 0).Err()
}

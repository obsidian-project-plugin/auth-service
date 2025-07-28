package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

type Provider[T any] struct {
	Client *redis.Client
	TTL    time.Duration
}

func NewProvider[T any](client *redis.Client, ttl time.Duration) *Provider[T] {
	return &Provider[T]{Client: client, TTL: ttl}
}

func (r *Provider[T]) WithCache(ctx context.Context, key string, block func() (T, error)) (T, error) {
	var zero T
	val, err := r.Client.Get(ctx, key).Result()
	if err == nil {
		var result T
		if err := json.Unmarshal([]byte(val), &result); err == nil {
			return result, nil
		}
	}
	if err != redis.Nil && err != nil {
		return zero, err
	}
	result, err := block()
	if err != nil {
		return zero, err
	}
	b, _ := json.Marshal(result)
	_ = r.Client.Set(ctx, key, b, r.TTL).Err()
	return result, nil
}

func (r *Provider[T]) MGet(ctx context.Context, keys []string) (map[string]*T, error) {
	res, err := r.Client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}
	result := make(map[string]*T, len(keys))
	for i, raw := range res {
		if raw == nil {
			result[keys[i]] = nil
			continue
		}
		var v T
		if err := json.Unmarshal([]byte(raw.(string)), &v); err == nil {
			result[keys[i]] = &v
		} else {
			result[keys[i]] = nil
		}
	}
	return result, nil
}

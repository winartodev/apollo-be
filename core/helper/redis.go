package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisUtil struct {
	client *redis.Client
}

func NewRedisUtil(client *redis.Client) (*RedisUtil, error) {
	return &RedisUtil{
		client: client,
	}, nil
}

func (r *RedisUtil) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

func (r *RedisUtil) Close() error {
	return r.client.Close()
}

// SetEx sets key-value with expiration, automatically marshalling the value
func (r *RedisUtil) SetEx(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("redisutil: failed to marshal value: %w", err)
	}

	return r.client.SetEx(ctx, key, jsonData, expiration).Err()
}

// Get gets a value by key and unmarshal it into the destination
func (r *RedisUtil) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("redisutil: failed to unmarshal data for key %s: %w", key, err)
	}

	return nil
}

// Delete removes one or more keys
func (r *RedisUtil) Delete(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (r *RedisUtil) Exists(ctx context.Context, key string) (bool, error) {
	count, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RedisUtil) IncrBy(ctx context.Context, key string, increment int64) (int64, error) {
	return r.client.IncrBy(ctx, key, increment).Result()
}

func (r *RedisUtil) DecrBy(ctx context.Context, key string, decrement int64) (int64, error) {
	return r.client.DecrBy(ctx, key, -decrement).Result()
}

func (r *RedisUtil) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.client.Expire(ctx, key, expiration).Err()
}

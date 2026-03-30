// Package cache provides a Redis-backed implementation of ports.Cache.
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache implements ports.Cache using Redis.
type RedisCache struct {
	client *redis.Client
}

// New creates a new RedisCache connected to the given URL.
func New(url string) (*RedisCache, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("cache: parse url: %w", err)
	}
	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("cache: ping redis: %w", err)
	}
	return &RedisCache{client: client}, nil
}

// Get retrieves the value of key. Returns ("", nil) if the key does not exist.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("cache: get %q: %w", key, err)
	}
	return val, nil
}

// Set stores key → value with the given TTL (0 = no expiry).
func (r *RedisCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if err := r.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("cache: set %q: %w", key, err)
	}
	return nil
}

// Delete removes one or more keys from the cache.
func (r *RedisCache) Delete(ctx context.Context, keys ...string) error {
	if err := r.client.Del(ctx, keys...).Err(); err != nil {
		return fmt.Errorf("cache: delete: %w", err)
	}
	return nil
}

// Exists reports whether a key exists.
func (r *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("cache: exists %q: %w", key, err)
	}
	return n > 0, nil
}

// Close closes the underlying Redis connection.
func (r *RedisCache) Close() error {
	return r.client.Close()
}

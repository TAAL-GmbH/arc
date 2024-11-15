package cache

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisStore is an implementation of CacheStore using Redis.
type RedisStore struct {
	client redis.UniversalClient
	ctx    context.Context
}

// NewRedisStore initializes a RedisStore.
func NewRedisStore(ctx context.Context, c redis.UniversalClient) *RedisStore {
	return &RedisStore{
		client: c,
		ctx:    ctx,
	}
}

// Get retrieves a value by key and hash (if given).
func (r *RedisStore) Get(hash *string, key string) ([]byte, error) {
	var result string
	var err error

	if hash == nil {
		result, err = r.client.Get(r.ctx, key).Result()
	} else {
		result, err = r.client.HGet(r.ctx, *hash, key).Result()
	}

	if errors.Is(err, redis.Nil) {
		return nil, ErrCacheNotFound
	} else if err != nil {
		return nil, errors.Join(ErrCacheFailedToGet, err)
	}

	return []byte(result), nil
}

// Set stores a value with a TTL for a specific hash (if given).
func (r *RedisStore) Set(hash *string, key string, value []byte, ttl time.Duration) error {
	var err error

	if hash == nil {
		err = r.client.Set(r.ctx, key, value, ttl).Err()
	} else {
		err = r.client.HSet(r.ctx, *hash, key, value).Err()
	}

	if err != nil {
		return errors.Join(ErrCacheFailedToSet, err)
	}

	return nil
}

// Del removes a value by key.
func (r *RedisStore) Del(hash *string, keys ...string) error {
	var result int64
	var err error

	if hash == nil {
		result, err = r.client.Del(r.ctx, keys...).Result()
	} else {
		result, err = r.client.HDel(r.ctx, *hash, keys...).Result()
	}

	if err != nil {
		return errors.Join(ErrCacheFailedToDel, err)
	}
	if result == 0 {
		return ErrCacheNotFound
	}
	return nil
}

// GetAllWithPrefix retrieves all key-value pairs that match a specific prefix.
func (r *RedisStore) GetAllWithPrefix(prefix string) (map[string][]byte, error) {
	var cursor uint64
	results := make(map[string][]byte)

	for {
		keys, newCursor, err := r.client.Scan(r.ctx, cursor, prefix+"*", 10).Result()
		if err != nil {
			return nil, errors.Join(ErrCacheFailedToScan, err)
		}

		for _, key := range keys {
			value, err := r.client.Get(r.ctx, key).Result()
			if errors.Is(err, redis.Nil) {
				// Key has been removed between SCAN and GET, skip it
				continue
			} else if err != nil {
				return nil, errors.Join(ErrCacheFailedToGet, err)
			}
			results[key] = []byte(value)
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}
	return results, nil
}

// GetAllForHash retrieves all key-value pairs for a specific hash.
func (r *RedisStore) GetAllForHash(hash string) (map[string][]byte, error) {
	values, err := r.client.HGetAll(r.ctx, hash).Result()
	if err != nil {
		return nil, errors.Join(ErrCacheFailedToGet, err)
	}

	result := make(map[string][]byte)
	for field, value := range values {
		result[field] = []byte(value)
	}
	return result, nil
}

// CountElementsForHash returns the number of elements in a hash.
func (r *RedisStore) CountElementsForHash(hash string) (int64, error) {
	count, err := r.client.HLen(r.ctx, hash).Result()
	if err != nil {
		return 0, errors.Join(ErrCacheFailedToGetCount, err)
	}
	return count, nil
}

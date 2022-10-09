package gin_cache_redis

import (
	"context"
	"errors"
	"github.com/chenyahui/gin-cache/persist"
	redisv9 "github.com/go-redis/redis/v9"
	"time"
)

type RedisV9Store struct {
	RedisClient *redisv9.Client
}

// NewRedisV9Store create a redis memory store with redis client
func NewRedisV9Store(redisClient *redisv9.Client) *RedisV9Store {
	return &RedisV9Store{
		RedisClient: redisClient,
	}
}

// Set put key value pair to redis, and expire after expireDuration
func (store *RedisV9Store) Set(key string, value interface{}, expire time.Duration) error {
	payload, err := persist.Serialize(value)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return store.RedisClient.Set(ctx, key, payload, expire).Err()
}

// Delete remove key in redis, do nothing if key doesn't exist
func (store *RedisV9Store) Delete(key string) error {
	ctx := context.Background()
	return store.RedisClient.Del(ctx, key).Err()
}

// Get retrieves an item from redis, if key doesn't exist, return ErrCacheMiss
func (store *RedisV9Store) Get(key string, value interface{}) error {
	ctx := context.TODO()
	payload, err := store.RedisClient.Get(ctx, key).Bytes()

	if errors.Is(err, redisv9.Nil) {
		return persist.ErrCacheMiss
	}

	if err != nil {
		return err
	}
	return persist.Deserialize(payload, value)
}

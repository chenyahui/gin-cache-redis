package gin_cache_redis

import (
	"context"
	"errors"
	"time"

	"github.com/chenyahui/gin-cache/persist"
	redisv9 "github.com/redis/go-redis/v9"
)

type RedisClusterStore struct {
	RedisClient *redisv9.ClusterClient `inject:""`
}

// NewRedisClusterStore create a redis memory store with redis client
func NewRedisClusterStore(redisClient *redisv9.ClusterClient) *RedisClusterStore {
	return &RedisClusterStore{
		RedisClient: redisClient,
	}
}

// Set put key value pair to redis, and expire after expireDuration
func (store *RedisClusterStore) Set(key string, value interface{}, expire time.Duration) error {
	payload, err := persist.Serialize(value)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return store.RedisClient.Set(ctx, key, payload, expire).Err()
}

// Delete remove key in redis, do nothing if key doesn't exist
func (store *RedisClusterStore) Delete(key string) error {
	ctx := context.Background()
	return store.RedisClient.Del(ctx, key).Err()
}

// Get retrieves an item from redis, if key doesn't exist, return ErrCacheMiss
func (store *RedisClusterStore) Get(key string, value interface{}) error {
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

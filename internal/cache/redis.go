package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/MostafaSensei106/Riko-Chan/config"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(cfg config.RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return r.client.Set(ctx, key, data, expiration).Err()
}

func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

func (r *RedisClient) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) ZAdd(ctx context.Context, key string, score float64, member interface{}) error {
	data, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("failed to marshal member: %w", err)
	}

	return r.client.ZAdd(ctx, key, &redis.Z{
		Score:  score,
		Member: string(data),
	}).Err()
}

func (r *RedisClient) ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error) {
	return r.client.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min: min,
		Max: max,
	}).Result()
}

func (r *RedisClient) ZRem(ctx context.Context, key string, member interface{}) error {
	data, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("failed to marshal member: %w", err)
	}

	return r.client.ZRem(ctx, key, string(data)).Err()
}

func (r *RedisClient) GetClient() *redis.Client {
	return r.client
}

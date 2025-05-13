package database

import (
	"context"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisStore struct {
	client *redis.Client
}

type RedisConfig struct {
	Addr     string
	Password string
	Db       int
}

func NewRedisStore(cfg *RedisConfig) *RedisStore {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Db,
	})

	return &RedisStore{
		client: client,
	}
}

func (r *RedisStore) Incr(ctx context.Context, key string, expirationSeconds int) (int, error) {
	pipe := r.client.TxPipeline()

	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, time.Duration(expirationSeconds)*time.Second)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return int(incr.Val()), nil
}

func (r *RedisStore) Get(ctx context.Context, key string) (int, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}

	count, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RedisStore) Set(ctx context.Context, key string, value int, expirationSeconds int) error {
	return r.client.Set(ctx, key, value, time.Duration(expirationSeconds)*time.Second).Err()
}

func (r *RedisStore) TTL(ctx context.Context, key string) (int, error) {
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, err
	}
	return int(ttl.Seconds()), nil
}

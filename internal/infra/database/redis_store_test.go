package database

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRedisStore(t *testing.T) {
	cfg := &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		Db:       0,
	}
	store := NewRedisStore(cfg)
	store.client.FlushDB(context.Background())

	t.Run("SetAndGet", func(t *testing.T) {
		err := store.Set(context.Background(), "testKey", 123, 10)
		require.NoError(t, err)

		val, err := store.Get(context.Background(), "testKey")
		require.NoError(t, err)
		require.Equal(t, 123, val)
	})

	t.Run("Incr", func(t *testing.T) {
		val, err := store.Incr(context.Background(), "testKey", 10)
		require.NoError(t, err)
		require.Equal(t, 124, val)
	})

	t.Run("TTL", func(t *testing.T) {
		ttl, err := store.TTL(context.Background(), "testKey")
		require.NoError(t, err)
		require.True(t, ttl > 0)
	})

	t.Run("GetNonexistentKey", func(t *testing.T) {
		val, err := store.Get(context.Background(), "doesNotExist")
		require.NoError(t, err)
		require.Equal(t, 0, val)
	})

	t.Run("SetExpiration", func(t *testing.T) {
		err := store.Set(context.Background(), "expireKey", 999, 1)
		require.NoError(t, err)
		time.Sleep(2 * time.Second)
		val, err := store.Get(context.Background(), "expireKey")
		require.NoError(t, err)
		require.Equal(t, 0, val)
	})
}

func TestRedisStore_NewRedisStore(t *testing.T) {
	cfg := &RedisConfig{
		Addr:     "localhost:6379",
		Password: "",
		Db:       0,
	}
	store := NewRedisStore(cfg)
	require.NotNil(t, store.client)

	_, err := store.client.Ping(context.Background()).Result()
	require.NoError(t, err)
}

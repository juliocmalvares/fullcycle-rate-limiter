package limiter

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MemoryStore struct {
	data  map[string]item
	mutex sync.Mutex
}

type item struct {
	value     int
	expiresAt time.Time
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]item),
	}
}

func (m *MemoryStore) Incr(ctx context.Context, key string, expirationSeconds int) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	now := time.Now()
	val, ok := m.data[key]

	if !ok || time.Now().After(val.expiresAt) {
		val = item{value: 0, expiresAt: now.Add(time.Duration(expirationSeconds) * time.Second)}
	}

	val.value++
	m.data[key] = val
	return val.value, nil
}

func (m *MemoryStore) Get(ctx context.Context, key string) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val, ok := m.data[key]
	if !ok || time.Now().After(val.expiresAt) {
		return 0, nil
	}
	return val.value, nil
}

func (m *MemoryStore) Set(ctx context.Context, key string, value int, expirationSeconds int) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.data[key] = item{
		value:     value,
		expiresAt: time.Now().Add(time.Duration(expirationSeconds) * time.Second),
	}
	return nil
}

func (m *MemoryStore) TTL(ctx context.Context, key string) (int, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val, ok := m.data[key]
	if !ok || time.Now().After(val.expiresAt) {
		return 0, nil
	}
	return int(time.Until(val.expiresAt).Seconds()), nil
}

func TestLimiter_ByIP(t *testing.T) {
	store := NewMemoryStore()
	limiter := NewLimiter(store, 5, 2) // 5 reqs, 2s window

	ctx := context.Background()
	ip := "192.168.0.1"

	for i := 0; i < 5; i++ {
		res, err := limiter.Check(ctx, ip, "")
		assert.NoError(t, err)
		assert.True(t, res.Allowed)
	}

	// 6th request should fail
	res, err := limiter.Check(ctx, ip, "")
	assert.NoError(t, err)
	assert.False(t, res.Allowed)
	assert.Equal(t, "you have reached the maximum number of requests or actions allowed within a certain time frame", res.Message)
}

func TestLimiter_ByTokenWithCustomLimit(t *testing.T) {
	store := NewMemoryStore()
	limiter := NewLimiter(store, 2, 2) // defaults: 2 reqs, 2s

	ctx := context.Background()
	token := "abc123"

	// Set custom config for token
	_ = store.Set(ctx, "config:token:abc123:limit", 3, 60)
	_ = store.Set(ctx, "config:token:abc123:expiration", 5, 60)

	for i := 0; i < 3; i++ {
		res, err := limiter.Check(ctx, "", token)
		assert.NoError(t, err)
		assert.True(t, res.Allowed)
	}

	// 4th request must be blocked
	res, err := limiter.Check(ctx, "", token)
	assert.NoError(t, err)
	assert.False(t, res.Allowed)
}

package database

import "context"

type DatabaseStore interface {
	Incr(ctx context.Context, key string, expirationSeconds int) (int, error)
	Get(ctx context.Context, key string) (int, error)
	TTL(ctx context.Context, key string) (int, error)
	Set(ctx context.Context, key string, value int, expirationSeconds int) error
}

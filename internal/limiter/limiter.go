package limiter

import (
	"context"
	"fmt"
	"posgoexpert-rate-limiter/internal/infra/database"
	"posgoexpert-rate-limiter/internal/logger"
)

type Limiter struct {
	DbStore           database.DatabaseStore
	DefaultLimit      int
	DefaultExpiration int
}

func NewLimiter(dbStore database.DatabaseStore, defaultLimit, defaultExpiration int) *Limiter {
	return &Limiter{
		DbStore:           dbStore,
		DefaultLimit:      defaultLimit,
		DefaultExpiration: defaultExpiration,
	}
}

type CheckResult struct {
	Allowed    bool
	Remaining  int
	TTLSeconds int
	Message    string
}

func (l *Limiter) Check(ctx context.Context, ip string, token string) (*CheckResult, error) {
	var key, configLimitKey, configExpKey string
	var limit, expiration int
	var err error
	logger.Logger.Info(fmt.Sprintf("Checking rate limit for IP: [%s], Token: [%s]", ip, token))
	if token != "" {
		key = fmt.Sprintf("limit:token:%s", token)
		configLimitKey = fmt.Sprintf("config:token:%s:limit", token)
		configExpKey = fmt.Sprintf("config:token:%s:expiration", token)

		limit, err = l.DbStore.Get(ctx, configLimitKey)
		if err != nil {
			return nil, err
		}
		if limit == 0 {
			limit = l.DefaultLimit
		}
		expiration, err = l.DbStore.Get(ctx, configExpKey)
		if err != nil {
			return nil, err
		}
		if expiration == 0 {
			expiration = l.DefaultExpiration
		}
		logger.Logger.Info(fmt.Sprintf("Rate limit for token [%s]: Limit [%d], Expiration [%d]", token, limit, expiration))
	} else if ip != "" {
		key = fmt.Sprintf("limit:ip:%s", ip)
		configLimitKey = fmt.Sprintf("config:ip:%s:limit", ip)
		configExpKey = fmt.Sprintf("config:ip:%s:expiration", ip)

		limit, err = l.DbStore.Get(ctx, configLimitKey)
		if err != nil {
			return nil, err
		}
		expiration, err = l.DbStore.Get(ctx, configExpKey)
		if err != nil {
			return nil, err
		}
		if limit == 0 {
			limit = l.DefaultLimit
		}
		if expiration == 0 {
			expiration = l.DefaultExpiration
		}
		logger.Logger.Info(fmt.Sprintf("Rate limit for IP [%s]: Limit [%d]", ip, limit))
	}

	count, err := l.DbStore.Incr(ctx, key, expiration)
	if err != nil {
		return nil, err
	}
	ttl, err := l.DbStore.TTL(ctx, key)
	if err != nil {
		return nil, err
	}
	remaining := limit - count
	logger.Logger.Info(fmt.Sprintf("Rate limit check: Key [%s], Count [%d], Remaining [%d], TTL [%d]", key, count, remaining, ttl))

	if remaining < 0 {
		return &CheckResult{
			Allowed:    false,
			Remaining:  0,
			TTLSeconds: ttl,
			Message:    "you have reached the maximum number of requests or actions allowed within a certain time frame",
		}, nil
	}

	return &CheckResult{
		Allowed:    true,
		Remaining:  remaining,
		TTLSeconds: ttl,
	}, nil
}

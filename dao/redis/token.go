package redis

import (
	"context"
	"time"
)

func BlacklistRefreshToken(ctx context.Context, jti string, expiresAt time.Time) error {
	if jti == "" {
		return nil
	}

	ttl := time.Until(expiresAt)
	if ttl <= 0 {
		return nil
	}

	key := getRedisKey(KeyRefreshTokenBlacklistPF + jti)
	return rdb.Set(ctx, key, "1", ttl).Err()
}

func IsRefreshTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	if jti == "" {
		return false, nil
	}

	key := getRedisKey(KeyRefreshTokenBlacklistPF + jti)
	n, err := rdb.Exists(ctx, key).Result()
	return n > 0, err
}

package repo

import (
	"CampusTake/internal/constants"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepo interface {
	SetBlacklist(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}

type tokenRepo struct {
	rdb *redis.Client
}

func NewTokenRepo(rdb *redis.Client) TokenRepo {
	return &tokenRepo{rdb: rdb}
}

func buildTokenBlacklistKey(token string) string {
	return fmt.Sprintf(constants.RedisKeyPrefixTokenBlacklist, token)
}

func (t *tokenRepo) SetBlacklist(ctx context.Context, token string, ttl time.Duration) error {
	key := buildTokenBlacklistKey(token)
	return t.rdb.Set(ctx, key, "1", ttl).Err()
}

func (t *tokenRepo) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	key := buildTokenBlacklistKey(token)
	val, err := t.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, nil
}

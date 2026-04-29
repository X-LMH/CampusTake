package repo

import (
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type VerifyCodeRepo interface {
	SetCode(ctx context.Context, phone, code string, ttl time.Duration) error
	GetCode(ctx context.Context, phone string) (string, error)
	DeleteCode(ctx context.Context, phone string) error
}

type verifyCodeRepo struct {
	rdb *redis.Client
}

func buildVerifyCodeKey(phone string) string {
	return fmt.Sprintf(enum.RedisKeyPrefixVerifyCode, phone)
}

func (v *verifyCodeRepo) SetCode(ctx context.Context, phone, code string, ttl time.Duration) error {
	key := buildVerifyCodeKey(phone)
	return v.rdb.Set(ctx, key, code, ttl).Err()
}

func (v *verifyCodeRepo) GetCode(ctx context.Context, phone string) (string, error) {
	key := buildVerifyCodeKey(phone)

	code, err := v.rdb.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", errx.ErrVerifyCodeNotFound
		}
		return "", err
	}

	return code, nil
}

func (v *verifyCodeRepo) DeleteCode(ctx context.Context, phone string) error {
	key := buildVerifyCodeKey(phone)

	res, err := v.rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	if res == 0 {
		return errx.ErrVerifyCodeNotFound
	}

	return nil
}

package impl

import (
	"CampusTake/internal/constants"
	errs "CampusTake/pkg/errors"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type VerifyCodeRepo interface {
	SetCode(ctx context.Context, phone, code string, ttl time.Duration) error
	GetCode(ctx context.Context, phone string) (string, error)
	DeleteCode(ctx context.Context, phone string) error
	SetVerifyToken(ctx context.Context, verifyToken constants.VerifyTokenType, ttl time.Duration) (string, error)
	GetVerifyToken(ctx context.Context, key string) (constants.VerifyTokenType, error)
	DeleteVerifyToken(ctx context.Context, key string) error
}

type verifyCodeRepo struct {
	RepoBase
}

func NewVerifyCodeRepo(db *gorm.DB, rdb redis.Cmdable) VerifyCodeRepo {
	return &verifyCodeRepo{
		RepoBase: RepoBase{
			db:  db,
			rdb: rdb,
		},
	}
}

func buildVerifyCodeKey(phone string) string {
	return fmt.Sprintf(constants.RedisKeyPrefixVerifyCode, phone)
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
			return "", errs.ErrVerifyCodeNotFound
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
		return errs.ErrVerifyCodeNotFound
	}

	return nil
}

func (v *verifyCodeRepo) SetVerifyToken(ctx context.Context, verifyToken constants.VerifyTokenType, ttl time.Duration) (string, error) {
	key := verifyToken.GenerateKey()

	// ✅ 修改点2：将 struct 序列化成 JSON
	data, err := json.Marshal(verifyToken)
	if err != nil {
		return "", err
	}

	// ❗修改点3：存入的是 []byte / string，而不是 struct
	err = v.rdb.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return "", err
	}

	// ✅ 修改点4（建议优化）：返回 token 而不是完整 key（可选）
	// 当前先保持你原逻辑不动
	return key, nil
}

func (v *verifyCodeRepo) GetVerifyToken(ctx context.Context, key string) (constants.VerifyTokenType, error) {
	var verifyToken constants.VerifyTokenType

	// ✅ 修改点5：用 Bytes() 获取原始数据（而不是 Scan）
	data, err := v.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return constants.VerifyTokenType{}, errs.ErrVerifyTokenNotFound
		}
		return constants.VerifyTokenType{}, err
	}

	// ✅ 修改点6：反序列化 JSON → struct
	err = json.Unmarshal(data, &verifyToken)
	if err != nil {
		return constants.VerifyTokenType{}, err
	}

	return verifyToken, nil
}

func (v *verifyCodeRepo) DeleteVerifyToken(ctx context.Context, key string) error {
	_, err := v.rdb.Del(ctx, key).Result()
	return err
}

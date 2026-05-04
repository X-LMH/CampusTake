package enum

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	RedisKeyPrefixVerifyCode     = "sms:code:%s"
	RedisKeyPrefixTokenBlacklist = "token:blacklist:%s"
	RedisKeyPrefixVerifyToken    = "verify:token:%s"
)

const (
	VerifyCodeTTL = 5 * time.Minute
	ResetTokenTTL = 10 * time.Minute
)

type VerifyTokenType struct {
	UserID int64
	Phone  string
}

func (v VerifyTokenType) GenerateKey() string {
	return fmt.Sprintf(RedisKeyPrefixVerifyToken, uuid.NewString())
}

package constants

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	RedisKeyPrefixVerifyCode         = "sms:code:%s"
	RedisKeyPrefixTokenBlacklist     = "token:blacklist:%s"
	RedisKeyPrefixVerifyToken        = "verify:token:%s"
	RedisKeyPrefixUser               = "user:%d"
	RedisKeyPrefixRiderProfileUserID = "rider:profile:user:%d"
	RedisKeyPrefixRiderProfileID     = "rider:profile:id:%d"
	RedisKeyPrefixOrderDetail        = "order:detail:%d"
	RedisKeyPrefixOrderGrabClaim     = "order:grab:claim:%d"
)

const (
	VerifyCodeTTL        = 5 * time.Minute
	ResetTokenTTL        = 10 * time.Minute
	UserCacheTTL         = 1 * time.Hour
	RiderProfileCacheTTL = 1 * time.Hour
	OrderDetailCacheTTL  = 15 * time.Minute
	OrderGrabClaimTTL    = 10 * time.Minute
)

type VerifyTokenType struct {
	UserID int64
	Phone  string
}

func (v VerifyTokenType) GenerateKey() string {
	return fmt.Sprintf(RedisKeyPrefixVerifyToken, uuid.NewString())
}

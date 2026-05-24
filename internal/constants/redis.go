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
	RedisKeyPrefixUserPhone          = "user:phone:%s"
	RedisKeyPrefixRiderProfileUserID = "rider:profile:user:%d"
	RedisKeyPrefixRiderProfileID     = "rider:profile:id:%d"
	RedisKeyPrefixOrderDetail        = "order:detail:%d"
	RedisKeyPrefixUserOrderList      = "order:list:user:%d:%d:%d:%d"
	RedisKeyPrefixAvailableOrderList = "order:list:available:%d:%d:%s:%s:%s:%s"
	RedisKeyPrefixOrderGrabClaim     = "order:grab:claim:%d"
	RedisKeyPrefixAddress            = "address:%d"
)

const (
	VerifyCodeTTL        = 5 * time.Minute
	ResetTokenTTL        = 10 * time.Minute
	UserCacheTTL         = 1 * time.Hour
	RiderProfileCacheTTL = 1 * time.Hour
	OrderDetailCacheTTL  = 15 * time.Minute
	OrderListCacheTTL    = 3 * time.Second
	OrderGrabClaimTTL    = 10 * time.Minute
	AddressCacheTTL      = 1 * time.Hour
)

type VerifyTokenType struct {
	UserID int64
	Phone  string
}

func (v VerifyTokenType) GenerateKey() string {
	return fmt.Sprintf(RedisKeyPrefixVerifyToken, uuid.NewString())
}

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

	EmptyOrderCacheValue = "__empty__"
	EmptyUserCacheValue  = "__empty__"

	OrderBloomKey = "bloom:order:id"
)

// 过期时间优化配置
const (
	VerifyCodeTTL        = 3 * time.Minute
	ResetTokenTTL        = 15 * time.Minute
	UserCacheTTL         = 2 * time.Hour
	RiderProfileCacheTTL = 2 * time.Hour
	OrderDetailCacheTTL  = 10 * time.Minute
	OrderListCacheTTL    = 2 * time.Second
	OrderGrabClaimTTL    = 2 * time.Minute
	AddressCacheTTL      = 2 * time.Hour
	OrderEmptyCacheTTL   = 1 * time.Minute
	UserEmptyCacheTTL    = 1 * time.Minute
	TokenBlacklistTTL    = 12 * time.Hour
)

type VerifyTokenType struct {
	UserID int64
	Phone  string
}

func (v VerifyTokenType) GenerateKey() string {
	return fmt.Sprintf(RedisKeyPrefixVerifyToken, uuid.NewString())
}

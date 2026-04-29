package enum

import "time"

const (
	RedisKeyPrefixVerifyCode     = "sms:code:%s"
	RedisKeyPrefixTokenBlacklist = "token:blacklist:%s"
)

const (
	VerifyCodeTTL = 5 * time.Minute
	ResetTokenTTL = 10 * time.Minute
)

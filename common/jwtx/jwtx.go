package jwtx

import (
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID int64         `json:"userID"`
	Role   enum.RoleType `json:"role,omitempty"`
	jwt.RegisteredClaims
}

type JwtConfig struct {
	SecretKey string
	Issuer    string
	Expire    time.Duration
}

// --- 底层通用逻辑 (私有，不对外暴露) ---

func generateGenericToken(cfg JwtConfig, userID int64, role enum.RoleType, subject string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.Issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.SecretKey))
}

// --- 上层语义化封装 (业务直接调用) ---

// GenerateToken 生成普通登录 Token (Access Token)
func GenerateToken(cfg JwtConfig, userID int64, role enum.RoleType) (string, error) {
	return generateGenericToken(cfg, userID, role, "access_token", cfg.Expire)
}

// GenerateResetToken 生成重置密码 Token (10分钟有效)
func GenerateResetToken(cfg JwtConfig, userID int64) (string, error) {
	return generateGenericToken(cfg, userID, 0, "reset_password_token", 10*time.Minute)
}

// GenerateChangePhoneToken 生成换绑手机 Token (5分钟有效)
// 这里就是你提到的未来可能用到的复用示例
func GenerateChangePhoneToken(cfg JwtConfig, userID int64) (string, error) {
	return generateGenericToken(cfg, userID, 0, "change_phone_token", 5*time.Minute)
}

// --- 解析逻辑 ---

// ParseToken 通用解析工具
func ParseToken(tokenString string, secretKey string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errx.NewParamError("token签名算法非法")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errx.NewCodeError(errx.TokenExpireError, "凭证已过期")
		}
		return nil, errx.NewCodeError(errx.TokenInvalidError, "凭证无效")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errx.NewCodeError(errx.TokenInvalidError, "凭证校验失败")
	}

	return claims, nil
}

// ParseResetToken 专门解析并验证重置密码凭证
func ParseResetToken(tokenString string, secretKey string) (*Claims, error) {
	claims, err := ParseToken(tokenString, secretKey)
	if err != nil {
		return nil, err
	}
	if claims.Subject != "reset_password_token" {
		return nil, errx.NewCodeError(errx.TokenInvalidError, "凭证用途错误")
	}
	return claims, nil
}

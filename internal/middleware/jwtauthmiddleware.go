// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package middleware

import (
	"CampusTake/internal/constants"
	errx2 "CampusTake/pkg/errors"
	"CampusTake/pkg/jwt"
	"CampusTake/pkg/response"
	"context"
	"net/http"
	"strings"
)

type JwtAuthMiddleware struct {
	SecretKey string
}

func NewJwtAuthMiddleware(secretKey string) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		SecretKey: secretKey,
	}
}

func (m *JwtAuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			response.Response(r, w, nil, errx2.NewCodeError(errx2.TokenMissingError, "未登录，请先登录"))
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Response(r, w, nil, errx2.NewCodeError(errx2.TokenInvalidError, "Authorization格式错误，应为 Bearer <token>"))
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := jwt.ParseToken(tokenString, m.SecretKey)
		if err != nil {
			response.Response(r, w, nil, err)
			return
		}

		ctx := context.WithValue(r.Context(), constants.CtxUserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, constants.CtxRoleKey, claims.Role)

		next(w, r.WithContext(ctx))
	}
}

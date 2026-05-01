package ctxx

import (
	"CampusTake/common/enum"
	"context"

	"CampusTake/common/errx"
)

func GetUserID(ctx context.Context) (int64, error) {
	v := ctx.Value(enum.CtxUserIDKey)
	userID, ok := v.(int64)
	if !ok || userID <= 0 {
		return 0, errx.NewCodeError(errx.TokenInvalidError, "登录状态无效")
	}
	return userID, nil
}

func MustUserID(ctx context.Context) int64 {
	userID, _ := GetUserID(ctx)
	return userID
}

func GetRole(ctx context.Context) enum.RoleType {
	v := ctx.Value(enum.CtxRoleKey)
	role, _ := v.(enum.RoleType)
	return role
}

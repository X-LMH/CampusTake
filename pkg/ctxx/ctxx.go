package ctxx

import (
	"CampusTake/internal/constants"
	"CampusTake/internal/enums"
	errx2 "CampusTake/pkg/errors"
	"context"
)

func GetUserID(ctx context.Context) (int64, error) {
	v := ctx.Value(constants.CtxUserIDKey)
	userID, ok := v.(int64)
	if !ok || userID <= 0 {
		return 0, errx2.NewCodeError(errx2.TokenInvalidError, "登录状态无效")
	}
	return userID, nil
}

func MustUserID(ctx context.Context) int64 {
	userID, _ := GetUserID(ctx)
	return userID
}

func GetRole(ctx context.Context) enums.RoleType {
	v := ctx.Value(constants.CtxRoleKey)
	role, _ := v.(enums.RoleType)
	return role
}

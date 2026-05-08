// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/internal/constants"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	errx2 "CampusTake/pkg/errors"
	"CampusTake/pkg/jwt"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordRequest) error {

	claims, err := jwt.ParseResetToken(req.ResetToken, l.svcCtx.JwtCfg.SecretKey)
	if err != nil {
		return err
	}

	// 2. 检查 Token 是否已在黑名单
	isBlack, err := l.svcCtx.Repo.Token().IsBlacklisted(l.ctx, req.ResetToken)
	if err != nil {
		return err
	}
	if isBlack {
		return errx2.NewCodeError(errx2.TokenInvalidError, "该重置链接已失效")
	}

	// 3. 修改数据库
	err = l.svcCtx.Repo.User().UpdatePasswordByID(l.ctx, claims.UserID, req.NewPassword)
	if err != nil {
		return err
	}

	// 4. 修改成功后，将 Token 拉黑
	_ = l.svcCtx.Repo.Token().SetBlacklist(l.ctx, req.ResetToken, constants.ResetTokenTTL)

	return nil
}

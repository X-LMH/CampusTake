// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"CampusTake/common/jwtx"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginWithVerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginWithVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginWithVerifyCodeLogic {
	return &LoginWithVerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginWithVerifyCodeLogic) LoginWithVerifyCode(req *types.LoginWithVerifyCodeRequest) (resp *types.LoginResponse, err error) {
	// 1. 验证验证码
	code, err := l.svcCtx.Repo.VerifyCode().GetCode(l.ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if code != req.VerifyCode {
		return nil, errx.ErrVerifyCodeWrong
	}

	user, err := l.svcCtx.Repo.User().GetByPhone(l.ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	// 用户被封禁
	if user.Status == enum.UserStatusDisabled {
		return nil, errx.ErrUserForbidden
	}

	token, err := jwtx.GenerateToken(l.svcCtx.JwtCfg, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	_ = l.svcCtx.Repo.VerifyCode().DeleteCode(l.ctx, req.Phone)

	return &types.LoginResponse{
		Token:    token,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

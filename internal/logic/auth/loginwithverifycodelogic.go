// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/internal/enums"
	"CampusTake/pkg/errors"
	"CampusTake/pkg/jwt"
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
	code, err := l.svcCtx.Repo.VerifyCode.GetCode(l.ctx, req.Phone)
	if err != nil {
		l.Errorf("验证码登录获取验证码失败，phone=%s，err=%v", req.Phone, err)
		return nil, err
	}
	if code != req.VerifyCode {
		l.Errorf("验证码登录验证码错误，phone=%s", req.Phone)
		return nil, errors.ErrVerifyCodeWrong
	}

	user, err := l.svcCtx.Repo.User.GetByPhone(l.ctx, req.Phone)
	if err != nil {
		l.Errorf("验证码登录查询用户失败，phone=%s，err=%v", req.Phone, err)
		return nil, err
	}

	// 用户被封禁
	if user.Status == enums.UserStatusDisabled {
		l.Errorf("验证码登录用户已被禁用，userID=%d，phone=%s", user.ID, req.Phone)
		return nil, errors.ErrUserForbidden
	}

	token, err := jwt.GenerateToken(l.svcCtx.JwtCfg, user.ID, user.Role)
	if err != nil {
		l.Errorf("验证码登录生成令牌失败，userID=%d，err=%v", user.ID, err)
		return nil, err
	}

	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, req.Phone)

	return &types.LoginResponse{
		Token:    token,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

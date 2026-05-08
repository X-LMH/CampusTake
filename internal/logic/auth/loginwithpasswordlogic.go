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

type LoginWithPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginWithPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginWithPasswordLogic {
	return &LoginWithPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginWithPasswordLogic) LoginWithPassword(req *types.LoginWithPasswordRequest) (*types.LoginResponse, error) {
	l.Logger.Debugf("LoginWithPassword request: %+v", req)
	// 根据手机号查询用户
	user, err := l.svcCtx.Repo.User().GetByPhone(l.ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	// 用户被封禁
	if user.Status == enums.UserStatusDisabled {
		return nil, errors.ErrUserForbidden
	}

	// 验证密码
	if user.Password != req.Password {
		return nil, errors.ErrPasswordWrong
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(l.svcCtx.JwtCfg, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Token:    token,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

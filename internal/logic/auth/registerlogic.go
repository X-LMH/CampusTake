// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"CampusTake/common/jwtx"
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"context"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (*types.RegisterResponse, error) {
	// 1. 验证验证码
	code, err := l.svcCtx.Repo.VerifyCode().GetCode(l.ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if code != req.VerifyCode {
		return nil, errx.ErrVerifyCodeWrong // ⚠ 这里建议不要用 PasswordWrong
	}

	// 2. 构造用户
	phoneSuffix := req.Phone
	if len(req.Phone) > 4 {
		phoneSuffix = req.Phone[len(req.Phone)-4:]
	}

	user := &model.User{
		Phone:    req.Phone,
		Password: req.Password,
		Nickname: "用户" + phoneSuffix,
		Avatar:   "base_avatar.png",
		Role:     enum.RoleUser,
	}

	// 3. 创建用户
	err = l.svcCtx.Repo.User().Create(l.ctx, user)
	if err != nil {
		if errors.Is(err, errx.ErrUserExist) {
			return nil, errx.ErrUserExist
		}
		return nil, err
	}

	// 5. 生成 token
	token, err := jwtx.GenerateToken(l.svcCtx.JwtCfg, user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	_ = l.svcCtx.Repo.VerifyCode().DeleteCode(l.ctx, req.Phone)

	return &types.RegisterResponse{
		Token:    token,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/errors"
	"CampusTake/pkg/jwt"
	"context"

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
	code, err := l.svcCtx.Repo.VerifyCode.GetCode(l.ctx, req.Phone)
	if err != nil {
		l.Errorf("注册获取验证码失败，phone=%s，err=%v", req.Phone, err)
		return nil, err
	}
	if code != req.VerifyCode {
		l.Errorf("注册验证码错误，phone=%s", req.Phone)
		return nil, errors.ErrVerifyCodeWrong // 这里建议不要用 PasswordWrong
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
		Avatar:   "default.png",
		Role:     enums.RoleUser,
		Status:   enums.UserStatusNormal,
	}

	// 3. 创建用户
	err = l.svcCtx.Repo.User.Create(l.ctx, user)
	if err != nil {
		l.Errorf("创建用户失败，phone=%s，err=%v", req.Phone, err)
		return nil, err
	}

	// 5. 生成 token
	token, err := jwt.GenerateToken(l.svcCtx.JwtCfg, user.ID, user.Role)
	if err != nil {
		l.Errorf("注册生成令牌失败，userID=%d，err=%v", user.ID, err)
		return nil, err
	}

	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, req.Phone)

	return &types.RegisterResponse{
		Token:    token,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}

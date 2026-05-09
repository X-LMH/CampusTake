// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/errors"
	"CampusTake/pkg/jwt"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type ForgetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewForgetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ForgetPasswordLogic {
	return &ForgetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ForgetPasswordLogic) ForgetPassword(req *types.ForgetPasswordRequest) (resp *types.ForgetPasswordResponse, err error) {
	// 1. 验证验证码
	code, err := l.svcCtx.Repo.VerifyCode.GetCode(l.ctx, req.Phone)
	if err != nil {
		return nil, err
	}
	if code != req.VerifyCode {
		return nil, errors.ErrVerifyCodeWrong
	}

	// 2. 获取用户信息
	user, err := l.svcCtx.Repo.User.GetByPhone(l.ctx, req.Phone)
	if err != nil {
		return nil, err
	}

	// 3. 设置忘记密码token
	resetToken, err := jwt.GenerateResetToken(l.svcCtx.JwtCfg, user.ID)
	if err != nil {
		return nil, err
	}

	// 4. 删除验证码
	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, req.Phone)

	return &types.ForgetPasswordResponse{
		ResetToken: resetToken,
	}, nil
}

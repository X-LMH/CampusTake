// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package phone

import (
	"CampusTake/internal/constants"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyNewPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewVerifyNewPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyNewPhoneLogic {
	return &VerifyNewPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *VerifyNewPhoneLogic) VerifyNewPhone(req *types.VerifyNewPhoneRequest) (resp *types.VerifyNewPhoneResponse, err error) {
	// ❗修改点1：建议先校验手机号是否已注册
	_, err = l.svcCtx.Repo.User.GetByPhone(l.ctx, req.NewPhone)
	if err == nil {
		l.Errorf("新手机号已被绑定，phone=%s", req.NewPhone)
		return nil, errors.ErrPhoneAlreadyBound
	}

	// 1. 验证码错误
	code, err := l.svcCtx.Repo.VerifyCode.GetCode(l.ctx, req.NewPhone)
	if err != nil {
		l.Errorf("获取新手机号验证码失败，phone=%s，err=%v", req.NewPhone, err)
		return nil, err
	}
	if code != req.Code {
		l.Errorf("新手机号验证码错误，phone=%s", req.NewPhone)
		return nil, errors.ErrVerifyCodeWrong
	}

	// 2. 生成verify token
	userID := ctxx.MustUserID(l.ctx)
	verifyToken := constants.VerifyTokenType{
		UserID: userID,
		Phone:  req.NewPhone,
	}
	token, err := l.svcCtx.Repo.VerifyCode.SetVerifyToken(l.ctx, verifyToken, constants.VerifyCodeTTL)
	if err != nil {
		l.Errorf("生成手机号变更令牌失败，userID=%d，phone=%s，err=%v", userID, req.NewPhone, err)
		return nil, err
	}

	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, req.NewPhone)

	return &types.VerifyNewPhoneResponse{
		VerifyToken: token,
	}, nil
}

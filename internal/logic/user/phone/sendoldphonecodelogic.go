// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package phone

import (
	"CampusTake/internal/auth"
	"CampusTake/internal/constants"
	"CampusTake/pkg/ctxx"
	"context"

	"CampusTake/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendOldPhoneCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendOldPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendOldPhoneCodeLogic {
	return &SendOldPhoneCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendOldPhoneCodeLogic) SendOldPhoneCode() error {
	userID := ctxx.MustUserID(l.ctx)
	user, err := l.svcCtx.Repo.User.GetByID(l.ctx, userID)
	if err != nil {
		l.Errorf("查询当前用户失败，userID=%d，err=%v", userID, err)
		return err
	}

	code, err := auth.Generate6DigitCode()
	if err != nil {
		l.Errorf("生成旧手机号验证码失败，userID=%d，err=%v", userID, err)
		return err
	}

	err = l.svcCtx.Repo.VerifyCode.SetCode(l.ctx, user.Phone, code, constants.VerifyCodeTTL)
	if err != nil {
		l.Errorf("保存旧手机号验证码失败，userID=%d，phone=%s，err=%v", userID, user.Phone, err)
		return err
	}
	// todo: 后续补齐发送验证码的逻辑
	return nil
}

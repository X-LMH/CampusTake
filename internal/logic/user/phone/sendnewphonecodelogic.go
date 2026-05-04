// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package phone

import (
	"CampusTake/common/ctxx"
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"CampusTake/common/utils"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendNewPhoneCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendNewPhoneCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNewPhoneCodeLogic {
	return &SendNewPhoneCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendNewPhoneCodeLogic) SendNewPhoneCode(req *types.SendNewPhoneCodeRequest) error {
	userID := ctxx.MustUserID(l.ctx)
	user, err := l.svcCtx.Repo.User().GetByID(l.ctx, userID)
	if err != nil {
		return err
	}
	if user.Phone == req.NewPhone {
		return errx.ErrPhoneSameWithOld // 原来名字也建议改
	}

	_, err = l.svcCtx.Repo.User().GetByPhone(l.ctx, req.NewPhone)
	if err == nil {
		return errx.ErrPhoneAlreadyBound
	}

	code, err := utils.Generate6DigitCode()
	if err != nil {
		return err
	}

	err = l.svcCtx.Repo.VerifyCode().SetCode(l.ctx, req.NewPhone, code, enum.VerifyCodeTTL)
	if err != nil {
		return err
	}

	// todo: 后续补齐发送验证码的逻辑

	return nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"CampusTake/common/utils"
	"context"
	"errors"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateVerifyCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateVerifyCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateVerifyCodeLogic {
	return &GenerateVerifyCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
func (l *GenerateVerifyCodeLogic) GenerateVerifyCode(req *types.GenerateVerifyCodeRequest) (*types.GenerateVerifyCodeResponse, error) {
	// 1. 频率控制：比如 60 秒内只能发一次
	code, err := l.svcCtx.Repo.VerifyCode().GetCode(l.ctx, req.Phone)
	if err == nil {
		return nil, errx.ErrVerifyCodeTooFrequent
	}
	if !errors.Is(err, errx.ErrVerifyCodeNotFound) {
		return nil, err
	}

	// 2. 生成验证码：严谨处理错误
	code, err = utils.Generate6DigitCode()
	if err != nil {
		return nil, err
	}

	// 3. 存入 Redis
	err = l.svcCtx.Repo.VerifyCode().SetCode(l.ctx, req.Phone, code, enum.VerifyCodeTTL)
	if err != nil {
		return nil, err
	}

	l.Debugf("Generated verify code %s for phone %s", code, req.Phone)

	// 4. 发送短信 (异步或同步)
	// 建议在 svcCtx 里集成一个 SMS 服务
	//go func() {
	//	// 这里记得传一个新的 context 或者处理好超时，不要直接用 l.ctx
	//	_ = l.svcCtx.SmsModel.Send(req.Phone, code)
	//}()

	return &types.GenerateVerifyCodeResponse{
		Code: code,
	}, nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package phone

import (
	"CampusTake/common/ctxx"
	"CampusTake/common/enum"
	"CampusTake/common/errx"
	"CampusTake/internal/repo"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePhoneLogic {
	return &ChangePhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePhoneLogic) ChangePhone(req *types.ChangePhoneRequest) error {
	userID := ctxx.MustUserID(l.ctx)
	user, err := l.svcCtx.Repo.User().GetByID(l.ctx, userID)
	if err != nil {
		return err
	}

	// 校验一下verifyToken是否正确并且合法
	verifyToken, err := l.svcCtx.Repo.VerifyCode().GetVerifyToken(l.ctx, req.VerifyToken)
	if err != nil {
		l.Infof("invalid verify token, user_id=%d, token=%s", userID, req.VerifyToken)
		return err
	}

	if verifyToken.UserID != userID {
		return errx.ErrUserPermissionDenied
	}

	switch req.VerifyType {
	case enum.ChangePhoneByPassword:
		if req.Password == "" {
			return errx.NewParamError("密码不能为空")
		}

		if req.Password != user.Password {
			return errx.ErrPasswordWrong
		}

	case enum.ChangePhoneByOldPhoneCode:
		if req.OldCode == "" {
			return errx.NewParamError("验证码不能为空")
		}

		code, err := l.svcCtx.Repo.VerifyCode().GetCode(l.ctx, user.Phone)
		if err != nil {
			return err
		}
		if code != req.OldCode {
			return errx.ErrVerifyCodeWrong
		}

	default:
		return errx.ErrInvalidParam
	}

	err = l.svcCtx.Repo.WithTx(l.ctx, func(r *repo.Repo) error {

		// 更新手机号（必须用 token 里的）
		err = r.User().UpdatePhoneByID(l.ctx, userID, verifyToken.Phone)
		if err != nil {
			return err
		}

		_ = r.VerifyCode().DeleteVerifyToken(l.ctx, req.VerifyToken)
		_ = r.VerifyCode().DeleteCode(l.ctx, verifyToken.Phone)
		_ = r.VerifyCode().DeleteCode(l.ctx, user.Phone)

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

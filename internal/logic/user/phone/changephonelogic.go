// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package phone

import (
	"CampusTake/internal/enums"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
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
	user, err := l.svcCtx.Repo.User.GetByID(l.ctx, userID)
	if err != nil {
		return err
	}

	// 校验一下verifyToken是否正确并且合法
	verifyToken, err := l.svcCtx.Repo.VerifyCode.GetVerifyToken(l.ctx, req.VerifyToken)
	if err != nil {
		l.Infof("invalid verify token, user_id=%d, token=%s", userID, req.VerifyToken)
		return err
	}

	if verifyToken.UserID != userID {
		return errs.ErrUserPermissionDenied
	}

	switch req.VerifyType {
	case enums.ChangePhoneByPassword:
		if req.Password == "" {
			return errs.NewParamError("密码不能为空")
		}

		if req.Password != user.Password {
			return errs.ErrPasswordWrong
		}

	case enums.ChangePhoneByOldPhoneCode:
		if req.OldCode == "" {
			return errs.NewParamError("验证码不能为空")
		}

		code, err := l.svcCtx.Repo.VerifyCode.GetCode(l.ctx, user.Phone)
		if err != nil {
			return err
		}
		if code != req.OldCode {
			return errs.ErrVerifyCodeWrong
		}

	default:
		return errs.ErrInvalidParam
	}

	err = l.svcCtx.Repo.User.UpdatePhoneByID(l.ctx, userID, verifyToken.Phone)
	if err != nil {
		return err
	}

	_ = l.svcCtx.Repo.VerifyCode.DeleteVerifyToken(l.ctx, req.VerifyToken)
	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, verifyToken.Phone)
	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, user.Phone)

	return nil
}

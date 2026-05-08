// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package auth

import (
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangePasswordLogic) ChangePassword(req *types.ChangePasswordRequest) error {
	// 新密码不能与旧密码相同
	if req.OldPassword == req.NewPassword {
		return errors.ErrPasswordNoChange
	}

	// 密码是否正确
	userID := ctxx.MustUserID(l.ctx)
	l.Debugf("ChangePassword userID: %d, req: %+v", userID, req)
	user, err := l.svcCtx.Repo.User().GetByID(l.ctx, userID)
	if err != nil {
		return err
	}
	if user.Password != req.OldPassword {
		return errors.ErrPasswordWrong
	}

	// 更新密码
	err = l.svcCtx.Repo.User().UpdatePasswordByID(l.ctx, userID, req.NewPassword)
	if err != nil {
		return err
	}
	return nil
}

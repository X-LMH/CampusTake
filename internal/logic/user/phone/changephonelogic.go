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
		l.Errorf("查询当前用户失败，userID=%d，err=%v", userID, err)
		return err
	}

	// 校验一下verifyToken是否正确并且合法
	verifyToken, err := l.svcCtx.Repo.VerifyCode.GetVerifyToken(l.ctx, req.VerifyToken)
	if err != nil {
		l.Errorf("获取换绑手机号验证码令牌失败，userID=%d，token=%s，err=%v", userID, req.VerifyToken, err)
		return err
	}

	if verifyToken.UserID != userID {
		l.Errorf("换绑手机号令牌不匹配，userID=%d，tokenUserID=%d，token=%s", userID, verifyToken.UserID, req.VerifyToken)
		return errs.ErrUserPermissionDenied
	}

	switch req.VerifyType {
	case enums.ChangePhoneByPassword:
		if req.Password == "" {
			l.Errorf("换绑手机号时密码不能为空，userID=%d", userID)
			return errs.NewParamError("密码不能为空")
		}

		if req.Password != user.Password {
			l.Errorf("换绑手机号密码错误，userID=%d", userID)
			return errs.ErrPasswordWrong
		}

	case enums.ChangePhoneByOldPhoneCode:
		if req.OldCode == "" {
			l.Errorf("换绑手机号验证码不能为空，userID=%d", userID)
			return errs.NewParamError("验证码不能为空")
		}

		code, err := l.svcCtx.Repo.VerifyCode.GetCode(l.ctx, user.Phone)
		if err != nil {
			l.Errorf("获取旧手机号验证码失败，userID=%d，phone=%s，err=%v", userID, user.Phone, err)
			return err
		}
		if code != req.OldCode {
			l.Errorf("旧手机号验证码错误，userID=%d，phone=%s", userID, user.Phone)
			return errs.ErrVerifyCodeWrong
		}

	default:
		l.Errorf("换绑手机号验证方式不合法，userID=%d，verifyType=%v", userID, req.VerifyType)
		return errs.ErrInvalidParam
	}

	err = l.svcCtx.Repo.User.UpdatePhoneByID(l.ctx, userID, verifyToken.Phone)
	if err != nil {
		l.Errorf("更新手机号失败，userID=%d，newPhone=%s，err=%v", userID, verifyToken.Phone, err)
		return err
	}

	_ = l.svcCtx.Repo.VerifyCode.DeleteVerifyToken(l.ctx, req.VerifyToken)
	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, verifyToken.Phone)
	_ = l.svcCtx.Repo.VerifyCode.DeleteCode(l.ctx, user.Phone)

	return nil
}

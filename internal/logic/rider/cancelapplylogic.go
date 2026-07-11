// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelApplyLogic {
	return &CancelApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelApplyLogic) CancelApply() error {
	userID := ctxx.MustUserID(l.ctx)

	// 1. 获取申请单
	profile, err := l.svcCtx.Repo.Rider.GetProfileByUserID(l.ctx, userID)
	if err != nil {
		l.Errorf("查询骑手申请记录失败，userID=%d，err=%v", userID, err)
		return err
	}

	// 2. 校验：只有待审核状态可以撤回
	if profile.AuditStatus != enums.RiderStatusPending {
		l.Errorf("骑手申请状态不允许撤回，userID=%d，status=%v", userID, profile.AuditStatus)
		return errors.ErrCannotCancelStatus
	}

	// 3. 更新状态为 已撤回 (3)
	return l.svcCtx.Repo.Rider.UpdateStatusByUserID(l.ctx, userID, enums.RiderStatusCancel)
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelAppealLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelAppealLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelAppealLogic {
	return &CancelAppealLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelAppealLogic) CancelAppeal(req *types.CancelAppealRequest) error {
	userID := ctxx.MustUserID(l.ctx)

	appeal, err := l.svcCtx.Repo.Appeal.GetByIDAndUserID(l.ctx, req.AppealID, userID)
	if err != nil {
		return err
	}

	// 只有待处理的申诉可以撤销
	if appeal.Status != enums.AppealStatusPending {
		return errs.ErrAppealCannotCancel
	}

	fromStatus := appeal.Status
	toStatus := enums.AppealStatusCancelled

	return l.svcCtx.Repo.Appeal.UpdateStatus(
		l.ctx,
		appeal.ID,
		fromStatus,
		toStatus,
	)
}

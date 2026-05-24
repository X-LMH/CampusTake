// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/repo"
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

	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		appeal, err := tx.Appeal.GetByIDAndUserID(l.ctx, req.AppealID, userID)
		if err != nil {
			l.Errorf("查询申诉失败 appealID=%d userID=%d err=%v", req.AppealID, userID, err)
			return err
		}

		// 仅允许待处理申诉撤销
		if appeal.Status != enums.AppealStatusPending {
			l.Errorf("申诉状态不允许撤销，appealID=%d，userID=%d，status=%v", req.AppealID, userID, appeal.Status)
			return errs.ErrAppealStatusInvalid
		}

		// 更新 appeal 表
		err = tx.Appeal.UpdateStatusByID(
			l.ctx,
			appeal.ID,
			enums.AppealStatusPending,
			enums.AppealStatusCancelled,
		)
		if err != nil {
			l.Errorf("更新申诉状态失败 appealID=%d err=%v", appeal.ID, err)
			return err
		}

		// 更新 order 表
		err = tx.Order.UpdateAppealStatusByID(
			l.ctx,
			appeal.OrderID,
			enums.OrderAppealStatusOngoing,
			enums.OrderAppealStatusNone,
		)
		if err != nil {
			l.Errorf("更新订单申诉状态失败 orderID=%d err=%v", appeal.OrderID, err)
			return err
		}

		return nil
	})
}

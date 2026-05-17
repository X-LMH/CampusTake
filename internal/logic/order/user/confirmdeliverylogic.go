// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"
	"time"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmDeliveryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmDeliveryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmDeliveryLogic {
	return &ConfirmDeliveryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmDeliveryLogic) ConfirmDelivery(
	req *types.ConfirmDeliveryRequest,
) error {

	userID := ctxx.MustUserID(l.ctx)

	order, err := l.svcCtx.Repo.Order.GetByIDAndUserID(
		l.ctx,
		req.OrderID,
		userID,
	)
	if err != nil {
		l.Errorf("查询订单失败，orderID=%d userID=%d err=%v", req.OrderID, userID, err)
		return err
	}

	// =========================
	// 只有已送达才能确认收货
	// =========================

	if order.Status != enums.OrderDelivered {
		return errs.ErrOrderStatusInvalid
	}

	at := time.Now()

	fromStatus := enums.OrderDelivered
	toStatus := enums.OrderCompleted

	return l.svcCtx.Repo.WithTx(
		l.ctx,
		func(tx *repo.RepoTx) error {

			// =========================
			// 更新订单状态
			// =========================

			err = tx.Order.UserUpdateStatusAndTime(
				l.ctx,
				req.OrderID,
				userID,
				fromStatus,
				toStatus,
				at,
				nil,
			)
			if err != nil {
				return err
			}

			// =========================
			// 创建订单日志
			// =========================

			log := &model.OrderLog{
				OrderID:      order.ID,
				FromStatus:   fromStatus,
				ToStatus:     toStatus,
				OperatorType: enums.OperatorTypeUser,
				OperatorID:   userID,
				Remark:       "用户确认收货",
				CreatedAt:    at,
			}

			return tx.Order.CreateLog(l.ctx, log)
		},
	)
}

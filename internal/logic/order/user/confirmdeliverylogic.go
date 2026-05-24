// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo"
	"CampusTake/internal/repo/query"
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

func (l *ConfirmDeliveryLogic) ConfirmDelivery(req *types.ConfirmDeliveryRequest) error {
	userID := ctxx.MustUserID(l.ctx)

	order, err := l.svcCtx.Repo.Order.GetByIDAndUserID(l.ctx, req.OrderID, userID)
	if err != nil {
		l.Errorf("query order failed, orderID=%d userID=%d err=%v", req.OrderID, userID, err)
		return err
	}

	if order.Status != enums.OrderDelivered {
		l.Errorf("order status does not allow confirm, orderID=%d userID=%d status=%v", req.OrderID, userID, order.Status)
		return errs.ErrOrderStatusInvalid
	}

	at := time.Now()
	fromStatus := enums.OrderDelivered
	toStatus := enums.OrderCompleted

	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		if err := tx.Order.UpdateStatusAndTime(
			l.ctx,
			query.OrderStatusUpdateQuery{
				OrderID: order.ID,
				UserID:  &userID,
			},
			fromStatus,
			toStatus,
			at,
			nil,
		); err != nil {
			l.Errorf("update confirm status failed, orderID=%d userID=%d err=%v", order.ID, userID, err)
			return err
		}

		log := &model.OrderLog{
			OrderID:      order.ID,
			FromStatus:   fromStatus,
			ToStatus:     toStatus,
			OperatorType: enums.OperatorTypeUser,
			OperatorID:   userID,
			Remark:       "User confirmed delivery",
			CreatedAt:    at,
		}
		if err := tx.Order.CreateLog(l.ctx, log); err != nil {
			l.Errorf("create confirm log failed, orderID=%d userID=%d err=%v", order.ID, userID, err)
			return err
		}

		if order.AppealStatus == enums.OrderAppealStatusApproved {
			return nil
		}
		if order.RiderID == nil {
			return errs.ErrOrderNoRider
		}

		if err := tx.Rider.IncrementCompletedOrderCount(l.ctx, *order.RiderID); err != nil {
			l.Errorf("increment rider completed count failed, orderID=%d riderID=%d err=%v", order.ID, *order.RiderID, err)
			return err
		}

		return nil
	})
}

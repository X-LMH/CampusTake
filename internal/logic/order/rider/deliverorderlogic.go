// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/mqs"
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

type DeliverOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeliverOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeliverOrderLogic {
	return &DeliverOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeliverOrderLogic) DeliverOrder(req *types.DeliverOrderRequest) error {
	userID := ctxx.MustUserID(l.ctx)

	deliveredAt := time.Now()

	fromStatus := enums.OrderDelivering
	toStatus := enums.OrderDelivered

	err := l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

		// 查询订单 + 权限校验
		order, err := tx.Order.GetByIDAndRiderID(
			l.ctx,
			req.OrderID,
			userID,
		)
		if err != nil {
			l.Errorf("查询骑手送达订单失败，orderID=%d，riderID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		// 状态校验
		if order.Status != fromStatus {
			l.Errorf("订单状态不允许送达，orderID=%d，riderID=%d，status=%v", req.OrderID, userID, order.Status)
			return errs.ErrOrderStatusInvalid
		}

		// 原子更新状态
		err = tx.Order.UpdateStatusAndTime(
			l.ctx,
			query.OrderStatusUpdateQuery{
				OrderID: req.OrderID,
				RiderID: &userID,
			},
			fromStatus,
			toStatus,
			deliveredAt,
			nil,
		)
		if err != nil {
			l.Errorf("更新骑手送达状态失败，orderID=%d，riderID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		// 写入订单日志
		log := &model.OrderLog{
			OrderID:      req.OrderID,
			FromStatus:   fromStatus,
			ToStatus:     toStatus,
			OperatorType: enums.OperatorTypeRider,
			OperatorID:   userID,
			Remark:       "骑手已送达",
			CreatedAt:    deliveredAt,
		}

		return tx.Order.CreateLog(l.ctx, log)
	})
	if err != nil {
		return err
	}

	err = mqs.PublishDelayConfirmOrder(l.svcCtx, req.OrderID)
	if err != nil {
		l.Errorf("发送延迟确认收货消息失败，orderID=%d err=%v", req.OrderID, err)
	}
	return nil
}

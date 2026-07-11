package rider

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

type PickupOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPickupOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PickupOrderLogic {
	return &PickupOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PickupOrderLogic) PickupOrder(req *types.PickupOrderRequest) error {

	userID := ctxx.MustUserID(l.ctx)

	deliveringAt := time.Now()

	fromStatus := enums.OrderAccepted
	toStatus := enums.OrderDelivering

	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

		// 查询订单 + 权限校验
		order, err := tx.Order.GetByIDAndRiderID(
			l.ctx,
			req.OrderID,
			userID,
		)
		if err != nil {
			l.Errorf("查询骑手取件订单失败，orderID=%d，riderID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		// 状态校验
		if order.Status != fromStatus {
			l.Errorf("订单状态不允许开始配送，orderID=%d，riderID=%d，status=%v", req.OrderID, userID, order.Status)
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
			deliveringAt,
			nil,
		)
		if err != nil {
			l.Errorf("更新骑手取件状态失败，orderID=%d，riderID=%d，err=%v", req.OrderID, userID, err)
			return err
		}

		// 写入订单日志
		log := &model.OrderLog{
			OrderID:      req.OrderID,
			FromStatus:   fromStatus,
			ToStatus:     toStatus,
			OperatorType: enums.OperatorTypeRider,
			OperatorID:   userID,
			Remark:       "骑手开始配送",
			CreatedAt:    deliveringAt,
		}

		return tx.Order.CreateLog(l.ctx, log)
	})
}

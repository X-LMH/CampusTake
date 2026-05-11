// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

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

	fromStatus := enums.OrderPickedUp
	toStatus := enums.OrderDelivered

	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {

		// 查询订单 + 权限校验
		order, err := tx.Order.GetByIDAndRiderID(
			l.ctx,
			req.OrderID,
			userID,
		)
		if err != nil {
			return err
		}

		// 状态校验
		if order.Status != fromStatus {
			return errs.ErrOrderStatusInvalid
		}

		// 原子更新状态
		err = tx.Order.RiderUpdateStatusAndTime(l.ctx, req.OrderID, userID, fromStatus, toStatus, deliveredAt, nil)
		if err != nil {
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
}

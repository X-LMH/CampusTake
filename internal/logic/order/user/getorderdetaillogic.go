// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/utils"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderDetailLogic {
	return &GetOrderDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderDetailLogic) GetOrderDetail(req *types.GetOrderDetailRequest) (resp *types.OrderItem, err error) {
	userID := ctxx.MustUserID(l.ctx)
	order, err := l.svcCtx.Repo.Order().GetByIDAndUserID(l.ctx, req.OrderID, userID)
	if err != nil {
		return nil, err
	}

	return &types.OrderItem{
		ID:                order.ID,
		OrderNo:           order.OrderNo,
		UserID:            order.UserID,
		RiderID:           order.RiderID,
		OrderType:         int8(order.OrderType),
		PickupAddressID:   order.PickupAddressID,
		DeliveryAddressID: order.DeliveryAddressID,
		RewardAmount:      order.RewardAmount,
		Status:            int8(order.Status),
		PaymentStatus:     int8(order.PaymentStatus),
		Remark:            order.Remark,
		CancelReason:      order.CancelReason,
		CreatedAt:         utils.FormatCreatedAt(order.CreatedAt),
		PaidAt:            utils.FormatTimePtr(order.PaidAt),
		AcceptedAt:        utils.FormatTimePtr(order.AcceptedAt),
		FinishedAt:        utils.FormatTimePtr(order.FinishedAt),
	}, nil
}

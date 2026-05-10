// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/pkg/ctxx"
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
	order, err := l.svcCtx.Repo.Order.GetByIDAndUserID(l.ctx, req.OrderID, userID)
	if err != nil {
		return nil, err
	}

	return types.ModelOrderToOrderDetail(order), nil
}

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

func (l *GetOrderDetailLogic) GetOrderDetail(req *types.GetOrderDetailRequest) (resp *types.GetOrderDetailResponse, err error) {
	userID := ctxx.MustUserID(l.ctx)

	// 1. 查询订单主表信息
	order, err := l.svcCtx.Repo.Order.GetByIDAndUserID(l.ctx, req.OrderID, userID)
	if err != nil {
		return nil, err
	}

	// 初始化返回体，默认 RiderInfo 为 nil
	resp = &types.GetOrderDetailResponse{
		Order:     types.ModelOrderToOrderItem(order),
		RiderInfo: nil,
	}

	// 2. 查询骑手信息
	if order.RiderID != nil && *order.RiderID > 0 {
		rider, err := l.svcCtx.Repo.Rider.GetProfileByRiderID(l.ctx, *order.RiderID)
		if err != nil {
			l.Errorf("查询骑手信息失败（已降级），riderID=%d orderID=%d err=%v", *order.RiderID, req.OrderID, err)
		} else if rider != nil {

			resp.RiderInfo = &types.RiderInfo{
				RiderID:             rider.ID,
				RealName:            maskRiderName(rider.RealName),
				RatingAvg:           rider.RatingAvg,
				RatingCount:         int(rider.RatingCount),
				CompletedOrderCount: int(rider.CompletedOrderCount),
			}
		}
	}

	return resp, nil
}

// 名字脱敏保持不变
func maskRiderName(name string) string {
	runes := []rune(name)
	if len(runes) <= 1 {
		return name
	}
	return string(runes[0]) + "*"
}

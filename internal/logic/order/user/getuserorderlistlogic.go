// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserOrderListLogic {
	return &GetUserOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserOrderListLogic) GetUserOrderList(req *types.GetUserOrderListRequest) (resp *types.GetUserOrderListResponse, err error) {
	userID := ctxx.MustUserID(l.ctx)
	// 直接查询订单列表
	status := enums.OrderStatus(req.Status)
	if !enums.IsOrderStatus(status) {
		status = enums.OrderDefault // 默认查询全部订单
	}
	pageResult, err := l.svcCtx.Repo.Order.GetListByUserID(l.ctx, userID, status, req.Page, req.Size)
	if err != nil {
		l.Errorf("查询用户订单列表失败, userID=%d, status=%d, err=%v", userID, status, err)
		return nil, err
	}

	orders, ok := pageResult.Records.([]*model.Order)
	if !ok {
		l.Error("用户订单分页数据类型断言失败")
		return nil, errors.ErrServiceError
	}

	list := make([]types.OrderItem, 0, len(orders))
	for _, item := range orders {
		orderItem := types.ModelOrderToOrderItem(item)
		list = append(list, orderItem)
	}

	return &types.GetUserOrderListResponse{
		Total: pageResult.Total,
		List:  list,
	}, nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/model"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvailableOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvailableOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailableOrderListLogic {
	return &GetAvailableOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvailableOrderListLogic) GetAvailableOrderList(req *types.GetAvailableOrderListRequest) (resp *types.GetAvailableOrderListResponse, err error) {
	orderClause, order := mapSort(req.SortBy, req.SortOrder)
	// 调用 repo
	pageResult, err := l.svcCtx.Repo.Order.GetAvailableForRider(
		l.ctx,
		req.Page,
		req.Size,
		orderClause,
		order,
		req.MinReward,
		req.MaxReward,
	)
	if err != nil {
		l.Errorf("查询可接单订单列表失败，page=%d，size=%d，err=%v", req.Page, req.Size, err)
		return nil, err
	}

	orders, ok := pageResult.Records.([]*model.Order)
	if !ok {
		l.Error("可接单订单列表分页数据类型断言失败")
		return nil, errors.ErrServiceError
	}

	list := make([]types.OrderItem, 0, len(orders))
	for _, item := range orders {
		if orderItem := types.ModelOrderToOrderItem(item); orderItem != nil {
			list = append(list, *orderItem)
		}
	}

	return &types.GetAvailableOrderListResponse{
		Total: pageResult.Total,
		List:  list,
	}, nil
}

func mapSort(sortBy, sortOrder string) (string, string) {
	var orderClause string
	switch sortBy {
	case "reward":
		orderClause = "reward_amount"
	case "createdAt":
		orderClause = "created_at"
	default:
		orderClause = "created_at"
	}

	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	return orderClause, sortOrder
}

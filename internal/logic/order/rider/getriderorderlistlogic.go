// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRiderOrderListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRiderOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRiderOrderListLogic {
	return &GetRiderOrderListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRiderOrderListLogic) GetRiderOrderList(req *types.GetRiderOrderListRequest) (resp *types.GetRiderOrderListResponse, err error) {
	userID := ctxx.MustUserID(l.ctx)

	pageResult, err := l.svcCtx.Repo.Order.GetListByRiderID(l.ctx, userID, enums.OrderStatus(req.Status), req.Page, req.Size)
	if err != nil {
		l.Errorf("查询骑手订单列表失败，userID=%d，status=%d，page=%d，size=%d，err=%v", userID, req.Status, req.Page, req.Size, err)
		return nil, err
	}

	orders, ok := pageResult.Records.([]*model.Order)
	if !ok {
		l.Error("骑手订单列表分页数据类型断言失败")
		return nil, errors.ErrServiceError
	}

	list := make([]types.OrderItem, 0, len(orders))
	for _, item := range orders {
		if orderItem := types.ModelOrderToOrderItem(item); orderItem != nil {
			list = append(list, *orderItem)
		}
	}
	return &types.GetRiderOrderListResponse{
		Total: pageResult.Total,
		List:  list,
	}, nil
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"context"

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
	// todo: add your logic here and delete this line

	return nil
}

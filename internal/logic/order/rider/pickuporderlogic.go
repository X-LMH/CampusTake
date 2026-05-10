// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"context"

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
	// todo: add your logic here and delete this line

	return nil
}

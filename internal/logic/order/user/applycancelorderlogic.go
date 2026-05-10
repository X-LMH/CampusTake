// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ApplyCancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyCancelOrderLogic {
	return &ApplyCancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyCancelOrderLogic) ApplyCancelOrder(req *types.ApplyCancelOrderRequest) error {

	return nil
}

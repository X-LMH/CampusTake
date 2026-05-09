// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogLogic {
	return &GetOrderLogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogLogic) GetOrderLog(req *types.GetOrderLogRequest) (resp *types.GetOrderLogResponse, err error) {
	// todo: add your logic here and delete this line

	return
}

// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReviewLogic {
	return &CreateReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateReviewLogic) CreateReview(req *types.CreateReviewRequest) error {
	// todo: add your logic here and delete this line

	return nil
}

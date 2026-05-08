// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAddressLogic {
	return &DeleteAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteAddressLogic) DeleteAddress(req *types.DeleteAddressRequest) error {
	userID := ctxx.MustUserID(l.ctx)

	// 直接调用 Repo 删除，Repo 内部已包含 userID 校验，防止越权
	err := l.svcCtx.Repo.Address().DeleteByID(l.ctx, req.AddressID, userID)
	if err != nil {
		return err // errors.ErrAddressNotFound 会在这里被抛出
	}

	return nil
}

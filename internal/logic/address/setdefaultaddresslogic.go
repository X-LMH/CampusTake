// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/common/ctxx"
	"CampusTake/internal/repo"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetDefaultAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultAddressLogic {
	return &SetDefaultAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultAddressLogic) SetDefaultAddress(req *types.SetDefaultAddressRequest) error {
	userID := ctxx.MustUserID(l.ctx)

	return l.svcCtx.Repo.WithTx(l.ctx, func(r *repo.Repo) error {
		if err := r.Address().ClearDefaultByUserID(l.ctx, userID); err != nil {
			return err
		}
		if err := r.Address().SetDefaultByID(l.ctx, req.AddressID, userID); err != nil {
			return err
		}
		return nil
	})
}

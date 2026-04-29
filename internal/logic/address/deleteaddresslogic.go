// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/common/ctxx"
	"CampusTake/common/errx"
	"CampusTake/internal/repo"
	"context"
	"errors"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

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

	return l.svcCtx.Repo.WithTx(l.ctx, func(r *repo.Repo) error {
		address, err := r.Address().GetByIDAndUserID(l.ctx, req.AddressID, userID)
		if err != nil {
			return err
		}

		if err := r.Address().DeleteByID(l.ctx, req.AddressID, userID); err != nil {
			return err
		}

		if address.IsDefault == 1 {
			nextAddress, err := r.Address().GetFirstByUserID(l.ctx, userID)
			if err != nil {
				if errors.Is(err, errx.ErrAddressNotFound) {
					return nil
				}
				return err
			}
			if err := r.Address().SetDefaultByID(l.ctx, nextAddress.ID, userID); err != nil {
				return err
			}
		}

		return nil
	})
}

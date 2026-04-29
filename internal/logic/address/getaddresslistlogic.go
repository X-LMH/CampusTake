// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/common/ctxx"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAddressListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressListLogic {
	return &GetAddressListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddressListLogic) GetAddressList() (resp []types.AddressItem, err error) {
	userID := ctxx.MustUserID(l.ctx)

	addressList, err := l.svcCtx.Repo.Address().GetListByUserID(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	resp = make([]types.AddressItem, len(addressList))
	for i, addr := range addressList {
		resp[i] = types.AddressItem{
			AddressID:    addr.ID,
			ContactName:  addr.ContactName,
			ContactPhone: addr.ContactPhone,
			Building:     addr.Building,
			Room:         addr.Room,
			IsDefault:    int8(addr.IsDefault),
		}
	}
	return resp, nil
}

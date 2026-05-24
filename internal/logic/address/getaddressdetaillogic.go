// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/pkg/ctxx"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAddressDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAddressDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAddressDetailLogic {
	return &GetAddressDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAddressDetailLogic) GetAddressDetail(req *types.GetAddressDetailRequest) (resp *types.AddressItem, err error) {
	userID := ctxx.MustUserID(l.ctx)
	addressID := req.AddressID

	address, err := l.svcCtx.Repo.Address.GetByIDAndUserID(l.ctx, addressID, userID)
	if err != nil {
		l.Errorf("查询地址详情失败，addressID=%d，userID=%d，err=%v", addressID, userID, err)
		return nil, err
	}

	return &types.AddressItem{
		AddressID:    address.ID,
		Type:         int8(address.Type),
		ContactName:  address.ContactName,
		ContactPhone: address.ContactPhone,
		Building:     address.Building,
		Room:         address.Room,
		Detail:       address.Detail,
		IsDefault:    int8(address.IsDefault),
	}, nil
}

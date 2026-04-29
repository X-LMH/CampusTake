// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/common/ctxx"
	"CampusTake/common/enum"
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddAddressLogic {
	return &AddAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddAddressLogic) AddAddress(req *types.AddAddressRequset) (resp *types.AddressItem, err error) {
	userID := ctxx.MustUserID(l.ctx)

	// 查看用户是否已有地址，若没有则新地址设为默认地址
	addresses, err := l.svcCtx.Repo.Address().GetListByUserID(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	var isDefault = enum.AddressIsDefault
	if len(addresses) > 0 {
		isDefault = enum.AddressNotDefault
	}

	// 创建地址
	address := &model.Address{
		UserID:       userID,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Building:     req.Building,
		Room:         req.Room,
		IsDefault:    isDefault,
	}
	err = l.svcCtx.Repo.Address().Create(l.ctx, address)
	if err != nil {
		return nil, err
	}

	return &types.AddressItem{
		AddressID:    address.ID,
		ContactName:  address.ContactName,
		ContactPhone: address.ContactPhone,
		Building:     address.Building,
		Room:         address.Room,
		IsDefault:    int8(address.IsDefault),
	}, nil
}

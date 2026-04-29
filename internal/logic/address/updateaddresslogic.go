// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/common/ctxx"
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAddressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAddressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAddressLogic {
	return &UpdateAddressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAddressLogic) UpdateAddress(req *types.UpdateAddressRequest) (resp *types.AddressItem, err error) {
	userID := ctxx.MustUserID(l.ctx)

	address := &model.Address{
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Building:     req.Building,
		Room:         req.Room,
	}

	// 更新地址
	if err := l.svcCtx.Repo.Address().UpdateByID(l.ctx, req.AddressID, userID, address); err != nil {
		return nil, err
	}

	// 查询更新后的地址
	updatedAddress, err := l.svcCtx.Repo.Address().GetByIDAndUserID(l.ctx, req.AddressID, userID)
	if err != nil {
		return nil, err
	}

	// 组装返回值
	resp = &types.AddressItem{
		AddressID:    updatedAddress.ID,
		ContactName:  updatedAddress.ContactName,
		ContactPhone: updatedAddress.ContactPhone,
		Building:     updatedAddress.Building,
		Room:         updatedAddress.Room,
		IsDefault:    int8(updatedAddress.IsDefault),
	}

	return resp, nil
}

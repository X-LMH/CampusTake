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
	addrType := enum.AddressType(req.Type)

	// 1. 查询该用户在该类型（收货/取件）下已有的地址数量
	// 建议在 repo 实现这个 Count 方法，比捞出整个 list 效率高得多
	count, err := l.svcCtx.Repo.Address().CountByUserIDAndType(l.ctx, userID, addrType)
	if err != nil {
		return nil, err
	}

	// 2. 逻辑：如果是该类型下的第一个地址，则设为默认
	isDefault := enum.AddressNotDefault
	if count == 0 {
		isDefault = enum.AddressIsDefault
	}

	// 3. 创建地址对象
	// 注意: Detail 字段来自请求，可为空
	address := &model.Address{
		UserID:       userID,
		Type:         addrType,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Building:     req.Building,
		Room:         req.Room,
		Detail:       req.Detail,
		IsDefault:    isDefault,
	}

	// 4. 写入数据库
	err = l.svcCtx.Repo.Address().Create(l.ctx, address)
	if err != nil {
		return nil, err
	}

	// 5. 返回结果（注意：address.ID 在 Create 执行后会被 GORM 自动回填）
	return &types.AddressItem{
		AddressID:    address.ID,
		Type:         int8(addrType),
		ContactName:  address.ContactName,
		ContactPhone: address.ContactPhone,
		Building:     address.Building,
		Room:         address.Room,
		Detail:       address.Detail,
		IsDefault:    int8(address.IsDefault),
	}, nil
}

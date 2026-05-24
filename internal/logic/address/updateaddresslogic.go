// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/internal/model"
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
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

	// 1. 构造更新模型
	address := &model.Address{
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Building:     req.Building,
		Room:         req.Room,
		Detail:       req.Detail,
	}

	// 2. 执行更新
	// 注意：Repo 内部应该使用 Updates(map[string]interface{}{...})
	// 以免 GORM 忽略掉零值字段更新（虽然地址字段通常不涉及零值问题）
	if err := l.svcCtx.Repo.Address.UpdateByID(l.ctx, req.AddressID, userID, address); err != nil {
		l.Errorf("更新地址失败，addressID=%d，userID=%d，err=%v", req.AddressID, userID, err)
		return nil, err
	}

	// 3. 查询更新后的完整对象（为了拿到 Type 和 IsDefault）
	updatedAddress, err := l.svcCtx.Repo.Address.GetByIDAndUserID(l.ctx, req.AddressID, userID)
	if err != nil {
		l.Errorf("查询更新后的地址失败，addressID=%d，userID=%d，err=%v", req.AddressID, userID, err)
		return nil, err
	}

	// 4. 组装返回值：确保 Type 字段不丢失
	return &types.AddressItem{
		AddressID:    updatedAddress.ID,
		Type:         int8(updatedAddress.Type), // 补齐这个字段
		ContactName:  updatedAddress.ContactName,
		ContactPhone: updatedAddress.ContactPhone,
		Building:     updatedAddress.Building,
		Room:         updatedAddress.Room,
		Detail:       updatedAddress.Detail,
		IsDefault:    int8(updatedAddress.IsDefault),
	}, nil
}

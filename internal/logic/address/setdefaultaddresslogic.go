// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package address

import (
	"CampusTake/internal/repo"
	"CampusTake/pkg/ctxx"
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

	// 1. 先查询该地址是否存在，并获取其类型 (Type)
	addr, err := l.svcCtx.Repo.Address.GetByIDAndUserID(l.ctx, req.AddressID, userID)
	if err != nil {
		l.Errorf("查询默认地址目标失败，addressID=%d，userID=%d，err=%v", req.AddressID, userID, err)
		return err // 这里底层已经返回了 errors.ErrAddressNotFound
	}

	// 2. 开启事务
	return l.svcCtx.Repo.WithTx(l.ctx, func(tx *repo.RepoTx) error {
		// 3. 只清空该用户下 “同类型” 的默认状态
		// 使用你 Repo 里已经定义好的 ClearDefaultByType
		if err := tx.Address.ClearDefaultByType(l.ctx, userID, addr.Type); err != nil {
			l.Errorf("清空默认地址失败，addressID=%d，userID=%d，地址类型=%d，err=%v", req.AddressID, userID, addr.Type, err)
			return err
		}

		// 4. 设置新的默认地址
		if err := tx.Address.SetDefaultByID(l.ctx, req.AddressID, userID); err != nil {
			l.Errorf("设置默认地址失败，addressID=%d，userID=%d，err=%v", req.AddressID, userID, err)
			return err
		}
		return nil
	})

}

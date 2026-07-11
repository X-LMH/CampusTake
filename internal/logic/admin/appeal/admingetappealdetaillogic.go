// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

import (
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetAppealDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetAppealDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetAppealDetailLogic {
	return &AdminGetAppealDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetAppealDetailLogic) AdminGetAppealDetail(req *types.GetAppealDetailRequest) (resp *types.GetAppealDetailResponse, err error) {

	// 1. 查询申诉
	appeal, err := l.svcCtx.Repo.Appeal.GetByID(l.ctx, req.AppealID)
	if err != nil {
		l.Errorf("管理员查询申诉详情失败 appealID=%d err=%v",
			req.AppealID, err)
		return nil, err
	}

	// 2. 查询最后一次处理记录
	handleModel, err := l.svcCtx.Repo.AppealHandle.GetLastHandleByAppealID(
		l.ctx,
		req.AppealID,
	)
	if err != nil {
		l.Errorf("管理员查询申诉处理记录失败 appealID=%d err=%v",
			req.AppealID, err)
		return nil, err
	}

	// 3. 返回
	return &types.GetAppealDetailResponse{
		Appeal: *types.ModelAppealToAppealItem(appeal),
		Handle: types.ModelAppealHandleToItem(handleModel),
	}, nil
}

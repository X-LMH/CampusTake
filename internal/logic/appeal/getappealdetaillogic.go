// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

import (
	"CampusTake/pkg/ctxx"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppealDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppealDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppealDetailLogic {
	return &GetAppealDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppealDetailLogic) GetAppealDetail(req *types.GetAppealDetailRequest) (resp *types.GetAppealDetailResponse, err error) {
	userID := ctxx.MustUserID(l.ctx)

	// 1. 查询申诉基础信息
	appeal, err := l.svcCtx.Repo.Appeal.GetByIDAndUserID(l.ctx, req.AppealID, userID)
	if err != nil {
		l.Errorf("查询申诉详情失败 appealID=%d userID=%d err=%v", req.AppealID, userID, err)
		return nil, err
	}

	// 2. 查询最近的处理记录
	handleModel, err := l.svcCtx.Repo.AppealHandle.GetLastHandleByAppealID(l.ctx, req.AppealID)
	if err != nil {
		l.Errorf("查询申诉处理详情失败 appealID=%d err=%v", req.AppealID, err)
		return nil, err
	}

	// 3. 组装返回对象
	return &types.GetAppealDetailResponse{
		Appeal: *types.ModelAppealToAppealItem(appeal),
		Handle: types.ModelAppealHandleToItem(handleModel),
	}, nil
}

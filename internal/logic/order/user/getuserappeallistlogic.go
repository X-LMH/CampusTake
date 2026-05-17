// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo/query"
	"CampusTake/pkg/ctxx"
	errs "CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserAppealListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserAppealListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserAppealListLogic {
	return &GetUserAppealListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserAppealListLogic) GetUserAppealList(req *types.GetUserAppealListRequest) (resp *types.GetUserAppealListResponse, err error) {
	userID := ctxx.MustUserID(l.ctx)

	// 1. 分页查询
	status := enums.AppealStatus(req.Status)
	pageQueryResult, err := l.svcCtx.Repo.Appeal.GetList(
		l.ctx,
		query.AppealQuery{
			Status:      &status,
			ApplicantID: &userID,
			Page:        req.Page,
			Size:        req.Size,
		},
	)
	if err != nil {
		return nil, err
	}

	// 2. 类型断言
	dbList, ok := pageQueryResult.Records.([]model.Appeal)
	if !ok {
		l.Errorf("GetUserAppealList type assertion failed, actual type: %T", pageQueryResult.Records)
		return nil, errs.ErrServiceError
	}

	// 3. 转换给前端的 List（保证即使为空，JSON 也是 [] 而不是 null）
	list := make([]types.AppealItem, 0, len(dbList))
	for _, appeal := range dbList {
		list = append(list, types.ModelAppealToAppealItem(appeal))
	}

	return &types.GetUserAppealListResponse{
		Total: pageQueryResult.Total,
		List:  list,
	}, nil
}

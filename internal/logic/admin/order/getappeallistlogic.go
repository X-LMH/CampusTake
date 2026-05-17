// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package order

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/internal/repo/query"
	errs "CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppealListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAppealListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppealListLogic {
	return &GetAppealListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAppealListLogic) GetAppealList(req *types.GetAppealListRequest) (resp *types.GetAppealListResponse, err error) {
	// 1. 分页查询
	status := enums.AppealStatus(req.Status)
	appealType := enums.AppealType(req.Type)

	pageQueryResult, err := l.svcCtx.Repo.Appeal.GetList(
		l.ctx,
		query.AppealQuery{
			Status:      &status,
			ApplicantID: &req.UserID,
			OrderID:     &req.OrderID,
			Type:        &appealType,
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
		l.Errorf("GetAppealList type assertion failed, actual type: %T", pageQueryResult.Records)
		return nil, errs.ErrServiceError
	}

	// 3. 转换给前端的 List（保证即使为空，JSON 也是 [] 而不是 null）
	list := make([]types.AppealItem, 0, len(dbList))
	for _, appeal := range dbList {
		list = append(list, types.ModelAppealToAppealItem(appeal))
	}

	return &types.GetAppealListResponse{
		Total: pageQueryResult.Total,
		List:  list,
	}, nil
}

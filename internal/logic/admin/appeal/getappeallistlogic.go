// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

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
	applicantRole := enums.AppealApplicantRole(req.ApplicantRole)

	pageQueryResult, err := l.svcCtx.Repo.Appeal.GetList(
		l.ctx,
		query.AppealListQuery{
			Status:        &status,
			ApplicantID:   &req.ApplicantID,
			OrderID:       &req.OrderID,
			ApplicantRole: &applicantRole,
			Page:          req.Page,
			Size:          req.Size,
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

	// 3. 转换给前端的 List
	list := make([]types.AppealItem, 0, len(dbList))
	for i := range dbList {
		list = append(list, *types.ModelAppealToAppealItem(&dbList[i]))
	}

	return &types.GetAppealListResponse{
		Total: pageQueryResult.Total,
		List:  list,
	}, nil
}

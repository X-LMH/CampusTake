// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package appeal

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

type GetMyAppealListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyAppealListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyAppealListLogic {
	return &GetMyAppealListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyAppealListLogic) GetMyAppealList(req *types.GetMyAppealListRequest) (resp *types.GetMyAppealListResponse, err error) {
	userID := ctxx.MustUserID(l.ctx)
	status := enums.AppealStatus(req.Status)

	// 1. 调用 Repo 获取分页结果 (Repo 返回的是 []model.Appeal)
	pageQueryResult, err := l.svcCtx.Repo.Appeal.GetList(
		l.ctx,
		query.AppealListQuery{
			Status:      &status,
			ApplicantID: &userID,
			Page:        req.Page,
			Size:        req.Size,
		},
	)
	if err != nil {
		l.Errorf("获取我的申诉列表失败 userID=%d err=%v", userID, err)
		return nil, err
	}

	// 2. 类型断言
	dbList, ok := pageQueryResult.Records.([]model.Appeal)
	if !ok {
		l.Errorf("GetMyAppealList type assertion failed, actual type: %T", pageQueryResult.Records)
		return nil, errs.ErrServiceError
	}

	// 3. 初始化值切片，容量与 dbList 一致（满足你响应体中的 []AppealItem）
	list := make([]types.AppealItem, 0, len(dbList))

	// 4. 遍历并进行类型转换
	for i := range dbList {
		// &dbList[i] 取出真正的 model 指针传入转换函数
		itemPtr := types.ModelAppealToAppealItem(&dbList[i])

		// 防空保护：只要转换出来的指针不为 nil，就通过 * 符号解引用变成值类型，追加进切片
		if itemPtr != nil {
			list = append(list, *itemPtr)
		}
	}

	// 5. 返回固定格式的响应体
	return &types.GetMyAppealListResponse{
		Total: pageQueryResult.Total,
		List:  list,
	}, nil
}

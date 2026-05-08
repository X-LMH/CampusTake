// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/internal/enums"
	"CampusTake/internal/model"
	"CampusTake/pkg/errors"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRiderApplyListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRiderApplyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRiderApplyListLogic {
	return &GetRiderApplyListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRiderApplyListLogic) GetRiderApplyList(req *types.GetRiderApplyListRequest) (resp *types.GetRiderApplyListResponse, err error) {
	// 1. 调用 Repo 获取封装好的 PageResult
	// 这里透传 req.Status, req.Page, req.PageSize
	pageResult, err := l.svcCtx.Repo.Rider().GetProfileList(l.ctx, enums.RiderAuditStatus(req.Status), req.Page, req.Size)
	if err != nil {
		l.Errorf("查询骑手申请列表失败, err: %v", err)
		return nil, err
	}

	// 2. 将 Records (interface{}) 断言回具体的 Model 切片
	list, ok := pageResult.Records.([]*model.RiderProfile)
	if !ok {
		l.Error("分页数据类型断言失败")
		return nil, errors.NewDefaultError("服务器开小差了，请稍后再试")
	}

	// 3. 执行 DTO (Data Transfer Object) 转换
	resList := make([]types.DetailApplyRiderInfo, 0, len(list))
	for _, item := range list {
		resList = append(resList, types.DetailApplyRiderInfo{
			ID:                 item.ID,
			UserID:             item.UserID,
			RealName:           item.RealName,
			StudentNo:          item.StudentNo,
			IDCardNo:           item.IDCardNo,
			DormitoryBuilding:  item.DormitoryBuilding,
			DormitoryRoom:      item.DormitoryRoom,
			CampusCardFrontURL: item.CampusCardFront,
			CampusCardBackURL:  item.CampusCardBack,
			Status:             int8(item.AuditStatus), // 状态码转换
			AuditRemark:        item.AuditRemark,       // 之前的拒绝理由
			UpdatedAt:          item.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 4. 返回 API 定义的响应格式
	return &types.GetRiderApplyListResponse{
		ApplyList: resList,
		Total:     pageResult.Total,
	}, nil
}

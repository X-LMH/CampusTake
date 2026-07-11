// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/pkg/ctxx"
	"context"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetApplyStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplyStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplyStatusLogic {
	return &GetApplyStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApplyStatusLogic) GetApplyStatus() (resp *types.BaseApplyRiderInfo, err error) {
	userID := ctxx.MustUserID(l.ctx)
	profile, err := l.svcCtx.Repo.Rider.GetProfileByUserID(l.ctx, userID)
	if err != nil {
		l.Errorf("查询骑手申请状态失败，userID=%d，err=%v", userID, err)
		return nil, err
	}

	campusCardFrontURL := l.svcCtx.Config.UploadConfig.UrlPrefix + profile.CampusCardFront
	campusCardBackURL := l.svcCtx.Config.UploadConfig.UrlPrefix + profile.CampusCardBack

	return &types.BaseApplyRiderInfo{
		Status:             profile.AuditStatus.String(),
		AuditRemark:        profile.AuditRemark,
		RealName:           profile.RealName,
		StudentNo:          profile.StudentNo,
		IDCardNo:           profile.IDCardNo,
		DormitoryBuilding:  profile.DormitoryBuilding,
		DormitoryRoom:      profile.DormitoryRoom,
		CampusCardFrontURL: campusCardFrontURL,
		CampusCardBackURL:  campusCardBackURL,
	}, nil
}

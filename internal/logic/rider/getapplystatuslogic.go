// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package rider

import (
	"CampusTake/common/ctxx"
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
	profile, err := l.svcCtx.Repo.Rider().GetProfileByUserID(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	campusCardFrontURL := l.svcCtx.Config.Upload.UrlPrefix + profile.CampusCardFront
	campusCardBackURL := l.svcCtx.Config.Upload.UrlPrefix + profile.CampusCardBack

	return &types.BaseApplyRiderInfo{
		Status:             int8(profile.AuditStatus),
		AuditRemark:        profile.AuditRemark,
		RealName:           profile.RealName,
		StudentNo:          profile.StudentNo,
		IdCardNo:           profile.IDCardNo,
		DormitoryBuilding:  profile.DormitoryBuilding,
		DormitoryRoom:      profile.DormitoryRoom,
		CampusCardFrontURL: campusCardFrontURL,
		CampusCardBackURL:  campusCardBackURL,
	}, nil
}

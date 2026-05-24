// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package profile

import (
	"CampusTake/pkg/ctxx"
	"context"
	"strings"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileRequest) (resp *types.UserProfileResponse, err error) {
	userID := ctxx.MustUserID(l.ctx)

	err = l.svcCtx.Repo.User.UpdateProfileByID(l.ctx, userID, req.Nickname, req.Gender)
	if err != nil {
		return nil, err
	}

	user, err := l.svcCtx.Repo.User.GetByID(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	urlPrefix := strings.TrimRight(l.svcCtx.Config.UploadConfig.UrlPrefix, "/")

	avatar := ""
	if user.Avatar != "" {
		avatar = urlPrefix + user.Avatar
	}

	return &types.UserProfileResponse{
		Phone:    user.Phone,
		Nickname: user.Nickname,
		Avatar:   avatar,
		Gender:   user.Gender,
	}, nil
}

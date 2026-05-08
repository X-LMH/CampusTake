package profile

import (
	"CampusTake/pkg/ctxx"
	"context"
	"strings"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile() (*types.UserProfileResponse, error) {
	userID := ctxx.MustUserID(l.ctx)

	user, err := l.svcCtx.Repo.User().GetByID(l.ctx, userID)
	if err != nil {
		return nil, err
	}

	// 拼接完整 URL
	urlPrefix := strings.TrimRight(l.svcCtx.Config.Upload.UrlPrefix, "/")
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

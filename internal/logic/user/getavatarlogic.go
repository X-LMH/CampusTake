// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"
	"os"
	"path/filepath"

	"CampusTake/internal/svc"
	"CampusTake/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvatarLogic {
	return &GetAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvatarLogic) GetAvatar(req *types.GetAvatarRequest) (string, error) {
	user, err := l.svcCtx.Repo.User().GetByID(l.ctx, req.UserID)
	if err != nil {
		return "", err
	}

	// 1️⃣ 没有头像 → 返回默认头像
	if user.Avatar == "" {
		return filepath.Join(l.svcCtx.Config.Upload.AvatarPath, "default.png"), nil
	}

	// 2️⃣ 从 URL 提取文件名
	fileName := filepath.Base(user.Avatar)
	if fileName == "." || fileName == "/" {
		return filepath.Join(l.svcCtx.Config.Upload.AvatarPath, "default.png"), nil
	}

	// 3️⃣ 拼接本地真实路径（⚠️ 这里只用 fileName）
	avatarPath := filepath.Join(l.svcCtx.Config.Upload.AvatarPath, fileName)

	// 4️⃣ 检查文件是否存在
	if _, err := os.Stat(avatarPath); err != nil {
		if os.IsNotExist(err) {
			logx.Infof("avatar not found, user_id=%d, path=%s", req.UserID, avatarPath)

			// 👉 推荐：返回默认头像，而不是报错
			return filepath.Join(l.svcCtx.Config.Upload.AvatarPath, "default.png"), nil
		}

		logx.Errorf("stat avatar failed: %v, user_id=%d, path=%s", err, req.UserID, avatarPath)
		return "", err
	}

	return avatarPath, nil
}

package profile

import (
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	"CampusTake/pkg/ctxx"
	"CampusTake/pkg/errors"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAvatarLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateAvatarLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAvatarLogic {
	return &UpdateAvatarLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateAvatarLogic) UpdateAvatar(file multipart.File, header *multipart.FileHeader) (*types.UpdateAvatarResponse, error) {

	// 1️. 用户ID
	userID := ctxx.MustUserID(l.ctx)

	// 2️. 限制大小（2MB）
	if header.Size > 2*1024*1024 {
		return nil, errors.ErrFileTooLarge
	}

	// 3️. 检测真实类型
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		return nil, errors.ErrFileFormatError
	}

	// ⚠️ 重置指针
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// 4️. 获取扩展名
	ext := getExtFromContentType(contentType)
	if ext == "" {
		return nil, errors.ErrFileFormatError
	}

	// 5️. 文件名（防缓存）
	filename := fmt.Sprintf("avatar_%d_%d%s", userID, time.Now().Unix(), ext)

	// 6️. 创建目录
	dir := l.svcCtx.Config.Upload.AvatarPath
	if err = os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// 7️. 保存文件
	filePath := filepath.Join(dir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return nil, err
	}

	// 8️. 生成 URL
	urlPrefix := strings.TrimRight(l.svcCtx.Config.Upload.UrlPrefix, "/")
	avatarUrl := urlPrefix + l.svcCtx.Config.Upload.AvatarPathPrefix + "/" + filename

	// 9️. 存数据库（只存相对路径）
	avatarDBUrl := l.svcCtx.Config.Upload.AvatarPathPrefix + "/" + filename

	if err = l.svcCtx.Repo.User().UpdateAvatarByID(l.ctx, userID, avatarDBUrl); err != nil {
		return nil, err
	}

	return &types.UpdateAvatarResponse{
		AvatarUrl: avatarUrl,
	}, nil
}

// 辅助函数
func getExtFromContentType(contentType string) string {
	switch contentType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	default:
		return ""
	}
}

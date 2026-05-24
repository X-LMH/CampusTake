package upload

import (
	"CampusTake/internal/svc"
	"CampusTake/internal/types"
	errs "CampusTake/pkg/errors"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadImageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadImageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadImageLogic {
	return &UploadImageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadImageLogic) UploadImage(r *http.Request, req *types.UploadImageRequest) (*types.UploadImageResponse, error) {

	// 1. 获取文件
	file, header, err := r.FormFile("file")
	if err != nil {
		logx.Errorf("获取上传文件失败，type=%s，err=%v", req.Type, err)
		return nil, err
	}
	defer func() { _ = file.Close() }()

	// 2. 校验图片
	if err = l.validateImage(header); err != nil {
		logx.Errorf("上传图片校验失败，filename=%s，type=%s，err=%v", header.Filename, req.Type, err)
		return nil, err
	}

	// 3. 获取上传目录
	savePath, accessPath, err := l.getUploadPath(req.Type)
	if err != nil {
		logx.Errorf("获取上传目录失败，type=%s，err=%v", req.Type, err)
		return nil, err
	}

	// 4. 创建目录
	if err = os.MkdirAll(savePath, os.ModePerm); err != nil {
		logx.Errorf("创建上传目录失败，path=%s，err=%v", savePath, err)
		return nil, err
	}

	// 5. 生成文件名
	ext := filepath.Ext(header.Filename)

	filename := fmt.Sprintf(
		"%d_%s%s",
		time.Now().Unix(),
		uuid.NewString(),
		ext,
	)

	fullPath := filepath.Join(savePath, filename)

	// 6. 保存文件
	if err = saveUploadedFile(file, fullPath); err != nil {
		logx.Errorf("保存上传文件失败，path=%s，err=%v", fullPath, err)
		return nil, err
	}

	// 7. 拼接 path
	path := fmt.Sprintf(
		"%s/%s",
		accessPath,
		filename,
	)

	// 8. 拼接完整URL
	url := fmt.Sprintf(
		"%s%s",
		l.svcCtx.Config.UploadConfig.UrlPrefix,
		path,
	)

	return &types.UploadImageResponse{
		Path: path,
		Url:  url,
	}, nil
}
func (l *UploadImageLogic) validateImage(header *multipart.FileHeader) error {

	// 限制大小
	maxSize := l.svcCtx.Config.UploadConfig.MaxImageSizeMB * 1024 * 1024

	if header.Size > maxSize {
		logx.Errorf("图片大小超出限制，filename=%s，size=%d，maxSize=%d", header.Filename, header.Size, maxSize)
		return errs.ErrImageTooLarge
	}

	// 校验后缀
	ext := strings.ToLower(filepath.Ext(header.Filename))

	allowExt := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
	}

	if !allowExt[ext] {
		logx.Errorf("图片格式不支持，filename=%s，ext=%s", header.Filename, ext)
		return errs.ErrImageFormatInvalid
	}

	return nil
}
func (l *UploadImageLogic) getUploadPath(uploadType string) (string, string, error) {

	switch uploadType {

	case "avatar":
		return l.svcCtx.Config.UploadConfig.Avatar.Path,
			l.svcCtx.Config.UploadConfig.Avatar.PathPrefix,
			nil

	case "appeal":
		return l.svcCtx.Config.UploadConfig.Appeal.Path,
			l.svcCtx.Config.UploadConfig.Appeal.PathPrefix,
			nil

	case "campus_card":
		return l.svcCtx.Config.UploadConfig.CampusCard.Path,
			l.svcCtx.Config.UploadConfig.CampusCard.PathPrefix,
			nil

	default:
		logx.Errorf("图片类型不支持，uploadType=%s", uploadType)
		return "", "", errs.ErrImageTypeInvalid
	}
}
func saveUploadedFile(file multipart.File, dst string) error {

	out, err := os.Create(dst)
	if err != nil {
		logx.Errorf("创建上传文件失败，path=%s，err=%v", dst, err)
		return err
	}
	defer func() { _ = out.Close() }()

	_, err = io.Copy(out, file)
	if err != nil {
		logx.Errorf("拷贝上传文件内容失败，path=%s，err=%v", dst, err)
	}

	return err
}

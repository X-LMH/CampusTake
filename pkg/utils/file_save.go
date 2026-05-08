package utils

import (
	"CampusTake/pkg/errors"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
)

const maxFileSize = 2 * 1024 * 1024

func SaveFileToLocal(file multipart.File, header *multipart.FileHeader, saveDir string, fileName string) (string, error) {
	if header.Size > maxFileSize {
		return "", errors.ErrFileTooLarge
	}

	// 【修改点】重置文件偏移量，防止在其他地方被读取过导致 Copy 失败
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return "", err
	}

	ext := path.Ext(header.Filename)
	finalName := fileName + ext
	finalPath := filepath.Join(saveDir, finalName) // 本地存储路径用 filepath

	dst, err := os.Create(finalPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return finalName, nil
}

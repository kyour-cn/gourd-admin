package services

import (
	"app/internal/modules/upload"
	"context"
	"mime/multipart"
	"path/filepath"
)

func NewFileService(ctx context.Context) *FileService {
	return &FileService{
		ctx: ctx,
	}
}

type FileService struct {
	ctx context.Context
}

func (s *FileService) CloudUpload(file multipart.File, handler *multipart.FileHeader, group string) (*upload.Output, error) {
	// 检查文件类型（可选）
	ext := filepath.Ext(handler.Filename)

	// 保存路径 按日期分目录，避免单目录文件过多
	savePath := upload.GenPath(group, ext)
	input := upload.Input{
		FileName: handler.Filename,
		Content:  file,
		Size:     handler.Size,
		Ext:      ext,
	}

	uploader, err := upload.GetUploader("")
	if err != nil {
		return nil, err
	}
	return uploader.Upload(s.ctx, input, savePath)
}

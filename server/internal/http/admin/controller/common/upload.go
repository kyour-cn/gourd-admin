package common

import (
	"app/internal/modules/common/upload"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Upload 上传
type Upload struct {
	Base //继承基础控制器
}

// Image 上传图片
func (c *Upload) Image(w http.ResponseWriter, r *http.Request) {
	// 限制上传大小：最大 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		_ = c.Fail(w, 500, "获取请求失败", err.Error())
		return
	}

	// 获取上传文件，字段名为 "file"
	file, handler, err := r.FormFile("file")
	if err != nil {
		_ = c.Fail(w, 501, "上传文件失败", err.Error())
		return
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)

	// 检查文件类型（可选）
	ext := filepath.Ext(handler.Filename)
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		_ = c.Fail(w, 501, "仅支持jpg/jpeg/png/gif格式", nil)
		return
	}

	// 保存路径
	savePath := filepath.Join("uploads", handler.Filename)

	input := upload.Input{
		FileName: handler.Filename,
		Content:  file,
		Size:     handler.Size,
	}

	uploader, err := upload.GetUploader("")
	if err != nil {
		_ = c.Fail(w, 502, "获取上传器失败", err.Error())
		return
	}

	output, err := uploader.Upload(r.Context(), input, "")
	if err != nil {
		_ = c.Fail(w, 502, "上传图片失败", err.Error())
		return
	}

	_ = c.Success(w, "上传图片成功", output)
}

// File 上传文件
func (c *Upload) File(w http.ResponseWriter, r *http.Request) {

	_ = c.Success(w, "上传图片成功", nil)
}

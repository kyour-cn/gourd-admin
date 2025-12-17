package controller

import (
	"app/internal/http/common/services"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"app/internal/modules/upload"
	"app/internal/orm/model"
	"app/internal/orm/query"
)

// Upload 上传
type Upload struct {
	Base //继承基础控制器
}

// Image 上传图片
func (c *Upload) Image(w http.ResponseWriter, r *http.Request) {
	claims, err := c.GetJwt(r)
	if err != nil {
		_ = c.Fail(w, 101, err.Error(), nil)
		return
	}

	// 限制上传大小：最大 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
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
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		_ = c.Fail(w, 501, "仅支持jpg/jpeg/png/gif格式", nil)
		return
	}

	// 检查文件类型是否为图片
	mimeType, err := upload.GetFileMimeType(file)
	if err != nil {
		_ = c.Fail(w, 502, "获取文件类型失败", err.Error())
		return
	}
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		_ = c.Fail(w, 501, "仅支持jpg/jpeg/png格式", nil)
		return
	}

	// 保存路径 按日期分目录，避免单目录文件过多
	service := services.NewFileService(r.Context())
	output, err := service.CloudUpload(file, handler, "images")
	if err != nil {
		_ = c.Fail(w, 502, "上传图片失败", err.Error())
		return
	}

	// 保存图片信息到数据库
	err = query.File.Create(&model.File{
		FileName:   handler.Filename,
		FileType:   mimeType, // 如 image/png
		FileExt:    ext,
		FileSize:   handler.Size,
		URL:        output.URL,
		FilePath:   output.Path,
		StorageID:  output.StorageID,
		StorageKey: output.Storage,
		HashMd5:    output.Hash,
		UserID:     claims.Sub,
	})
	if err != nil {
		_ = c.Fail(w, 503, "保存图片信息失败", err.Error())
		return
	}
	_ = c.Success(w, "上传图片成功", output)
}

// File 上传文件
func (c *Upload) File(w http.ResponseWriter, r *http.Request) {
	claims, err := c.GetJwt(r)
	if err != nil {
		_ = c.Fail(w, 101, err.Error(), nil)
		return
	}

	// 限制上传大小：最大 10MB
	if err := r.ParseMultipartForm(10 << 20); err != nil {
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
	//if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
	//	_ = c.Fail(w, 501, "仅支持jpg/jpeg/png/gif格式", nil)
	//	return
	//}

	// 检查文件类型是否允许上传
	mimeType, err := upload.GetFileMimeType(file)
	if err != nil {
		_ = c.Fail(w, 502, "获取文件类型失败", err.Error())
		return
	}
	//if mimeType != "image/jpeg" && mimeType != "image/png" {
	//	_ = c.Fail(w, 501, "仅支持jpg/jpeg/png格式", nil)
	//	return
	//}

	// 保存路径 按日期分目录，避免单目录文件过多
	service := services.NewFileService(r.Context())
	output, err := service.CloudUpload(file, handler, "files")
	if err != nil {
		_ = c.Fail(w, 502, "上传图片失败", err.Error())
		return
	}

	// 保存图片信息到数据库
	err = query.File.Create(&model.File{
		FileName:   handler.Filename,
		FileType:   mimeType, // 如 image/png
		FileExt:    ext,
		FileSize:   handler.Size,
		URL:        output.URL,
		FilePath:   output.Path,
		StorageID:  output.StorageID,
		StorageKey: output.Storage,
		HashMd5:    output.Hash,
		UserID:     claims.Sub,
	})
	if err != nil {
		_ = c.Fail(w, 503, "保存图片信息失败", err.Error())
		return
	}
	_ = c.Success(w, "上传图片成功", output)
}

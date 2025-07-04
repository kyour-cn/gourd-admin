package common

import (
	"app/internal/modules/common/upload"
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

// Upload 上传
type Upload struct {
	Base //继承基础控制器
}

// Image 上传图片
func (c *Upload) Image(w http.ResponseWriter, r *http.Request) {
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
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
		_ = c.Fail(w, 501, "仅支持jpg/jpeg/png/gif格式", nil)
		return
	}

	// 保存路径 按日期分目录，避免单目录文件过多
	dir := filepath.Join("uploads", time.Now().Format("200601/02"))

	// 生成文件名，避免重复
	fileName := handler.Filename + uuid.New().String()

	// 计算文件名的MD5值
	md5Hash := md5.Sum([]byte(fileName))
	md5String := hex.EncodeToString(md5Hash[:])
	savePath := filepath.ToSlash(filepath.Join(dir, md5String+ext))

	input := upload.Input{
		FileName: handler.Filename,
		Content:  file,
		Size:     handler.Size,
		Ext:      ext,
	}

	uploader, err := upload.GetUploader("")
	if err != nil {
		_ = c.Fail(w, 502, "获取上传器失败", err.Error())
		return
	}

	// 开始上传
	output, err := uploader.Upload(r.Context(), input, savePath)
	if err != nil {
		_ = c.Fail(w, 502, "上传图片失败", err.Error())
		return
	}

	_ = c.Success(w, "上传图片成功", output)
}

// File 上传文件
func (c *Upload) File(w http.ResponseWriter, r *http.Request) {
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
	//if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
	//	_ = c.Fail(w, 501, "仅支持jpg/jpeg/png/gif格式", nil)
	//	return
	//}

	// 保存路径 按日期分目录，避免单目录文件过多
	dir := filepath.Join("uploads", time.Now().Format("200601/02"))

	// 生成文件名，避免重复
	fileName := handler.Filename + uuid.New().String()

	// 计算文件名的MD5值
	md5Hash := md5.Sum([]byte(fileName))
	md5String := hex.EncodeToString(md5Hash[:])
	savePath := filepath.ToSlash(filepath.Join(dir, md5String+ext))

	input := upload.Input{
		FileName: handler.Filename,
		Content:  file,
		Size:     handler.Size,
		Ext:      ext,
	}

	uploader, err := upload.GetUploader("")
	if err != nil {
		_ = c.Fail(w, 502, "获取上传器失败", err.Error())
		return
	}

	// 开始上传
	output, err := uploader.Upload(r.Context(), input, savePath)
	if err != nil {
		_ = c.Fail(w, 502, "上传图片失败", err.Error())
		return
	}

	_ = c.Success(w, "上传文件成功", output)
}

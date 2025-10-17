package upload

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"gorm.io/gen"

	"app/internal/orm/model"
	"app/internal/orm/query"
)

// Input 上传请求参数
type Input struct {
	FileName string    // 文件名（带后缀）
	Content  io.Reader // 文件内容
	Size     int64     // 文件大小
	Ext      string    // 文件后缀
}

// Output 上传返回结果
type Output struct {
	URL       string `json:"url"`       // 可访问链接
	Path      string `json:"path"`      // 存储路径
	Ext       string `json:"ext"`       // 存储路径
	FileName  string `json:"fileName"`  // 文件名（带后缀）
	Hash      string `json:"hash"`      // 文件hash（如md5）
	Storage   string `json:"storage"`   // 存储类型
	StorageID int64  `json:"storageID"` // 存储类型
}

// Uploader 接口定义了上传和删除文件的方法
type Uploader interface {
	Upload(ctx context.Context, input Input, savePath string) (*Output, error)
	Delete(ctx context.Context, path string) error
	GetStorageModel() *model.FileStorage // 获取存储模型
}

// GetUploader 获取上传器
// key: 上传器标识，默认为空时取默认
func GetUploader(key string) (Uploader, error) {

	qfs := query.FileStorage

	var condition = []gen.Condition{
		qfs.Status.Eq(1),
	}

	if key == "" {
		// 取默认
		condition = append(condition, qfs.IsDefault.Eq(1))
	} else {
		condition = append(condition, qfs.Key.Eq(key))
	}

	storage, err := query.FileStorage.
		Where(condition...).
		First()
	if err != nil {
		return nil, fmt.Errorf("获取上传器失败: %w", err)
	}

	switch storage.Key {
	case "local":
		return NewLocalUploader(storage), nil
	default:
		return nil, fmt.Errorf("不支持的上传类型: %s", storage.Key)
	}
}

// GenPath 生成存储路径
func GenPath(group string, ext string) string {
	root := "uploads" // 根目录
	if group != "" {
		root = filepath.Join(root, group)
	}
	t := time.Now()
	// 保存路径 按日期分目录，避免单目录文件过多
	dir := filepath.Join(root, t.Format("200601/02"))

	fileName := t.Format("0405") + uuid.New().String() + ext
	savePath := filepath.ToSlash(filepath.Join(dir, fileName))
	return savePath
}

// GetFileMimeType 获取文件的MIME类型
func GetFileMimeType(file multipart.File) (string, error) {
	_, _ = file.Seek(0, 0)
	// 读取前512字节来检测文件类型
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %w", err)
	}
	_, _ = file.Seek(0, 0)
	return http.DetectContentType(buf), nil
}

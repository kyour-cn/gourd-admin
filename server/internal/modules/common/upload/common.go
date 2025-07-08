package upload

import (
	"app/internal/orm/model"
	"app/internal/orm/query"
	"context"
	"fmt"
	"gorm.io/gen"
	"io"
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
	URL      string `json:"url"`      // 可访问链接
	Path     string `json:"path"`     // 存储路径
	FileName string `json:"fileName"` // 文件名（带后缀）
	Hash     string `json:"hash"`     // 文件hash（如md5）
	Storage  string `json:"storage"`  // 存储类型
}

// Uploader 接口定义了上传和删除文件的方法
type Uploader interface {
	Upload(ctx context.Context, input Input, savePath string) (*Output, error)
	Delete(path string) error
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

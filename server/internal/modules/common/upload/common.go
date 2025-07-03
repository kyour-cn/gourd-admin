package upload

import (
	"app/internal/orm/query"
	"fmt"
	"gorm.io/gen"
)

// GetUploader 获取默认上传器
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
		return NewLocalUploader(), nil
	default:
		return nil, fmt.Errorf("不支持的上传类型: %s", storage.Key)
	}
}

package upload

import (
	"context"
	"io"
)

// Input 上传请求参数
type Input struct {
	FileName string    // 文件名（带后缀）
	Content  io.Reader // 文件内容
	Size     int64     // 文件大小
	MimeType string    // 文件类型
}

// Output 上传返回结果
type Output struct {
	URL     string // 可访问链接
	Path    string // 存储路径
	Hash    string // 文件hash（如md5）
	Storage string // 存储类型
}

// Uploader 接口定义了上传和删除文件的方法
type Uploader interface {
	Upload(ctx context.Context, input Input, savePath string) (*Output, error)
	Delete(path string) error
}

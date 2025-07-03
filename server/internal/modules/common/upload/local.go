package upload

import (
	"context"
	"io"
	"os"
)

func NewLocalUploader() Uploader {
	return &LocalUploader{
		StoreKey: "local", // 本地存储的唯一标识符
	}
}

type LocalUploader struct {
	StoreKey string // 存储的唯一标识符
}

func (l LocalUploader) Upload(_ context.Context, input Input, savePath string) (*Output, error) {
	// 模拟本地上传逻辑
	// 这里可以使用os包将文件写入到指定路径
	// 例如：os.Write
	file, err := os.Create(savePath)
	if err != nil {
		return nil, err
	}
	// 拷贝文件内容
	if _, err := io.Copy(file, input.Content); err != nil {
		return nil, err
	}

	return &Output{
		URL:     "/" + savePath,
		Path:    savePath,
		Storage: "local",
	}, nil
}

func (l LocalUploader) Delete(path string) error {

	return nil
}

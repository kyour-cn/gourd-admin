package upload

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"app/internal/config"
	"app/internal/orm/model"
)

func NewLocalUploader(storage *model.FileStorage) Uploader {
	// 获取配置中静态文件目录
	baseDir := "./web/"
	conf, err := config.GetHttpConfig()
	if err == nil {
		baseDir = conf.Static
	}

	return &LocalUploader{
		StoreKey: "local", // 本地存储的唯一标识符
		BaseDir:  baseDir, // 本地存储的基础目录
		storage:  storage,
	}
}

type LocalUploader struct {
	StoreKey string // 存储的唯一标识符
	BaseDir  string // 本地存储的根目录
	storage  *model.FileStorage
}

func (u LocalUploader) GetStorageModel() *model.FileStorage {
	return u.storage
}

func (u LocalUploader) Upload(_ context.Context, input Input, savePath string) (*Output, error) {
	// 获取目录部分创建所有目录
	dir := filepath.Dir(u.BaseDir + savePath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// 创建文件
	outFile, err := os.Create(u.BaseDir + savePath)
	if err != nil {
		return nil, err
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	// 创建 md5 哈希器
	hash := md5.New()
	reader := io.TeeReader(input.Content, hash)

	// 拷贝文件内容
	if _, err := io.Copy(outFile, reader); err != nil {
		return nil, err
	}

	return &Output{
		URL:       savePath,
		Path:      savePath,
		FileName:  input.FileName,
		Storage:   u.StoreKey,
		StorageID: u.storage.ID,
		Hash:      hex.EncodeToString(hash.Sum(nil)),
	}, nil
}

func (u LocalUploader) Delete(_ context.Context, path string) error {
	fullPath := u.BaseDir + path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // 文件不存在，直接返回
	}
	if err := os.Remove(fullPath); err != nil {
		return err
	}

	return nil
}

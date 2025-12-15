package upload

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"app/internal/orm/model"
)

type QiniuConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Bucket    string `json:"bucket"`
	Endpoint  string `json:"endpoint"`
	Domain    string `json:"domain"`
}

type QiniuResponse struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

func NewQiniuUploader(storage *model.FileStorage) Uploader {
	return &QiniuUploader{
		StoreKey: "qiniu", // 本地存储的唯一标识符
		storage:  storage,
	}
}

type QiniuUploader struct {
	StoreKey string // 存储的唯一标识符
	storage  *model.FileStorage
}

func (u QiniuUploader) GetStorageModel() *model.FileStorage {
	return u.storage
}

func (u QiniuUploader) Upload(_ context.Context, input Input, savePath string) (*Output, error) {
	if u.storage.Config == nil {
		return nil, fmt.Errorf("七牛存储配置为空")
	}
	// 开始上传
	conf := QiniuConfig{}
	if err := json.Unmarshal([]byte(*u.storage.Config), &conf); err != nil {
		return nil, fmt.Errorf("解析七牛配置失败: %w", err)
	}
	res, err := u.uploadFileToCloud(conf, input.Content, savePath)
	if err != nil {
		return nil, err
	}

	return &Output{
		URL:       conf.Domain + savePath,
		Path:      savePath,
		FileName:  input.FileName,
		Storage:   u.StoreKey,
		StorageID: u.storage.ID,
		Hash:      res.Hash,
	}, nil
}

func (u QiniuUploader) Delete(_ context.Context, path string) error {
	if u.storage.Config == nil {
		return fmt.Errorf("七牛存储配置为空")
	}
	// 删除文件
	conf := QiniuConfig{}
	if err := json.Unmarshal([]byte(*u.storage.Config), &conf); err != nil {
		return fmt.Errorf("解析七牛配置失败: %w", err)
	}
	err := u.deleteFile(conf, conf.Bucket, path)
	if err != nil {
		return err
	}
	return nil
}

// 生成七牛上传凭证
func (u QiniuUploader) makeUploadToken(bucket, key string, ak, sk string, expire time.Duration) (string, error) {
	// 1. 构造上传策略
	deadline := time.Now().Add(expire).Unix()
	putPolicy := map[string]interface{}{
		"scope":    bucket + ":" + key, // 上传到指定空间且覆盖同名文件
		"deadline": deadline,
	}
	// 2. JSON 并 Base64 URL-safe 编码
	policyJSON, err := json.Marshal(putPolicy)
	if err != nil {
		return "", err
	}
	encodedPolicy := base64.URLEncoding.EncodeToString(policyJSON)
	// 3. 使用 SK 对 encodedPolicy 做 HMAC-SHA1 签名
	mac := hmac.New(sha1.New, []byte(sk))
	mac.Write([]byte(encodedPolicy))
	sign := mac.Sum(nil)
	encodedSign := base64.URLEncoding.EncodeToString(sign)
	// 4. 拼接最终 Token： AK:encodedSign:encodedPolicy
	token := fmt.Sprintf("%s:%s:%s", ak, encodedSign, encodedPolicy)
	return token, nil
}

// 上传文件到七牛
func (u QiniuUploader) uploadFileToCloud(conf QiniuConfig, src io.Reader, key string) (*QiniuResponse, error) {
	// 生成上传凭证，设定 1 小时过期
	token, err := u.makeUploadToken(conf.Bucket, key, conf.AccessKey, conf.SecretKey, time.Hour)
	if err != nil {
		return nil, err
	}

	// 构造 multipart/form-data 请求体
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	// token 字段
	if err := writer.WriteField("token", token); err != nil {
		return nil, err
	}
	// key 字段：对象存储中的名字
	if err := writer.WriteField("key", key); err != nil {
		return nil, err
	}
	// file 字段：文件内容
	part, err := writer.CreateFormFile("file", "")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, src); err != nil {
		return nil, err
	}
	// 关闭 writer，以写入结尾 boundary
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// 发送 HTTP POST 请求
	req, err := http.NewRequest("POST", "https://"+conf.Endpoint, &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	// 解析响应
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upload failed, status: %s, body: %s", resp.Status, string(respBody))
	}

	res := &QiniuResponse{}
	if err := json.Unmarshal(respBody, res); err != nil {
		return nil, fmt.Errorf("解析七牛上传响应失败: %w", err)
	}

	return res, nil
}

func (u QiniuUploader) makeAuthorization(conf QiniuConfig, method, path, host, contentType string, body []byte) string {
	signingStr := fmt.Sprintf("%s %s\nHost: %s", method, path, host)

	if contentType != "" {
		signingStr += fmt.Sprintf("\nContent-Type: %s", contentType)
	}

	signingStr += "\n\n" // 两个换行
	if len(body) > 0 && contentType != "application/octet-stream" {
		signingStr += string(body)
	}

	mac := hmac.New(sha1.New, []byte(conf.SecretKey))
	mac.Write([]byte(signingStr))
	sign := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	return fmt.Sprintf("Qiniu %s:%s", conf.AccessKey, sign)
}

// 删除文件：bucket 是空间名，key 是文件名
func (u QiniuUploader) deleteFile(conf QiniuConfig, bucket, key string) error {
	// Step 1: 计算 EncodedEntryURI
	entryURI := fmt.Sprintf("%s:%s", bucket, key)
	encodedEntryURI := base64.URLEncoding.EncodeToString([]byte(entryURI))

	// Step 2: 构造请求信息
	method := "POST"
	path := fmt.Sprintf("/delete/%s", encodedEntryURI)
	host := "rs.qiniuapi.com"
	contentType := "application/x-www-form-urlencoded"
	urlStr := "http://" + host + path
	var body []byte

	// Step 3: 构造请求
	req, err := http.NewRequest(method, urlStr, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Host", host)
	req.Header.Set("Content-Type", contentType)

	// Step 4: 构造签名
	auth := u.makeAuthorization(conf, method, path, host, contentType, body)
	req.Header.Set("Authorization", auth)

	// Step 5: 发起请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("删除失败，状态码：%d", resp.StatusCode)
}

package common

import "net/http"

// Upload 上传
type Upload struct {
	Base //继承基础控制器
}

// Image 上传图片
func (c *Upload) Image(w http.ResponseWriter, r *http.Request) {

}

// File 上传文件
func (c *Upload) File(w http.ResponseWriter, r *http.Request) {

}

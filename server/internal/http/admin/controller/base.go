package controller

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strconv"
)

// Base 基础控制器
// 所有控制器都应继承此控制器，可以在此控制器中定义公共方法
type Base struct {
}

// Success 成功时响应
func (*Base) Success(w http.ResponseWriter, message string, data any) (err error) {
	if message == "" {
		message = "success"
	}
	res := map[string]any{
		"code":    0,
		"data":    data,
		"message": message,
	}

	str, _ := json.Marshal(res)
	_, err = w.Write(str)
	return
}

// Fail 失败响应
func (*Base) Fail(w http.ResponseWriter, code int, message string, data any) (err error) {
	if message == "" {
		message = "fail"
	}
	res := map[string]any{
		"code":    code,
		"data":    data,
		"message": message,
	}

	str, _ := json.Marshal(res)
	_, err = w.Write(str)
	return
}

// JsonReqUnmarshal 解析json请求参数
func (*Base) JsonReqUnmarshal(r *http.Request, req any) error {
	// 解析json参数
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	// 校验参数
	return validator.New().Struct(req)
}

// PageParam 获取分页参数
func (*Base) PageParam(r *http.Request, defaultPage int, defaultSize int) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page == 0 {
		page = defaultPage
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize == 0 {
		pageSize = defaultSize
	}
	return page, pageSize
}

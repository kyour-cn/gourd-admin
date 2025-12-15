package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
	json "github.com/json-iterator/go"

	"app/internal/http/common/dto"
)

// Base 基础控制器
// 所有控制器都应继承此控制器，可以在此控制器中定义公共方法
type Base struct{}

type Response struct {
	Code    int    `json:"code"`
	Data    any    `json:"data"`
	Message string `json:"message"`
}

var (
	decoder = form.NewDecoder()
)

// Success 成功时响应
func (*Base) Success(w http.ResponseWriter, message string, data any) (err error) {
	if message == "" {
		message = "success"
	}
	res := Response{
		Code:    0,
		Data:    data,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	str, _ := json.Marshal(res)
	_, err = w.Write(str)
	return
}

// Fail 失败响应
func (*Base) Fail(w http.ResponseWriter, code int, message string, data any) (err error) {
	if message == "" {
		message = "fail"
	}
	res := Response{
		Code:    code,
		Data:    data,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	str, _ := json.Marshal(res)
	_, err = w.Write(str)
	return
}

// JsonReqUnmarshal 解析json请求参数
func (*Base) JsonReqUnmarshal(r *http.Request, req any) error {
	// 解析json参数
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	v := validate.Struct(req)
	zhcn.Register(v)

	// 执行验证
	if v.Validate() {
		return nil
	}

	// 校验参数失败 并返回第一个错误
	for _, errs := range v.Errors.All() {
		for _, err := range errs {
			return errors.New(err)
		}
	}
	return nil
}

// QueryReqUnmarshal 解析GET请求参数
func (*Base) QueryReqUnmarshal(r *http.Request, req any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	if err := decoder.Decode(req, r.Form); err != nil {
		return err
	}

	v := validate.Struct(req)
	zhcn.Register(v)

	// 执行验证
	if v.Validate() {
		return nil
	}

	// 校验参数失败 并返回第一个错误
	for _, errs := range v.Errors.All() {
		for _, err := range errs {
			return errors.New(err)
		}
	}
	return nil
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

// GetJwt 从请求上下文中获取JWT信息
func (*Base) GetJwt(r *http.Request) (*dto.UserClaims, error) {
	jwtValue := r.Context().Value("jwt")
	if _, ok := jwtValue.(dto.UserClaims); !ok {
		return nil, errors.New("获取登录信息失败")
	}
	claims := jwtValue.(dto.UserClaims)
	return &claims, nil
}

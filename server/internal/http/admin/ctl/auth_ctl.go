package ctl

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/user"
	"gourd/internal/util/captcha"
	"net/http"
)

// AuthCtl 用户控制器
type AuthCtl struct {
	common.BaseCtl //继承基础控制器
}

// Captcha 获取验证码
func (ctl *AuthCtl) Captcha(w http.ResponseWriter, r *http.Request) {

	data, err := captcha.GenerateSlide()
	if err != nil {
		_ = ctl.Fail(w, 1, "生成验证码失败："+err.Error(), nil)
		return
	}

	_ = ctl.Success(w, "", data)

}

type LoginReq struct {
	Username   string `json:"username" validate:"required,min=5,max=20"`
	Password   string `json:"password" validate:"required,min=6,max=32"`
	CaptchaKey string `json:"captcha_key" validate:"required"`
	Md5        bool   `json:"md5"`
	Point      struct {
		X int64 `json:"x"`
		Y int64 `json:"y"`
	}
}

type LoginResp struct {
	Token    string     `json:"token"`
	UserInfo model.User `json:"userInfo"`
}

// Login 登录
func (ctl *AuthCtl) Login(w http.ResponseWriter, r *http.Request) {

	req := LoginReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		_ = ctl.Fail(w, 101, "请求参数异常："+err.Error(), nil)
		return
	}

	err = validator.New().Struct(req)
	if err != nil {
		_ = ctl.Fail(w, 102, "请求参数错误", err.Error())
		return
	}

	// 验证码校验
	if !captcha.VerifySlide(req.CaptchaKey, req.Point.X, req.Point.Y) {
		_ = ctl.Fail(w, 103, "验证失败", nil)
		return
	}

	res := LoginResp{}

	_ = ctl.Success(w, "", res)
}

// Info 获取用户信息
func (ctl *AuthCtl) Info(w http.ResponseWriter, r *http.Request) {

	repository := user.Repository{
		Ctx: r.Context(),
	}

	qu := query.User

	userData, _ := repository.Query().
		Where(qu.ID.Eq(1)).
		Select(
			qu.ID,
		).
		First()

	// 响应结果
	_ = ctl.Success(w, "", userData)
}

// Add 创建用户
func (ctl *AuthCtl) Add(w http.ResponseWriter, r *http.Request) {

	repository := user.Repository{
		Ctx: r.Context(),
	}

	_ = repository.Begin()

	userData := model.User{
		Username: "go_create",
	}

	err := repository.Create(&userData)
	if err != nil {
		_ = ctl.Fail(w, 1, "添加失败："+err.Error(), nil)
		_ = repository.Rollback()
		return
	}

	_ = repository.Commit()

	// 响应结果
	_ = ctl.Success(w, "", userData)
}

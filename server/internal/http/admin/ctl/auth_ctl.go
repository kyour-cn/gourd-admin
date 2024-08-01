package ctl

import (
	"crypto/md5"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"gourd/internal/http/admin/common"
	"gourd/internal/http/admin/service"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories"
	"gourd/internal/util/captcha"
	"io"
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
	Token    string      `json:"token"`
	UserInfo *model.User `json:"userInfo"`
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

	if !req.Md5 {
		hash := md5.New()
		_, _ = io.WriteString(hash, req.Password)
		req.Password = string(hash.Sum(nil))
	}

	userData, err := service.LoginUser(req.Username, req.Password)
	if err != nil {
		_ = ctl.Fail(w, 104, "账号或密码错误", err)
		return
	}

	// 生成token
	token, err := service.GenerateToken(userData)
	if err != nil {
		_ = ctl.Fail(w, 105, "生成token失败", err.Error())
		return
	}

	res := LoginResp{
		Token:    token,
		UserInfo: userData,
	}

	_ = ctl.Success(w, "", res)
}

func (ctl *AuthCtl) GetMenu(w http.ResponseWriter, r *http.Request) {
	// 获取jwt
	token := r.Context().Value("jwt").(jwt.MapClaims)

	userId := int32(token["id"].(float64))

	ru := repositories.User{
		Ctx: r.Context(),
	}

	uq := query.User

	userInfo, err := ru.Query().
		Where(uq.ID.Eq(userId)).
		First()
	if err != nil {
		_ = ctl.Fail(w, 101, "获取用户信息失败", err.Error())
		return
	}

	menus, err := service.GetMenu(userInfo)
	if err != nil {
		_ = ctl.Fail(w, 102, "获取菜单失败", err.Error())
		return
	}

	_ = ctl.Success(w, "", menus)
}

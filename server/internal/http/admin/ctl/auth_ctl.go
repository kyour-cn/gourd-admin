package ctl

import (
	"crypto/md5"
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
func (c *AuthCtl) Captcha(w http.ResponseWriter, r *http.Request) {

	data, err := captcha.GenerateSlide()
	if err != nil {
		_ = c.Fail(w, 1, "生成验证码失败："+err.Error(), nil)
		return
	}

	_ = c.Success(w, "", data)

}

// Login 登录
func (c *AuthCtl) Login(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		Username   string `json:"username" validate:"required,min=5,max=20"`
		Password   string `json:"password" validate:"required,min=6,max=32"`
		CaptchaKey string `json:"captcha_key" validate:"required"`
		Md5        bool   `json:"md5"`
		Point      struct {
			X int64 `json:"x"`
			Y int64 `json:"y"`
		}
	}

	type Resp struct {
		Token    string      `json:"token"`
		UserInfo *model.User `json:"userInfo"`
	}

	req := Req{}
	err := c.JsonReqUnmarshal(r, &req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	// 验证码校验
	if !captcha.VerifySlide(req.CaptchaKey, req.Point.X, req.Point.Y) {
		_ = c.Fail(w, 102, "验证失败", nil)
		return
	}

	if !req.Md5 {
		hash := md5.New()
		_, _ = io.WriteString(hash, req.Password)
		req.Password = string(hash.Sum(nil))
	}

	userData, err := service.LoginUser(req.Username, req.Password)
	if err != nil {
		_ = c.Fail(w, 103, "账号或密码错误", err)
		return
	}
	if userData.Status != 1 {
		_ = c.Fail(w, 104, "账号异常或被锁定", err)
		return
	}

	rr := repositories.Role{}
	roleData, err := rr.Query().
		Where(
			query.Role.ID.Eq(userData.RoleID),
			query.Role.Status.Eq(1),
		).
		First()
	if err != nil {
		_ = c.Fail(w, 104, "账号角色异常,请联系管理员", err)
		return
	}

	// 生成token
	tokenData := map[string]any{
		"id":   userData.ID,
		"role": userData.RoleID,
		"app":  roleData.AppID,
	}
	token, err := service.GenerateToken(tokenData)
	if err != nil {
		_ = c.Fail(w, 105, "生成token失败", err.Error())
		return
	}

	res := Resp{
		Token:    token,
		UserInfo: userData,
	}

	_ = c.Success(w, "", res)
}

func (c *AuthCtl) GetMenu(w http.ResponseWriter, r *http.Request) {
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
		_ = c.Fail(w, 101, "获取用户信息失败", err.Error())
		return
	}

	menus, err := service.GetMenu(userInfo)
	if err != nil {
		_ = c.Fail(w, 102, "获取菜单失败", err.Error())
		return
	}

	_ = c.Success(w, "", menus)
}

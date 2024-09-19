package ctl

import (
	"crypto/md5"
	"encoding/hex"
	"gourd/internal/config"
	"gourd/internal/http/admin/common"
	"gourd/internal/http/admin/service"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/util/captcha"
	"net/http"
)

// AuthCtl 用户控制器
type AuthCtl struct {
	common.BaseCtl //继承基础控制器
}

// Captcha 获取验证码
func (c *AuthCtl) Captcha(w http.ResponseWriter, _ *http.Request) {
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
		Expire   int64       `json:"expire"`
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
		hash := md5.Sum([]byte(req.Password))
		req.Password = hex.EncodeToString(hash[:])
	}

	userData, err := service.LoginUser(r.Context(), req.Username, req.Password)
	if err != nil {
		_ = c.Fail(w, 103, "登录失败："+err.Error(), "")
		return
	}
	if userData.Status != 1 {
		_ = c.Fail(w, 104, "账号异常或被锁定", err)
		return
	}

	roleData, err := query.Role.WithContext(r.Context()).
		Where(
			query.Role.ID.Eq(userData.RoleID),
			query.Role.Status.Eq(1),
		).
		First()
	if err != nil {
		_ = c.Fail(w, 104, "账号角色异常,请联系管理员", err)
		return
	}

	jwtConf, err := config.GetJwtConfig()
	if err != nil {
		_ = c.Fail(w, 104, "token配置异常,请联系管理员", err)
	}
	// 生成token
	claims := service.UserClaims{
		Uid:    userData.ID,
		RoleId: userData.RoleID,
		AppId:  roleData.AppID,
	}
	token, err := service.GenerateToken(claims)
	if err != nil {
		_ = c.Fail(w, 105, "生成token失败", err.Error())
		return
	}

	res := Resp{
		Token:    token,
		UserInfo: userData,
		Expire:   jwtConf.Expire,
	}

	_ = c.Success(w, "", res)
}

func (c *AuthCtl) GetMenu(w http.ResponseWriter, r *http.Request) {
	// 获取jwt并解析
	jwtValue := r.Context().Value("jwt")
	if _, ok := jwtValue.(service.UserClaims); !ok {
		_ = c.Fail(w, 101, "获取登录信息失败", "jwt信息不正确")
		return
	}
	claims := jwtValue.(service.UserClaims)

	uq := query.User

	userInfo, err := uq.WithContext(r.Context()).
		Where(uq.ID.Eq(claims.Uid)).
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

package common

import (
	"app/internal/config"
	"app/internal/modules/auth"
	"app/internal/modules/dblog"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/captcha"
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

// Auth 用户控制器
type Auth struct {
	Base //继承基础控制器
}

// Captcha 获取验证码
func (c *Auth) Captcha(w http.ResponseWriter, _ *http.Request) {
	data, err := captcha.GenerateSlide()
	if err != nil {
		_ = c.Fail(w, 1, "生成验证码失败："+err.Error(), nil)
		return
	}

	_ = c.Success(w, "", data)
}

// Login 登录
func (c *Auth) Login(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Username   string `json:"username" validate:"required,min=5,max=20"`
		Password   string `json:"password" validate:"required,min=6,max=32"`
		CaptchaKey string `json:"captcha_key" validate:"required"`
		Md5        bool   `json:"md5"`
		Point      struct {
			X int `json:"x"`
			Y int `json:"y"`
		}
	}

	req := Req{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
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

	// 登录
	userData, err := auth.LoginUser(r.Context(), req.Username, req.Password)
	if err != nil {
		_ = c.Fail(w, 103, "登录失败："+err.Error(), "")
		return
	}
	if userData.Status != 1 {
		_ = c.Fail(w, 104, "账号异常或被锁定", err)
		return
	}

	// 创建token
	jwtConf, err := config.GetJwtConfig()
	if err != nil {
		_ = c.Fail(w, 104, "token配置异常,请联系管理员", err)
	}
	// 生成token
	claims := auth.UserClaims{
		Sub:  userData.ID,
		Name: userData.Nickname,
	}
	token, err := auth.GenerateToken(claims)
	if err != nil {
		_ = c.Fail(w, 105, "生成token失败", err.Error())
		return
	}

	// 记录登录日志
	_ = dblog.New("login").
		WithModel(&model.Log{
			RequestUserID: userData.ID,
			RequestUser:   userData.Nickname,
		}).
		WithRequest(r).
		Write("登录后台", "")

	res := struct {
		Token    string      `json:"token"`
		UserInfo *model.User `json:"userInfo"`
		Expire   int64       `json:"expire"`
	}{
		Token:    token,
		UserInfo: userData,
		Expire:   jwtConf.Expire,
	}

	_ = c.Success(w, "", res)
}

func (c *Auth) GetMenu(w http.ResponseWriter, r *http.Request) {
	// 获取jwt并解析
	jwtValue := r.Context().Value("jwt")
	if _, ok := jwtValue.(auth.UserClaims); !ok {
		_ = c.Fail(w, 101, "获取登录信息失败", "jwt信息不正确")
		return
	}
	claims := jwtValue.(auth.UserClaims)

	// TODO: APPID应由前端指定
	var appId int32 = 1

	uq := query.User

	userInfo, err := uq.WithContext(r.Context()).
		Preload(
			query.User.UserRole,
			query.User.UserRole.Role,
		).
		Where(uq.ID.Eq(claims.Sub)).
		First()
	if err != nil {
		_ = c.Fail(w, 101, "获取用户信息失败", err.Error())
		return
	}

	menus, err := auth.GetMenu(userInfo, appId)
	if err != nil {
		_ = c.Fail(w, 102, "获取菜单失败", err.Error())
		return
	}

	permissions, err := auth.GetPermission(userInfo, appId)
	if err != nil {
		_ = c.Fail(w, 102, "获取权限失败", err.Error())
		return
	}

	res := struct {
		Menu        auth.MenuTreeArr `json:"menu"`
		Permissions []string         `json:"permissions"`
	}{
		Menu:        menus,
		Permissions: permissions,
	}

	_ = c.Success(w, "", res)
}

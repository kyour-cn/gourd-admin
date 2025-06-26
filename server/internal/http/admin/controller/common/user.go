package common

import (
	"app/internal/modules/auth"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

// User 用户控制器
type User struct {
	Base //继承基础控制器
}

// ResetPassword 重置密码
func (c *User) ResetPassword(w http.ResponseWriter, r *http.Request) {
	jwtValue := r.Context().Value("jwt")
	if _, ok := jwtValue.(auth.UserClaims); !ok {
		_ = c.Fail(w, 101, "获取登录信息失败", "jwt信息不正确")
		return
	}
	claims := jwtValue.(auth.UserClaims)

	type Req struct {
		UserPassword       string `json:"user_password" validate:"required,min=6,max=32"`
		NewPassword        string `json:"new_password" validate:"required,min=6,max=32"`
		ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=6,max=32"`
	}

	req := Req{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}
	if req.NewPassword != req.ConfirmNewPassword {
		_ = c.Fail(w, 102, "新密码和确认密码不一致", nil)
		return
	}

	hash := md5.Sum([]byte(req.UserPassword))
	checkPass := hex.EncodeToString(hash[:])

	// 查询旧密码是否正确
	check, _ := query.User.
		Where(
			query.User.ID.Eq(claims.Sub),
			query.User.Password.Eq(checkPass),
		).
		Select(query.User.ID).
		First()
	if check == nil {
		_ = c.Fail(w, 103, "旧密码不正确", nil)
		return
	}

	// 新密码加密
	newHash := md5.Sum([]byte(req.NewPassword))
	newPass := hex.EncodeToString(newHash[:])

	// 更新密码
	_, err := query.User.WithContext(r.Context()).
		Where(query.User.ID.Eq(claims.Sub)).
		Select(query.User.Password).
		Updates(model.User{
			Password: newPass,
		})
	if err != nil {
		_ = c.Fail(w, 1, "密码更新失败", err.Error())
		return
	}

	_ = c.Success(w, "密码重置成功", nil)
}

package controller

import (
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

// Info 用户信息查询、编辑
func (c *User) Info(w http.ResponseWriter, r *http.Request) {
	claims, err := c.GetJwt(r)
	if err != nil {
		_ = c.Fail(w, 101, err.Error(), "")
		return
	}

	// 查询用户信息
	qu := query.User
	user, err := qu.WithContext(r.Context()).
		Where(qu.ID.Eq(claims.Sub)).
		Select(qu.Username, qu.Nickname).
		First()
	if err != nil {
		_ = c.Fail(w, 1, "查询用户信息失败", err.Error())
		return
	}

	// 如果是POST请求
	if r.Method == http.MethodPost {
		req := struct {
			Nickname string `json:"nickname" validate:"required,min=2,max=20"`
		}{}
		if err := c.JsonReqUnmarshal(r, &req); err != nil {
			_ = c.Fail(w, 101, "请求参数异常", err.Error())
			return
		}
		// 更新用户信息
		_, err = qu.WithContext(r.Context()).
			Where(qu.ID.Eq(claims.Sub)).
			Select(qu.Nickname).
			Updates(&model.User{
				Nickname: req.Nickname,
			})
		if err != nil {
			_ = c.Fail(w, 1, "更新用户信息失败", err.Error())
			return
		}

		_ = c.Success(w, "编辑成功", "")
	} else {
		_ = c.Success(w, "获取用户信息成功", user)
	}
}

// ResetPassword 重置密码
func (c *User) ResetPassword(w http.ResponseWriter, r *http.Request) {
	claims, err := c.GetJwt(r)
	if err != nil {
		_ = c.Fail(w, 101, err.Error(), "")
		return
	}

	req := struct {
		UserPassword       string `json:"user_password" validate:"required,min=6,max=32"`
		NewPassword        string `json:"new_password" validate:"required,min=6,max=32"`
		ConfirmNewPassword string `json:"confirm_new_password" validate:"required,min=6,max=32"`
	}{}
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
	_, err = query.User.WithContext(r.Context()).
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

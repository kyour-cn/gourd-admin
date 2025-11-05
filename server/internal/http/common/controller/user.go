package controller

import (
	"app/internal/http/common/dto"
	"app/internal/http/common/services"
	"net/http"
)

// User 用户中心
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

	service := services.NewUserService(r.Context())

	// 查询用户信息
	user, err := service.GetInfo(claims)
	if err != nil {
		_ = c.Fail(w, 1, "查询用户信息失败", err.Error())
		return
	}

	// 如果是POST请求
	if r.Method == http.MethodPost {
		req := dto.UserUpdateNameReq{
			Claims: claims,
		}
		if err := c.JsonReqUnmarshal(r, &req); err != nil {
			_ = c.Fail(w, 101, "请求参数异常", err.Error())
			return
		}
		err := service.UpdateInfo(req)
		if err != nil {
			_ = c.Fail(w, 1, "更新失败", err.Error())
			return
		}

		_ = c.Success(w, "更新成功", nil)
	} else {
		_ = c.Success(w, "", user)
	}
}

// ResetPassword 重置密码
func (c *User) ResetPassword(w http.ResponseWriter, r *http.Request) {
	claims, err := c.GetJwt(r)
	if err != nil {
		_ = c.Fail(w, 101, err.Error(), "")
		return
	}

	req := dto.UserResetPasswordReq{
		Claims: claims,
	}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}
	service := services.NewUserService(r.Context())

	err = service.ResetPassword(req)
	if err != nil {
		_ = c.Fail(w, 500, err.Error(), nil)
		return
	}
	_ = c.Success(w, "密码重置成功", nil)
}

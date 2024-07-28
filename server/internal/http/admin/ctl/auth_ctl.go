package ctl

import (
	"encoding/json"
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/user"
	"net/http"
)

// AuthCtl 用户控制器
type AuthCtl struct {
	common.BaseCtl //继承基础控制器
}

type LoginResp struct {
	Token    string     `json:"token"`
	UserInfo model.User `json:"userInfo"`
}

// Login 登录
func (ctl *AuthCtl) Login(w http.ResponseWriter, r *http.Request) {

	// 创建一个结构体来存储解析后的JSON数据
	var requestData map[string]interface{}

	// 解析请求体中的JSON数据
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// TODO: 获取用户名和密码
	username := requestData["username"].(string)
	password := requestData["password"].(string)

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

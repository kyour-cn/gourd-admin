package system

import (
	"net/http"

	"app/internal/http/admin/dto"
	"app/internal/http/admin/services"
	"app/internal/http/common/controller"
)

// User 用户控制器
type User struct {
	controller.Base //继承基础控制器
}

func (c *User) List(w http.ResponseWriter, r *http.Request) {
	req := &dto.UserListReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	service := services.NewUserService(r.Context())
	res, err := service.GetList(req)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}

func (c *User) Export(w http.ResponseWriter, r *http.Request) {
	req := &dto.UserExportReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	service := services.NewUserService(r.Context())
	err := service.Export(req)
	if err != nil {
		_ = c.Fail(w, 500, "导出失败", err.Error())
		return
	}
	_ = c.Success(w, "", nil)
}

func (c *User) Add(w http.ResponseWriter, r *http.Request) {
	req := &dto.UserCreateReq{}
	if err := c.JsonReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	service := services.NewUserService(r.Context())
	if err := service.Create(req); err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}
	_ = c.Success(w, "success", nil)
}

func (c *User) Edit(w http.ResponseWriter, r *http.Request) {
	req := &dto.UserUpdateReq{}
	if err := c.JsonReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	service := services.NewUserService(r.Context())
	res, err := service.Update(req)
	if err != nil {
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}
	_ = c.Success(w, "success", res)
}

func (c *User) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []uint32 `json:"ids"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	service := services.NewUserService(r.Context())
	res, err := service.Delete(req.Ids)
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}
	_ = c.Success(w, "success", res)
}

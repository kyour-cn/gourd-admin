package system

import (
	"net/http"

	"app/internal/http/admin/dto"
	"app/internal/http/admin/services"
	"app/internal/http/common/controller"
)

// Role 用户控制器
type Role struct {
	controller.Base //继承基础控制器
}

func (c *Role) List(w http.ResponseWriter, r *http.Request) {
	req := &dto.RoleListReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	service := services.NewRoleService(r.Context())
	res, err := service.GetList(req)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}

func (c *Role) Add(w http.ResponseWriter, r *http.Request) {
	req := &dto.RoleCreateReq{}
	if err := c.JsonReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), nil)
		return
	}

	service := services.NewRoleService(r.Context())
	if err := service.Create(req); err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}
	_ = c.Success(w, "success", nil)
}

func (c *Role) Edit(w http.ResponseWriter, r *http.Request) {
	req := &dto.RoleUpdateReq{}
	if err := c.JsonReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), nil)
		return
	}
	// 从查询参数读取 type（permission 或 空）
	req.Type = r.URL.Query().Get("type")

	service := services.NewRoleService(r.Context())
	res, err := service.Update(req)
	if err != nil {
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}

func (c *Role) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []uint32 `json:"ids"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), nil)
		return
	}

	service := services.NewRoleService(r.Context())
	res, err := service.Delete(req.Ids)
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}
	_ = c.Success(w, "", res)
}

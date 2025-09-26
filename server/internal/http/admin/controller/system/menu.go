package system

import (
	"net/http"

	"app/internal/http/admin/dto"
	"app/internal/http/admin/services"
	"app/internal/http/common/controller"
	"app/internal/modules/admin/auth"
)

// Menu 用户控制器
type Menu struct {
	controller.Base //继承基础控制器
}

func (c *Menu) List(w http.ResponseWriter, r *http.Request) {
	req := &dto.MenuListReq{}
	if err := c.QueryReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常："+err.Error(), "")
		return
	}

	menus, err := auth.GetMenuFormApp(req.AppId)
	if err != nil {
		return
	}

	_ = c.Success(w, "", menus)
}

func (c *Menu) Add(w http.ResponseWriter, r *http.Request) {
	req := &dto.MenuCreateReq{}
	if err := c.JsonReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	service := services.NewMenuService(r.Context())
	menu, err := service.Create(req)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}
	_ = c.Success(w, "success", menu)
}

func (c *Menu) Edit(w http.ResponseWriter, r *http.Request) {
	req := &dto.MenuUpdateReq{}
	if err := c.JsonReqUnmarshal(r, req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	service := services.NewMenuService(r.Context())
	err := service.Update(req)
	if err != nil {
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}
	_ = c.Success(w, "success", nil)
}

// Delete 删除分类
func (c *Menu) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []int64 `json:"ids"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	service := services.NewMenuService(r.Context())
	err := service.Delete(req.Ids)
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}
	_ = c.Success(w, "success", nil)
}

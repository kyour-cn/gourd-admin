package ctl

import (
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories"
	"net/http"
	"strconv"
)

// AppCtl 用户控制器
type AppCtl struct {
	common.BaseCtl //继承基础控制器
}

func (c *AppCtl) List(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
	type Res struct {
		Rows  []*model.App `json:"rows"`
		Total int64        `json:"total"`
	}

	// 获取参数
	req := Req{
		Page:     1,
		PageSize: 10,
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page != 0 {
		req.Page = page
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize != 0 {
		req.PageSize = pageSize
	}

	ra := repositories.App{
		Ctx: r.Context(),
	}

	// 查询列表
	list, count, err := ra.Query().
		FindByPage((req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}

	res := Res{
		Rows:  list,
		Total: count,
	}

	_ = c.Success(w, "", res)
}

func (c *AppCtl) Add(w http.ResponseWriter, r *http.Request) {
	_ = c.Success(w, "success", nil)
}

func (c *AppCtl) Edit(w http.ResponseWriter, r *http.Request) {
	_ = c.Success(w, "success", nil)
}

func (c *AppCtl) Delete(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Ids []int32 `json:"ids"`
	}

	rm := repositories.App{
		Ctx: r.Context(),
	}

	req := Req{}
	err := c.JsonReqUnmarshal(r, &req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err = rm.Query().Where(query.App.ID.In(req.Ids...)).Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

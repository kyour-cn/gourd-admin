package ctl

import (
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"net/http"
)

// AppCtl 用户控制器
type AppCtl struct {
	common.BaseCtl //继承基础控制器
}

func (c *AppCtl) List(w http.ResponseWriter, r *http.Request) {
	type Res struct {
		Rows  []*model.App `json:"rows"`
		Total int64        `json:"total"`
	}

	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	// 查询列表
	list, count, err := query.App.WithContext(r.Context()).
		FindByPage((page-1)*pageSize, pageSize)
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
	req := &model.App{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	err := query.App.WithContext(r.Context()).Create(req)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *AppCtl) Edit(w http.ResponseWriter, r *http.Request) {
	req := &model.App{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	qm := query.App

	_, err := qm.WithContext(r.Context()).
		Where(query.App.ID.Eq(req.ID)).
		Select(
			qm.Name,
			qm.Key,
			qm.Remark,
			qm.Status,
			qm.Sort,
		).
		Updates(req)
	if err != nil {
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *AppCtl) Delete(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Ids []int32 `json:"ids"`
	}

	req := Req{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err := query.App.WithContext(r.Context()).
		Where(query.App.ID.In(req.Ids...)).
		Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

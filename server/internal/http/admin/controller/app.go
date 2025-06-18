package controller

import (
	"app/internal/orm/model"
	"app/internal/orm/query"
	"net/http"
)

// App 用户控制器
type App struct {
	Base //继承基础控制器
}

func (c *App) List(w http.ResponseWriter, r *http.Request) {
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

func (c *App) Add(w http.ResponseWriter, r *http.Request) {
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

func (c *App) Edit(w http.ResponseWriter, r *http.Request) {
	req := model.App{}
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

func (c *App) Delete(w http.ResponseWriter, r *http.Request) {
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

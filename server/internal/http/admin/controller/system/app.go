package system

import (
	"app/internal/http/common/controller"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"net/http"
)

// App 用户控制器
type App struct {
	controller.Base //继承基础控制器
}

func (c *App) List(w http.ResponseWriter, r *http.Request) {
	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	// 查询列表
	list, count, err := query.App.WithContext(r.Context()).
		FindByPage((page-1)*pageSize, pageSize)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}

	res := struct {
		Rows  []*model.App `json:"rows"`
		Total int64        `json:"total"`
	}{
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

	qa := query.App

	_, err := qa.WithContext(r.Context()).
		Where(query.App.ID.Eq(req.ID)).
		Select(qa.Name, qa.Key, qa.Remark, qa.Status, qa.Sort).
		Updates(req)
	if err != nil {
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *App) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []int32 `json:"ids"`
	}{}
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

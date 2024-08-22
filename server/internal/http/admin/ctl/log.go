package ctl

import (
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"net/http"
)

// LogCtl 用户控制器
type LogCtl struct {
	common.BaseCtl //继承基础控制器
}

func (c *LogCtl) LevelList(w http.ResponseWriter, r *http.Request) {

	type Res struct {
		Rows  []*model.LogLevel `json:"rows"`
		Total int64             `json:"total"`
	}

	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	// 查询列表
	list, count, err := query.LogLevel.WithContext(r.Context()).
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

func (c *LogCtl) List(w http.ResponseWriter, r *http.Request) {

	type Res struct {
		Rows  []*model.Log `json:"rows"`
		Total int64        `json:"total"`
	}

	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	// 查询列表
	list, count, err := query.Log.WithContext(r.Context()).
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

func (c *LogCtl) LogStat(w http.ResponseWriter, r *http.Request) {

	//TODO: 日志统计
	_ = c.Success(w, "", nil)
}

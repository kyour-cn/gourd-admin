package ctl

import (
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"net/http"
	"time"
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
	type Res struct {
		Days []string             `json:"days"`
		Rows []*model.LogStatView `json:"rows"`
	}
	params := r.URL.Query()

	// 获取参数
	startTimeStr, endTimeStr := params.Get("start_time"), params.Get("end_time")
	if startTimeStr == "" || endTimeStr == "" {
		_ = c.Fail(w, 1, "时间不能为空", nil)
		return
	}

	startTime, err1 := time.Parse(time.DateTime, startTimeStr)
	entTime, err2 := time.Parse(time.DateTime, endTimeStr)
	if err1 != nil || err2 != nil {
		_ = c.Fail(w, 101, "时间格式异常", nil)
		return
	}

	// 生成时间列表
	days := GenerateDays(startTime, entTime, "2006-01-02")

	// 查询日志数量
	list, _ := query.LogStatView.Where(query.LogStatView.Date.Between(
		startTime.Format(time.DateOnly),
		entTime.Format(time.DateOnly),
	)).Find()

	_ = c.Success(w, "", Res{
		Days: days,
		Rows: list,
	})
}

func GenerateDays(startDate, endDate time.Time, format string) []string {
	var days []string
	for current := startDate; !current.After(endDate); current = current.AddDate(0, 0, 1) {
		days = append(days, current.Format(format))
	}
	return days
}

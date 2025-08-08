package system

import (
	"app/internal/http/common/controller"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gen"
)

// Log 用户控制器
type Log struct {
	controller.Base //继承基础控制器
}

// TypeList 日志类型列表
func (c *Log) TypeList(w http.ResponseWriter, r *http.Request) {
	page, pageSize := c.PageParam(r, 1, 10)

	// 查询列表
	list, count, err := query.LogType.WithContext(r.Context()).
		FindByPage((page-1)*pageSize, pageSize)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}

	res := struct {
		Rows  []*model.LogType `json:"rows"`
		Total int64            `json:"total"`
	}{
		Rows:  list,
		Total: count,
	}

	_ = c.Success(w, "", res)
}

// List 日志列表
func (c *Log) List(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	var condition []gen.Condition

	// 时间筛选
	startTimeStr, endTimeStr := params.Get("start_time"), params.Get("end_time")
	if startTimeStr == "" || endTimeStr == "" {
		_ = c.Fail(w, 1, "时间范围不能为空", nil)
		return
	}
	startTime, err1 := time.Parse(time.DateTime, startTimeStr)
	entTime, err2 := time.Parse(time.DateTime, endTimeStr)
	if err1 != nil || err2 != nil {
		_ = c.Fail(w, 101, "时间格式异常", nil)
		return
	}
	condition = append(condition, query.Log.CreateTime.Between(startTime.Unix(), entTime.Unix()))

	// 类型筛选
	logType := params.Get("type_id")
	if logType != "" {
		typeId, _ := strconv.Atoi(logType)
		condition = append(condition, query.Log.TypeID.Eq(int64(typeId)))
	}

	// 查询列表
	list, count, err := query.Log.WithContext(r.Context()).
		Where(condition...).
		Order(query.Log.ID.Desc()).
		FindByPage((page-1)*pageSize, pageSize)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}

	res := struct {
		Rows  []*model.Log `json:"rows"`
		Total int64        `json:"total"`
	}{
		Rows:  list,
		Total: count,
	}

	_ = c.Success(w, "", res)
}

// LogStat 日志统计
func (c *Log) LogStat(w http.ResponseWriter, r *http.Request) {
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
	days := c.generateDays(startTime, entTime, time.DateOnly)

	type LogStat struct {
		Date     string `gorm:"column:date" json:"date"`
		Count_   int64  `gorm:"column:count;not null" json:"count"`
		TypeName string `gorm:"column:type_name;not null;comment:日志级别名称" json:"type_name"`
		TypeID   int32  `gorm:"column:type_id;not null;comment:日志级别 <10为系统日志" json:"type_id"`
	}
	var logRows []*LogStat

	db := query.App.UnderlyingDB()

	// 查询日志数量
	rows := db.Table(query.Log.TableName()).
		Select(
			"date_format(from_unixtime(`create_time`), '%Y-%m-%d') AS `date`",
			"count(0) AS `count`",
			"type_name",
			"type_id",
		).
		Where("`create_time` BETWEEN ? AND ?", startTime.Unix(), entTime.Unix()).
		Group("date, type_name, type_id").
		Find(&logRows)
	if rows.Error != nil {
		_ = c.Fail(w, 500, "查询日志统计失败", rows.Error.Error())
		return
	}

	_ = c.Success(w, "", struct {
		Days []string   `json:"days"`
		Rows []*LogStat `json:"rows"`
	}{
		Days: days,
		Rows: logRows,
	})
}

// generateDays 时间列表生成
func (c *Log) generateDays(startDate, endDate time.Time, format string) []string {
	var days []string
	for current := startDate; !current.After(endDate); current = current.AddDate(0, 0, 1) {
		days = append(days, current.Format(format))
	}
	return days
}

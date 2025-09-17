package services

import (
	"app/internal/http/admin/dto"
	"app/internal/orm/query"
	"context"
	"fmt"
	"time"

	"gorm.io/gen"
)

func NewLogTypeService(ctx context.Context) *LogTypeService {
	return &LogTypeService{
		ctx: ctx,
	}
}

type LogTypeService struct {
	ctx context.Context
}

func (s *LogTypeService) GetTypeList(req *dto.LogTypeListReq) (*dto.PageListReq, error) {
	q := query.LogType
	var conds []gen.Condition

	list, count, err := q.WithContext(s.ctx).
		Where(conds...).
		FindByPage((req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PageListReq{
		Rows:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *LogTypeService) GetList(req *dto.LogListReq) (*dto.PageListReq, error) {
	q := query.Log
	var conds []gen.Condition

	// 类型筛选
	if req.TypeId != 0 {
		conds = append(conds, query.Log.TypeID.Eq(req.TypeId))
	}

	if req.StartTime == "" || req.EndTime == "" {
		return nil, fmt.Errorf("时间范围不能为空")
	}
	startTime, err1 := time.Parse(time.DateTime, req.StartTime)
	entTime, err2 := time.Parse(time.DateTime, req.EndTime)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("时间格式异常")
	}
	conds = append(conds, query.Log.CreatedAt.Between(startTime, entTime))

	list, count, err := q.WithContext(s.ctx).
		Where(conds...).
		FindByPage((req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PageListReq{
		Rows:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// GetLogStat 日志统计
func (s *LogTypeService) GetLogStat(req *dto.LogStatReq) (any, error) {
	if req.StartTime == "" || req.EndTime == "" {
		return nil, fmt.Errorf("时间范围不能为空")
	}
	startTime, err1 := time.Parse(time.DateTime, req.StartTime)
	entTime, err2 := time.Parse(time.DateTime, req.EndTime)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("时间格式异常")
	}

	// 生成时间列表
	days := s.generateDays(startTime, entTime, time.DateOnly)

	type LogStat struct {
		Date     string `gorm:"column:date" json:"date"`
		Count    int64  `gorm:"column:count;not null" json:"count"`
		TypeName string `gorm:"column:type_name;not null;comment:日志级别名称" json:"type_name"`
		TypeID   int32  `gorm:"column:type_id;not null;comment:日志级别 <10为系统日志" json:"type_id"`
	}
	var logRows []*LogStat

	db := query.App.UnderlyingDB()

	// 查询日志数量
	rows := db.Table(query.Log.TableName()).
		Select(
			"date_format(`created_at`, '%Y-%m-%d') AS `date`",
			"count(*) AS `count`",
			"type_name",
			"type_id",
		).
		Where("`created_at` BETWEEN ? AND ?", req.StartTime, req.EndTime).
		Group("date, type_name, type_id").
		Find(&logRows)
	if rows.Error != nil {
		return nil, fmt.Errorf("查询日志统计失败")
	}
	return map[string]any{
		"days": days,
		"rows": logRows,
	}, nil
}

// generateDays 时间列表生成
func (s *LogTypeService) generateDays(startDate, endDate time.Time, format string) []string {
	var days []string
	for current := startDate; !current.After(endDate); current = current.AddDate(0, 0, 1) {
		days = append(days, current.Format(format))
	}
	return days
}

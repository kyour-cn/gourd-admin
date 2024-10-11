package service

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"net/http"
)

type Log struct {
	Ctx   context.Context
	Type  *model.LogType
	Model *model.Log
}

// NewLog create a new log instance
func NewLog() *Log {
	log := &Log{
		Model: &model.Log{},
		Ctx:   context.Background(),
	}
	return log
}

// WithRequest add request information to log
func (l *Log) WithRequest(r *http.Request) *Log {
	if r == nil {
		return l
	}

	// 取出请求参数
	l.Model.RequestIP = r.RemoteAddr
	l.Model.RequestSource = r.URL.String()

	ctx := r.Context()
	l.Ctx = ctx

	// 取出用户信息
	jwtValue := ctx.Value("jwt")
	if claims, ok := jwtValue.(UserClaims); ok {
		// 取出uid
		l.Model.RequestUserID = claims.Uid
		// 取出用户名称
		l.Model.RequestUser = claims.Uname
	}

	return l
}

// WithType add log type to log
func (l *Log) WithType(t model.LogType) *Log {
	l.Type = &t
	return l
}

// WithTypeLabel add log type by label to log
func (l *Log) WithTypeLabel(label string) *Log {

	//TODO: 优化查询, 避免每次都查询数据库

	lt := query.LogType
	// 根据label查找类型
	logType, err := lt.Where(lt.Label.Eq(label)).First()
	if err != nil {
		// 如果找不到类型，则默认debug类型
		logType, _ = lt.First()
	}
	l.Type = logType

	l.Model.TypeID = logType.ID
	l.Model.TypeName = logType.Name
	if logType.AppID > 0 {
		l.Model.AppID = logType.AppID
	}

	return l
}

// Write log to database
func (l *Log) Write(title string, value string) error {
	if title != "" {
		l.Model.Title = title
	}
	if value != "" {
		l.Model.Value = value
	}

	return query.Log.WithContext(l.Ctx).Create(l.Model)
}

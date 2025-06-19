package dblog

import (
	"app/internal/modles/auth"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"context"
	"fmt"
	"net/http"
)

type Log struct {
	Ctx   context.Context
	Type  *model.LogType
	Model *model.Log
}

// New create a new log instance
func New() *Log {
	log := &Log{
		Model: &model.Log{},
		Ctx:   context.Background(),
	}
	return log
}

// WithModel add model to log
func (l *Log) WithModel(model *model.Log) *Log {
	l.Model = model
	return l
}

// WithRequest add request information to log
func (l *Log) WithRequest(r *http.Request) *Log {
	if r == nil {
		return l
	}

	// 取出请求参数
	l.Model.RequestIP = r.RemoteAddr
	l.Model.RequestSource = fmt.Sprintf("[%s] %s", r.Method, r.URL.String())

	ctx := r.Context()
	l.Ctx = ctx

	// 取出用户信息
	jwtClaims := ctx.Value("jwt")
	if claims, ok := jwtClaims.(auth.UserClaims); ok {
		// 取出uid
		l.Model.RequestUserID = claims.Sub
		// 取出用户名称
		l.Model.RequestUser = claims.Name
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

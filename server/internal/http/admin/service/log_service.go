package service

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
)

type Log struct {
	Level *model.LogLevel
	ctx   *context.Context
}

// Write log to database
func (l *Log) Write(title string, value string) {
	data := &model.Log{
		Title: title,
		Value: value,
	}

	if l.ctx != nil {
		// read uid from context
		ctx := *l.ctx
		jwtValue := ctx.Value("jwt")
		if claims, ok := jwtValue.(UserClaims); ok {
			// 取出uid
			data.RequestUserID = claims.Uid

		}

		//if jwt != nil  {
		//    data.UID = ctx.Value("uid").(int64)
		//}
	}

	err := query.Log.Create(data)
	if err != nil {
		return
	}
}

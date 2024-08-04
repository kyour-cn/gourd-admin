package repositories

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/base"
)

type App struct {
	base.Repository
	Ctx context.Context
}

func (r *App) Query() query.IAppDo {
	if r.Tx != nil {
		return r.Tx.App.WithContext(r.Ctx)
	}
	return query.App.WithContext(r.Ctx)
}

func (r *App) Create(ud *model.App) error {
	userQ := r.Query()
	return userQ.Create(ud)
}

package repositories

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/base"
)

type Menu struct {
	base.Repository
	Ctx context.Context
}

func (r *Menu) Query() query.IMenuDo {
	if r.Tx != nil {
		return r.Tx.Menu.WithContext(r.Ctx)
	}
	return query.Menu.WithContext(r.Ctx)
}

func (r *Menu) Create(ud *model.Menu) error {
	userQ := r.Query()
	// TODO: add more fields
	return userQ.Create(ud)
}

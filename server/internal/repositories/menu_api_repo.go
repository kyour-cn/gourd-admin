package repositories

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/base"
)

type MenuApi struct {
	base.Repository
	Ctx context.Context
}

func (r *MenuApi) Query() query.IMenuAPIDo {
	if r.Tx != nil {
		return r.Tx.MenuAPI.WithContext(r.Ctx)
	}
	return query.MenuAPI.WithContext(r.Ctx)
}

func (r *MenuApi) Create(ud *model.MenuAPI) error {
	q := r.Query()
	return q.Create(ud)
}

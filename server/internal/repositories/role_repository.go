package repositories

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/base"
)

type Role struct {
	base.Repository
	Ctx context.Context
}

func (r *Role) Query() query.IRoleDo {
	if r.Tx != nil {
		return r.Tx.Role.WithContext(r.Ctx)
	}
	return query.Role.WithContext(r.Ctx)
}

func (r *Role) Create(ud *model.Role) error {
	userQ := r.Query()
	// TODO: add more fields
	return userQ.Create(ud)
}

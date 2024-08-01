package repositories

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/base"
)

type Rule struct {
	base.Repository
	Ctx context.Context
}

func (r *Rule) Query() query.IRuleDo {
	if r.Tx != nil {
		return r.Tx.Rule.WithContext(r.Ctx)
	}
	return query.Rule.WithContext(r.Ctx)
}

func (r *Rule) Create(ud *model.Rule) error {
	userQ := r.Query()
	// TODO: add more fields
	return userQ.Create(ud)
}

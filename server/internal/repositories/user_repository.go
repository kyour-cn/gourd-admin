package repositories

import (
	"context"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/base"
)

type User struct {
	base.Repository
	Ctx context.Context
}

func (r *User) Query() query.IUserDo {
	if r.Tx != nil {
		return r.Tx.User.WithContext(r.Ctx)
	}
	return query.User.WithContext(r.Ctx)
}

func (r *User) Create(ud *model.User) error {
	userQ := r.Query()
	// TODO: add more fields
	return userQ.Create(ud)
}

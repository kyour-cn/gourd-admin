package dto

import "app/internal/orm/model"

type UserListReq struct {
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1,lte=500"`
	Keyword  string `form:"keyword"`
}

type UserCreateReq struct {
	model.User
	Roles []int64 `json:"roles"`
}

type UserUpdateReq struct {
	model.User
	ID    int64   `json:"id" validate:"gt=0"`
	Roles []int64 `json:"roles"`
}

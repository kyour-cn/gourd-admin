package dto

import "app/internal/orm/model"

type UserListReq struct {
	Page     int    `form:"page" validate:"min:1" label:"分页"`
	PageSize int    `form:"page_size" validate:"min:1|max:500" label:"每页数量"`
	Keyword  string `form:"keyword"`
}

type UserExportReq struct {
	Keyword string `form:"keyword"`
}

type UserCreateReq struct {
	model.User
	Roles []int32 `json:"roles"`
}

type UserUpdateReq struct {
	model.User
	ID    int32   `json:"id" validate:"gt:0"`
	Roles []int32 `json:"roles"`
}

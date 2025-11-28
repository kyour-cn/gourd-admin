package dto

import "app/internal/orm/model"

type AppListReq struct {
	Page     int    `form:"page" validate:"min:1" label:"分页"`
	PageSize int    `form:"page_size" validate:"min:1|max:500" label:"每页数量"`
	Keyword  string `form:"keyword"`
}

type AppCreateReq struct {
	model.App
	ID int64 `json:"id" validate:"eq:0"`
}

type AppUpdateReq struct {
	model.App
	ID int64 `json:"id" validate:"gt:0"`
}

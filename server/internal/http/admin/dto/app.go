package dto

import "app/internal/orm/model"

type AppListReq struct {
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1,lte=500"`
	Keyword  string `form:"keyword"`
}

type AppCreateReq struct {
	model.App
	ID int64 `json:"id" validate:"eq=0"`
}

type AppUpdateReq struct {
	model.App
	ID int64 `json:"id" validate:"gt=0"`
}

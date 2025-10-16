package dto

import "app/internal/http/common/services"

type MenuListReq struct {
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1,lte=1000"`
	AppId    int64  `form:"app_id" validate:"gt=0"`
	Keyword  string `form:"keyword"`
}

type MenuCreateReq struct {
	Pid       int64             `json:"pid"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Component string            `json:"component"`
	Meta      services.MenuMate `json:"meta"`
	AppId     int64             `json:"app_id"`
}

type MenuUpdateReq struct {
	Id        int64             `json:"id"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Component string            `json:"component"`
	Sort      int32             `json:"sort"`
	Meta      services.MenuMate `json:"meta"`
	AppId     int64             `json:"appId"`
	Pid       int64             `json:"pid"`
	Status    int32             `json:"status"`
	ApiList   []struct {
		Path string `json:"path"`
		Tag  string `json:"tag"`
	} `json:"apiList"`
}

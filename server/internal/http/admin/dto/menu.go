package dto

import "app/internal/http/common/services"

type MenuListReq struct {
	Page     int    `form:"page" validate:"min:1" label:"分页"`
	PageSize int    `form:"page_size" validate:"min:1|max:1000"`
	AppId    int32  `form:"app_id" validate:"gt:0"`
	Keyword  string `form:"keyword"`
}

type MenuCreateReq struct {
	Pid       int32             `json:"pid"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Component string            `json:"component"`
	Meta      services.MenuMate `json:"meta"`
	AppId     int32             `json:"app_id"`
}

type MenuUpdateReq struct {
	Id        int32             `json:"id"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Component string            `json:"component"`
	Sort      int32             `json:"sort"`
	Meta      services.MenuMate `json:"meta"`
	AppId     int32             `json:"appId"`
	Pid       int32             `json:"pid"`
	Status    int32             `json:"status"`
	ApiList   []struct {
		Path string `json:"path"`
		Tag  string `json:"tag"`
	} `json:"apiList"`
}

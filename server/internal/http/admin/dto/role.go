package dto

type RoleListReq struct {
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1,lte=500"`
	AppId    int64  `form:"app_id"`
	Ids      string `form:"ids"` // 逗号分隔的ID列表
}

type RoleCreateReq struct {
	AppID       int64  `json:"app_id" validate:"gt=0"`
	Name        string `json:"name" validate:"required"`
	Sort        int32  `json:"sort"`
	Status      int32  `json:"status"` // 0 禁用 1 启用
	Remark      string `json:"remark"`
	IsAdmin     int32  `json:"is_admin"`
	Rules       string `json:"rules"` // 权限ID 列表
	RulesCheckd string `json:"rules_checkd"`
}

type RoleUpdateReq struct {
	Type        string `json:"type"`
	ID          int64  `json:"id" validate:"gt=0"`
	AppID       int64  `json:"app_id" validate:"gt=0"`
	Name        string `json:"name" validate:"required"`
	Sort        int32  `json:"sort"`
	Status      int32  `json:"status"` // 0 禁用 1 启用
	Remark      string `json:"remark"`
	IsAdmin     int32  `json:"is_admin"`
	Rules       string `json:"rules"` // 权限ID 列表
	RulesCheckd string `json:"rules_checkd"`
}

package dto

type RoleListReq struct {
	Page     int    `form:"page" validate:"min:1" label:"分页"`
	PageSize int    `form:"page_size" validate:"min:1|max:500" label:"每页数量"`
	AppId    uint32 `form:"app_id"`
	Ids      string `form:"ids"` // 逗号分隔的ID列表
	Name     string `form:"name"`
}

type RoleCreateReq struct {
	AppID   uint32 `json:"app_id" validate:"gt:0"`
	Name    string `json:"name" validate:"required" label:"角色名称"`
	Sort    uint32 `json:"sort"`
	Status  uint32 `json:"status"` // 0 禁用 1 启用
	Remark  string `json:"remark"`
	IsAdmin uint32 `json:"is_admin"`
}

type RoleUpdateReq struct {
	Type         string `json:"type"`
	ID           uint32 `json:"id" validate:"gt:0"`
	AppID        uint32 `json:"app_id"`
	Name         string `json:"name"`
	Sort         uint32 `json:"sort"`
	Status       uint32 `json:"status"` // 0 禁用 1 启用
	Remark       string `json:"remark"`
	IsAdmin      uint32 `json:"is_admin"`
	Rules        string `json:"rules"` // 权限ID 列表
	RulesChecked string `json:"rules_checked"`
}

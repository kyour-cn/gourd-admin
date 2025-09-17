package dto

type PageListReq struct {
	Rows     any   `json:"rows"`      // 页码
	Total    int64 `json:"total"`     // 每页数量
	Page     int   `json:"page"`      // 总数
	PageSize int   `json:"page_size"` // 每页数量
}

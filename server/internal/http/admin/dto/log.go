package dto

type LogTypeListReq struct {
	Page     int `form:"page" validate:"gte=1"`
	PageSize int `form:"page_size" validate:"gte=1,lte=500"`
}

type LogListReq struct {
	Page      int    `form:"page" validate:"gte=1"`
	PageSize  int    `form:"page_size" validate:"gte=1,lte=500"`
	TypeId    int64  `form:"type_id"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}

type LogStatReq struct {
	StartTime string `form:"start_time" validate:"required"`
	EndTime   string `form:"end_time" validate:"required"`
}

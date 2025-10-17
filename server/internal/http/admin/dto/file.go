package dto

import (
	"app/internal/modules/common/auth"
	"app/internal/orm/model"
	"mime/multipart"
)

type FileListReq struct {
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1,lte=500"`
	Keyword  string `form:"keyword"`
}

type FileUploadReq struct {
	Claims     *auth.UserClaims
	File       multipart.File
	FileHeader *multipart.FileHeader
	MenuId     int64
}

type FileUpdateReq struct {
	model.File
	ID int64 `json:"id" validate:"gt=0"`
}

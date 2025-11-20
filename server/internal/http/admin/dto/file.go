package dto

import (
	"mime/multipart"

	"app/internal/http/common/dto"
	"app/internal/orm/model"
)

type FileListReq struct {
	Page     int    `form:"page" validate:"gte=1"`
	PageSize int    `form:"page_size" validate:"gte=1,lte=500"`
	Keyword  string `form:"keyword"`
}

type FileUploadReq struct {
	Claims     *dto.UserClaims
	File       multipart.File
	FileHeader *multipart.FileHeader
	MenuId     int64
}

type FileUpdateReq struct {
	model.File
	ID int64 `json:"id" validate:"gt=0"`
}

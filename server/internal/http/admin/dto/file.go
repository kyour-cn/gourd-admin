package dto

import (
	"mime/multipart"

	"app/internal/http/common/dto"
	"app/internal/orm/model"
)

type FileMenuAddReq struct {
	Name string `json:"name" validate:"required"`
}

type FileListReq struct {
	Page     int    `form:"page" validate:"min:1" label:"分页"`
	PageSize int    `form:"page_size" validate:"min:1|max:500" label:"每页数量"`
	MenuID   int32  `form:"menu_id"`
	Keyword  string `form:"keyword"`
}

type FileUploadReq struct {
	Claims     *dto.UserClaims
	File       multipart.File
	FileHeader *multipart.FileHeader
	MenuID     int32
}

type FileUpdateReq struct {
	model.File
	ID int32 `json:"id" validate:"gt:0"`
}

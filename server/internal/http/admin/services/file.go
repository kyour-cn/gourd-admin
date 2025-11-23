package services

import (
	"app/internal/http/common/services"
	"context"
	"path/filepath"

	"gorm.io/gen"

	"app/internal/http/admin/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
)

func NewFileService(ctx context.Context) *FileService {
	return &FileService{
		ctx: ctx,
	}
}

type FileService struct {
	ctx context.Context
}

func (s *FileService) GetMenuList() ([]*model.FileMenu, error) {
	return query.FileMenu.WithContext(s.ctx).Find()
}

func (s *FileService) AddMenu(req *dto.FileMenuAddReq) error {
	menu := &model.FileMenu{
		Name: req.Name,
	}
	err := query.FileMenu.WithContext(s.ctx).Create(menu)
	if err != nil {
		return err
	}
	return nil
}

func (s *FileService) DeleteMenu(id int32) error {
	_, err := query.FileMenu.WithContext(s.ctx).
		Where(query.FileMenu.ID.Eq(id)).
		Delete()
	if err != nil {
		return err
	}
	return nil
}

func (s *FileService) GetList(req *dto.FileListReq) (*dto.PageListReq, error) {
	q := query.File
	var conds []gen.Condition

	// 关键词搜索
	if req.Keyword != "" {
		conds = append(conds, q.FileName.Like("%"+req.Keyword+"%"))
	}

	// 菜单ID搜索
	if req.MenuID > 0 {
		conds = append(conds, q.MenuID.Eq(req.MenuID))
	}

	list, count, err := q.WithContext(s.ctx).
		Where(conds...).
		Order(q.ID.Desc()).
		FindByPage((req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PageListReq{
		Rows:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *FileService) Upload(req *dto.FileUploadReq) (*model.File, error) {
	q := query.File

	fileName := req.FileHeader.Filename

	// 保存路径 按日期分目录，避免单目录文件过多
	service := services.NewFileService(s.ctx)
	output, err := service.CloudUpload(req.File, req.FileHeader, "files")
	if err != nil {
		return nil, err
	}

	file := &model.File{
		FileSize:   req.FileHeader.Size,
		FileType:   req.FileHeader.Header.Get("Content-Type"),
		FileName:   fileName,
		FileExt:    filepath.Ext(fileName),
		URL:        output.URL,
		FilePath:   output.Path,
		StorageID:  output.StorageID,
		StorageKey: output.Storage,
		HashMd5:    output.Hash,
		UserID:     req.Claims.Sub,
		MenuID:     req.MenuId,
	}

	err = q.WithContext(s.ctx).Create(file)
	return file, err
}

//func (s *FileService) Update(req *dto.FileUpdateReq) (gen.ResultInfo, error) {
//	q := query.File
//
//	return q.WithContext(s.ctx).
//		Where(q.ID.Eq(req.ID)).
//		//Select(q.Name, q.Key, q.Remark, q.Status, q.Sort).
//		Updates(&req.File)
//}

func (s *FileService) Delete(ids []int64) (gen.ResultInfo, error) {
	q := query.File

	return q.WithContext(s.ctx).
		Where(q.ID.In(ids...)).
		Delete()
}

package services

import (
	"context"

	"gorm.io/gen"

	"app/internal/http/admin/dto"
	"app/internal/orm/query"
)

func NewAppService(ctx context.Context) *AppService {
	return &AppService{
		ctx: ctx,
	}
}

type AppService struct {
	ctx context.Context
}

func (s *AppService) GetList(req *dto.AppListReq) (*dto.PageListReq, error) {
	q := query.App
	var conds []gen.Condition

	// 关键词搜索
	if req.Keyword != "" {
		conds = append(conds, q.Name.Like("%"+req.Keyword+"%"))
	}

	list, count, err := q.WithContext(s.ctx).
		Where(conds...).
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

func (s *AppService) Create(req *dto.AppCreateReq) error {
	q := query.App

	return q.WithContext(s.ctx).
		Create(&req.App)
}

func (s *AppService) Update(req *dto.AppUpdateReq) (gen.ResultInfo, error) {
	q := query.App

	return q.WithContext(s.ctx).
		Where(q.ID.Eq(req.ID)).
		Select(q.Name, q.Key, q.Remark, q.Status, q.Sort).
		Updates(&req.App)
}

func (s *AppService) Delete(ids []int64) (gen.ResultInfo, error) {
	q := query.App

	return q.WithContext(s.ctx).
		Where(q.ID.In(ids...)).
		Delete()
}

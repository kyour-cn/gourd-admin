package services

import (
	"context"
	"strconv"
	"strings"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"app/internal/http/admin/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
)

func NewRoleService(ctx context.Context) *RoleService {
	return &RoleService{ctx: ctx}
}

type RoleService struct {
	ctx context.Context
}

func (s *RoleService) GetList(req *dto.RoleListReq) (*dto.PageListRes, error) {
	q := query.Role
	var conds []gen.Condition

	// 关键词搜索
	if req.Name != "" {
		conds = append(conds, q.Name.Like("%"+req.Name+"%"))
	}

	if req.AppId > 0 {
		conds = append(conds, q.AppID.Eq(req.AppId))
	}

	if req.Ids != "" {
		var idSlice []uint32
		for _, v := range strings.Split(req.Ids, ",") {
			if v == "" {
				continue
			}
			n, _ := strconv.Atoi(v)
			idSlice = append(idSlice, uint32(n))
		}
		if len(idSlice) > 0 {
			conds = append(conds, q.ID.In(idSlice...))
		}
	}

	list, count, err := q.WithContext(s.ctx).
		Where(conds...).
		Order(q.AppID.Asc(), q.Sort.Asc()).
		FindByPage((req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PageListRes{
		Rows:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *RoleService) Create(req *dto.RoleCreateReq) error {
	q := query.Role

	role := &model.Role{
		AppID:   req.AppID,
		Name:    req.Name,
		Status:  req.Status,
		Remark:  req.Remark,
		IsAdmin: req.IsAdmin,
		Sort:    req.Sort,
	}

	return q.WithContext(s.ctx).
		Select(q.AppID, q.Name, q.Status, q.Remark, q.IsAdmin, q.Sort).
		Create(role)
}

func (s *RoleService) Update(req *dto.RoleUpdateReq) (gen.ResultInfo, error) {
	q := query.Role
	var fields []field.Expr

	// 适配不同更新场景
	if req.Type == "permission" {
		fields = append(fields, q.Rules, q.RulesChecked)
	} else {
		fields = append(fields, q.Name, q.Remark, q.Status, q.IsAdmin, q.Sort)
	}
	return q.WithContext(s.ctx).
		Where(q.ID.Eq(req.ID)).
		Select(fields...).
		Updates(&model.Role{
			IsAdmin:      req.IsAdmin,
			Name:         req.Name,
			Sort:         req.Sort,
			Remark:       req.Remark,
			Status:       req.Status,
			Rules:        req.Rules,
			RulesChecked: req.RulesChecked,
		})
}

func (s *RoleService) Delete(ids []uint32) (*gen.ResultInfo, error) {
	tx := query.Q.Begin()
	re, err := tx.Role.WithContext(s.ctx).
		Where(query.Role.ID.In(ids...)).
		Delete()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return &re, err
	}

	// 删除用户角色关联
	_, err = tx.UserRole.WithContext(s.ctx).
		Where(query.UserRole.RoleID.In(ids...)).
		Delete()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, err
		}
		return &re, err
	}
	return &re, tx.Commit()
}

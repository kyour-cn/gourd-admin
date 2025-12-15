package services

import (
	"context"
	"encoding/json"
	"fmt"

	"app/internal/http/admin/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
)

func NewMenuService(ctx context.Context) *MenuService {
	return &MenuService{
		ctx: ctx,
	}
}

type MenuService struct {
	ctx context.Context
}

func (s *MenuService) Create(req *dto.MenuCreateReq) (*model.Menu, error) {
	q := query.Menu

	mate, _ := json.Marshal(req.Meta)
	mateStr := string(mate)
	menu := &model.Menu{
		AppID:     req.AppId,
		Pid:       req.Pid,
		Name:      req.Name,
		Title:     req.Meta.Title,
		Type:      req.Meta.Type,
		Path:      req.Path,
		Component: req.Component,
		Sort:      0,
		Meta:      &mateStr,
	}

	err := q.WithContext(s.ctx).Create(menu)
	if err != nil {
		return nil, fmt.Errorf("创建失败: %w", err)
	}

	return menu, nil
}

func (s *MenuService) Update(req *dto.MenuUpdateReq) error {
	tx := query.Q.Begin()

	mate, _ := json.Marshal(req.Meta)
	_, err := tx.Menu.WithContext(s.ctx).
		Where(query.Menu.ID.Eq(req.Id)).
		Updates(map[string]any{
			"name":      req.Name,
			"title":     req.Meta.Title,
			"type":      req.Meta.Type,
			"path":      req.Path,
			"component": req.Component,
			"sort":      req.Sort,
			"meta":      mate,
			"pid":       req.Pid,
		})
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("更新失败: %w", err)
	}

	//更新菜单API权限
	_, err = tx.MenuAPI.WithContext(s.ctx).
		Where(query.MenuAPI.MenuID.Eq(req.Id)).
		Delete()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return fmt.Errorf("更新失败: %w", err)
	}

	for _, api := range req.ApiList {
		err = tx.MenuAPI.Create(&model.MenuAPI{
			MenuID: req.Id,
			AppID:  req.AppId,
			Path:   api.Path,
			Tag:    api.Tag,
		})
		if err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return fmt.Errorf("更新失败: %w", err)
		}
	}

	return tx.Commit()
}

// 递归获取所有子分类ID
func (s *MenuService) getAllSubMenuIDs(ids []uint32) ([]uint32, error) {
	q := query.Menu
	var allIDs = make([]uint32, 0)
	var stack = make([]uint32, len(ids))
	copy(stack, ids)
	for len(stack) > 0 {
		currentID := stack[0]
		stack = stack[1:]
		allIDs = append(allIDs, currentID)
		// 查找当前ID的所有子分类
		children, err := q.WithContext(s.ctx).Where(q.Pid.Eq(currentID)).Find()
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			stack = append(stack, child.ID)
		}
	}
	return allIDs, nil
}

func (s *MenuService) Delete(ids []uint32) error {
	q := query.Menu

	// 递归查找所有需要删除的分类ID（包括子分类）
	allIDs, err := s.getAllSubMenuIDs(ids)
	if err != nil {
		return err
	}
	_, err = q.WithContext(s.ctx).
		Where(q.ID.In(allIDs...)).
		Delete()
	if err != nil {
		return err
	}

	_, err = query.Menu.WithContext(s.ctx).
		Where(q.ID.In(ids...)).
		Delete()

	return err
}

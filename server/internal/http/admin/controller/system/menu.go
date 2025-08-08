package system

import (
	"app/internal/http/common/controller"
	"app/internal/modules/admin/auth"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

// Menu 用户控制器
type Menu struct {
	controller.Base //继承基础控制器
}

func (c *Menu) List(w http.ResponseWriter, r *http.Request) {
	req := struct {
		AppId int64 `json:"app_id"`
	}{
		AppId: 1,
	}
	appId, _ := strconv.Atoi(r.URL.Query().Get("app_id"))
	if appId != 0 {
		req.AppId = int64(appId)
	}

	menus, err := auth.GetMenuFormApp(req.AppId)
	if err != nil {
		return
	}

	_ = c.Success(w, "", menus)
}

func (c *Menu) Add(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Pid       int64         `json:"pid"`
		Name      string        `json:"name"`
		Path      string        `json:"path"`
		Component string        `json:"component"`
		Meta      auth.MenuMate `json:"meta"`
		AppId     int64         `json:"app_id"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	mate, _ := json.Marshal(req.Meta)
	data := &model.Menu{
		AppID:     req.AppId,
		Pid:       req.Pid,
		Name:      req.Name,
		Title:     req.Meta.Title,
		Type:      req.Meta.Type,
		Path:      req.Path,
		Component: req.Component,
		Status:    1,
		Sort:      0,
		Meta:      string(mate),
	}

	err := query.Menu.WithContext(r.Context()).Create(data)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	_ = c.Success(w, "success", data)
}

func (c *Menu) Edit(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Id        int64         `json:"id"`
		Name      string        `json:"name"`
		Path      string        `json:"path"`
		Component string        `json:"component"`
		Sort      int32         `json:"sort"`
		Meta      auth.MenuMate `json:"meta"`
		AppId     int64         `json:"appId"`
		Pid       int64         `json:"pid"`
		ApiList   []struct {
			Path string `json:"path"`
			Tag  string `json:"tag"`
		} `json:"apiList"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	mate, _ := json.Marshal(req.Meta)
	_, err := query.Menu.WithContext(r.Context()).
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
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}

	qma := query.MenuAPI

	//更新菜单API权限
	_, _ = qma.WithContext(r.Context()).
		Where(qma.MenuID.Eq(req.Id)).
		Delete()

	for _, api := range req.ApiList {
		_ = qma.Create(&model.MenuAPI{
			MenuID: req.Id,
			AppID:  req.AppId,
			Path:   api.Path,
			Tag:    api.Tag,
		})
	}

	_ = c.Success(w, "success", nil)
}

// 递归获取所有子分类ID
func (c *Menu) getAllSubMenuIDs(ctx context.Context, ids []int64) ([]int64, error) {
	q := query.Menu
	var allIDs = make([]int64, 0)
	var stack = make([]int64, len(ids))
	copy(stack, ids)
	for len(stack) > 0 {
		currentID := stack[0]
		stack = stack[1:]
		allIDs = append(allIDs, currentID)
		// 查找当前ID的所有子分类
		children, err := q.WithContext(ctx).Where(q.Pid.Eq(currentID)).Find()
		if err != nil {
			return nil, err
		}
		for _, child := range children {
			stack = append(stack, child.ID)
		}
	}
	return allIDs, nil
}

// Delete 删除分类
func (c *Menu) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []int64 `json:"ids"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}
	q := query.Menu

	// 递归查找所有需要删除的分类ID（包括子分类）
	allIDs, err := c.getAllSubMenuIDs(r.Context(), req.Ids)
	if err != nil {
		_ = c.Fail(w, 1, "查找子分类失败", err.Error())
		return
	}
	_, err = q.WithContext(r.Context()).
		Where(q.ID.In(allIDs...)).
		Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

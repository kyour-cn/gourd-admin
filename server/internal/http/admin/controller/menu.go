package controller

import (
	"app/internal/modles/auth"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"encoding/json"
	"net/http"
	"strconv"
)

// Menu 用户控制器
type Menu struct {
	Base //继承基础控制器
}

func (c *Menu) List(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		AppId int32 `json:"app_id"`
	}
	// 获取参数
	req := Req{
		AppId: 1,
	}
	appId, _ := strconv.Atoi(r.URL.Query().Get("app_id"))
	if appId != 0 {
		req.AppId = int32(appId)
	}

	menus, err := auth.GetMenuFormApp(req.AppId)
	if err != nil {
		return
	}

	_ = c.Success(w, "", menus)
}

func (c *Menu) Add(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		ParentId  int32         `json:"parentId"`
		Name      string        `json:"name"`
		Path      string        `json:"path"`
		Component string        `json:"component"`
		Meta      auth.MenuMate `json:"meta"`
		AppId     int32         `json:"app_id"`
	}
	// 获取参数
	req := Req{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	mate, _ := json.Marshal(req.Meta)
	data := &model.Menu{
		AppID:     req.AppId,
		Pid:       req.ParentId,
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
	type Req struct {
		Id        int32         `json:"id"`
		Name      string        `json:"name"`
		Path      string        `json:"path"`
		Component string        `json:"component"`
		Sort      int32         `json:"sort"`
		Meta      auth.MenuMate `json:"meta"`
		AppId     int32         `json:"appId"`
		ApiList   []struct {
			Path string `json:"path"`
			Tag  string `json:"tag"`
		} `json:"apiList"`
	}

	// 获取参数
	req := Req{}
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

func (c *Menu) Delete(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Ids []int32 `json:"ids"`
	}

	req := Req{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err := query.Menu.WithContext(r.Context()).
		Where(query.Menu.ID.In(req.Ids...)).
		Or(query.Menu.Pid.In(req.Ids...)). //删除子菜单
		Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

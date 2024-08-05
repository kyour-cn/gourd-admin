package ctl

import (
	"encoding/json"
	"gourd/internal/http/admin/common"
	"gourd/internal/http/admin/service"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories"
	"net/http"
	"strconv"
)

// MenuCtl 用户控制器
type MenuCtl struct {
	common.BaseCtl //继承基础控制器
}

func (c *MenuCtl) List(w http.ResponseWriter, r *http.Request) {

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

	menus, err := service.GetMenuFormApp(req.AppId)
	if err != nil {
		return
	}

	_ = c.Success(w, "", menus)
}

type MenuMate struct {
	Title            string `json:"title"`
	Icon             string `json:"icon"`
	Active           string `json:"active"`
	Color            string `json:"color"`
	Type             string `json:"type"`
	Fullpage         bool   `json:"fullpage"`
	Tag              string `json:"tag"`
	Hidden           bool   `json:"hidden"`
	HiddenBreadcrumb bool   `json:"hiddenBreadcrumb"`
}

func (c *MenuCtl) Add(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		ParentId  int32    `json:"parentId"`
		Name      string   `json:"name"`
		Path      string   `json:"path"`
		Component string   `json:"component"`
		Meta      MenuMate `json:"meta"`
		AppId     int32    `json:"app_id"`
	}
	// 获取参数
	req := Req{}
	err := c.JsonReqUnmarshal(r, &req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	rm := repositories.Menu{
		Ctx: r.Context(),
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

	err = rm.Create(data)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	_ = c.Success(w, "success", data)
}

func (c *MenuCtl) Edit(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Id        int32    `json:"id"`
		Name      string   `json:"name"`
		Path      string   `json:"path"`
		Component string   `json:"component"`
		Sort      int32    `json:"sort"`
		Meta      MenuMate `json:"meta"`
		AppId     int32    `json:"appId"`
		ApiList   []struct {
			Path string `json:"path"`
			Tag  string `json:"tag"`
		} `json:"apiList"`
	}

	// 获取参数
	req := Req{}
	err := c.JsonReqUnmarshal(r, &req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	rm := repositories.Menu{
		Ctx: r.Context(),
	}

	mate, _ := json.Marshal(req.Meta)
	_, err = rm.Query().
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

	//更新菜单API权限
	rmApi := repositories.MenuApi{
		Ctx: r.Context(),
	}

	_, _ = rmApi.Query().
		Where(query.MenuAPI.MenuID.Eq(req.Id)).
		Delete()

	for _, api := range req.ApiList {
		_ = rmApi.Create(&model.MenuAPI{
			MenuID: req.Id,
			AppID:  req.AppId,
			Path:   api.Path,
			Tag:    api.Tag,
		})
	}

	_ = c.Success(w, "success", nil)
}

func (c *MenuCtl) Delete(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Ids []int32 `json:"ids"`
	}

	rm := repositories.Menu{
		Ctx: r.Context(),
	}

	req := Req{}
	err := c.JsonReqUnmarshal(r, &req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err = rm.Query().Where(query.Menu.ID.In(req.Ids...)).Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

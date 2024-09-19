package ctl

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"net/http"
	"strconv"
)

// RoleCtl 用户控制器
type RoleCtl struct {
	common.BaseCtl //继承基础控制器
}

func (c *RoleCtl) List(w http.ResponseWriter, r *http.Request) {

	type Res struct {
		Rows  []*model.Role `json:"rows"`
		Total int64         `json:"total"`
	}

	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	// 获取参数
	appId, _ := strconv.Atoi(r.URL.Query().Get("app_id"))

	var conditions []gen.Condition
	if appId > 0 {
		conditions = append(conditions, query.Role.AppID.Eq(int32(appId)))
	}

	// 查询列表
	list, count, err := query.Role.WithContext(r.Context()).
		Preload(query.Role.App).
		Where(conditions...).
		Order(query.Role.AppID.Asc(), query.Role.Sort.Asc()).
		FindByPage((page-1)*pageSize, pageSize)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}

	res := Res{
		Rows:  list,
		Total: count,
	}

	_ = c.Success(w, "", res)
}

func (c *RoleCtl) Add(w http.ResponseWriter, r *http.Request) {
	req := &model.Role{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	err := query.Role.WithContext(r.Context()).Create(req)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *RoleCtl) Edit(w http.ResponseWriter, r *http.Request) {
	req := &model.Role{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	qm := query.Role

	var fields []field.Expr
	if r.URL.Query().Get("type") == "permission" {
		// 权限编辑
		fields = append(fields,
			qm.Rules,
			qm.RulesCheckd,
		)
	} else {
		fields = append(fields,
			qm.IsAdmin,
			qm.Name,
			qm.Remark,
			qm.Status,
			qm.Sort,
		)
	}

	_, err := query.Role.WithContext(r.Context()).
		Where(query.Role.ID.Eq(req.ID)).
		Select(fields...).
		Updates(req)
	if err != nil {
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *RoleCtl) Delete(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Ids []int32 `json:"ids"`
	}

	req := Req{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err := query.Role.WithContext(r.Context()).
		Where(query.Role.ID.In(req.Ids...)).
		Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

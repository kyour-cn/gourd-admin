package ctl

import (
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories"
	"net/http"
	"strconv"
)

// RoleCtl 用户控制器
type RoleCtl struct {
	common.BaseCtl //继承基础控制器
}

func (c *RoleCtl) List(w http.ResponseWriter, r *http.Request) {

	type Req struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
		AppId    int `json:"app_id"`
	}
	type Res struct {
		Rows  []*model.Role `json:"rows"`
		Total int64         `json:"total"`
	}

	// 获取参数
	req := Req{}
	req.Page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	req.PageSize, _ = strconv.Atoi(r.URL.Query().Get("page_size"))
	req.AppId, _ = strconv.Atoi(r.URL.Query().Get("app_id"))

	ra := repositories.Role{
		Ctx: r.Context(),
	}

	var conditions []gen.Condition
	if req.AppId > 0 {
		conditions = append(conditions, query.Role.AppID.Eq(int32(req.AppId)))
	}

	// 查询列表
	list, count, err := ra.Query().
		Where(conditions...).
		FindByPage((req.Page-1)*req.PageSize, req.PageSize)
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
	err := c.JsonReqUnmarshal(r, req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	rm := repositories.Role{
		Ctx: r.Context(),
	}

	err = rm.Create(req)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *RoleCtl) Edit(w http.ResponseWriter, r *http.Request) {
	req := &model.Role{}
	err := c.JsonReqUnmarshal(r, req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	rm := repositories.Role{
		Ctx: r.Context(),
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

	_, err = rm.Query().
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

	rm := repositories.Role{
		Ctx: r.Context(),
	}

	req := Req{}
	err := c.JsonReqUnmarshal(r, &req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err = rm.Query().Where(query.Role.ID.In(req.Ids...)).Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

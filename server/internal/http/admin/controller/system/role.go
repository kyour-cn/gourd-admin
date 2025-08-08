package system

import (
	"app/internal/http/common/controller"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

// Role 用户控制器
type Role struct {
	controller.Base //继承基础控制器
}

func (c *Role) List(w http.ResponseWriter, r *http.Request) {
	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	var conditions []gen.Condition

	// 筛选指定应用
	appId, _ := strconv.Atoi(r.URL.Query().Get("app_id"))
	if appId > 0 {
		conditions = append(conditions, query.Role.AppID.Eq(int64(appId)))
	}

	// 筛选指定id列表
	ids := r.URL.Query().Get("ids")
	if ids != "" {
		idSlice := make([]int64, 0)
		for _, v := range strings.Split(ids, ",") {
			num, _ := strconv.Atoi(v)
			idSlice = append(idSlice, int64(num))
		}
		conditions = append(conditions, query.Role.ID.In(idSlice...))
	}

	// 查询列表
	list, count, err := query.Role.WithContext(r.Context()).
		//Preload(query.Role.App).
		Where(conditions...).
		Order(query.Role.AppID.Asc(), query.Role.Sort.Asc()).
		FindByPage((page-1)*pageSize, pageSize)
	if err != nil {
		_ = c.Fail(w, 500, "获取列表失败", err.Error())
		return
	}

	res := struct {
		Rows  []*model.Role `json:"rows"`
		Total int64         `json:"total"`
	}{
		Rows:  list,
		Total: count,
	}

	_ = c.Success(w, "", res)
}

func (c *Role) Add(w http.ResponseWriter, r *http.Request) {
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

func (c *Role) Edit(w http.ResponseWriter, r *http.Request) {
	req := model.Role{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	qm := query.Role

	var fields []field.Expr
	if r.URL.Query().Get("type") == "permission" {
		// 权限编辑
		fields = append(fields,
			qm.Rules, qm.RulesCheckd,
		)
	} else {
		fields = append(fields,
			qm.IsAdmin, qm.Name, qm.Remark, qm.Status, qm.Sort,
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

func (c *Role) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []int64 `json:"ids"`
	}{}
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

package ctl

import (
	"crypto/md5"
	"encoding/hex"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gourd/internal/http/admin/common"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories"
	"net/http"
	"strconv"
)

// UserCtl 用户控制器
type UserCtl struct {
	common.BaseCtl //继承基础控制器
}

func (c *UserCtl) List(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
		Keyword  string `json:"keyword"`
	}
	type Res struct {
		Rows  []*model.User `json:"rows"`
		Total int64         `json:"total"`
	}

	// 获取参数
	req := Req{
		Page:     1,
		PageSize: 10,
	}
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page != 0 {
		req.Page = page
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize != 0 {
		req.PageSize = pageSize
	}

	var conditions []gen.Condition

	qu := query.User

	req.Keyword = r.URL.Query().Get("keyword")
	if req.Keyword != "" {
		conditions = append(conditions, qu.Where(
			qu.Where(
				qu.Nickname.Like("%"+req.Keyword+"%"),
			).Or(
				qu.Nickname.Like("%"+req.Keyword+"%"),
			).Or(
				qu.Mobile.Like("%"+req.Keyword+"%"),
			),
		))
	}

	ru := repositories.User{
		Ctx: r.Context(),
	}

	// 查询列表
	list, count, err := ru.Query().
		Preload(
			query.User.Role.Select(
				query.Role.ID,
				query.Role.Name,
			),
		).
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

func (c *UserCtl) Add(w http.ResponseWriter, r *http.Request) {
	req := &model.User{}
	err := c.JsonReqUnmarshal(r, req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	ru := repositories.User{
		Ctx: r.Context(),
	}

	err = ru.Create(req)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *UserCtl) Edit(w http.ResponseWriter, r *http.Request) {
	req := &model.User{}
	err := c.JsonReqUnmarshal(r, req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	ru := repositories.User{
		Ctx: r.Context(),
	}
	qu := query.User

	fields := []field.Expr{
		qu.Nickname,
		qu.Username,
		qu.Mobile,
		qu.Avatar,
		qu.Status,
		qu.RoleID,
	}

	// 如果密码不为空，则加密后更新密码
	if req.Password != "" {
		hash := md5.Sum([]byte(req.Password))
		req.Password = hex.EncodeToString(hash[:])
		fields = append(fields, qu.Password)
	}

	_, err = ru.Query().
		Where(query.User.ID.Eq(req.ID)).
		Select(fields...).
		Updates(req)
	if err != nil {
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *UserCtl) Delete(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Ids []int32 `json:"ids"`
	}

	rm := repositories.User{
		Ctx: r.Context(),
	}

	req := Req{}
	err := c.JsonReqUnmarshal(r, &req)
	if err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err = rm.Query().Where(query.User.ID.In(req.Ids...)).Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

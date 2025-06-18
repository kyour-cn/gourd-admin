package controller

import (
	"app/internal/orm/model"
	"app/internal/orm/query"
	"crypto/md5"
	"encoding/hex"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"net/http"
)

// User 用户控制器
type User struct {
	Base //继承基础控制器
}

func (c *User) List(w http.ResponseWriter, r *http.Request) {
	type Res struct {
		Rows  []*model.User `json:"rows"`
		Total int64         `json:"total"`
	}

	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	var conditions []gen.Condition

	qu := query.User

	keyword := r.URL.Query().Get("keyword")
	if keyword != "" {
		conditions = append(conditions, qu.Where(
			qu.Where(
				qu.Nickname.Like("%"+keyword+"%"),
			).Or(
				qu.Nickname.Like("%"+keyword+"%"),
			).Or(
				qu.Mobile.Like("%"+keyword+"%"),
			),
		))
	}

	// 查询列表
	list, count, err := query.User.WithContext(r.Context()).
		Preload(
			query.User.Role.Select(
				query.Role.ID,
				query.Role.Name,
			),
		).
		Where(conditions...).
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

func (c *User) Add(w http.ResponseWriter, r *http.Request) {
	req := &model.User{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	// 如果密码不为空，则加密后更新密码
	if req.Password != "" {
		hash := md5.Sum([]byte(req.Password))
		req.Password = hex.EncodeToString(hash[:])
	}

	err := query.User.WithContext(r.Context()).Create(req)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *User) Edit(w http.ResponseWriter, r *http.Request) {
	req := model.User{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
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

	_, err := query.User.WithContext(r.Context()).
		Where(query.User.ID.Eq(req.ID)).
		Select(fields...).
		Updates(req)
	if err != nil {
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *User) Delete(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Ids []int32 `json:"ids"`
	}

	req := Req{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	_, err := query.User.WithContext(r.Context()).
		Where(query.User.ID.In(req.Ids...)).
		Delete()
	if err != nil {
		_ = c.Fail(w, 1, "删除失败", err.Error())
		return
	}

	_ = c.Success(w, "success", nil)
}

package system

import (
	"app/internal/http/common/controller"
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
	controller.Base //继承基础控制器
}

func (c *User) List(w http.ResponseWriter, r *http.Request) {
	// 分页参数
	page, pageSize := c.PageParam(r, 1, 10)

	var conditions []gen.Condition
	qu := query.User

	keyword := r.URL.Query().Get("keyword")
	if keyword != "" {
		conditions = append(conditions, qu.Where(
			qu.Where(qu.Username.Like("%"+keyword+"%")).
				Or(qu.Nickname.Like("%"+keyword+"%")),
		))
	}

	// 查询列表
	list, count, err := query.User.WithContext(r.Context()).
		Preload(
			query.User.UserRole,
			query.User.UserRole.Role.Select(
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

	res := struct {
		Rows  []*model.User `json:"rows"`
		Total int64         `json:"total"`
	}{
		Rows:  list,
		Total: count,
	}

	_ = c.Success(w, "", res)
}

func (c *User) Add(w http.ResponseWriter, r *http.Request) {
	req := struct {
		model.User
		Roles []int32 `json:"roles"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	// 如果密码不为空，则加密后更新密码
	if req.Password != "" {
		hash := md5.Sum([]byte(req.Password))
		req.Password = hex.EncodeToString(hash[:])
	}

	err := query.User.WithContext(r.Context()).Create(&req.User)
	if err != nil {
		_ = c.Fail(w, 1, "创建失败", err.Error())
		return
	}

	// 新增用户角色
	err = c.UpdateRole(req.ID, req.Roles)
	if err != nil {
		_ = c.Fail(w, 1, "更新角色失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *User) Edit(w http.ResponseWriter, r *http.Request) {
	req := struct {
		model.User
		Roles []int32 `json:"roles"`
	}{}
	if err := c.JsonReqUnmarshal(r, &req); err != nil {
		_ = c.Fail(w, 101, "请求参数异常", err.Error())
		return
	}

	qu := query.User

	fields := []field.Expr{
		qu.Nickname, qu.Username, qu.Avatar, qu.Status,
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
		Updates(req.User)
	if err != nil {
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}

	// 更新用户角色
	err = c.UpdateRole(req.ID, req.Roles)
	if err != nil {
		_ = c.Fail(w, 1, "更新失败", err.Error())
		return
	}

	_ = c.Success(w, "success", req)
}

func (c *User) UpdateRole(userID int32, roleIDs []int32) error {
	// 查出原本的角色
	oldRoles, err := query.UserRole.
		Where(query.UserRole.UserID.Eq(userID)).
		Find()
	if err != nil {
		return err
	}

	// 构建新旧角色对比映射
	oldRoleMap := make(map[int32]bool)
	for _, role := range oldRoles {
		oldRoleMap[role.RoleID] = true
	}

	newRoleMap := make(map[int32]bool)
	for _, roleID := range roleIDs {
		newRoleMap[roleID] = true
	}

	// 删除失效角色
	for roleID := range oldRoleMap {
		if !newRoleMap[roleID] {
			_, err := query.UserRole.Where(query.UserRole.UserID.Eq(userID), query.UserRole.RoleID.Eq(roleID)).Delete()
			if err != nil {
				return err
			}
		}
	}

	// 新增新增角色
	for _, roleID := range roleIDs {
		if oldRoleMap[roleID] {
			continue
		}
		err := query.UserRole.Create(&model.UserRole{
			UserID: userID,
			RoleID: roleID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *User) Delete(w http.ResponseWriter, r *http.Request) {
	req := struct {
		Ids []int32 `json:"ids"`
	}{}
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

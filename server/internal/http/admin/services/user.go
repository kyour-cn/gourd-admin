package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"

	"github.com/go-gourd/gourd/event"
	"gorm.io/gen"
	"gorm.io/gen/field"

	"app/internal/http/admin/dto"
	comDto "app/internal/http/common/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
)

func NewUserService(ctx context.Context) *UserService {
	return &UserService{ctx: ctx}
}

type UserService struct {
	ctx context.Context
}

func (s *UserService) GetList(req *dto.UserListReq) (*dto.PageListRes, error) {
	q := query.User
	var conds []gen.Condition

	// 关键字搜索：用户名或昵称
	if req.Keyword != "" {
		conds = append(conds, q.Where(
			q.Where(q.Username.Like("%"+req.Keyword+"%")).
				Or(q.Nickname.Like("%"+req.Keyword+"%")),
		))
	}

	list, count, err := q.WithContext(s.ctx).
		Preload(
			query.User.UserRole,
			query.User.UserRole.Role.Select(query.Role.ID, query.Role.Name),
		).
		Where(conds...).
		FindByPage((req.Page-1)*req.PageSize, req.PageSize)
	if err != nil {
		return nil, err
	}

	return &dto.PageListRes{
		Rows:     list,
		Total:    count,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (s *UserService) Export(req *dto.UserExportReq) error {
	jwtValue := s.ctx.Value("jwt")
	if _, ok := jwtValue.(comDto.UserClaims); !ok {
		return errors.New("获取登录信息失败")
	}
	claims := jwtValue.(comDto.UserClaims)

	content, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = query.Task.WithContext(s.ctx).Create(&model.Task{
		Title:   "用户列表导出",
		Group_:  "user",
		Label:   "user",
		UserID:  claims.Sub,
		Type:    "export",
		Content: string(content),
	})

	// 触发任务运行事件
	event.Trigger("task.run", context.Background())

	return err
}

func (s *UserService) Create(req *dto.UserCreateReq) error {
	// 密码处理
	if req.Password != "" {
		hash := md5.Sum([]byte(req.Password))
		req.Password = hex.EncodeToString(hash[:])
	}

	// 查询用户是否存在
	_, err := query.User.WithContext(s.ctx).
		Where(query.User.Username.Eq(req.Username)).
		First()
	if err == nil {
		return errors.New("用户名已存在")
	}

	return query.Q.Transaction(func(tx *query.Query) error {
		err = tx.User.WithContext(s.ctx).Create(&req.User)
		if err != nil {
			return err
		}

		// 新增用户角色
		err = s.updateRole(tx, req.ID, req.Roles)
		if err != nil {
			return err
		}
		return nil
	})
}

func (s *UserService) Update(req *dto.UserUpdateReq) (any, error) {
	qu := query.User

	fields := []field.Expr{
		qu.Nickname, qu.Username, qu.Avatar, qu.Status,
	}
	// 密码更新（如提供）
	if req.Password != "" {
		hash := md5.Sum([]byte(req.Password))
		req.Password = hex.EncodeToString(hash[:])
		fields = append(fields, qu.Password)
	}

	var res gen.ResultInfo

	// 更新用户信息
	err := query.Q.Transaction(func(tx *query.Query) (err error) {
		res, err = tx.User.WithContext(s.ctx).
			Where(qu.ID.Eq(req.ID)).
			Select(fields...).
			Updates(req.User)
		if err != nil {
			return err
		}

		// 更新用户角色
		err = s.updateRole(tx, req.ID, req.Roles)
		if err != nil {
			return err
		}

		return nil
	})
	return res, err
}

func (s *UserService) Delete(ids []int32) (gen.ResultInfo, error) {
	return query.User.WithContext(s.ctx).
		Where(query.User.ID.In(ids...)).
		Delete()
}

// updateRole 差异更新用户角色
func (s *UserService) updateRole(tx *query.Query, userID int32, roleIDs []int32) error {
	q := query.UserRole

	// 原有角色
	oldRoles, err := tx.UserRole.WithContext(s.ctx).
		Where(q.UserID.Eq(userID)).
		Find()
	if err != nil {
		return err
	}

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
			_, err := tx.UserRole.WithContext(s.ctx).
				Where(q.UserID.Eq(userID), q.RoleID.Eq(roleID)).
				Delete()
			if err != nil {
				return err
			}
		}
	}

	// 新增缺失角色
	for _, roleID := range roleIDs {
		if oldRoleMap[roleID] {
			continue
		}
		err := tx.UserRole.WithContext(s.ctx).
			Create(&model.UserRole{UserID: userID, RoleID: roleID})
		if err != nil {
			return err
		}
	}
	return nil
}

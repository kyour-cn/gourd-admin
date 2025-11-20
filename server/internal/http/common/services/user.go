package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"app/internal/http/common/dto"
	"app/internal/orm/model"
	"app/internal/orm/query"
)

func NewUserService(ctx context.Context) *UserService {
	return &UserService{
		ctx: ctx,
	}
}

type UserService struct {
	ctx context.Context
}

func (s *UserService) GetInfo(claims *dto.UserClaims) (*model.User, error) {
	qu := query.User
	return qu.WithContext(s.ctx).
		Where(qu.ID.Eq(claims.Sub)).
		Select(qu.Username, qu.Nickname).
		First()
}

func (s *UserService) UpdateInfo(req dto.UserUpdateNameReq) error {
	qu := query.User
	_, err := qu.WithContext(s.ctx).
		Where(qu.ID.Eq(req.Claims.Sub)).
		Select(qu.Nickname).
		Updates(&model.User{
			Nickname: req.Nickname,
		})
	return err
}

func (s *UserService) ResetPassword(req dto.UserResetPasswordReq) error {
	if req.NewPassword != req.ConfirmNewPassword {
		return errors.New("新密码和确认密码不一致")
	}

	hash := md5.Sum([]byte(req.UserPassword))
	checkPass := hex.EncodeToString(hash[:])

	// 查询旧密码是否正确
	check, _ := query.User.
		Where(
			query.User.ID.Eq(req.Claims.Sub),
			query.User.Password.Eq(checkPass),
		).
		Select(query.User.ID).
		First()
	if check == nil {
		return errors.New("旧密码不正确")
	}

	// 新密码加密
	newHash := md5.Sum([]byte(req.NewPassword))
	newPass := hex.EncodeToString(newHash[:])

	// 更新密码
	_, err := query.User.WithContext(s.ctx).
		Where(query.User.ID.Eq(req.Claims.Sub)).
		Select(query.User.Password).
		Updates(model.User{
			Password: newPass,
		})
	if err != nil {
		return errors.New("密码更新失败")
	}

	return nil
}

func (s *UserService) GetTaskList(claims *dto.UserClaims) ([]*model.Task, error) {
	q := query.Task
	return q.WithContext(s.ctx).
		Where(q.UserID.Eq(claims.Sub)).
		Select(q.ID, q.Title, q.Content, q.Status, q.StatusName, q.Type, q.CreatedAt, q.UpdatedAt).
		Find()
}

package service

import (
	"github.com/golang-jwt/jwt/v5"
	"gourd/internal/config"
	"gourd/internal/orm/model"
	"gourd/internal/orm/query"
	"gourd/internal/repositories/user"
	"time"
)

func LoginUser(username string, password string) (*model.User, error) {

	// TODO: 登录频率限制锁

	ur := user.Repository{}
	uq := query.User

	// 查询用户
	userModel, err := ur.Query().
		Where(
			uq.Username.Eq(username),
			uq.Password.Eq(password),
		).
		Select(
			uq.ID,
			uq.Realname,
			uq.Username,
			uq.Mobile,
			uq.Avatar,
			uq.RegisterTime,
			uq.Status,
			uq.RoleID,
		).
		First()
	if err != nil {
		return nil, err
	}

	return userModel, nil
}

// GenerateToken 生成token
func GenerateToken(user *model.User) (string, error) {

	conf, err := config.GetJwtConfig()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iss": "gourd_admin",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Second * time.Duration(conf.Expire)).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(conf.Secret))
}

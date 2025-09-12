package auth

import (
	"app/internal/config"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/cache"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Sub  int64  `json:"sub"`
	Name string `json:"name"`
}

// LoginUser 登录用户
func LoginUser(_ context.Context, username string, password string) (*model.User, error) {
	// 登录频率限制锁 10秒
	key := "login_lock:" + username
	val, ok := cache.Get(key)
	if !ok {
		val = 0
	}
	failures, _ := val.(int)
	// 10秒登录失败次数超过3次，禁止登录
	if failures > 3 {
		return nil, errors.New("登录失败次数过多，请稍后再试")
	}

	uq := query.User

	// 查询用户
	userModel, err := uq.
		Preload(uq.UserRole, uq.UserRole.Role, uq.UserRole.Role.App).
		Where(
			uq.Username.Eq(username),
			uq.Password.Eq(password),
		).
		Select(uq.ID, uq.Nickname, uq.Username, uq.Avatar, uq.CreatedAt, uq.Status).
		First()
	if err != nil {
		// 登录失败次数+1
		cache.Set(key, failures+1, 10*time.Second)
		return nil, errors.New("用户名或密码错误")
	}
	cache.Delete(key)

	return userModel, nil
}

// GenerateToken 生成token
func GenerateToken(claims UserClaims) (string, error) {
	// 读取配置
	conf, err := config.GetJwtConfig()
	if err != nil {
		return "", err
	}

	// 设置签署时间和过期时间
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(conf.Expire)))

	// 使用HS256算法签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(conf.Secret))
}

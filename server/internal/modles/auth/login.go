package auth

import (
	"app/internal/config"
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/redisutil"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

// UserClaims TODO： 待完善调整
type UserClaims struct {
	jwt.RegisteredClaims
	Sub   int32  `json:"sub"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	AppId int32  `json:"app_id"` // 移除改为用户多角色
}

// LoginUser 登录用户
func LoginUser(ctx context.Context, username string, password string) (*model.User, error) {

	rdb, err := redisutil.GetRedis(ctx)
	if err != nil {
		return nil, errors.New("redis连接失败")
	}

	// 登录频率限制锁 10秒
	key := "login_lock:" + username
	val, _ := rdb.Get(ctx, key).Result()
	failures, _ := strconv.Atoi(val)
	// 10秒登录失败次数超过3次，禁止登录
	if failures > 3 {
		return nil, errors.New("登录失败次数过多，请稍后再试")
	}

	uq := query.User

	// 查询用户
	userModel, err := uq.
		Where(
			uq.Username.Eq(username),
			uq.Password.Eq(password),
		).
		Select(
			uq.ID,
			uq.Nickname,
			uq.Username,
			uq.Mobile,
			uq.Avatar,
			uq.CreateTime,
			uq.Status,
			uq.RoleID,
		).
		First()
	if err != nil {

		// 登录失败次数+1
		rdb.Incr(ctx, key)
		rdb.Expire(ctx, key, 10*time.Second)

		return nil, errors.New("用户名或密码错误")
	}

	rdb.Del(ctx, key)

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

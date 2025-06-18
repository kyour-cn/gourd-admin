package auth

import (
	"app/internal/orm/model"
	"app/internal/orm/query"
	"app/internal/util/redisutil"
	"context"
	"errors"
	"strconv"
	"time"
)

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

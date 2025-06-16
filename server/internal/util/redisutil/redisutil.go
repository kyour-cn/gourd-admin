package redisutil

import (
	"app/internal/config"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"strconv"
)

// RedisClient 是 Redis 客户端实例
var RedisClient *redis.Client

// GetRedis 初始化redis
func GetRedis(ctx context.Context) (client *redis.Client, err error) {
	if RedisClient != nil {
		return RedisClient, nil
	}

	// 自动初始化redis
	client, err = InitRedis(ctx)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func InitRedis(ctx context.Context) (client *redis.Client, err error) {

	dbConf, err := config.GetDBConfig("redis")
	if err != nil {
		return nil, errors.New("redis config is nil")
	}
	db, err := strconv.Atoi(dbConf.Database)
	if err != nil {
		return nil, errors.New("redis config database error")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     dbConf.Host + ":" + strconv.Itoa(dbConf.Port),
		Password: dbConf.Pass,
		DB:       db,
	})
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	RedisClient = client
	return client, nil
}

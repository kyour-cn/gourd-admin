package util

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gourd/internal/config"
	"gourd/internal/orm/query"
	"log/slog"
	"strconv"
	"time"
)

// 使用自定义logger接管gorm日志
type dbLogWriter struct{}

func (w dbLogWriter) Printf(format string, args ...any) {
	slog.Warn(fmt.Sprintf(format, args...))
}

// InitDatabase 初始化数据库连接
func InitDatabase() error {

	// 连接数据库
	dbConf, err := config.GetDBConfig("mysql")
	if err != nil {
		return errors.New("mysql config is nil")
	}

	// 连接数据库
	dsn := dbConf.GenerateDsn()
	gormConfig := &gorm.Config{
		// 替换默认日志
		Logger: logger.New(
			dbLogWriter{},
			logger.Config{
				SlowThreshold:             time.Duration(dbConf.SlowLogTime) * time.Millisecond, // 慢 SQL 阈值
				LogLevel:                  logger.Warn,                                          // 日志级别
				IgnoreRecordNotFoundError: true,                                                 // 忽略记录未找到错误
				Colorful:                  true,                                                 // 禁用彩色打印
			},
		),
	}
	mysqlDb, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return err
	}

	// 设置默认查询器
	query.SetDefault(mysqlDb)

	return nil
}

var redisConn *redis.Client

// GetRedis 初始化redis
func GetRedis() (client *redis.Client, err error) {
	if redisConn != nil {
		return redisConn, nil
	}

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
	redisConn = client
	return client, nil
}

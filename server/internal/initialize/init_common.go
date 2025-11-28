package initialize

import (
	"app/internal/util/cache"
	"context"
)

func InitCommon(ctx context.Context) {

	// 初始化日志
	err := InitLog()
	if err != nil {
		panic(err)
	}

	// 初始化缓存
	cache.InitDefaultCache(ctx)
}

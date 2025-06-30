package cache

import (
	"context"
	"sync"
	"time"
)

var (
	commonCache *Cache
	once        sync.Once
)

// InitCommonCache 初始化并定期执行 GC
func InitCommonCache(ctx context.Context) {
	once.Do(func() {
		commonCache = NewCache()

		go func() {
			ticker := time.NewTicker(1 * time.Minute) // 可配置的 GC 间隔
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					commonCache.GC()
				case <-ctx.Done():
					return
				}
			}
		}()
	})
}

func GetCommonCache() *Cache {
	InitCommonCache(context.Background()) // 保证已初始化
	return commonCache
}

func Set(key string, value any, duration time.Duration) {
	GetCommonCache().Set(key, value, duration)
}

func Get(key string) (any, bool) {
	return GetCommonCache().Get(key)
}

func Delete(key string) {
	GetCommonCache().Delete(key)
}

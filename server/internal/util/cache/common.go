package cache

import (
	"context"
	"sync"
	"time"
)

var (
	defaultCache *Cache
	once         sync.Once
)

// InitDefaultCache 初始化并定期执行 GC
func InitDefaultCache(ctx context.Context) {
	once.Do(func() {
		defaultCache = NewCache()
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					defaultCache.GC()
				case <-ctx.Done():
					return
				}
			}
		}()
	})
}

func GetDefaultCache() *Cache {
	InitDefaultCache(context.Background()) // 保证已初始化
	return defaultCache
}

func Set(key string, value any, duration time.Duration) {
	GetDefaultCache().Set(key, value, duration)
}

func Get(key string) (any, bool) {
	return GetDefaultCache().Get(key)
}

func Delete(key string) {
	GetDefaultCache().Delete(key)
}

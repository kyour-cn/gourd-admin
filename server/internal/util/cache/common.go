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
			ticker := time.NewTicker(30 * time.Second)
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

// Remember 尝试从默认缓存获取，未命中则使用 fn 生成并写入（泛型函数，方法无类型形参）
func Remember[T any](key string, duration time.Duration, fn func() (T, error)) (T, error) {
	c := GetDefaultCache()

	// 快路径：从缓存读取并断言类型
	if v, ok := c.Get(key); ok {
		if tv, ok2 := v.(T); ok2 {
			return tv, nil
		}
	}
	// 生成新值
	v, err := fn()
	if err != nil {
		var zero T
		return zero, err
	}
	// 写入缓存
	c.Set(key, v, duration)
	return v, nil
}

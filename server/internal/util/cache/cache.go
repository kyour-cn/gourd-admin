package cache

import (
	"sync"
	"time"
)

// Item 表示一个缓存项
type Item struct {
	Value      any
	Expiration int64 // Unix timestamp (毫秒)
}

// IsExpired 判断是否过期
func (item Item) IsExpired() bool {
	return item.Expiration > 0 && time.Now().UnixMilli() > item.Expiration
}

// Cache 是一个简易内存缓存
type Cache struct {
	data map[string]Item
	mu   sync.RWMutex
}

// NewCache 创建一个新缓存实例
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]Item),
	}
}

// Set 设置一个键值对，duration为有效期。如果为0则永不过期
func (c *Cache) Set(key string, value any, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiration int64
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixMilli()
	}

	c.data[key] = Item{
		Value:      value,
		Expiration: expiration,
	}
}

// Get 获取一个值，若存在且未过期则返回
func (c *Cache) Get(key string) (any, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.data[key]
	if !found || item.IsExpired() {
		return nil, false
	}
	return item.Value, true
}

// Delete 删除一个键
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// Clear 清空所有缓存
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data = make(map[string]Item)
}

// GC 清理过期的项
func (c *Cache) GC() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixMilli()
	for k, v := range c.data {
		if v.Expiration > 0 && now > v.Expiration {
			delete(c.data, k)
		}
	}
}

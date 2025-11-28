package cache

import (
	"sync"
	"time"
)

// cacheItem 缓存项
type cacheItem struct {
	value      string
	expiration time.Time
}

// LocalCache 本地内存缓存（LRU）
type LocalCache struct {
	items   map[string]*cacheItem
	mu      sync.RWMutex
	maxSize int
	ttl     time.Duration
}

// NewLocalCache 创建本地缓存
func NewLocalCache(maxSize int, ttl time.Duration) *LocalCache {
	lc := &LocalCache{
		items:   make(map[string]*cacheItem),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// 启动定时清理过期数据
	go lc.cleanupExpired()

	return lc
}

// Get 获取缓存
func (lc *LocalCache) Get(key string) (string, bool) {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	item, exists := lc.items[key]
	if !exists {
		return "", false
	}

	// 检查是否过期
	if time.Now().After(item.expiration) {
		return "", false
	}

	return item.value, true
}

// Set 设置缓存
func (lc *LocalCache) Set(key, value string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// 如果超过最大容量，随机删除一些旧数据
	if len(lc.items) >= lc.maxSize {
		// 简单实现：删除10%的数据
		count := 0
		target := lc.maxSize / 10
		for k := range lc.items {
			delete(lc.items, k)
			count++
			if count >= target {
				break
			}
		}
	}

	lc.items[key] = &cacheItem{
		value:      value,
		expiration: time.Now().Add(lc.ttl),
	}
}

// Delete 删除缓存
func (lc *LocalCache) Delete(key string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	delete(lc.items, key)
}

// Size 获取缓存大小
func (lc *LocalCache) Size() int {
	lc.mu.RLock()
	defer lc.mu.RUnlock()
	return len(lc.items)
}

// Clear 清空缓存
func (lc *LocalCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.items = make(map[string]*cacheItem)
}

// cleanupExpired 定期清理过期数据
func (lc *LocalCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		lc.mu.Lock()
		now := time.Now()
		for key, item := range lc.items {
			if now.After(item.expiration) {
				delete(lc.items, key)
			}
		}
		lc.mu.Unlock()
	}
}

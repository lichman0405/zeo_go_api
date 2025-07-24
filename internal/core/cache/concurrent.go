package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"zeo-api/internal/config"
)

type Cache struct {
	shards  []shard
	config  *config.CacheConfig
	baseDir string
	mu      sync.RWMutex
}

type shard struct {
	items map[string]cacheItem
	mu    sync.RWMutex
}

type cacheItem struct {
	Data     map[string][]byte
	Created  time.Time
	HitCount int64
}

func NewCache(cfg *config.CacheConfig, baseDir string) *Cache {
	shards := make([]shard, cfg.Shards)
	for i := range shards {
		shards[i] = shard{
			items: make(map[string]cacheItem),
		}
	}
	return &Cache{
		shards:  shards,
		config:  cfg,
		baseDir: baseDir,
	}
}

func (c *Cache) Get(key string) (map[string][]byte, bool) {
	shard := c.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	item, exists := shard.items[key]
	if !exists {
		return nil, false
	}

	if c.config.TTL > 0 && time.Since(item.Created) > c.config.TTL {
		delete(shard.items, key)
		return nil, false
	}

	item.HitCount++
	shard.items[key] = item
	return item.Data, true
}

func (c *Cache) Set(key string, data map[string][]byte) {
	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.items[key] = cacheItem{
		Data:    data,
		Created: time.Now(),
	}
}

func (c *Cache) Delete(key string) {
	shard := c.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

func (c *Cache) getShard(key string) *shard {
	hash := sha256.Sum256([]byte(key))
	index := int(hash[0]) % len(c.shards)
	return &c.shards[index]
}

func GenerateCacheKey(filePath string, args []string) string {
	h := sha256.New()
	h.Write([]byte(filePath))
	for _, arg := range args {
		h.Write([]byte(arg))
	}
	return hex.EncodeToString(h.Sum(nil))
}

func (c *Cache) ClearExpired() {
	for i := range c.shards {
		shard := &c.shards[i]
		shard.mu.Lock()
		for key, item := range shard.items {
			if c.config.TTL > 0 && time.Since(item.Created) > c.config.TTL {
				delete(shard.items, key)
			}
		}
		shard.mu.Unlock()
	}
}

func (c *Cache) Stats() (total, hits int64) {
	for i := range c.shards {
		shard := &c.shards[i]
		shard.mu.RLock()
		total += int64(len(shard.items))
		for _, item := range shard.items {
			hits += item.HitCount
		}
		shard.mu.RUnlock()
	}
	return
}

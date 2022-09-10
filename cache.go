package mycache

import (
	"sync"

	"github.com/myyppp/mycache/lru"
)

type cache struct {
	mu          sync.Mutex
	lru         *lru.Cache[string, ByteView]
	cacheOption lru.CacheOption
}

func (c *cache) set(k string, v ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 延迟初始化
	if c.lru == nil {
		c.lru = lru.New[string, ByteView](c.cacheOption)
	}
	c.lru.Set(k, v)
}

func (c *cache) get(k string) (v ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	v, ok = c.lru.Get(k)

	return
}

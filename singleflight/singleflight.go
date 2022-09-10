package singleflight

import "sync"

// call 正在进行中或已结束的请求
type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

// Group 管理不同 key 的请求（call）
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 针对相同的 key，无论 Do 被调用多少次，fn 只被调用一次
func (g *Group) Do(key string, fn func() (any, error)) (any, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}

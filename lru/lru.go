package lru

import "github.com/myyppp/mycache/list"

// Cache lru 缓存
type Cache[K comparable, V any] struct {
	ll      *list.List[entry[K, V]]
	items   map[K]*list.Element[entry[K, V]]
	options *options
}

// entry 双向链表的数据类型，保存 key 是为了在淘汰队首节点时便于从 map 中删除
type entry[K comparable, V any] struct {
	key   K
	value V
}

// New 初始化一个新的 lru 缓存
func New[K comparable, V any](cacheOption ...CacheOption) *Cache[K, V] {
	c := &Cache[K, V]{
		ll:      list.NewList[entry[K, V]](),
		items:   make(map[K]*list.Element[entry[K, V]]),
		options: defaultOptions(),
	}

	for _, o := range cacheOption {
		o.apply(c.options)
	}

	return c
}

// Len 返回缓存中 items 的数量
func (c *Cache[K, V]) Len() int {
	return c.ll.Len()
}

// Set
func (c *Cache[K, V]) Set(k K, v V) {
	if e, ok := c.items[k]; ok {
		e.Value.value = v
		c.ll.MoveToFront(e)
		return
	}

	e := c.ll.PushFront(entry[K, V]{k, v})
	c.items[k] = e
	if c.ll.Len() > c.options.capactiy {
		c.deleteElement(c.ll.Back())
	}
}

// Get 从缓存中获取 item，更新最近使用
func (c *Cache[K, V]) Get(k K) (v V, ok bool) {
	e, ok := c.items[k]
	if !ok {
		return
	}

	c.ll.MoveToFront(e)
	return e.Value.value, true
}

// Peek 从缓存中获取 item 但不更新最近使用
func (c *Cache[K, V]) Peek(k K) (v V, ok bool) {
	e, ok := c.items[k]
	if !ok {
		return
	}

	return e.Value.value, true
}

// Delete 从缓存中删除指定的 item
func (c *Cache[K, V]) Delete(k K) bool {
	e, ok := c.items[k]
	if !ok {
		return false
	}

	c.deleteElement(e)
	return true
}

// Flush 从缓存中删除所有的 itmes
func (c *Cache[K, V]) Flush() {
	c.ll.Init()
	c.items = make(map[K]*list.Element[entry[K, V]])
}

func (c *Cache[K, V]) deleteElement(e *list.Element[entry[K, V]]) {
	delete(c.items, e.Value.key)
	c.ll.Remove(e)
}

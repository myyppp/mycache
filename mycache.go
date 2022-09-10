package mycache

import (
	"fmt"
	"log"
	"sync"

	"github.com/myyppp/mycache/lru"
	"github.com/myyppp/mycache/singleflight"
)

// Group 缓存命名空间，每个 Group 拥有唯一的名称 name
type Group struct {
	name      string
	getter    Getter // 缓存未命中时获取数据源的回调
	mainCache cache  // 并发缓存
	peers     PeerPicker

	loader *singleflight.Group // 确保每个 key 只请求一次
}

// Getter 缓存不存在时，调用该回调函数获取数据
type Getter interface {
	Get(k string) ([]byte, error) // 回调函数
}

// GetterFunc 使用函数实现 Getter 接口
type GetterFunc func(k string) ([]byte, error)

// Get 实现 Getter 接口
func (f GetterFunc) Get(k string) ([]byte, error) {
	return f(k)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup 创建一个新的 Group 实例
func NewGroup(name string, cahceCapacity int, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheOption: lru.WithCapacity(cahceCapacity)},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup 返回 NewGroup 创建的 Group。使用只读锁，不涉及变量的写操作
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// RegisterPeers 注册一个 PeerPicker
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// Get 从 mainCache 中查找缓存，如果存在返回缓存值
// 不存在，调用 load方法，load 调用 getLocally，
// getLocally 调用用户回调函数 g.getter.Get() 获取数据源
// 并将数据添加到缓存中
func (g *Group) Get(k string) (ByteView, error) {
	if k == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(k); ok {
		log.Println("cache hit")
		return v, nil
	}

	return g.load(k)
}

func (g *Group) load(k string) (v ByteView, err error) {
	// Do 只执行一次
	view, err := g.loader.Do(k, func() (any, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PeerPick(k); ok {
				if v, err = g.getFromPeer(peer, k); err == nil {
					return v, err
				}
				log.Println("[myCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(k)
	})

	if err == nil {
		return view.(ByteView), nil
	}

	return
}

// getLocally 本地获取缓存
func (g *Group) getLocally(k string) (ByteView, error) {
	bytes, err := g.getter.Get(k)
	if err != nil {
		return ByteView{}, err
	}

	v := ByteView{b: cloneBytes(bytes)}
	g.populateCache(k, v)
	return v, nil
}

// getFromPeer 分布式场景下获取缓存
func (g *Group) getFromPeer(peer PeerGetter, k string) (ByteView, error) {
	bytes, err := peer.Get(g.name, k)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}

// populateCache 将源数据添加到缓存中
func (g *Group) populateCache(k string, v ByteView) {
	g.mainCache.set(k, v)
}

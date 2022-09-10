package mycache

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/myyppp/mycache/consistenthash"
)

const (
	defaultBasePath = "/_mycache/"
	defaultReplicas = 50
)

// HTTPPool 节点间 HTTP 通信的核心数据结构
type HTTPPool struct {
	self     string // base url，主机名/ip 和端口号
	basePath string // 节点间通信的前缀

	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter // 一个远程节点对应一个 httpGetter
}

// NewHTTPPool
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log 日志信息
func (p *HTTPPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP 处理所有 http 请求
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path: " + r.URL.Path)
	}
	p.Log("%s %s", r.Method, r.URL.Path)
	// 访问路径格式 /<basepath>/<groupname>/<key> required
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice())
}

// Set 实例化了一致性哈希算法，并添加传入的节点
// 并为每个节点创建一个 HTTP 客户端 httpGetter
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{baseURL: peer + p.basePath}
	}
}

// PeerPick HTTPPool 实现 PeerPicker 接口，返回节点对应的 HTTP 客户端
func (p *HTTPPool) PeerPick(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if peer := p.peers.Get(key); peer != "" && peer != p.self {
		p.Log("Pick peer %s", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

var _ PeerPicker = (*HTTPPool)(nil)

type httpGetter struct {
	baseURL string // 将要访问远程节点的地址
}

// httpGetter 实现 PeerGetter 接口
func (h *httpGetter) Get(group string, k string) ([]byte, error) {
	u := fmt.Sprintf(
		"%v%v%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(k),
	)
	res, err := http.Get(u) // 获取返回值
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned: %v", res.Status)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %v", err)
	}

	return bytes, nil
}

var _ PeerGetter = (*httpGetter)(nil)

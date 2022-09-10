package mycache

// PeerPicker 根据传入的 key 选择相应的节点
type PeerPicker interface {
	PeerPick(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 从 group 中获取缓存，对应于流程中 HTTP 的客户端
type PeerGetter interface {
	Get(group string, k string) ([]byte, error)
}

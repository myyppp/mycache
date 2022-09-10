package lru

const (
	DefaultCapacity = 10000 // 默认的 cache 容量
)

type options struct {
	capactiy int
}

// CacheOption 配置 lru 缓存
type CacheOption interface {
	apply(*options)
}

// CacheOptionFunc 包装一个函数来实现 CacheOption 接口
type CacheOptionFunc func(o *options)

func (f CacheOptionFunc) apply(o *options) {
	f(o)
}

// 功能选项模式
func WithCapacity(capactiy int) CacheOption {
	return CacheOptionFunc(func(o *options) {
		o.capactiy = capactiy
	})
}

func defaultOptions() *options {
	return &options{
		capactiy: DefaultCapacity,
	}
}

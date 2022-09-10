package mycache

// ByteView 用于表示缓存值
type ByteView struct {
	b []byte // []byte 为了能够支持任意的数据类型的存储
}

// Len
func (v ByteView) Len() int {
	return len(v.b)
}

// String
func (v ByteView) String() string {
	return string(v.b)
}

// ByteSlice 返回一个拷贝，防止缓存值被篡改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

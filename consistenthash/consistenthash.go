package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash
type Hash func(data []byte) uint32

// Map 包含所有的 hashed 的 keys
type Map struct {
	hash     Hash
	replicas int            // 虚拟节点倍数
	keys     []int          // 哈希环，sorted
	hashMap  map[int]string // 虚拟节点到真实节点的映射，k 虚拟节点的 hash，v 真实节点
}

// New 允许自定义虚拟节点倍数和 Hash 函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 添加节点
func (m *Map) Add(keys ...string) {
	for _, k := range keys {
		for i := 0; i < m.replicas; i++ {
			// 通过添加编号的方式区分不同的虚拟节点
			// m.hash() 计算虚拟节点的哈希值
			hash := int(m.hash([]byte(strconv.Itoa(i) + k)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = k
		}
	}
	sort.Ints(m.keys) // 对环上的哈希值进行排序
}

// Get 选择节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return hash <= m.keys[i]
	})

	// idx == len(m.keys) 时，选择 m.keys[0]
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

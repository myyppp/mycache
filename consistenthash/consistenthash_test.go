package consistenthash_test

import (
	"strconv"
	"testing"

	"github.com/myyppp/mycache/consistenthash"
)

func TestHashing(t *testing.T) {
	hash := consistenthash.New(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})

	// 虚拟节点为 2, 4, 6, 12, 14, 16, 22, 24, 26
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}

	// 8, 18, 28
	hash.Add("8")

	// 27 现在应该映射到节点 8
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s", k, v)
		}
	}
}

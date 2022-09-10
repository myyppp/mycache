package singleflight_test

import (
	"testing"

	"github.com/myyppp/mycache/singleflight"
)

func TestDo(t *testing.T) {
	var g singleflight.Group
	v, err := g.Do("key", func() (any, error) {
		return "bar", nil
	})

	if v != "bar" || err != nil {
		t.Errorf("Do v = %v, error = %v", v, err)
	}
}

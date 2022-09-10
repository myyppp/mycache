package mycache_test

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/myyppp/mycache"
)

func TestGetter(t *testing.T) {
	f := mycache.GetterFunc(func(k string) ([]byte, error) {
		return []byte(k), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	g := mycache.NewGroup("scores", 2<<10, mycache.GetterFunc(
		func(k string) ([]byte, error) {
			log.Println("SlowDB search key", k)
			if v, ok := db[k]; ok {
				if _, ok := db[k]; !ok {
					loadCounts[k] = 0
				}
				loadCounts[k]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", k)
		},
	))

	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", v)
		}
		if _, err := g.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // cache hit
	}

	if view, err := g.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}

func TestHTTPPool(t *testing.T) {
	mycache.NewGroup("scores", 2<<10, mycache.GetterFunc(
		func(k string) ([]byte, error) {
			log.Println("SlowDB search key", k)
			if v, ok := db[k]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", k)
		},
	))

	addr := "localhost:9999"
	peers := mycache.NewHTTPPool(addr)
	log.Println("geecache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

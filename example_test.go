package ttlcache_test

import (
	"fmt"
	"time"

	. "github.com/damoon/ttlcache"
)

type stringCache struct {
	*Cache
}

func (cache *stringCache) Get(k string) (int, bool) {
	cached, found := cache.GetUnsafe(k)
	if !found {
		return 0, false
	}
	return cached.(int), true
}

func (cache *stringCache) Set(k string, v int) {
	cache.SetUnsafe(k, v)
}

func ExampleCache() {
	c := &stringCache{NewCache(1 * time.Second)}
	c.Set("1", 2)

	fmt.Println(c.Get("1"))
	// Output: 2 true
}

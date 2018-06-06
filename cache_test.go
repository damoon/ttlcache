package ttlcache_test

import (
	"testing"
	"time"

	. "github.com/damoon/ttlcache"
)

func TestGet(t *testing.T) {
	cache := NewCache(time.Second)

	data, exists := cache.Get("hello")
	if exists {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set("hello", "world")
	data, exists = cache.Get("hello")
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data != "world" {
		t.Errorf("Expected cache to return `world` for `hello`")
	}
}

func TestInt(t *testing.T) {
	cache := NewCache(time.Second)

	data, exists := cache.Get(1)
	if exists {
		t.Errorf("Expected empty cache to return no data")
	}

	cache.Set(1, 2)
	data, exists = cache.Get(1)
	if !exists {
		t.Errorf("Expected cache to return data for `hello`")
	}
	if data != 2 {
		t.Errorf("Expected cache to return `world` for `hello`")
	}
}

func TestExpiration(t *testing.T) {
	cache := NewCache(time.Second)

	cache.Set("x", "1")
	cache.Set("y", "z")
	cache.Set("z", "3")

	item, exists := cache.Get("x")
	if !exists || item != "1" {
		t.Errorf("Expected `x` to not have expired after 500ms")
	}
	_, exists = cache.Get("y")
	if !exists {
		t.Errorf("Expected `y` to not have expired")
	}
	_, exists = cache.Get("z")
	if !exists {
		t.Errorf("Expected `z` to not have expired")
	}
	count := cache.Count()
	if count != 3 {
		t.Errorf("Expected cache to contain 3 item")
	}

	timer := cache.StartCleanupTimer(time.Second)
	defer timer.Stop()

	<-time.After(500 * time.Millisecond)
	item, exists = cache.Get("x")
	if !exists || item != "1" {
		t.Errorf("Expected `x` to not have expired after 500ms")
	}

	<-time.After(600 * time.Millisecond)

	item, exists = cache.Get("x")
	if !exists || item != "1" {
		t.Errorf("Expected `x` to not have expired")
	}
	_, exists = cache.Get("y")
	if exists {
		t.Errorf("Expected `y` to have expired")
	}
	_, exists = cache.Get("z")
	if exists {
		t.Errorf("Expected `z` to have expired")
	}
	count = cache.Count()
	if count != 1 {
		t.Errorf("Expected cache to contain 1 item")
	}

	<-time.After(2 * time.Second)
	_, exists = cache.Get("x")
	if exists {
		t.Errorf("Expected `x` to have expired")
	}
	count = cache.Count()
	if count != 0 {
		t.Errorf("Expected cache to be empty")
	}
}

package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/oarkflow/pkg/fastdb"
	"github.com/oarkflow/pkg/str"
)

// Cache is a thread safe cache.
type Cache[K comparable, V any] struct {
	cache Interface[K, *Item[K, V]]
	// mu is used to do lock in some method process.
	mu      sync.Mutex
	janitor *janitor
	store   *fastdb.DB
	persist bool
	bucket  string
}

// New creates a new thread safe Cache.
// The janitor will not be stopped which is created by this function. If you
// want to stop the janitor gracefully, You should use the `NewContext` function
// instead of this.
//
// There are several Cache replacement policies available with you specified any options.
func New[K comparable, V any](opts ...Option[K, V]) *Cache[K, V] {
	return NewContext(context.Background(), opts...)
}

// NewContext creates a new thread safe Cache with context.
// This function will be stopped by an internal janitor when the context is canceled.
//
// There are several Cache replacement policies available with you specified any options.
func NewContext[K comparable, V any](ctx context.Context, opts ...Option[K, V]) *Cache[K, V] {
	o := newOptions[K, V]()
	for _, optFunc := range opts {
		optFunc(o)
	}
	store, _ := fastdb.New(fastdb.Config{
		StorageType: fastdb.DiskStorage,
		Path:        "./data",
	})
	cache := &Cache[K, V]{
		cache:   o.cache,
		janitor: newJanitor(ctx, o.janitorInterval),
		store:   store,
		persist: o.persist,
		bucket:  o.bucket,
	}
	cache.janitor.run(cache.DeleteExpired)
	return cache
}

// Get looks up a key's value from the cache.
func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.cache.Get(key)

	if !ok {
		i, o := c.store.Get(c.bucket, fmt.Sprintf("%v", key))
		ok = o
		if !ok {
			return
		}
		json.Unmarshal(i, &value)
		c.Set(key, value)
		return
	}
	// Returns nil if the item has been expired.
	// Do not delete here and leave it to an external process such as Janitor.
	if item.Expired() {
		return value, false
	}

	return item.Value, true
}

// DeleteExpired all expired items from the cache.
func (c *Cache[K, V]) DeleteExpired() {
	c.mu.Lock()
	keys := c.cache.Keys()
	c.mu.Unlock()
	for _, key := range keys {
		c.mu.Lock()
		// if is expired, delete it and return nil instead
		item, ok := c.cache.Get(key)
		if ok && item.Expired() {
			c.cache.Delete(key)
		}
		c.mu.Unlock()
	}
}

// Set sets a value to the cache with key. replacing any existing value.
func (c *Cache[K, V]) Set(key K, val V, opts ...ItemOption) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := newItem(key, val, opts...)
	c.cache.Set(key, item)
	if c.persist {
		c.Persist(key, item)
	}
}

// Keys returns the keys of the cache. the order is relied on algorithms.
func (c *Cache[K, V]) Keys() []K {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Keys()
}

// Delete deletes the item with provided key from the cache.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Delete(key)
}

// Contains reports whether key is within cache.
func (c *Cache[K, V]) Contains(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.cache.Get(key)
	return ok
}

func (c *Cache[K, V]) Persist(k K, val *Item[K, V]) error {
	key := fmt.Sprintf("%v", k)
	_, err := c.store.Del(c.bucket, key)
	if err != nil {
		return err
	}
	value := fmt.Sprintf("%v", val.Value)
	return c.store.Set(c.bucket, key, str.ToByte(value))
}

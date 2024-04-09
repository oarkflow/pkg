package cache

import (
	"time"

	"github.com/oarkflow/pkg/cache/policy/clock"
	"github.com/oarkflow/pkg/cache/policy/fifo"
	"github.com/oarkflow/pkg/cache/policy/lfu"
	"github.com/oarkflow/pkg/cache/policy/lru"
	"github.com/oarkflow/pkg/cache/policy/mru"
	"github.com/oarkflow/pkg/cache/policy/simple"
)

// Interface is a common-cache interface.
type Interface[K comparable, V any] interface {
	// Get looks up a key's value from the cache.
	Get(key K) (value V, ok bool)
	// Set sets a value to the cache with key. replacing any existing value.
	Set(key K, val V)
	// Keys returns the keys of the cache. The order is relied on algorithms.
	Keys() []K
	// Delete deletes the item with provided key from the cache.
	Delete(key K)
}

var (
	_ = []Interface[struct{}, any]{
		(*simple.Cache[struct{}, any])(nil),
		(*lru.Cache[struct{}, any])(nil),
		(*lfu.Cache[struct{}, any])(nil),
		(*fifo.Cache[struct{}, any])(nil),
		(*mru.Cache[struct{}, any])(nil),
		(*clock.Cache[struct{}, any])(nil),
	}
)

// Item is an item
type Item[K comparable, V any] struct {
	Key                   K
	Value                 V
	Expiration            time.Time
	InitialReferenceCount int
}

// Expired returns true if the item has expired.
func (item *Item[K, V]) Expired() bool {
	if item.Expiration.IsZero() {
		return false
	}
	return nowFunc().After(item.Expiration)
}

// GetReferenceCount returns reference count to be used when setting
// the cache item for the first time.
func (item *Item[K, V]) GetReferenceCount() int {
	return item.InitialReferenceCount
}

var nowFunc = time.Now

// ItemOption is an option for cache item.
type ItemOption func(*itemOptions)

type itemOptions struct {
	expiration     time.Time // default none
	referenceCount int
}

// WithExpiration is an option to set expiration time for any items.
// If the expiration is zero or negative value, it treats as w/o expiration.
func WithExpiration(exp time.Duration) ItemOption {
	return func(o *itemOptions) {
		o.expiration = nowFunc().Add(exp)
	}
}

// WithReferenceCount is an option to set reference count for any items.
// This option is only applicable to cache policies that have a reference count (e.g., Clock, LFU).
// referenceCount specifies the reference count value to set for the cache item.
//
// the default is 1.
func WithReferenceCount(referenceCount int) ItemOption {
	return func(o *itemOptions) {
		o.referenceCount = referenceCount
	}
}

// newItem creates a new item with specified any options.
func newItem[K comparable, V any](key K, val V, opts ...ItemOption) *Item[K, V] {
	o := new(itemOptions)
	for _, optFunc := range opts {
		optFunc(o)
	}
	return &Item[K, V]{
		Key:                   key,
		Value:                 val,
		Expiration:            o.expiration,
		InitialReferenceCount: o.referenceCount,
	}
}

// Option is an option for cache.
type Option[K comparable, V any] func(*options[K, V])

type options[K comparable, V any] struct {
	cache           Interface[K, *Item[K, V]]
	janitorInterval time.Duration
	persist         bool
	bucket          string
}

func newOptions[K comparable, V any]() *options[K, V] {
	return &options[K, V]{
		cache:           simple.NewCache[K, *Item[K, V]](),
		janitorInterval: time.Minute,
	}
}

// AsLRU is an option to make a new Cache as LRU algorithm.
func AsLRU[K comparable, V any](opts ...lru.Option) Option[K, V] {
	return func(o *options[K, V]) {
		o.cache = lru.NewCache[K, *Item[K, V]](opts...)
	}
}

// AsLFU is an option to make a new Cache as LFU algorithm.
func AsLFU[K comparable, V any](opts ...lfu.Option) Option[K, V] {
	return func(o *options[K, V]) {
		o.cache = lfu.NewCache[K, *Item[K, V]](opts...)
	}
}

// AsFIFO is an option to make a new Cache as FIFO algorithm.
func AsFIFO[K comparable, V any](opts ...fifo.Option) Option[K, V] {
	return func(o *options[K, V]) {
		o.cache = fifo.NewCache[K, *Item[K, V]](opts...)
	}
}

// AsMRU is an option to make a new Cache as MRU algorithm.
func AsMRU[K comparable, V any](opts ...mru.Option) Option[K, V] {
	return func(o *options[K, V]) {
		o.cache = mru.NewCache[K, *Item[K, V]](opts...)
	}
}

// AsClock is an option to make a new Cache as clock algorithm.
func AsClock[K comparable, V any](opts ...clock.Option) Option[K, V] {
	return func(o *options[K, V]) {
		o.cache = clock.NewCache[K, *Item[K, V]](opts...)
	}
}

// WithJanitorInterval is an option to specify how often cache should delete expired items.
//
// Default is 1 minute.
func WithJanitorInterval[K comparable, V any](d time.Duration) Option[K, V] {
	return func(o *options[K, V]) {
		o.janitorInterval = d
	}
}

func WithPersist[K comparable, V any](persist bool, bucket string) Option[K, V] {
	return func(o *options[K, V]) {
		o.persist = persist
		o.bucket = bucket
	}
}

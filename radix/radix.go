package radix

import (
	"strings"

	"github.com/oarkflow/pkg/maps"
)

// Node represents a node in the Radix Trie
type Node struct {
	children maps.IMap[string, *Node]
	bitmask  int
}

// NewNode creates a new Radix Trie node
func NewNode() *Node {
	return &Node{
		children: maps.New[string, *Node](),
	}
}

// Trie represents a Radix Trie data structure
type Trie struct {
	root  *Node
	cache maps.IMap[string, int]
}

// New initializes a new Radix Trie
func New() *Trie {
	return &Trie{
		root:  NewNode(),
		cache: maps.New[string, int](),
	}
}

// InsertPermission inserts a permission into the Radix Trie
func (rt *Trie) InsertPermission(keys []string, bitmask int) {
	node := rt.root
	for _, key := range keys {
		if n, ok := node.children.Get(key); !ok || n == nil {
			node.children.Set(key, NewNode())
		}
		if n, ok := node.children.Get(key); ok {
			node = n
		}
	}
	node.bitmask = bitmask
}

// HasPermission checks if a role has permission for a specific URL and method
func (rt *Trie) HasPermission(keys []string, method int) bool {
	cacheKey := strings.Join(keys, "|")
	if bitmask, ok := rt.cache.Get(cacheKey); ok {
		return bitmask&method != 0
	}
	node := rt.root
	for _, key := range keys {
		if n, ok := node.children.Get(key); !ok || n == nil {
			rt.cache.Set(cacheKey, 0)
			return false
		}
		if n, ok := node.children.Get(key); ok {
			node = n
		}
	}
	rt.cache.Set(cacheKey, node.bitmask)
	return node.bitmask&method != 0
}

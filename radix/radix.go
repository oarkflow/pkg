package radix

import (
	"github.com/oarkflow/pkg/maps"
)

// Node represents a node in the Radix Trie
type Node struct {
	children maps.IMap[string, *Node]
	bitmask  int
}

// AddChild adds a node to the Radix Trie
func (n *Node) AddChild(key string, child *Node) {
	n.children.Set(key, child)
}

// RemoveChild removes a node from the Radix Trie
func (n *Node) RemoveChild(key string) {
	n.children.Del(key)
}

func (n *Node) Children(key string) (*Node, bool) {
	return n.children.Get(key)
}

// NewNode creates a new Radix Trie node
func NewNode() *Node {
	return &Node{
		children: maps.New[string, *Node](),
	}
}

// Trie represents a Radix Trie data structure
type Trie struct {
	root *Node
	// cache maps.IMap[string, int]
}

// New initializes a new Radix Trie
func New() *Trie {
	return &Trie{
		root: NewNode(),
	}
}

// InsertPermission inserts a permission into the Radix Trie
func (rt *Trie) InsertPermission(keys []string, bitmask int) {
	node := rt.root
	for _, key := range keys {
		n, _ := node.Children(key)
		if n == nil {
			n = NewNode()
			node.AddChild(key, n)
		}
		node = n
	}
	node.bitmask = bitmask
}

// HasPermission checks if a role has permission for a specific URL and method
func (rt *Trie) HasPermission(keys []string, method int) bool {
	node := rt.root
	for _, key := range keys {
		n, _ := node.Children(key)
		if n == nil {
			return false
		}
		node = n
	}
	return node.bitmask&method != 0
}

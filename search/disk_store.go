package search

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shamaton/msgpack/v2"
)

// KVStore defines the interface for our key-value store
type KVStore[K comparable, V any] interface {
	Set(key K, value V) error
	Get(key K) (V, bool)
	Del(key K) error
}

// DiskStore implements KVStore using files on disk
type DiskStore[K comparable, V any] struct {
	basePath string
	compress bool
}

// NewDiskStore creates a new DiskStore instance
func NewDiskStore[K comparable, V any](basePath string, compress bool) (KVStore[K, V], error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}
	return &DiskStore[K, V]{basePath: basePath, compress: compress}, nil
}

// Set stores a key-value pair on disk
func (s *DiskStore[K, V]) Set(key K, value V) error {
	fileName := filepath.Join(s.basePath, fmt.Sprintf("%v.json", key))
	jsonData, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	if s.compress {
		compressed, err := Compress(jsonData)
		if err != nil {
			return err
		}

		return os.WriteFile(fileName, compressed, 0644)
	}
	return os.WriteFile(fileName, jsonData, 0644)
}

// Get retrieves a value for a given key from disk
func (s *DiskStore[K, V]) Get(key K) (V, bool) {
	var err error
	fileName := filepath.Join(s.basePath, fmt.Sprintf("%v.json", key))
	_, err = os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return *new(V), false
		}
		return *new(V), false
	}
	var value V
	if s.compress {
		compressedData, err := os.ReadFile(fileName)
		if err != nil {
			return *new(V), false
		}
		jsonData, err := Decompress(compressedData)
		err = msgpack.Unmarshal(jsonData, &value)
	} else {
		file, err := os.ReadFile(fileName)
		if err != nil {
			return *new(V), false
		}
		err = msgpack.Unmarshal(file, &value)
	}
	if err != nil {
		return *new(V), false
	}
	return value, true
}

// Del removes a key-value pair from disk
func (s *DiskStore[K, V]) Del(key K) error {
	fileName := filepath.Join(s.basePath, fmt.Sprintf("%v", key))
	return os.Remove(fileName)
}

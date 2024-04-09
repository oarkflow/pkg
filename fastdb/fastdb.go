package fastdb

/* ------------------------------- Imports --------------------------- */

import (
	"errors"
	"fmt"
	"sync"

	"github.com/oarkflow/pkg/fastdb/persist"
)

/* ---------------------- Constants/Types/Variables ------------------ */

// DB represents a collection of key-value pairs that persist on disk or memory.
type DB struct {
	aof  *persist.AOF
	keys map[string]map[string][]byte
	mu   sync.RWMutex
}

type Storage string

const (
	MemoryStorage Storage = ":memory:"
	DiskStorage   Storage = ":disk:"
)

type Config struct {
	StorageType Storage
	Path        string
	Filename    string
	SyncTime    int
}

func New(cfg ...Config) (*DB, error) {
	keys := make(map[string]map[string][]byte)
	var (
		aof    *persist.AOF
		err    error
		config Config
	)
	if len(cfg) > 0 {
		config = cfg[0]
	}
	if config.Path == "" {
		config.Path = "./data"
	}
	if config.Filename == "" {
		config.Filename = "fast.db"
	}
	if config.SyncTime == 0 {
		config.SyncTime = 100
	}

	switch config.StorageType {
	case DiskStorage:
		aof, keys, err = persist.New(config.Path, config.Filename, config.SyncTime)
		return &DB{aof: aof, keys: keys}, err
	default:
		return &DB{aof: aof, keys: keys}, err
	}
}

// Optimize optimizes the file to reflect the latest state.
func (fdb *DB) Optimize() error {
	fdb.mu.Lock()
	defer fdb.mu.Unlock()

	var err error

	err = fdb.aof.Optimize(fdb.keys)
	if err != nil {
		err = fmt.Errorf("defrag error: %w", err)
	}

	return err
}

// Del deletes one map value in a bucket.
func (fdb *DB) Del(bucket string, key string) (bool, error) {
	var err error

	fdb.mu.Lock()
	defer fdb.mu.Unlock()

	// bucket exists?
	_, found := fdb.keys[bucket]
	if !found {
		return found, nil
	}

	// key exists in bucket?
	_, found = fdb.keys[bucket][key]
	if !found {
		return found, nil
	}

	if fdb.aof != nil {
		lines := "del\n" + bucket + "_" + key + "\n"

		err = fdb.aof.Write(lines)
		if err != nil {
			return false, fmt.Errorf("del->write error: %w", err)
		}
	}

	delete(fdb.keys[bucket], key)

	if len(fdb.keys[bucket]) == 0 {
		delete(fdb.keys, bucket)
	}

	return true, nil
}

// Get returns one map value from a bucket.
func (fdb *DB) Get(bucket string, key string) ([]byte, bool) {
	fdb.mu.RLock()
	defer fdb.mu.RUnlock()

	data, ok := fdb.keys[bucket][key]

	return data, ok
}

// GetAll returns all map values from a bucket.
func (fdb *DB) GetAll(bucket string) (map[string][]byte, error) {
	fdb.mu.RLock()
	defer fdb.mu.RUnlock()

	bmap, found := fdb.keys[bucket]
	if !found {
		return nil, errors.New("bucket not found")
	}

	return bmap, nil
}

// Info returns info about the storage.
func (fdb *DB) Info() string {
	count := 0
	for i := range fdb.keys {
		count += len(fdb.keys[i])
	}

	return fmt.Sprintf("%d record(s) in %d bucket(s)", count, len(fdb.keys))
}

// Set stores one map value in a bucket.
func (fdb *DB) Set(bucket string, key string, value []byte) error {
	fdb.mu.Lock()
	defer fdb.mu.Unlock()

	if fdb.aof != nil {
		lines := "set\n" + bucket + "_" + key + "\n" + string(value) + "\n"

		err := fdb.aof.Write(lines)
		if err != nil {
			return fmt.Errorf("sel->write error: %w", err)
		}
	}
	if fdb.keys == nil {
		fdb.keys = make(map[string]map[string][]byte)
	}
	_, found := fdb.keys[bucket]
	if !found {
		fdb.keys[bucket] = map[string][]byte{}
	}

	fdb.keys[bucket][key] = value

	return nil
}

// Close closes the database.
func (fdb *DB) Close() error {
	if fdb.aof != nil {
		fdb.mu.Lock()
		defer fdb.mu.Unlock()

		err := fdb.aof.Close()
		if err != nil {
			return fmt.Errorf("close error: %w", err)
		}
	}

	fdb.keys = make(map[string]map[string][]byte)

	return nil
}

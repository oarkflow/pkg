package konf

import (
	"os"
	"sync"

	json "github.com/bytedance/sonic"

	"github.com/oarkflow/pkg/flat"
)

type Konf struct {
	data    map[string]any
	mu      *sync.RWMutex
	file    string
	autoEnv bool
}

func New(file string, autoEnv bool) *Konf {
	return &Konf{
		data:    make(map[string]any),
		mu:      &sync.RWMutex{},
		file:    file,
		autoEnv: autoEnv,
	}
}

func (c *Konf) Add(key string, value any) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
}

func (c *Konf) Get(key string) any {
	return c.data[key]
}

func (c *Konf) Read() error {
	file, err := os.ReadFile(c.file)
	if err != nil {
		return err
	}
	var cfg map[string]any
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		return err
	}
	mp, err := flat.Flatten(cfg, nil)
	if err != nil {
		return err
	}
	flat.MergeMap(mp, c.data)
	if c.autoEnv {
		c.data = flat.MapReplace(mp, flat.FlattenEnv(nil))
	} else {
		c.data = mp
	}
	return nil
}

func (c *Konf) All() map[string]any {
	return c.data
}

func (c *Konf) Unmarshal(dst any) error {
	if c.data == nil {
		err := c.Read()
		if err != nil {
			return err
		}
	}
	merged, err := flat.Unflatten(c.data, nil)
	if err != nil {
		return err
	}
	bt, err := json.Marshal(merged)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, dst)
}

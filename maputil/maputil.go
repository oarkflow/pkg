package maputil

import (
	"encoding/json"
)

func CopyMap(m map[string]any) map[string]any {
	cp := make(map[string]any)
	for k, v := range m {
		vm, ok := v.(map[string]any)
		if ok {
			cp[k] = CopyMap(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

func ToMap[T any, V any](data T) (V, error) {
	var v V
	bt, err := json.Marshal(data)
	if err != nil {
		return v, err
	}
	err = json.Unmarshal(bt, &v)
	return v, err
}

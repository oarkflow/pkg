package sliceutil

import (
	"fmt"
	"reflect"
	"strings"
)

func ConcatSlices[T any](slices ...[]T) []T {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	result := make([]T, totalLen)
	var i int
	for _, s := range slices {
		i += copy(result[i:], s)
	}
	return result
}

// GroupBy is a generic function that groups items by a specified key.
func GroupBy[T any](items []T, fieldNames ...string) map[string][]T {
	var t T
	kind := reflect.ValueOf(t).Kind()
	grouped := make(map[string][]T)
	keys := make([]string, len(fieldNames))
	for _, item := range items {
		switch kind {
		case reflect.Struct:
			itemValue := reflect.ValueOf(item)
			for i, fieldName := range fieldNames {
				fieldValue := itemValue.FieldByName(fieldName)
				keys[i] = fmt.Sprint(fieldValue.Interface())
			}
		case reflect.Map:
			if m, ok := any(item).(map[string]any); ok {
				for i, fieldName := range fieldNames {
					if val, ok := m[fieldName]; ok {
						keys[i] = fmt.Sprint(val)
					}
				}
			}
		}

		k := strings.Join(keys, "_")
		if _, ok := grouped[k]; !ok {
			grouped[k] = make([]T, 0, len(items)/len(keys))
		}
		grouped[k] = append(grouped[k], item)
	}
	return grouped
}

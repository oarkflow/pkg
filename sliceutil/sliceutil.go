package sliceutil

import (
	"github.com/samber/lo"
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

func FilterMapByFieldValue(data []map[string]any, field string, value any) []map[string]any {
	return lo.FilterMap[map[string]any, map[string]any](data, func(item map[string]any, index int) (map[string]any, bool) {
		if val, ok := item[field]; ok {
			if val == value {
				return item, true
			}
		}
		return nil, false
	})
}

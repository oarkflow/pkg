package search

import (
	"strconv"
	"strings"
)

func isEqual(dataVal, val any) bool {
	switch val := val.(type) {
	case string:
		switch gtVal := dataVal.(type) {
		case string:
			return strings.EqualFold(val, gtVal)
		}
		return false
	case int:
		switch gtVal := dataVal.(type) {
		case int:
			return val == gtVal
		case uint:
			return val == int(gtVal)
		case float64:
			return float64(val) == gtVal
		}
		return false
	case float64:
		switch gtVal := dataVal.(type) {
		case int:
			return val == float64(gtVal)
		case uint:
			return val == float64(gtVal)
		case float64:
			return val == gtVal
		}
		return false
	case bool:
		switch gtVal := dataVal.(type) {
		case bool:
			return val == gtVal
		case string:
			v, err := strconv.ParseBool(gtVal)
			if err != nil {
				return false
			}
			return val == v
		}
		return false
	}
	return false
}

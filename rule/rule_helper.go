package rule

func searchDeeplyNestedSlice(nestedSlice []interface{}, targetSlice []interface{}) bool {
	targetMap := make(map[interface{}]struct{})
	for _, target := range targetSlice {
		targetMap[target] = struct{}{}
	}

	for _, element := range nestedSlice {
		switch v := element.(type) {
		case []interface{}:
			if searchDeeplyNestedSlice(v, targetSlice) {
				return true
			}
		default:
			if _, found := targetMap[v]; found {
				return true
			}
		}
	}
	return false
}

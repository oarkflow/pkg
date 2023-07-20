package rule

// searchDeeplyNestedSlice searches for a target slice in a nested slice.
// It returns true if any of the target slice elements are found in the nested slice
// or in any of its nested slices. Otherwise, it returns false.
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

// flattenSlice flattens a nested slice into a single slice.
func flattenSlice(slice []interface{}) []interface{} {
	var result []interface{}
	for _, element := range slice {
		switch element := element.(type) {
		case []interface{}:
			result = append(result, flattenSlice(element)...)
		default:
			result = append(result, element)
		}
	}
	return result
}

// sumIntSlice sums up all the elements in a slice and returns the result.
func sumIntSlice(slice []any) int {
	var sum int
	for _, element := range slice {
		switch element := element.(type) {
		case int:
			sum += element
		case float64:
			sum += int(element)
		}
	}
	return sum
}

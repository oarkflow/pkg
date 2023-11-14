package algo

// FindMatchInArray return StartIndexElementWithPartOfPattern, StopIndexElementWithPartOfPattern, ShiftBeginningPatternInFirstElement, err
func FindMatchInArray(arr []string, pattern string) (int, int, int, error) {

	strL := ""
	for _, item := range arr {
		strL += item
	}
	shift, err := KMPSearch(strL, pattern)
	if err != nil {
		return 0, 0, 0, err
	}

	globalSize := 0
	index := 0
	for i, item := range arr {
		globalSize += len(item)
		index = i
		if shift < globalSize {
			break
		}
	}
	LenOfFirstPart := len(arr[index]) - (globalSize - shift)

	shiftSize := globalSize
	shiftIndex := index
	if globalSize < shift+len(pattern) {
		for _, item := range arr[index+1:] {

			shiftSize += len(item)
			shiftIndex += 1
			if shift+len(pattern) < shiftSize+1 {
				break
			}
		}
	}
	return index, shiftIndex, LenOfFirstPart, nil
}

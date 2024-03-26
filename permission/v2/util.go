package v2

func MatchResource(value, pattern string) bool {
	var i, j int
	for i < len(value) && value[i] == ' ' {
		i++
	}
	for j < len(pattern) && pattern[j] == ' ' {
		j++
	}
	for ; i < len(value) && j < len(pattern); i, j = i+1, j+1 {
		if pattern[j] == '*' {
			for i < len(value) {
				i++
			}
			break
		}
		if pattern[j] == ':' {
			for j < len(pattern) && pattern[j] != '/' && pattern[j] != '*' {
				j++
			}
			for i < len(value) && value[i] != '/' && value[i] != ' ' {
				i++
			}
			continue
		}
		if pattern[j] != value[i] {
			return false
		}
	}
	return (i == len(value) && j == len(pattern)) || (j == len(pattern) && pattern[j-1] == '*')
}

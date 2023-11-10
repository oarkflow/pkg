package docx

import (
	"regexp"
	"strings"
)

var (
	funcPattern = regexp.MustCompile(`([^,]+\(.+?\))|([^,]+)`)
	argPattern  = regexp.MustCompile(`([^,]+\(.+?\))|([^,]+)`)
	strReplacer = strings.NewReplacer("{", "", "}", "", "$$", `"`, `'`, `"`)
	numPattern  = regexp.MustCompile("^-?\\d*(\\.\\d+)?$")
)

func IsFunction(str string) bool {
	return str != "" && strings.Contains(str, "(") && strings.Contains(str, ")")
}

func GetFunction(str string) [][]string {
	return argPattern.FindAllStringSubmatch(str, -1)
}

func GetArguments(str string) [][]string {
	return argPattern.FindAllStringSubmatch(str, -1)
}

func IsNumber(str string) bool {
	return str != "" && numPattern.MatchString(str)
}

package permission

import "strings"

func GetDomain(args ...interface{}) (interface{}, error) {
	key := args[0].(string)
	if key == "" {
		return "*", nil
	}
	domainWithWork := strings.Split(key, ":")
	return domainWithWork[0], nil
}

func GetWork(args ...interface{}) (interface{}, error) {
	key := args[0].(string)
	if key == "" {
		return "*", nil
	}
	domainWithWork := strings.Split(key, ":")
	if len(domainWithWork) == 2 {
		return domainWithWork[1], nil
	}
	return "*", nil
}

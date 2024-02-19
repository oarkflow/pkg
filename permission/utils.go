package permission

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/casbin/casbin/v2/model"
	"github.com/oarkflow/govaluate"
	"gopkg.in/yaml.v3"

	"github.com/oarkflow/pkg/flat"
	"github.com/oarkflow/pkg/str"
)

func roleModel() interface{} {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, dom, obj, act")                                                                                                                                                                                        // [request_definition]
	m.AddDef("p", "p", "sub, dom, obj, act")                                                                                                                                                                                        // [policy_definition]
	m.AddDef("g", "g", "_, _, _")                                                                                                                                                                                                   // [role_definition]
	m.AddDef("g2", "g2", "_, _")                                                                                                                                                                                                    // [role_definition]
	m.AddDef("e", "e", "some(where (p.eft == allow))")                                                                                                                                                                              // [policy_effect]
	m.AddDef("m", "m", "((g(r.sub, p.sub, get_domain(r.dom)) && get_domain(r.dom) == get_domain(p.dom) && globMatch(get_work(r.dom), get_work(p.dom))) || g2(r.dom, p.dom)) && globMatch(r.obj, p.obj) && globMatch(r.act, p.act)") // [matchers]
	return m
}

func getMapFromString(text string) (map[string]any, error) {
	var t map[string]any
	st := str.ToByte(text)
	err := json.Unmarshal(st, &t)
	if err != nil {
		err = yaml.Unmarshal(st, &t)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

var (
	sliceRe = regexp.MustCompile(`\b\w+\.\d+\b`)
)

func IsMatch(item map[string]any, fields map[string]any) (bool, error) {
	i, err := flat.Flatten(item, nil)
	if err != nil {
		return false, nil
	}
	f, err := flat.Flatten(fields, nil)
	if err != nil {
		return false, nil
	}
	for k, v := range f {
		vt := fmt.Sprintf("%v", v)
		if vt == "*" {
			continue
		}
		if t, ok := i[k]; !ok {
			if sliceRe.MatchString(k) {
				iKey := k[:strings.Index(k, ".")]
				iVal := fields[iKey]
				switch reflect.TypeOf(iVal).Kind() {
				case reflect.Slice:
					tc := fmt.Sprintf("%v", item[iKey])
					var strSlice []string
					s := reflect.ValueOf(iVal)
					for i := 0; i < s.Len(); i++ {
						strSlice = append(strSlice, fmt.Sprintf("%v", s.Index(i)))
					}
					return str.Contains(strSlice, tc), nil
				}
			}
			return false, nil // errors.New("field not found: " + k)
		} else {
			if !reflect.DeepEqual(v, t) {
				return false, nil
			}
		}
	}
	return true, nil
}

var CasFunc = map[string]govaluate.ExpressionFunction{
	"isMatch": func(args ...interface{}) (interface{}, error) {
		switch attributes := args[0].(type) {
		case map[string]any:
			switch condition := args[1].(type) {
			case map[string]any:
				return IsMatch(attributes, condition)
			case string:
				t, err := getMapFromString(condition)
				if err != nil || len(t) == 0 {
					return true, nil
				}
				return IsMatch(attributes, t)
			}
		case string:
			if attributes != "" {
				attr, err := getMapFromString(attributes)
				if err != nil {
					return false, nil
				}
				switch condition := args[1].(type) {
				case map[string]any:
					return IsMatch(attr, condition)
				case string:
					t, err := getMapFromString(condition)
					if err != nil {
						return true, nil
					}
					if len(t) == 0 {
						return str.EqualFold(attributes, condition), nil
					}
					return IsMatch(attr, t)
				}
			}
		}
		return false, nil
	},
}

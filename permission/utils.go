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
		switch requestData := args[0].(type) {
		case map[string]any:
			switch policyCondition := args[1].(type) {
			case map[string]any:
				return IsMatch(requestData, policyCondition)
			case string:
				t, err := getMapFromString(policyCondition)
				if err != nil || len(t) == 0 {
					return true, nil
				}
				return IsMatch(requestData, t)
			}
		case string:
			if requestData != "" {
				attr, err := getMapFromString(requestData)
				if err != nil {
					return false, nil
				}
				switch policyCondition := args[1].(type) {
				case map[string]any:
					return IsMatch(attr, policyCondition)
				case string:
					if policyCondition == "" {
						return true, nil
					}
					t, err := getMapFromString(policyCondition)
					if err != nil {
						return true, nil
					}
					if len(t) == 0 {
						return str.EqualFold(requestData, policyCondition), nil
					}
					return IsMatch(attr, t)
				}
			} else {
				attr := make(map[string]any)
				switch policyCondition := args[1].(type) {
				case map[string]any:
					return IsMatch(attr, policyCondition)
				case string:
					if policyCondition == "" {
						return true, nil
					}
					t, err := getMapFromString(policyCondition)
					if err != nil {
						return true, nil
					}
					if len(t) == 0 {
						return str.EqualFold(requestData, policyCondition), nil
					}
					return IsMatch(attr, t)
				}
			}
		}
		return false, nil
	},
	"relatedDomain": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 || Instance == nil {
			return args[0], nil
		}
		reqDomain := args[0]
		ds := Instance.GetFilteredNamedGroupingPolicy("g", 0, args[1].(string))
		for _, dGroup := range ds {
			if len(dGroup) == 4 {
				d := strings.TrimSpace(dGroup[2])
				if d == "*" {
					return reqDomain, nil
				}
				if dGroup[3] == "true" {
					if d != "" {
						domains := Instance.GetRelatedDomains(d)
						if str.Contains(domains, reqDomain.(string)) {
							return d, nil
						}
					}
				}
			}
		}
		return reqDomain, nil
	},
	// inArray is used to parse the policy entities and match the request entity.
	"inArray": func(args ...interface{}) (interface{}, error) {
		entity := args[0].(string)
		entity = strings.TrimSpace(entity)
		if entity == "" {
			return false, nil
		}
		policyEntities := args[1].(string)
		// We assume that the entities in the policy are separated by commas.
		for _, policyEntity := range strings.Split(policyEntities, ",") {
			matched, err := regexp.MatchString(policyEntity, entity)
			if err != nil {
				return false, err
			}
			if matched {
				return true, nil
			}
		}
		return false, nil
	},

	// hasAccess is used to parse the policy entities and match the request entity.
	"hasAccess": func(args ...interface{}) (interface{}, error) {
		if len(args) != 2 || Instance == nil {
			return true, nil
		}
		sub := args[0]
		entity := args[1]
		ds := Instance.GetFilteredNamedGroupingPolicy("g", 0, sub.(string))
		for _, dGroup := range ds {
			if len(dGroup) == 5 {
				d := strings.TrimSpace(dGroup[4])
				if d == "*" || d == "" {
					return true, nil
				}
				splitted := strings.Split(d, ",")
				for _, policyEntity := range splitted {
					policyEntity = strings.TrimSpace(policyEntity)
					matched, err := regexp.MatchString(policyEntity, entity.(string))
					if err != nil {
						return false, nil
					}
					if matched {
						return true, nil
					}
				}
			}
		}
		return false, nil
	},
}

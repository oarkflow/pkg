package rule

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/oarkflow/pkg/dipper"
	"github.com/oarkflow/pkg/evaluate"
	"github.com/oarkflow/pkg/sjson"
)

func (condition *Condition) filterMap(data Data) any {
	// use sjson to get the Value
	dataJson, err := json.Marshal(data)
	if err != nil {
		return false
	}
	if condition.Filter.Key == "" {
		return nil
	}
	if condition.Filter.LookupData == nil && condition.Filter.LookupHandler == nil {
		return nil
	} else if condition.Filter.LookupData == nil {
		condition.Filter.LookupData = condition.Filter.LookupHandler()
	}
	if condition.Filter.LookupData != nil {
		lookupData := condition.Filter.LookupData
		if condition.Filter.Condition != "" {
			c := condition.Filter.Condition
			tags := unique(re.FindAllString(c, -1))
			for _, tag := range tags {
				if strings.Contains(tag, "[data.") {
					dField := strings.ReplaceAll(strings.ReplaceAll(tag, "[data.", ""), "]", "")
					v := sjson.Get(string(dataJson), dField)
					if v.Type == sjson.JSON {
						c = strings.ReplaceAll(c, tag, fmt.Sprintf("%v", v))
					} else {
						c = strings.ReplaceAll(c, tag, fmt.Sprintf("'%v'", v))
					}
				}
			}
			eval, err := evaluate.Parse(c)
			if err == nil {
				switch d := condition.Filter.LookupData.(type) {
				case []map[string]any:
					var filteredLookupData []map[string]any
					for _, dRow := range d {
						param := evaluate.NewEvalParams(dRow)
						rs, err := eval.Eval(param)
						if err == nil {
							if rs.(bool) {
								filteredLookupData = append(filteredLookupData, dRow)
							}
						}
					}
					lookupData = filteredLookupData
				case []any:
					var filteredLookupData []map[string]any
					for _, t := range d {
						switch dRow := t.(type) {
						case map[string]any:
							param := evaluate.NewEvalParams(dRow)
							rs, err := eval.Eval(param)
							if err == nil {
								if rs.(bool) {
									filteredLookupData = append(filteredLookupData, dRow)
								}
							}
						}
					}
					lookupData = filteredLookupData
				}
			}
		}
		lookupJSON, err := json.Marshal(lookupData)
		if err != nil {
			return false
		}
		lookupData = sjson.GetBytes(lookupJSON, condition.Filter.Key).Value()
		switch lookupData := lookupData.(type) {
		case []any:
			if len(lookupData) > 0 {
				d := dipper.FilterSlice(data, condition.Field, lookupData)
				if dipper.Error(d) == nil {
					if strings.Contains(condition.Field, ".[].") {
						p := strings.Split(condition.Field, ".[].")
						left := p[0]
						if left != "" {
							dipper.Set(data, left, d)
						}
					} else {
						dipper.Set(data, condition.Field, d)
					}
				}
			} else {
				if strings.Contains(condition.Field, ".[].") {
					p := strings.Split(condition.Field, ".[].")
					left := p[0]
					if left != "" {
						dipper.Set(data, left, nil)
					}
				} else {
					dipper.Set(data, condition.Field, nil)
				}
			}

		case []string:
			if len(lookupData) > 0 {
				var t []any
				for _, a := range lookupData {
					t = append(t, a)
				}
				d := dipper.FilterSlice(data, condition.Field, t)
				if dipper.Error(d) == nil {
					if strings.Contains(condition.Field, ".[].") {
						p := strings.Split(condition.Field, ".[].")
						left := p[0]
						if left != "" {
							dipper.Set(data, left, d)
						}
					} else {
						dipper.Set(data, condition.Field, d)
					}
				}
			} else {
				if strings.Contains(condition.Field, ".[].") {
					p := strings.Split(condition.Field, ".[].")
					left := p[0]
					if left != "" {
						dipper.Set(data, left, nil)
					}
				} else {
					dipper.Set(data, condition.Field, nil)
				}
			}
		}
		return lookupData
	}
	return nil
}

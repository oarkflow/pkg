package str

import (
	"encoding/json"
	"fmt"
	"regexp"
)

var re = regexp.MustCompile("%<([a-zA-Z0-9_]+)>[.0-9]*[xsvTtbcdoqXxUeEfFgGp]")

// Printf support named format
func Printf(format string, params any) {
	f, p := parse(format, GetMapValues(params))
	fmt.Printf(f, p...)
}

// Printfln support named format
func Printfln(format string, params any) {
	f, p := parse(format, GetMapValues(params))
	fmt.Printf(f, p...)
	fmt.Println()
}

// Sprintf support named format
func Sprintf(format string, params any) string {
	values := GetMapValues(params)
	f, p := parse(format, values)
	return fmt.Sprintf(f, p...)
}

// Sprintfln support named format
func Sprintfln(format string, params any) string {
	f, p := parse(format, GetMapValues(params))
	return fmt.Sprintf(f+"\n", p...)
}

func parse(format string, params map[string]any) (string, []any) {
	f, n := reformat(format)
	p := make([]any, len(n))
	for i, v := range n {
		p[i] = params[v]
	}
	return f, p
}

func reformat(f string) (string, []string) {
	i := re.FindAllStringSubmatchIndex(f, -1)

	ord := []string{}
	pair := []int{0}
	for _, v := range i {
		ord = append(ord, f[v[2]:v[3]])
		pair = append(pair, v[2]-1)
		pair = append(pair, v[3]+1)
	}
	pair = append(pair, len(f))
	plen := len(pair)

	out := ""
	for n := 0; n < plen; n += 2 {
		out += f[pair[n]:pair[n+1]]
	}

	return out, ord
}

// GetMapValues convert interface to map[string]any
func GetMapValues(input any) map[string]any {
	var output = map[string]any{}
	switch input := input.(type) {
	case []byte:
		err := json.Unmarshal(input, &output)
		if err != nil {
			panic(err)
		}
	default:
		bt, err := json.Marshal(input)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(bt, &output)
		if err != nil {
			panic(err)
		}
	}

	// Get values form pointers into map
	for k, val := range output {
		switch typedValue := val.(type) {
		case *float64:
		case *float32:
		case *int64:
		case *int16:
		case *int:
			if typedValue == nil {
				output[k] = 0
			} else {
				output[k] = *typedValue
			}
			break
		case *string:
			if typedValue == nil {
				output[k] = ""
			} else {
				output[k] = *typedValue
			}
			break
		case *bool:
			if typedValue == nil || !*typedValue {
				output[k] = "false"
			} else {
				output[k] = "true"
			}
			break
		case bool:
			if typedValue {
				output[k] = "true"
			} else {
				output[k] = "false"
			}
			break
		default:
			output[k] = val
		}
	}
	return output
}

package gender

import (
	"github.com/oarkflow/pkg/str"
)

var words = map[string]string{
	"mr":   "male",
	"mr.":  "male",
	"he":   "male",
	"his":  "male",
	"Mr":   "male",
	"Mr.":  "male",
	"He":   "male",
	"His":  "male",
	"she":  "female",
	"her":  "female",
	"miss": "female",
	"ms":   "female",
	"ms.":  "female",
	"She":  "female",
	"Her":  "female",
	"Miss": "female",
	"Ms":   "female",
	"Ms.":  "female",
}

var converter = map[string]any{
	"he":   "she",
	"his":  "her",
	"He":   "She",
	"His":  "Her",
	"she":  "he",
	"her":  "his",
	"miss": "mr",
	"ms":   "mr",
	"ms.":  "mr.",
	"She":  "He",
	"Her":  "His",
	"Miss": "Mr",
	"Ms":   "Mr",
	"Ms.":  "Mr.",
	"mr":   map[string]string{"married": "ms", "single": "miss"},
	"Mr":   map[string]string{"married": "Ms", "single": "Miss"},
	"mr.":  map[string]string{"married": "ms.", "single": "miss"},
	"Mr.":  map[string]string{"married": "Ms.", "single": "Miss"},
}

func Convert(word, gender string, married ...bool) string {
	if w, exists := words[word]; exists {
		if str.ToLower(w) == str.ToLower(gender) {
			return word
		}
		e, exist := converter[word]
		if !exist {
			return word
		}
		switch e := e.(type) {
		case string:
			return e
		case map[string]string:
			if len(married) > 0 {
				return e["married"]
			} else {
				return e["single"]
			}
		}
		return word
	}
	return word
}

package series

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type stringElement struct {
	e   string
	nan bool
}

// force stringElement struct to implement Element interface
var _ Element = (*stringElement)(nil)

func (e *stringElement) Set(value interface{}) {
	switch val := value.(type) {
	case string:
		e.SetString(val)
	case int:
		e.SetInt(val)
	case float64:
		e.SetFloat(val)
	case bool:
		e.SetBool(val)
	case Element:
		e.SetElement(val)
	default:
		e.nan = true
	}
}

func (e *stringElement) SetElement(val Element) {
	e.nan = val.IsNA()
	e.e = val.String()
}

func (e *stringElement) SetBool(val bool) {
	e.nan = false
	if val {
		e.e = "true"
	} else {
		e.e = "false"
	}
}
func (e *stringElement) SetFloat(val float64) {
	if math.IsNaN(val) {
		e.nan = true
	} else {
		e.nan = false
	}
	e.e = strconv.FormatFloat(val, 'f', 6, 64)
}

func (e *stringElement) SetInt(val int) {
	e.nan = false
	e.e = strconv.Itoa(val)
}

func (e *stringElement) SetString(val string) {
	e.e = val
	if e.e == NaN {
		e.nan = true
	} else {
		e.nan = false
	}
}

func (e stringElement) Copy() Element {
	if e.IsNA() {
		return &stringElement{"", true}
	}
	return &stringElement{e.e, false}
}

func (e stringElement) IsNA() bool {
	return e.nan
}

func (e stringElement) Type() Type {
	return String
}

func (e stringElement) Val() ElementValue {
	if e.IsNA() {
		return nil
	}
	return string(e.e)
}

func (e stringElement) String() string {
	if e.IsNA() {
		return NaN
	}
	return string(e.e)
}

func (e stringElement) Int() (int, error) {
	if e.IsNA() {
		return 0, fmt.Errorf("can't convert NaN to int")
	}
	return strconv.Atoi(e.e)
}

func (e stringElement) Float() float64 {
	if e.IsNA() {
		return math.NaN()
	}
	f, err := strconv.ParseFloat(e.e, 64)
	if err != nil {
		return math.NaN()
	}
	return f
}

func (e stringElement) Bool() (bool, error) {
	if e.IsNA() {
		return false, fmt.Errorf("can't convert NaN to bool")
	}
	switch strings.ToLower(e.e) {
	case "true", "t", "1":
		return true, nil
	case "false", "f", "0":
		return false, nil
	}
	return false, fmt.Errorf("can't convert String \"%v\" to bool", e.e)
}

func (e stringElement) Eq(elem Element) bool {
	if e.IsNA() || elem.IsNA() {
		return false
	}
	return e.e == elem.String()
}

func (e stringElement) Neq(elem Element) bool {
	if e.IsNA() || elem.IsNA() {
		return false
	}
	return e.e != elem.String()
}

func (e stringElement) Less(elem Element) bool {
	if e.IsNA() || elem.IsNA() {
		return false
	}
	return e.e < elem.String()
}

func (e stringElement) LessEq(elem Element) bool {
	if e.IsNA() || elem.IsNA() {
		return false
	}
	return e.e <= elem.String()
}

func (e stringElement) Greater(elem Element) bool {
	if e.IsNA() || elem.IsNA() {
		return false
	}
	return e.e > elem.String()
}

func (e stringElement) GreaterEq(elem Element) bool {
	if e.IsNA() || elem.IsNA() {
		return false
	}
	return e.e >= elem.String()
}

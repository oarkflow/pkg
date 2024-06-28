package rule

import (
	"encoding/json"
	"regexp"

	"github.com/oarkflow/xid"

	"github.com/oarkflow/pkg/jet"
	"github.com/oarkflow/pkg/maputil"
	"github.com/oarkflow/pkg/str"
)

type Filter struct {
	LookupData    any        `json:"lookup_data"`
	LookupHandler func() any `json:"-"`
	Key           string     `json:"key"`
	Condition     string     `json:"condition"`
	LookupSource  string     `json:"lookup_source"`
}

type Expr struct {
	Value string `json:"value"`
}

type (
	Data      any
	Condition struct {
		Filter       Filter            `json:"filter"`
		Value        any               `json:"value"`
		Key          string            `json:"key"`
		ConditionKey string            `json:"condition_key"`
		Field        string            `json:"field"`
		Operator     ConditionOperator `json:"operator"`
	}
)

var re = regexp.MustCompile("\\[([^\\[\\]]*)\\]")

func unique(s []string) []string {
	inResult := make(map[string]bool)
	var result []string
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func NewCondition(field string, operator ConditionOperator, value any, filter ...Filter) *Condition {
	var f Filter
	if len(filter) > 0 {
		f = filter[0]
	}
	return &Condition{
		Field:    field,
		Operator: operator,
		Value:    value,
		Filter:   f,
	}
}

type CallbackFn func(data Data) any

type Node interface {
	Apply(d Data, callback ...CallbackFn) any
}

type Conditions struct {
	Operator  JoinOperator `json:"operator,omitempty"`
	id        string
	Condition []*Condition `json:"condition,omitempty"`
	Reverse   bool         `json:"reverse"`
}

type Response struct {
	Data      Data
	Processed bool
	Result    bool
}

func (node *Conditions) Apply(d Data) Response {
	var nodeResult bool
	switch node.Operator {
	case AND:
		nodeResult = true
		for _, condition := range node.Condition {
			nodeResult = nodeResult && condition.Validate(d)
		}
		break
	case OR:
		nodeResult = false
		for _, condition := range node.Condition {
			nodeResult = nodeResult || condition.Validate(d)
		}
		break
	}
	if node.Reverse {
		nodeResult = !nodeResult
	}
	response := Response{
		Processed: true,
		Result:    nodeResult,
	}
	if nodeResult {
		response.Data = d
	}
	return response
}

type Join struct {
	Left     *Group       `json:"left,omitempty"`
	Operator JoinOperator `json:"operator,omitempty"`
	Right    *Group       `json:"right,omitempty"`
	id       string
}

func (join *Join) Apply(d Data) Response {
	leftResponse := join.Left.Apply(d)
	rightResponse := join.Right.Apply(d)
	var joinResult bool
	switch join.Operator {
	case AND:
		joinResult = leftResponse.Result && rightResponse.Result
		break
	case OR:
		joinResult = leftResponse.Result || rightResponse.Result
		break
	}
	response := Response{
		Processed: true,
		Result:    joinResult,
	}
	if joinResult {
		response.Data = d
	}
	return response
}

type Group struct {
	Left     *Conditions  `json:"left,omitempty"`
	Operator JoinOperator `json:"operator,omitempty"`
	Right    *Conditions  `json:"right,omitempty"`
	id       string
}

func (group *Group) Apply(d Data) Response {
	resultLeft := group.Left.Apply(d)
	resultRight := group.Right.Apply(d)
	var groupResult bool
	switch group.Operator {
	case AND:
		groupResult = resultLeft.Result && resultRight.Result
		break
	case OR:
		groupResult = resultLeft.Result || resultRight.Result
		break
	}
	response := Response{
		Processed: true,
		Result:    groupResult,
	}
	if groupResult {
		response.Data = d
	}
	return response
}

type Option struct {
	ID          string `json:"id"`
	ErrorMsg    string `json:"error_msg"`
	ErrorAction string `json:"error_action"` // warning message, restrict, restrict + warning message
}

type ErrorResponse struct {
	ErrorMsg    string `json:"error_msg"`
	ErrorAction string `json:"error_action"` // warning message, restrict, restrict + warning message
}

func (e *ErrorResponse) Error() string {
	bt, _ := json.Marshal(e)
	return str.FromByte(bt)
}

type Rule struct {
	successHandler CallbackFn
	ID             string        `json:"id,omitempty"`
	ErrorMsg       string        `json:"error_msg"`
	ErrorAction    string        `json:"error_action"`
	Conditions     []*Conditions `json:"conditions"`
	Groups         []*Group      `json:"groups"`
	Joins          []*Join       `json:"joins"`
}

func New(id ...string) *Rule {
	rule := &Rule{}
	if len(id) > 0 {
		rule.ID = id[0]
	} else {
		rule.ID = xid.New().String()
	}
	return rule
}

func (r *Rule) addNode(operator JoinOperator, condition ...*Condition) *Conditions {
	node := &Conditions{
		Condition: condition,
		Operator:  operator,
		id:        xid.New().String(),
	}
	r.Conditions = append(r.Conditions, node)
	return node
}

func (r *Rule) And(condition ...*Condition) *Conditions {
	return r.addNode(AND, condition...)
}

func (r *Rule) Or(condition ...*Condition) *Conditions {
	return r.addNode(OR, condition...)
}

func (r *Rule) Not(condition ...*Condition) *Conditions {
	return r.addNode(NOT, condition...)
}

func (r *Rule) Group(left *Conditions, operator JoinOperator, right *Conditions) *Group {
	group := &Group{
		Left:     left,
		Operator: operator,
		Right:    right,
		id:       xid.New().String(),
	}
	r.Groups = append(r.Groups, group)
	return group
}

func (r *Rule) Join(left *Group, operator JoinOperator, right *Group) *Join {
	join := &Join{
		Left:     left,
		Operator: operator,
		Right:    right,
		id:       xid.New().String(),
	}
	r.Joins = append(r.Joins, join)
	return join
}

func (r *Rule) OnSuccess(handler CallbackFn) {
	r.successHandler = handler
}

func (r *Rule) apply(d Data) Data {
	result := r.Validate(d)
	if !result {
		return nil
	}
	return d
}

func (r *Rule) Validate(d Data) bool {
	var result, n, g, j bool
	for i, node := range r.Conditions {
		if len(node.Condition) == 0 {
			continue
		}
		if i == 0 && node.Operator == AND {
			n = true
		} else if i == 0 && node.Operator == OR {
			n = false
		}
		response := node.Apply(d)
		switch node.Operator {
		case AND:
			n = n && response.Result
			break
		case OR:
			n = n || response.Result
			break
		}
	}
	if len(r.Groups) == 0 {
		result = n
	}
	for i, group := range r.Groups {
		if i == 0 && group.Operator == AND {
			g = true
		} else if i == 0 && group.Operator == OR {
			g = false
		}
		response := group.Apply(d)
		switch group.Operator {
		case AND:
			g = g && response.Result
			break
		case OR:
			g = g || response.Result
			break
		}
	}
	if len(r.Groups) > 0 && len(r.Joins) == 0 {
		result = g
	}
	for i, join := range r.Joins {
		if i == 0 && join.Operator == AND {
			j = true
		} else if i == 0 && join.Operator == OR {
			j = false
		}
		response := join.Apply(d)
		switch join.Operator {
		case AND:
			j = j && response.Result
			break
		case OR:
			j = j || response.Result
			break
		}
	}
	if len(r.Joins) > 0 {
		result = j
	}
	return result
}

func (r *Rule) Apply(d Data, callback ...CallbackFn) (any, error) {
	defaultCallbackFn := func(data Data) any {
		return data
	}

	if len(callback) > 0 {
		defaultCallbackFn = callback[0]
	}
	switch d := d.(type) {
	case map[string]any:
		dt := maputil.CopyMap(d)
		rt := r.apply(dt)
		if rt == nil && r.ErrorAction != "" {
			errorMsg, _ := jet.Parse(r.ErrorMsg, d)
			return nil, &ErrorResponse{
				ErrorMsg:    errorMsg,
				ErrorAction: r.ErrorAction,
			}
		}
		if rt != nil && r.successHandler != nil {
			return r.successHandler(rt), nil
		}
		return defaultCallbackFn(rt), nil
	case []map[string]any:
		var data []map[string]any
		for _, line := range d {
			l := maputil.CopyMap(line)
			result := r.apply(l)
			if result != nil {
				data = append(data, line)
			}
		}
		if len(data) == 0 && r.ErrorAction != "" {
			errorMsg, _ := jet.Parse(r.ErrorMsg, d)
			return nil, &ErrorResponse{
				ErrorMsg:    errorMsg,
				ErrorAction: r.ErrorAction,
			}
		}
		if len(data) > 0 && r.successHandler != nil {
			return r.successHandler(data), nil
		}
		return defaultCallbackFn(data), nil
	case []any:
		var data []map[string]any
		for _, line := range d {
			switch line := line.(type) {
			case map[string]any:
				l := maputil.CopyMap(line)
				result := r.apply(l)
				if result != nil {
					data = append(data, line)
				}
			}
		}
		if len(data) == 0 && r.ErrorAction != "" {
			errorMsg, _ := jet.Parse(r.ErrorMsg, d)
			return nil, &ErrorResponse{
				ErrorMsg:    errorMsg,
				ErrorAction: r.ErrorAction,
			}
		}
		if len(data) > 0 && r.successHandler != nil {
			return r.successHandler(data), nil
		}
		return defaultCallbackFn(data), nil
	}
	return defaultCallbackFn(nil), nil
}

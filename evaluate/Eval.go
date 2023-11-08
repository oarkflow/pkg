package evaluate

import (
	"fmt"
	"sync"
)

type EvalParams struct {
	Variables map[string]interface{}
	Operators map[string]Operator
}

func (expr ExprNode) Eval(params EvalParams) (interface{}, error) {
	switch expr.Type {
	case NodeTypeLiteral:
		return expr.Value, nil
	case NodeTypeVariable:
		/*
			var value any
				var ok bool
				if strings.Contains(expr.Name, ".") {
					bt, _ := json.Marshal(params.Variables)
					rs := sjson.GetBytes(bt, expr.Name)
					if !rs.Exists() {
						value = nil
						ok = false
					} else {
						value = rs.Value()
						ok = true
					}
				} else {
					value, ok = params.Variables[expr.Name]
				}
		*/
		value, ok := params.Variables[expr.Name]
		if !ok {
			return nil, fmt.Errorf("variable undefined: %v [pos=%d; len=%d]", expr.Name, expr.SourcePos, expr.SourceLen)
		}

		// Check if var is a node that can be Eval'd
		node, nodeType := value.(ExprNode)
		if !nodeType {
			return value, nil
		}

		for _, v := range node.Vars() {
			if v == expr.Name {
				return nil, fmt.Errorf("variable can not refer to itself: %v [pos=%d; len=%d]", expr.Name, expr.SourcePos, expr.SourceLen)
			}
		}
		return node.Eval(params)
	case NodeTypeOperator:
		operator, ok := params.Operators[expr.Name]
		if !ok {
			return nil, fmt.Errorf("operator undefined: %v [pos=%d; len=%d]", expr.Name, expr.SourcePos, expr.SourceLen)
		}
		return operator(EvalContext{params: params, expr: expr})
	}
	return nil, fmt.Errorf("bad expr type: %v", expr)
}

type Operators struct {
	builtin map[string]Operator
	mu      *sync.RWMutex
}

func (o *Operators) Add(key string, ops Operator) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.builtin[key] = ops
}

var defaultOperators *Operators

func AddCustomOperator(key string, ops Operator) {
	defaultOperators.Add(key, ops)
}

func init() {
	defaultOperators = &Operators{
		builtin: BuiltinOperators(),
		mu:      &sync.RWMutex{},
	}
}

func NewEvalParams(variables map[string]interface{}) EvalParams {
	return EvalParams{
		Variables: variables,
		Operators: defaultOperators.builtin,
	}
}

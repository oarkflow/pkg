package rule

type JoinOperator string

const (
	AND JoinOperator = "AND"
	OR  JoinOperator = "OR"
	NOT JoinOperator = "NOT"
)

// Creating rule to work with data of type map[string]any

type ConditionOperator string

const (
	EQ          ConditionOperator = "eq"
	NEQ         ConditionOperator = "neq"
	GT          ConditionOperator = "gt"
	LT          ConditionOperator = "lt"
	GTE         ConditionOperator = "gte"
	LTE         ConditionOperator = "lte"
	EqCount     ConditionOperator = "eq_count"
	NeqCount    ConditionOperator = "neq_count"
	GtCount     ConditionOperator = "gt_count"
	LtCount     ConditionOperator = "lt_count"
	GteCount    ConditionOperator = "gte_count"
	LteCount    ConditionOperator = "lte_count"
	BETWEEN     ConditionOperator = "between"
	IN          ConditionOperator = "in"
	NotIn       ConditionOperator = "not_in"
	CONTAINS    ConditionOperator = "contains"
	NotContains ConditionOperator = "not_contains"
	StartsWith  ConditionOperator = "starts_with"
	EndsWith    ConditionOperator = "ends_with"
	NotZero     ConditionOperator = "not_zero"
	IsZero      ConditionOperator = "is_zero"
	IsNull      ConditionOperator = "is_null"
	NotNull     ConditionOperator = "not_null"
)

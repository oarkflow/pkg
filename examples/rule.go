package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/oarkflow/pkg/evaluate"
	"github.com/oarkflow/pkg/rule"
	"github.com/oarkflow/pkg/timeutil"
)

var data1 = map[string]any{
	"patient": map[string]any{
		"first_name": "John",
		"gender":     "male",
		"salary":     25000,
		"dob":        "1989-04-19",
	},
	"cpt": map[string]any{
		"code": "code1",
	},
}

var data2 = map[string]any{
	"patient": map[string]any{
		"first_name": "Michelle",
		"gender":     "female",
	},
	"cpt": map[string]any{
		"code": "code3",
	},
}

var requestData = []byte(`
{
    "patient": {
        "dob": "2003-04-10"
    },
  "coding": [
    {
      "em": {
          "code": "123",
          "encounter_uid": 1,
          "work_item_uid": 2, 
          "billing_provider": "Test provider",
          "resident_provider": "Test Resident Provider"
      },
      "cpt": [
          {
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          },
          {
              "code": "OBS011",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          },
          {
              "code": "OBS011",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          }
      ]
    },
    {
      "em": {
          "code": "123",
          "encounter_uid": 1,
          "work_item_uid": 2, 
          "billing_provider": "Test provider",
          "resident_provider": "Test Resident Provider"
      },
      "cpt": [
          {
              "code": "OBS01",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          },
          {
              "code": "OBS011",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          },
          {
              "code": "OBS011",
              "billing_provider": "Test provider",
              "resident_provider": "Test Resident Provider"
          }
      ]
    }
  ]
}
`)

var jsonSchema = []byte(`
[
	{
		"conditions": [
			{
				"condition": [
					{
						"key": "check-blacklist-cpt-by-age",
						"field": "coding.#.cpt.#.code",
						"operator": "in",
						"filter": {
							"key": ".[].code",
							"lookup_data": [
								{
									"code": "OBS01",
									"min_age": 18,
									"max_age": 54
								},
								{
									"code": "JUS01",
									"min_age": 34,
									"max_age": 54
								}
							],
							"condition": "age([data.patient.dob]) >= min_age && age([data.patient.dob]) <= max_age"
						}
					}
				],
				"operator": "AND",
				"reverse": true
			}
		],
		"error_msg": "Invalid CPT/ICD Code by age",
		"error_action": "DENY"
	},
	{
		"groups": [
			{
				"left": {
					"condition": [
						{
							"key": "check-greater-than-two-obs",
							"field": "coding.#.cpt.#.code",
							"operator": "gte_count",
							"value": "2",
							"filter": {
								"key": ".[].code",
								"lookup_data": [
									{
										"code": "OBS01",
										"no_charge": "1"
									},
									{
										"code": "JUS01",
										"no_charge": "1"
									}
								]
							}
						}
					],
					"operator": "AND",
					"reverse": true
				},
				"operator": "OR",
				"right": {
					"condition": [
						{
							"key": "check-one-obs-cpt",
							"field": "coding.#.cpt.#.code",
							"operator": "eq_count",
							"value": "1",
							"filter": {
								"key": ".[].code",
								"lookup_data": [
									{
										"code": "OBS01",
										"no_charge": "1"
									},
									{
										"code": "JUS01",
										"no_charge": "1"
									}
								]
							}
						},
						{
							"key": "check-atleast-one-obs-cpt-in-em",
							"field": "em.code",
							"operator": "in",
							"value": [
								"OBS01",
								"JUS01"
							]
						}
					],
					"operator": "AND",
					"reverse": true
				}
			}
		],
		"error_msg": "Invalid Code",
		"error_action": "DENY"
	}
]
`)

func builtinAge(ctx evaluate.EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(1); err != nil {
		return 0, err
	}
	left, err := ctx.Arg(0)
	if err != nil {
		return nil, err
	}
	t, err := timeutil.ParseTime(left)
	if err != nil {
		return nil, err
	}
	return timeutil.CalculateToNow(t), err
}

func twoConditionsWithAndOp() {
	// Applying expression rules:
	// 1) When using rule to count values, use expression in filter->condition
	// 2) When checking for values within specific lookup sources, use expression as value in condition
	filter := rule.Filter{
		Key: ".[].salary",
		LookupData: []map[string]any{
			{
				"title":   "Min Salary",
				"salary":  12000,
				"min_age": 18,
				"max_age": 20,
			},
			{
				"title":   "Avg Salary",
				"salary":  13000,
				"min_age": 21,
				"max_age": 30,
			},
			{
				"title":   "Max Salary",
				"salary":  25000,
				"min_age": 31,
				"max_age": 40,
			},
		},
		Condition: "age([data.patient.dob]) >= min_age && age([data.patient.dob]) <= max_age",
	}
	// c1 := rule.NewCondition("patient.salary", rule.IN, rule.Expr{Value: "age([data.patient.dob]) >= min_age && age([data.patient.dob]) <= max_age"}, filter)
	c1 := rule.NewCondition("patient.salary", rule.GteCount, 1, filter)
	c2 := rule.NewCondition("cpt.code", rule.IN, []string{"code1", "code2"})
	r1 := rule.New()
	r1.And(c1, c2)
	fmt.Println(r1.Apply(data1))
}

func groupConditions() {
	c1 := rule.NewCondition("patient.gender", rule.EQ, "male")
	c2 := rule.NewCondition("cpt.code", rule.IN, []string{"code1", "code2"})

	c3 := rule.NewCondition("patient.gender", rule.EQ, "female")
	c4 := rule.NewCondition("cpt.code", rule.IN, []string{"code13"})

	r1 := rule.New()
	var r2 rule.Rule
	node1 := r1.And(c1, c2) // condition group
	node2 := r1.And(c3, c4) // condition group
	r1.Group(node1, rule.OR, node2)
	bt, _ := json.Marshal(r1)
	fmt.Println(string(bt))
	json.Unmarshal(bt, &r2)
	// fmt.Println(r1.Apply(data2))
	fmt.Println(r2.Apply(data2))
}

func jsonConditions() {
	var rules []*rule.Rule
	var data map[string]any
	err := json.Unmarshal(jsonSchema, &rules)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(requestData, &data)
	if err != nil {
		panic(err)
	}
	for _, r := range rules {
		d, err := r.Apply(data)
		if err != nil {
			panic(err)
		}
		fmt.Println(d)
	}
}

func main() {
	evaluate.AddCustomOperator("age", builtinAge)
	// twoConditionsWithAndOp()
	// groupConditions()
	start := time.Now()
	jsonConditions()
	fmt.Printf("%s", time.Since(start))
}

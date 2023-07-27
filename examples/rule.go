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
  "company_id": 1,
  "work_item_client_refs": [
    {
      "work_item_id": 40,
      "client_ref": "SOMETHING"
    },
    {
      "work_item_id": 41,
      "client_ref": "41SOMETHING"
    },
    {
      "work_item_id": 66,
      "client_ref": null
    },
    {
      "work_item_id": 48
    }
  ],
  "first_name": "John",
  "middle_name": "P",
  "last_name": "Doeh",
  "provider_type_id": 1,
  "title": null,
  "display_name": "Doeh, John P - MD",
  "provider_lov": "Doeh, John P - MD",
  "npi": 123498765,
  "provider_email": "jpdoeh@gmail.com",
  "valid_for_report": true
}
`)

var requestData3 = []byte(`
{
  "patient_header": {
    "patient_status": "Inpatient",
    "admit_date": null,
    "injury_date": null,
    "status": "COMPLETE"
  },
  "coding": [
    {
      "dos": "2020/01/01",
      "details": {
        "pro": {
          "em": {
            "em_level": "993851",
            "billing_provider": "Burnett, Brett",
            "em_downcode": false,
            "shared": false
          },
          "downcode": [],
          "special": [],
          "cpt": [],
          "hcpcs": []
        },
        "dx": {
          "pro": [
            {
              "dx_reason": "Principal",
              "code": "L89013"
            },
            {
              "dx_reason": "Admitting"
            },
            {
              "dx_reason": "Additional"
            },
            {
              "dx_reason": "Additional"
            }
          ]
        },
        "cdi": {
          "pro": []
        },
        "notes": []
      }
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
    "error_msg": "Client reference is required for some work items provided.",
    "error_action": "DENY",
    "conditions": [
      {
        "operator": "AND",
        "condition": [
          {
            "filter": {
              "lookup_data": [
                {
                  "work_item_id": 48,
                  "client_ref_req_ind": "false"
                },
                {
                  "work_item_id": 49,
                  "client_ref_req_ind": "false"
                },
                {
                  "work_item_id": 65,
                  "client_ref_req_ind": "false"
                },
                {
                  "work_item_id": 66,
                  "client_ref_req_ind": "true"
                },
                {
                  "work_item_id": 145,
                  "client_ref_req_ind": "false"
                }
              ],
              "key": "#.work_item_id",
              "condition": "client_ref_req_ind == 'true' && work_item_id in [data.work_item_client_refs.#(client_ref==~false)#.work_item_id]",
              "lookup_source": "vw_wi_client_ref"
            },
            "value": null,
            "key": "check-client-ref-for-work-item-wi",
            "condition_key": "",
            "field": "work_item_client_refs.#.work_item_id",
            "operator": "in"
          }
        ],
        "reverse": true
      }
    ],
    "groups": null,
    "joins": null
  },
{
		"groups": [
			{
				"left": {
					"condition": [
	{
		"key": "check-greater-than-two-obs-in-cpt-pro",
		"field": "coding.#.details.pro.cpt.#.procedure_num",
		"operator": "gte_count",
		"value": "2",
		"filter": {
			"key": ".[].code",
			"lookup_data": [],
			"condition": "charge_type == 'ED_PROFEE' && [data.request_param.wid] == string(work_item_id)"
		}
	}
					],
					"operator": "AND",
					"reverse": true
				},
				"operator": "AND",
				"right": {
					"condition": [
	{
		"key": "check-one-obs-in-cpt-pro",
		"field": "coding.#.details.pro.cpt.#.procedure_num",
		"operator": "eq_count",
		"value": "1",
		"filter": {
			"key": ".[].code",
  "lookup_data": [],
			"condition": "charge_type == 'ED_PROFEE' && [data.request_param.wid] == string(work_item_id)"
		}
	},
	{
		"key": "check-obs-in-em-pro",
		"field": "coding.#.details.pro.em.em_level",
		"operator": "in",
		"filter": {
			"key": ".[].code",
      "lookup_data": [],
			"condition": "charge_type == 'ED_PROFEE' && [data.request_param.wid] == string(work_item_id)"
		}
	}
					],
					"operator": "AND",
					"reverse": true
				}
			}
		],
		"error_msg": "Multiple OBS codes found on same DOS. Please review and correct coding.",
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
	evaluate.AddCustomOperator("string", ToString)
	// twoConditionsWithAndOp()
	// groupConditions()
	start := time.Now()
	jsonConditions()
	fmt.Printf("%s", time.Since(start))
}

// ToString converts the given value to a string.
func ToString(ctx evaluate.EvalContext) (interface{}, error) {
	if err := ctx.CheckArgCount(1); err != nil {
		return nil, err
	}
	left, err := ctx.Arg(0)
	if err != nil {
		return nil, err
	}
	return fmt.Sprintf("%v", left), nil
}

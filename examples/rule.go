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
    "patient_header": {
        "disch_disp": null,
        "transfer_dest": null,
        "patient_status": null,
        "admit_date": null,
        "injury_date": null,
        "lmp_date": null,
        "status": "IN_PROGRESS",
        "message": null,
        "patient_dob": "2019-05-02",
        "patient_sex": "M"
    },
    "coding": [
        {
            "dos": "2020/01/01",
            "details": {
                "pro": {
                    "em": {
                        "em_modifier1": "8",
                        "em_downcode": false,
                        "shared": false
                    },
                    "downcode": [],
                    "special": [],
                    "cpt": [
                        {
                            "procedure_num": "AN65450",
                            "procedure_qty": 1,
                            "billing_provider": null,
                            "secondary_provider": null
                        }
                    ],
                    "hcpcs": []
                },
                "fac": {
                    "em": null,
                    "special": [],
                    "cpt": [],
                    "hcpcs": []
                },
                "dx": {
                    "pro": [],
                    "fac": []
                },
                "cdi": {
                    "pro": [],
                    "fac": []
                },
                "notes": []
            }
        }
    ]
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
[{},{"joins":null,"groups":null,"error_msg":"Invalid CPT Code by age","conditions":[{"reverse":true,"operator":"OR","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-cpt-by-age"}]}],"error_action":"DENY"},{"joins":null,"groups":null,"error_msg":"Invalid CPT Code by gender","conditions":
[{"reverse":true,"operator":"OR","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-cpt-by-gender"}]}],"error_action":"DENY"},{"joins":null,"groups":null,"error_msg":"Invalid CPT Code by age","conditions":[{"reverse":true,"operator":"OR","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-cpt-by-age-pro"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-cpt-by-age-fac"}]}],"error_action":"DENY"},{"joins":null,"groups":null,"error_msg":"Invalid ICD10 Code by gender","conditions":[{"reverse":true,"operator":"OR","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-icd10-by-gender-pro"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-icd10-by-gender-fac"}]}],"error_action":"DENY"},{"joins":null,"groups":null,"error_msg":"Invalid CPT Code by gender","conditions":[{"reverse":true,"operator":"OR","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-cpt-by-gender-pro"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-cpt-by-gender-fac"}]}],"error_action":"DENY"},{"joins":null,"groups":null,"error_msg":"Invalid ICD10 Code by age","conditions":[{"reverse":true,"operator":"OR","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-icd10-by-age-pro"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-blacklist-icd10-by-age-fac"}]}],"error_action":"DENY"},{"joins":null,"groups":[{"left":{"reverse":true,"operator":"AND","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-greater-than-two-obs-in-cpt-pro"}]},"right":{"reverse":true,"operator":"AND","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-one-obs-in-cpt-pro"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-obs-in-em-pro"}]},"operator":"AND"}],"error_msg":"Multiple OBS codes found on same DOS. Please review and correct coding.","conditions":null,"error_action":"DENY"},{"joins":null,"groups":null,"error_msg":"Two E/M codes were identified on same DOS. Please validate this is NOT a Medicare patient before completing this encounter.","conditions":[{"reverse":true,"operator":"AND","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-one-obs-in-cpt-fac"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-no-obs-in-em-fac"}]}],"error_action":"WARN"},{"joins":null,"groups":null,"error_msg":"Two E/M codes were identified on same DOS. Please validate this is NOT a Medicare patient before completing this encounter.","conditions":[{"reverse":true,"operator":"AND","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-one-obs-in-cpt-pro"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-no-obs-in-em-pro"}]}],"error_action":"WARN"},{"joins":null,"groups":[{"left":{"reverse":true,"operator":"AND","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-greater-than-two-obs-in-cpt--fac"}]},"right":{"reverse":true,"operator":"AND","condition":[{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-one-obs-in-cpt-fac"},{"key":"","field":"","value":null,"filter":{"key":"","condition":"","lookup_data":null,"lookup_source":""},"operator":"","condition_key":"check-obs-in-em-fac"}]},"operator":"AND"}],"error_msg":"Multiple OBS codes found on same DOS. Please review and correct coding.","conditions":null,"error_action":"DENY"}]
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

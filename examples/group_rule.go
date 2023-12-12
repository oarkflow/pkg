package main

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/oarkflow/pkg/maputil"
	"github.com/oarkflow/pkg/rule"
	"github.com/oarkflow/pkg/utils"
)

func main() {
	Group1()
	Group2()
	Group3()
}

type Route interface {
	ToRule(exceptFields ...string) *rule.PriorityRule
}

type ProviderRoute struct {
	Provider             string `gorm:"provider" json:"provider" form:"provider" query:"provider"`
	RouteType            string `gorm:"route_type" json:"route_type" form:"route_type" query:"route_type"` // D W OTP
	SourceAddrType       string `gorm:"column:source_addr_type;size:20;" json:"source_addr_type" form:"source_addr_type" query:"source_addr_type"`
	SourceAddr           string `gorm:"column:source_addr;size:20;" json:"source_addr" form:"source_addr" query:"source_addr"`
	MessageCountOperator string `gorm:"message_count_operator" json:"message_count_operator" form:"message_count_operator" query:"message_count_operator"`
	CountryCode          string `gorm:"country_code" json:"country_code" form:"country_code" query:"country_code"`
	CarrierCode          string `gorm:"carrier_code" json:"carrier_code" form:"carrier_code" query:"carrier_code"`
}

func prepareConditions(route Route, exceptFields []string) *rule.PriorityRule {
	priority := 0
	var conditions []*rule.Condition
	data, err := maputil.ToMap[any, map[string]any](route)
	if err != nil {
		return nil
	}
	for key, val := range data {
		if !slices.Contains(exceptFields, key) {
			if !utils.IsZeroVal(val) {
				priority++
				conditions = append(conditions, &rule.Condition{
					Field:    key,
					Operator: rule.EQ,
					Value:    val,
				})
			}
		}
	}
	r := rule.New()
	r.OnSuccess(func(data rule.Data) any {
		return route
	})
	r.And(conditions...)
	return &rule.PriorityRule{
		Rule:     r,
		Priority: priority,
	}
}

func (route *ProviderRoute) ToRule(exceptFields ...string) *rule.PriorityRule {
	return prepareConditions(route, exceptFields)
}

type UserRoute struct {
	ProviderType   string `gorm:"provider_type" json:"provider_type" form:"provider_type" query:"provider_type"`
	RouteType      string `gorm:"route_type" json:"route_type" form:"route_type" query:"route_type"` // D W OTP
	RequestType    string `gorm:"request_type" json:"request_type" form:"request_type" query:"request_type"`
	SourceAddrType string `gorm:"source_addr_type" json:"source_addr_type" form:"source_addr_type" query:"source_addr_type"`
	SourceAddr     string `gorm:"source_addr" json:"source_addr" form:"source_addr" query:"source_addr"`
	CountryCode    string `gorm:"country_code" json:"country_code" form:"country_code" query:"country_code"`
	CarrierCode    string `gorm:"carrier_code" json:"carrier_code" form:"carrier_code" query:"carrier_code"`
	Provider       string `gorm:"provider" json:"provider" form:"provider" query:"provider"`
	UserID         uint   `gorm:"user_id" json:"user_id" form:"user_id" query:"user_id"`
}

func (route *UserRoute) ToRule(exceptFields ...string) *rule.PriorityRule {
	return prepareConditions(route, exceptFields)
}

func SearchProvider(data any, routes []Route) (any, error) {
	var mp map[string]any
	dt, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(dt, &mp)
	if err != nil {
		panic(err)
	}
	var routePriorities []*rule.PriorityRule
	for _, providerRoute := range routes {
		routePriorities = append(routePriorities, providerRoute.ToRule("provider"))
	}
	ruleGroup := rule.NewRuleGroup(rule.Config{
		Rules:    routePriorities,
		Priority: rule.HighestPriority,
	})
	return ruleGroup.Apply(mp)
}

var providerRoutes = []Route{
	&UserRoute{
		Provider:    "test3",
		CountryCode: "IN",
		UserID:      1,
	},
	&ProviderRoute{
		Provider:       "routee",
		RouteType:      "W",
		SourceAddrType: "longcode",
		CountryCode:    "NP",
	},
	&ProviderRoute{
		Provider:    "routee1",
		CountryCode: "NP",
	},
	&ProviderRoute{
		Provider:       "test1",
		RouteType:      "W",
		SourceAddrType: "shortcode",
		CountryCode:    "IN",
	},
	&ProviderRoute{
		Provider:    "test2",
		CountryCode: "IN",
	},
}

var userRoute1 = UserRoute{
	CountryCode: "NP",
	UserID:      1,
}

var userRoute2 = UserRoute{
	CountryCode: "IN",
	UserID:      1,
}

var userRoute3 = UserRoute{
	CountryCode: "IN",
	UserID:      2,
}

func Group1() {
	route, err := SearchProvider(userRoute1, providerRoutes)
	if err != nil {
		panic(err)
	}
	fmt.Println(route)
}

func Group2() {
	route, err := SearchProvider(userRoute2, providerRoutes)
	if err != nil {
		panic(err)
	}
	fmt.Println(route)
}

func Group3() {
	route, err := SearchProvider(userRoute3, providerRoutes)
	if err != nil {
		panic(err)
	}
	fmt.Println(route)
}

package main

import (
	"fmt"

	"github.com/oarkflow/pkg/pattern"
)

/*
	type ProviderRoute struct {
		Provider             string `gorm:"provider" json:"provider" form:"provider" query:"provider"`
		RouteType            string `gorm:"route_type" json:"route_type" form:"route_type" query:"route_type"` // D W OTP
		SourceAddrType       string `gorm:"column:source_addr_type;size:20;" json:"source_addr_type" form:"source_addr_type" query:"source_addr_type"`
		SourceAddr           string `gorm:"column:source_addr;size:20;" json:"source_addr" form:"source_addr" query:"source_addr"`
		MessageCountOperator string `gorm:"message_count_operator" json:"message_count_operator" form:"message_count_operator" query:"message_count_operator"`
		CountryCode          string `gorm:"country_code" json:"country_code" form:"country_code" query:"country_code"`
		CarrierCode          string `gorm:"carrier_code" json:"carrier_code" form:"carrier_code" query:"carrier_code"`
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

	var providerRoutes = []*ProviderRoute{
		{
			Provider:    "test3",
			CountryCode: "IN",
		},
		{
			Provider:       "routee",
			RouteType:      "W",
			SourceAddrType: "longcode",
			CountryCode:    "NP",
		},
		{
			Provider:    "routee1",
			CountryCode: "NP",
		},
		{
			Provider:       "test1",
			RouteType:      "W",
			SourceAddrType: "shortcode",
			CountryCode:    "IN",
		},
		{
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
*/
func main() {
	a := 3
	b := 115
	result, err := pattern.
		Match(a, b).
		Case(func(args ...any) (any, error) {
			fmt.Println(args)
			return 5, nil
		}, 3, 15).
		Case(func(args ...any) (any, error) {
			fmt.Println(args)
			return 4, nil
		}, pattern.EXISTS, pattern.ANY).
		Default(func(args ...any) (any, error) {
			fmt.Println(args)
			return 2, nil
		}).
		Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(result, err)
}

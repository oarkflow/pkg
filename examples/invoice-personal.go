package main

import (
	"github.com/oarkflow/pkg/invoice"
)

func main() {
	inv, err := invoice.New(&invoice.Config{
		Business: &invoice.Business{
			Details: invoice.Contact{
				Name:      "Sujit Prasad Baniya",
				Address1:  "Battishputali",
				Logo:      "new-logo.png",
				City:      "Kathmandu",
				ZipCode:   "9",
				Country:   "Nepal",
				Telephone: "+9779856034616",
				Email:     "s.baniya.np@gmail.com",
			},
			BankDetail: &invoice.BankDetail{
				AccountName:   "SUJIT PRASAD BANIYA",
				AccountNumber: "05-2984342-53",
				BankName:      "STANDARD CHARTERED BANK NEPAL LIMITED",
				BankAddress:   "LAZIMPAT, KATHMANDU NEPAL",
				SwiftCode:     "SCBLNPKA",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	detail := &invoice.Detail{
		InvoiceNumber: "2",
		Currency:      "$",
		PaymentTerms:  "1",
		Date:          "2024-06-14",
		DueDays:       1,
		TaxRate:       "0",
		Customer: invoice.Customer{
			Details: invoice.Contact{
				Name:     "Boxen Labs SRL",
				Address1: "Str Grigore Mora 38",
				Address2: "Bucharest 011888",
				City:     "Bucharest",
				State:    "Bucharest",
				ZipCode:  "011888",
				Country:  "Romania",
			},
		},
		Items: []invoice.Item{
			{
				Description: "Software Development Service (From Github Issues ETA) till June 14th 2024",
				Quantity:    22.5,
				UnitPrice:   30,
				Currency:    "$",
			},
		},
	}
	err = inv.Create(detail).RenderToFile("Software Development Service Invoice June 14th 2024.pdf")
	if err != nil {
		panic(err)
	}
}

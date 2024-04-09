package main

import (
	"github.com/oarkflow/pkg/invoice"
)

func main() {
	inv, err := invoice.New(&invoice.Config{
		Business: &invoice.Business{
			Details: invoice.Contact{
				Name:      "Orgware Construct Pvt. Ltd.",
				Address1:  "Prachin Marg, Old Baneshwor",
				Logo:      "logo.jpg",
				City:      "Kathmandu",
				ZipCode:   "10",
				Country:   "Nepal",
				Telephone: "+977-1-4497653",
				Email:     "info@orgwareconstruct.com",
			},
			BankDetail: &invoice.BankDetail{
				AccountName:   "ORGWARE CONSTRUCT PVT. LTD.",
				AccountNumber: "08001010007253",
				BankName:      "GLOBAL IME BANK LIMITED",
				BankAddress:   "KAMALADI, 28",
				SwiftCode:     "GLBBNPKA",
			},
			ContactName:  "Sujit Baniya",
			ContactEmail: "s.baniya.np@gmail.com",
			ContactPhone: "+977-9856034616",
		},
		Secret: "r28GYSTF9oJeiXvnkIDLLqu8RGWg3VUS",
	})
	if err != nil {
		panic(err)
	}
	detail := &invoice.Detail{
		InvoiceNumber: "7",
		Currency:      "$",
		PaymentTerms:  "1",
		Date:          "2024-03-25",
		DueDays:       5,
		TaxRate:       "0",
		Customer: invoice.Customer{
			Details: invoice.Contact{
				Name:     "Edelberg + Associates",
				Address1: "1205 Johnson Ferry Rd.",
				Address2: "Suite 136-356",
				City:     "Marietta",
				State:    "GA",
				ZipCode:  "30068",
				Country:  "US",
			},
		},
		Items: []invoice.Item{
			{
				Description: "CARE 2.0 Development and Support for the month of April 2024 (Advance) ",
				Quantity:    1,
				UnitPrice:   7500,
				Currency:    "$",
			},
			{
				Description: "CARE 2.0 Development and Support for the month of March 2024 ",
				Quantity:    1,
				UnitPrice:   1250,
				Currency:    "$",
			},
		},
		/*Transactions: []invoice.Transaction{
			{
				Description:   "CARE 2.0 Development and Support for the month of May 2023 ",
				PaymentMethod: "Wire Transfer",
				Quantity:      1,
				UnitPrice:     7000,
				Currency:      "$",
			},
		},*/
	}
	err = inv.Create(detail).RenderToFile("CARE 2.0 Invoice #7 for March 2024.pdf")
	if err != nil {
		panic(err)
	}
}

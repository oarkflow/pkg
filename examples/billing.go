package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/oarkflow/pkg/pdfs"
	"github.com/oarkflow/pkg/pdfs/color"
	"github.com/oarkflow/pkg/pdfs/consts"
	"github.com/oarkflow/pkg/pdfs/props"
)

var primaryColor = color.Hex2RGB("012A4A")
var emptyRow = []string{"", "", "", "", ""}

func main() {
	begin := time.Now()
	m := pdfs.New(consts.Portrait, consts.A4)
	m.SetPageMargins(10, 15, 10)
	// m.SetBorder(true)

	header, content := getContent()
	if len(content) <= 10 {
		rowsToFill := 10 - len(content)
		for i := 0; i < rowsToFill; i++ {
			content = append(content, emptyRow)
		}
	}
	setupInvoiceDetails(m)

	m.TableList(header, content, props.TableList{
		ContentProp: props.TableListContent{
			GridSizes: []uint{1, 5, 2, 2, 2},
			Color:     *primaryColor,
			Align:     []consts.Align{consts.Left, consts.Left, consts.Right, consts.Right, consts.Right},
		},
		HeaderProp: props.TableListContent{
			GridSizes:       []uint{1, 5, 2, 2, 2},
			Align:           []consts.Align{consts.Left, consts.Left, consts.Right, consts.Right, consts.Right},
			Color:           color.NewWhite(),
			BackgroundColor: color.NewBlack(),
		},
		ContentStyles: []consts.Style{
			consts.Normal,
			consts.Bold,
			consts.Italic,
		},
		VerticalHeaderPadding:    2,
		HorizontalHeaderPadding:  2,
		VerticalContentPadding:   2,
		HorizontalContentPadding: 2,
		AlternatedBackground:     color.Hex2RGB("eaeaea"),
	})

	m.Row(5, func() {
		m.Col(2, func() {
			m.Text("Notes/Remarks:", props.Text{
				Top:   5,
				Size:  10,
				Color: *primaryColor,
			})
		})
		m.ColSpace(6)
		m.Col(2, func() {
			m.Text("Sub Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
		m.Col(2, func() {
			m.Text("$12,500.00", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
	m.Row(5, func() {
		m.ColSpace(8)
		m.Col(2, func() {
			m.Text("Tax (0%):", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
		m.Col(2, func() {
			m.Text("$0.00", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
	m.Row(5, func() {
		m.ColSpace(8)
		m.Col(2, func() {
			m.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  14,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
		m.Col(2, func() {
			m.Text("$12,500.00", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  14,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
	m.Row(15, func() {

	})
	m.Row(5, func() {
		m.ColSpace(8)
		m.Col(4, func() {
			m.Text("Thank you for your business!", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
	err := m.OutputFileAndClose("sample1.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))
}

func setupInvoiceDetails(m pdfs.IPdf) {
	setPDFHeader(m)
	setPDFFooter(m)
	m.SetAliasNbPages("{nb}")
	m.SetFirstPageNb(1)
}

func setPDFFooter(m pdfs.IPdf) {
	m.RegisterFooter(func() {
		m.Row(40, func() {
			m.Col(4, func() {
				m.Text("Payment Details:", props.Text{
					Top:   4,
					Align: consts.Left,
					Size:  12,
					Color: color.Color{
						Red:   124,
						Green: 124,
						Blue:  124,
					},
				})
				m.Text("BENEFICIARY NAME:", props.Text{
					Top:   10,
					Align: consts.Left,
					Size:  10,
				})
				m.Text("BENEFICIARY ACCOUNT NUMBER:", props.Text{
					Top:   15,
					Align: consts.Left,
					Size:  10,
				})
				m.Text("BANK NAME:", props.Text{
					Top:   20,
					Align: consts.Left,
					Size:  10,
				})
				m.Text("BANK ADDRESS:", props.Text{
					Top:   25,
					Align: consts.Left,
					Size:  10,
				})
				m.Text("BANK SWIFT CODE:", props.Text{
					Top:   30,
					Align: consts.Left,
					Size:  10,
				})
			})

			m.Col(4, func() {
				m.Text("", props.Text{
					Top:   4,
					Align: consts.Left,
					Size:  12,
					Color: color.Color{
						Red:   124,
						Green: 124,
						Blue:  124,
					},
				})
				m.Text("ORGWARE CONSTRUCT PVT. LTD.", props.Text{
					Top:   10,
					Align: consts.Left,
					Style: consts.Bold,
					Size:  10,
					Color: *primaryColor,
				})
				m.Text("08001010007253", props.Text{
					Top:   15,
					Align: consts.Left,
					Style: consts.Bold,
					Size:  10,
					Color: *primaryColor,
				})
				m.Text("GLOBAL IME BANK LIMITED", props.Text{
					Top:   20,
					Align: consts.Left,
					Style: consts.Bold,
					Size:  10,
					Color: *primaryColor,
				})
				m.Text("KAMALADI, 28", props.Text{
					Top:   25,
					Align: consts.Left,
					Style: consts.Bold,
					Size:  10,
					Color: *primaryColor,
				})
				m.Text("GLBBNPKA", props.Text{
					Top:   30,
					Align: consts.Left,
					Style: consts.Bold,
					Size:  10,
					Color: *primaryColor,
				})
			})
			m.Col(4, func() {
				m.Text("Contact Information:", props.Text{
					Top:   4,
					Align: consts.Right,
					Size:  12,
					Color: color.Color{
						Red:   124,
						Green: 124,
						Blue:  124,
					},
				})
				m.Text("Sujit Baniya", props.Text{
					Top:   10,
					Align: consts.Right,
					Size:  10,
				})
				m.Text("s.baniya.np@gmail.com", props.Text{
					Top:   15,
					Align: consts.Right,
					Size:  10,
				})
				m.Text("+977-9856034616", props.Text{
					Top:   20,
					Align: consts.Right,
					Size:  10,
				})
			})
		})
		m.Line(1)
		m.Row(3, func() {

		})
		m.Row(7, func() {
			m.Col(2, func() {
				byteSlices, err := os.ReadFile("logo.jpg")
				if err != nil {
					fmt.Println("Got error while opening file:", err)
					os.Exit(1)
				}
				base64image := base64.StdEncoding.EncodeToString(byteSlices)
				_ = m.Base64Image(base64image, consts.Jpg, props.Rect{
					Percent: 100,
				})
			})
			m.Col(4, func() {
				m.Text("Orgware Construct Pvt. Ltd.", props.Text{
					Top:   0,
					Align: consts.Right,
					Size:  10,
					Style: consts.Bold,
					Color: *primaryColor,
				})
				m.Text("Prachin Marg, Old Baneshwor", props.Text{
					Top:   4,
					Align: consts.Right,
					Size:  8,
				})
			})
			m.Col(4, func() {
				m.Text("Tel: +977-1-4497653", props.Text{
					Top:   0,
					Align: consts.Right,
					Size:  8,
				})
				m.Text("info@orgwareconstruct.com", props.Text{
					Top:   4,
					Align: consts.Right,
					Size:  8,
				})
			})
			m.Col(2, func() {
				m.Text(strconv.Itoa(m.GetCurrentPage())+"/{nb}", props.Text{
					Align: consts.Right,
					Size:  8,
				})
			})
		})
	})
}

func getContent() ([]string, [][]string) {
	header := []string{"SN", "Item", "Quantity", "Rate", "Amount"}

	contents := [][]string{
		{"1", "CARE 2.0 Development and Support for the month of April 2023", "1", "$12,500", "$12,500.00"},
		{"2", "CARE 2.0 Development and Support for the month of March 2023", "1", "$12,500", "$12,500.00"},
	}
	return header, contents
}

func setPDFHeader(m pdfs.IPdf) {
	byteSlices, err := os.ReadFile("logo.jpg")
	if err != nil {
		fmt.Println("Got error while opening file:", err)
		os.Exit(1)
	}
	base64image := base64.StdEncoding.EncodeToString(byteSlices)
	m.RegisterHeader(func() {
		m.Row(20, func() {
			m.Col(3, func() {
				_ = m.Base64Image(base64image, consts.Jpg, props.Rect{
					Percent: 100,
				})
			})

			m.ColSpace(6)

			m.Col(3, func() {
				m.Text("INVOICE", props.Text{
					Align: consts.Right,
					Size:  28,
				})
				m.Text("#4", props.Text{
					Top:   11,
					Align: consts.Right,
					Size:  18,
					Color: *color.Hex2RGB("a4a4a4"),
				})
			})
		})

		m.Line(1.0,
			props.Line{
				Color: *primaryColor,
			},
		)
		invoiceDetails(m)

	})
}

func invoiceDetails(m pdfs.IPdf) {
	m.Row(40, func() {
		m.Col(4, func() {
			m.Text("Bill From:", props.Text{
				Top:   4,
				Align: consts.Left,
				Size:  12,
				Color: color.Color{
					Red:   124,
					Green: 124,
					Blue:  124,
				},
			})
			m.Text("Orgware Construct Pvt. Ltd.", props.Text{
				Top:   10,
				Align: consts.Left,
				Size:  10,
				Style: consts.Bold,
				Color: *primaryColor,
			})

			m.Text("Prachin Marg, Old Baneshwor", props.Text{
				Size: 10,
				Top:  15,
			})
			m.Text("Kathmandu - 10, Nepal", props.Text{
				Size: 10,
				Top:  20,
			})
			m.Text("Tel: +977-1-4497653", props.Text{
				Size: 10,
				Top:  25,
			})
			m.Text("info@orgwareconstruct.com", props.Text{
				Size: 10,
				Top:  30,
			})
		})
		m.Col(3, func() {
			m.Text("Bill To:", props.Text{
				Top:   4,
				Align: consts.Left,
				Size:  12,
				Color: color.Color{
					Red:   124,
					Green: 124,
					Blue:  124,
				},
			})
			m.Text("Edelberg + Associates", props.Text{
				Top:   10,
				Align: consts.Left,
				Size:  10,
				Style: consts.Bold,
				Color: *primaryColor,
			})
			m.Text("1205 Johnson Ferry Rd.", props.Text{
				Top:   15,
				Align: consts.Left,
				Size:  10,
			})
			m.Text("Suite 136-356", props.Text{
				Top:   20,
				Align: consts.Left,
				Size:  10,
			})
			m.Text("Marietta, GA 30068, US", props.Text{
				Top:   25,
				Align: consts.Left,
				Size:  10,
			})
		})
		m.Col(3, func() {
			m.Text("Invoice Detail:", props.Text{
				Top:  4,
				Size: 12,
				Color: color.Color{
					Red:   124,
					Green: 124,
					Blue:  124,
				},
			})
			m.Text("Date:", props.Text{
				Top:  10,
				Size: 10,
			})
			m.Text("Payment Terms:", props.Text{
				Top:  15,
				Size: 10,
			})
			m.Text("Due Date:", props.Text{
				Top:  20,
				Size: 10,
			})
			m.Text("Balance Due:", props.Text{
				Top:   25,
				Size:  16,
				Style: consts.Bold,
				Color: *primaryColor,
			})
		})
		m.Col(2, func() {
			m.Text("", props.Text{
				Top: 4,
			})
			m.Text("Apr 25, 2023", props.Text{
				Top:   10,
				Align: consts.Right,
				Size:  10,
				Style: consts.Bold,
				Color: *primaryColor,
			})
			m.Text("1", props.Text{
				Top:   15,
				Align: consts.Right,
				Size:  10,
				Style: consts.Bold,
				Color: *primaryColor,
			})
			m.Text("Apr 28, 2023", props.Text{
				Top:   20,
				Align: consts.Right,
				Size:  10,
				Style: consts.Bold,
				Color: *primaryColor,
			})
			m.Text("$12,500", props.Text{
				Top:   25,
				Align: consts.Right,
				Size:  16,
				Style: consts.Bold,
				Color: *primaryColor,
			})
		})
	})
	m.Line(1.0,
		props.Line{
			Color: *primaryColor,
		},
	)
	m.Row(10, func() {

	})
}

package invoice

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/skip2/go-qrcode"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	_ "github.com/skip2/go-qrcode"

	"github.com/oarkflow/pkg/checksum"
	"github.com/oarkflow/pkg/decimal"
	"github.com/oarkflow/pkg/pdfs"
	"github.com/oarkflow/pkg/pdfs/color"
	"github.com/oarkflow/pkg/pdfs/consts"
	"github.com/oarkflow/pkg/pdfs/props"
	"github.com/oarkflow/pkg/str"
	"github.com/oarkflow/pkg/timeutil"
)

// FloatStr takes a float and gives back a monetary, human-formatted
// value.
var r = regexp.MustCompile("-?[0-9,]+.[0-9]{2}")
var primaryColor = color.Hex2RGB("012A4A")
var emptyRow = []string{"", "", "", "", ""}

func FloatStr(f float64) string {
	roundedFloat := math.Round(f*100) / 100
	p := message.NewPrinter(language.English)
	results := r.FindAllString(p.Sprintf("%f", roundedFloat), 1)

	if len(results) < 1 {
		panic("got some ridiculous number that has no decimals")
	}

	return results[0]
}

type Contact struct {
	Name      string `json:"name"`
	Address1  string `json:"address1"`
	Logo      string `json:"logo"`
	Address2  string `json:"address2"`
	City      string `json:"city"`
	State     string `json:"state"`
	ZipCode   string `json:"zipCode"`
	Country   string `json:"country"`
	Telephone string `json:"telephone"`
	Email     string `json:"email"`
}

type BankDetail struct {
	AccountName   string `json:"name"`
	AccountNumber string `json:"account_number"`
	BankName      string `json:"bank_name"`
	BankAddress   string `json:"bank_address"`
	SwiftCode     string `json:"swift_code"`
}

type Esewa struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Business struct {
	Details      Contact     `json:"details"`
	BankDetail   *BankDetail `json:"bank_detail"`
	Esewa        *Esewa      `json:"esewa"`
	PayPal       string      `json:"paypal"`
	ContactName  string      `json:"contact_name"`
	ContactEmail string      `json:"contact_email"`
	ContactPhone string      `json:"contact_phone"`
}

type Customer struct {
	Details Contact `json:"details"`
}

type Item struct {
	Description string  `yaml:"description" json:"description"`
	Currency    string  `yaml:"currency" json:"currency"`
	Quantity    float64 `yaml:"quantity" json:"quantity"`
	UnitPrice   float64 `yaml:"unit_price" json:"unit_price"`
}

type Transaction struct {
	Description   string  `yaml:"description" json:"description"`
	PaymentMethod string  `yaml:"payment_method" json:"payment_method"`
	Currency      string  `yaml:"currency" json:"currency"`
	Quantity      float64 `yaml:"quantity" json:"quantity"`
	UnitPrice     float64 `yaml:"unit_price" json:"unit_price"`
}

func (b *Item) Total() float64 {
	return decimal.NewFromFloat(b.UnitPrice).Mul(decimal.NewFromFloat(b.Quantity)).InexactFloat64()
}

func (b *Transaction) Total() float64 {
	return decimal.NewFromFloat(b.UnitPrice).Mul(decimal.NewFromFloat(b.Quantity)).InexactFloat64()
}

func (b *Transaction) Strings(index string) []string {
	return []string{
		index,
		b.PaymentMethod,
		b.Description,
		strconv.FormatFloat(b.Quantity, 'f', 2, 64),
		FloatStr(b.UnitPrice),
		FloatStr(b.Total()),
	}
}

func (b *Item) Strings(index string) []string {
	return []string{
		index,
		b.Description,
		strconv.FormatFloat(b.Quantity, 'f', 2, 64),
		FloatStr(b.UnitPrice),
		FloatStr(b.Total()),
	}
}

type Detail struct {
	Department        string        `yaml:"department" json:"department"`
	InvoiceNumber     string        `yaml:"invoice_number" json:"invoice_number"`
	InvoiceURL        string        `json:"invoice_url"`
	Currency          string        `yaml:"currency" json:"currency"`
	PaymentTerms      string        `yaml:"payment_terms" json:"payment_terms"`
	DueDate           string        `yaml:"due_date" json:"due_date"`
	DueDays           int           `json:"due_days"`
	Date              string        `yaml:"date" json:"date"`
	Status            string        `yaml:"status" json:"status"`
	GeneratedDate     string        `yaml:"generated_date" json:"generated_date"`
	SubTotalAmount    float64       `json:"sub_total_amount"`
	TransactionAmount float64       `yaml:"transaction_amount" json:"transaction_amount"`
	TotalAmount       float64       `yaml:"total_amount" json:"total_amount"`
	TaxRate           string        `json:"tax_rate"`
	Note              string        `json:"note"`
	TaxAmount         float64       `json:"tax_amount"`
	UseExactDate      bool          `yaml:"use_exact_date" json:"use_exact_date"`
	Customer          Customer      `yaml:"customer" json:"customer"`
	Items             []Item        `yaml:"items" json:"items"`
	Transactions      []Transaction `json:"transactions" yaml:"transactions"`
	items             [][]string
	transactions      [][]string
}

func (c *Detail) Filename() string {
	return fmt.Sprintf("Invoice-%s-%s.pdf", c.InvoiceNumber, c.Date)
}

func (c *Detail) Header() []string {
	return []string{"SN", "Item", "Qty", "Rate (" + c.Currency + ")", "Amount (" + c.Currency + ")"}
}

func (c *Detail) TransactionHeader() []string {
	return []string{"SN", "Payment Method", "Item", "Qty", "Rate (" + c.Currency + ")", "Amount (" + c.Currency + ")"}
}

func (c *Detail) Total() float64 {
	var total float64
	for _, item := range c.Items {
		total += item.Total()
	}
	c.SubTotalAmount = total
	if c.TaxRate != "" {
		totalDecimal := decimal.NewFromFloat(total)
		taxAmount := decimal.RequireFromString(c.TaxRate).Div(decimal.NewFromFloat(100)).Mul(totalDecimal)
		c.TaxAmount = taxAmount.InexactFloat64()
		total = totalDecimal.Add(taxAmount).InexactFloat64()
	}
	return total
}

func (c *Detail) TransactionTotal() float64 {
	var total float64
	for _, item := range c.Transactions {
		total += item.Total()
	}
	if c.TaxRate != "" {
		totalDecimal := decimal.NewFromFloat(total)
		taxAmount := decimal.RequireFromString(c.TaxRate).Div(decimal.NewFromFloat(100)).Mul(totalDecimal)
		c.TaxAmount = taxAmount.InexactFloat64()
		total = totalDecimal.Add(taxAmount).InexactFloat64()
	}
	return total
}

type Margin struct {
	Left  float64 `json:"left,omitempty"`
	Top   float64 `json:"top"`
	Right float64 `json:"right"`
}

type Config struct {
	Business    *Business          `json:"business"`
	Margin      *Margin            `json:"margin"`
	Orientation consts.Orientation `json:"orientation"`
	PaperSize   consts.PageSize    `json:"paper_size"`
	Secret      string             `json:"secret"`
}

type logo struct {
	ext         consts.Extension
	base64Image string
}

type Invoice struct {
	config *Config
	engine pdfs.IPdf
	logo   *logo
}

func New(config *Config) (*Invoice, error) {
	if config.Orientation == "" {
		config.Orientation = consts.Portrait
	}
	if config.PaperSize == "" {
		config.PaperSize = consts.A4
	}
	if config.Margin == nil {
		config.Margin = &Margin{
			Left:  10,
			Top:   15,
			Right: 10,
		}
	}
	m := pdfs.New(config.Orientation, config.PaperSize)
	m.SetPageMargins(config.Margin.Left, config.Margin.Top, config.Margin.Right)
	m.SetCreationDate(time.Now())
	m.SetAuthor(config.Business.Details.Name, true)
	m.SetCreator(config.Business.Details.Name, true)
	invoice := &Invoice{config: config, engine: m}
	err := invoice.getEncodedLogo()
	if err != nil {
		return nil, err
	}
	return invoice, nil
}

func (i *Invoice) Create(detail *Detail) *Invoice {
	if detail.Date == "" {
		detail.Date = time.Now().Format(time.DateOnly)
	}
	if detail.DueDate == "" && detail.DueDays == 0 {
		detail.DueDate = detail.Date
	} else if detail.DueDays > 0 {
		date, err := timeutil.ParseTime(detail.Date)
		if err == nil {
			detail.DueDate = date.Add(time.Duration(detail.DueDays) * 24 * time.Hour).Format(time.DateOnly)
		} else {
			detail.DueDate = detail.Date
		}
	} else {
		currentDate, err := timeutil.ParseTime(detail.Date)
		if err != nil {
			detail.DueDate = time.Now().Format(time.DateOnly)
		}
		date, _, err := timeutil.Parse(detail.DueDate, currentDate)
		if err != nil {
			detail.DueDate = time.Now().Format(time.DateOnly)
		} else {
			detail.DueDate = date.Format(time.DateOnly)
		}
	}
	var transactions [][]string
	for i, item := range detail.Transactions {
		transactions = append(transactions, item.Strings(strconv.Itoa(i+1)))
	}

	detail.transactions = transactions
	detail.TransactionAmount = detail.TransactionTotal()
	detail.TotalAmount = detail.Total() - detail.TransactionAmount
	var items [][]string
	for i, item := range detail.Items {
		items = append(items, item.Strings(strconv.Itoa(i+1)))
	}
	if len(items) <= 10 && len(transactions) == 0 {
		rowsToFill := 10 - len(items)
		for i := 0; i < rowsToFill; i++ {
			items = append(items, emptyRow)
		}
	}
	detail.items = items
	i.init(detail)
	return i
}

func (i *Invoice) String(detail *Detail) string {
	return str.FromByte(i.Byte(detail))
}

func (i *Invoice) Byte(detail *Detail) []byte {
	if i.config.Business.Esewa != nil {
		qrData := map[string]any{
			"amount": detail.TotalAmount,
		}
		qrData["eSewa_id"] = i.config.Business.Esewa.ID
		qrData["name"] = i.config.Business.Esewa.Name
		bt, _ := json.Marshal(qrData)
		return bt
	}
	if i.config.Business.BankDetail != nil {
		qrData := map[string]any{
			"amount": detail.TotalAmount,
			"number": detail.InvoiceNumber,
			"date":   detail.Date,
			"url":    detail.InvoiceURL,
			"from":   i.config.Business.Details.Name,
			"to":     detail.Customer.Details.Name,
			"items":  len(detail.items),
		}
		qrData["bankCode"] = i.config.Business.BankDetail.SwiftCode
		qrData["accountName"] = i.config.Business.BankDetail.AccountName
		qrData["accountNumber"] = i.config.Business.BankDetail.AccountNumber
		qrData["remarks"] = "Invoice #" + detail.InvoiceNumber
		bt, _ := json.Marshal(qrData)
		return bt
	}
	return nil
}

func (i *Invoice) RenderToFile(outFileName ...string) error {
	var file string
	if len(outFileName) > 0 {
		file = outFileName[0]
	} else {
		file = fmt.Sprintf("Invoice-%d-%s.pdf", 1, time.Now().Format(time.DateOnly))
	}
	return i.engine.OutputFileAndClose(file)
}

func (i *Invoice) Render() (bytes.Buffer, error) {
	return i.engine.Output()
}

func (i *Invoice) init(detail *Detail) {
	i.prepareHeader(detail)
	i.prepareFooter(detail)
	i.engine.SetAliasNbPages("{nb}")
	i.engine.SetFirstPageNb(1)
	i.prepareItemTable(detail)
}

func (i *Invoice) prepareItemTable(detail *Detail) {
	i.engine.TableList(detail.Header(), detail.items, props.TableList{
		ContentProp: props.TableListContent{
			GridSizes: []uint{1, 6, 1, 2, 2},
			Color:     *primaryColor,
			Align:     []consts.Align{consts.Left, consts.Left, consts.Right, consts.Right, consts.Right},
		},
		HeaderProp: props.TableListContent{
			GridSizes:       []uint{1, 6, 1, 2, 2},
			Align:           []consts.Align{consts.Left, consts.Left, consts.Right, consts.Right, consts.Right},
			Color:           *primaryColor,
			BackgroundColor: *color.Hex2RGB("E9F2FB"),
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
	})
	if len(detail.Transactions) > 0 {
		i.engine.Row(13, func() {
			i.engine.Col(2, func() {
				i.engine.Text("Transactions:", props.Text{
					Top:   8,
					Size:  10,
					Color: *primaryColor,
					Style: consts.Bold,
				})
			})
		})
		i.prepareTransactionTable(detail)
	}
	i.prepareSummary(detail)
}

func (i *Invoice) prepareTransactionTable(detail *Detail) {
	i.engine.TableList(detail.TransactionHeader(), detail.transactions, props.TableList{
		ContentProp: props.TableListContent{
			GridSizes: []uint{1, 2, 4, 1, 2, 2},
			Color:     *primaryColor,
			Align:     []consts.Align{consts.Left, consts.Left, consts.Left, consts.Right, consts.Right, consts.Right},
		},
		HeaderProp: props.TableListContent{
			GridSizes:       []uint{1, 2, 4, 1, 2, 2},
			Align:           []consts.Align{consts.Left, consts.Left, consts.Left, consts.Right, consts.Right, consts.Right},
			Color:           *primaryColor,
			BackgroundColor: *color.Hex2RGB("E9F2FB"),
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
	})
}

func (i *Invoice) prepareSummary(detail *Detail) {
	i.engine.Row(5, func() {
		i.engine.Col(2, func() {
			i.engine.Text("Notes/Remarks:", props.Text{
				Top:   5,
				Size:  10,
				Color: *primaryColor,
			})
		})
		i.engine.ColSpace(6)
		i.engine.Col(2, func() {
			i.engine.Text("Sub Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
		i.engine.Col(2, func() {
			i.engine.Text(detail.Currency+FloatStr(detail.SubTotalAmount), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
	i.engine.Row(5, func() {
		i.engine.Col(6, func() {
			i.engine.Text(detail.Note, props.Text{
				Top:   5,
				Size:  10,
				Style: consts.Bold,
				Color: *primaryColor,
			})
		})
		i.engine.ColSpace(2)
		i.engine.Col(2, func() {
			i.engine.Text(fmt.Sprintf("Tax (%s%%):", detail.TaxRate), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
		i.engine.Col(2, func() {
			i.engine.Text(detail.Currency+FloatStr(detail.TaxAmount), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
	if len(detail.Transactions) > 0 {
		i.engine.Row(5, func() {
			i.engine.ColSpace(8)
			i.engine.Col(2, func() {
				i.engine.Text("Last Transactions:", props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  10,
					Align: consts.Right,
					Color: *primaryColor,
				})
			})
			i.engine.Col(2, func() {
				i.engine.Text(fmt.Sprintf("(%s%s)", detail.Currency, FloatStr(detail.TransactionAmount)), props.Text{
					Top:   5,
					Style: consts.Bold,
					Size:  10,
					Align: consts.Right,
					Color: *primaryColor,
				})
			})
		})
	}
	i.engine.Row(5, func() {
		i.engine.ColSpace(8)
		i.engine.Col(2, func() {
			i.engine.Text("Total:", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  14,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
		i.engine.Col(2, func() {
			i.engine.Text(detail.Currency+FloatStr(detail.TotalAmount), props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  14,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
	i.engine.Row(15, func() {

	})
	i.engine.Row(5, func() {
		i.engine.ColSpace(8)
		i.engine.Col(4, func() {
			i.engine.Text("Thank you for your business with us!", props.Text{
				Top:   5,
				Style: consts.Bold,
				Size:  10,
				Align: consts.Right,
				Color: *primaryColor,
			})
		})
	})
}

var supportedImageTypes = []string{"image/jpeg", "image/jpg", "image/png", "image/x-png"}

func getLogo(logo string) (data []byte, ext string, err error) {
	if strings.HasPrefix(logo, "http") || strings.HasPrefix(logo, "https") {
		var r *http.Response
		r, err = http.Get(logo)
		if err != nil {
			return
		}
		defer r.Body.Close()

		fname := path.Base(logo)
		ext = strings.ReplaceAll(filepath.Ext(fname), ".", "")
		if !str.Contains(supportedImageTypes, r.Header.Get("Content-Type")) {
			err = errors.New("Unsupported image type: " + ext)
			return
		}
		data, err = io.ReadAll(r.Body)
	} else {
		data, err = os.ReadFile(logo)
		if err != nil {
			return
		}

		ext = strings.ReplaceAll(filepath.Ext(logo), ".", "")
		if !str.Contains(supportedImageTypes, mime.TypeByExtension("."+ext)) {
			err = errors.New("Unsupported local image type: " + ext)
			return
		}
	}
	return
}

func (i *Invoice) getEncodedLogo() error {
	byteSlices, ext, err := getLogo(i.config.Business.Details.Logo)
	if err != nil {
		return err
	}
	i.logo = &logo{
		ext:         consts.Extension(ext),
		base64Image: base64.StdEncoding.EncodeToString(byteSlices),
	}
	return nil
}

func (i *Invoice) prepareHeader(detail *Detail) {
	i.engine.RegisterHeader(func() {
		i.engine.Row(22, func() {
			if i.logo != nil {
				i.engine.Col(4, func() {
					_ = i.engine.Base64Image(i.logo.base64Image, i.logo.ext, props.Rect{
						Percent: 100,
					})
				})
				i.engine.Col(3, func() {
					i.engine.Text("Outstanding Amount", props.Text{
						Align: consts.Center,
						Size:  10,
						Top:   2,
						Color: *primaryColor,
					})
					i.engine.Text(detail.Currency+FloatStr(detail.TotalAmount), props.Text{
						Align: consts.Center,
						Size:  20,
						Top:   6,
						Style: consts.Bold,
						Color: *primaryColor,
					})
				})
				i.engine.Col(3, func() {
					i.engine.Text("Date: "+detail.Date+"  Terms:"+detail.PaymentTerms, props.Text{
						Align: consts.Right,
						Size:  10,
						Top:   2,
						Color: *primaryColor,
					})
					i.engine.Text("Invoice No. "+detail.InvoiceNumber, props.Text{
						Align: consts.Right,
						Size:  20,
						Top:   6,
						Style: consts.Bold,
						Color: *primaryColor,
					})
					i.engine.Text("Due Date: "+detail.DueDate, props.Text{
						Align: consts.Right,
						Size:  10,
						Top:   15,
						Color: *primaryColor,
					})
				})
			} else {
				i.engine.ColSpace(6)
			}

			if detail.InvoiceURL == "" {
				detail.InvoiceURL = "https://orgwareconstruct.com"
			}
			qrData := i.String(detail)
			png, _ := qrcode.Encode(qrData, qrcode.Medium, 1024)
			i.engine.Col(2, func() {
				i.engine.Base64Image(base64.StdEncoding.EncodeToString(png), "png", props.Rect{
					Percent: 100,
					Left:    10,
				})

			})
		})

		/*i.engine.Line(1.0,
			props.Line{
				Color: *primaryColor,
			},
		)*/
		i.prepareDetail(detail)
	})
}

func (i *Invoice) prepareBusinessDetail(detail Contact, align consts.Align) {
	i.engine.Text(detail.Name, props.Text{
		Top:   10,
		Align: align,
		Size:  10,
		Style: consts.Bold,
		Color: *primaryColor,
	})

	i.engine.Text(detail.Address1, props.Text{
		Size:  10,
		Top:   15,
		Align: align,
		Color: *primaryColor,
	})
	top := 20.0
	if detail.Address2 != "" {
		i.engine.Text(detail.Address2, props.Text{
			Size:  10,
			Align: align,
			Top:   top,
			Color: *primaryColor,
		})
		top += 5

	}

	i.engine.Text(fmt.Sprintf("%s, %s %s, %s", detail.City, detail.State, detail.ZipCode, detail.Country), props.Text{
		Size:  10,
		Align: align,
		Top:   top,
		Color: *primaryColor,
	})
	top += 5
	if detail.Telephone != "" {
		i.engine.Text("Tel: "+detail.Telephone, props.Text{
			Size:  10,
			Top:   top,
			Align: align,
			Color: *primaryColor,
		})
	}
	top += 5
	if detail.Email != "" {
		i.engine.Text(detail.Email, props.Text{
			Size:  10,
			Top:   top,
			Align: align,
			Color: *primaryColor,
		})
	}
}

func (i *Invoice) prepareDetail(detail *Detail) {
	businessDetail := i.config.Business.Details
	customerDetail := detail.Customer.Details
	i.engine.Row(37, func() {
		i.engine.Col(6, func() {
			i.engine.Text("From", props.Text{
				Top:   4,
				Align: consts.Left,
				Size:  12,
				Color: color.Color{
					Red:   124,
					Green: 124,
					Blue:  124,
				},
			})
			i.prepareBusinessDetail(businessDetail, consts.Left)
		})
		i.engine.Col(6, func() {
			i.engine.Text("To", props.Text{
				Top:   4,
				Align: consts.Right,
				Size:  12,
				Color: color.Color{
					Red:   124,
					Green: 124,
					Blue:  124,
				},
			})
			i.prepareBusinessDetail(customerDetail, consts.Right)
		})
	})
	/*i.engine.Line(1.0,
		props.Line{
			Color: *primaryColor,
		},
	)*/
	i.engine.Row(5, func() {})
}

func (i *Invoice) prepareFooter(detail *Detail) {
	businessDetails := i.config.Business.Details
	paypal := i.config.Business.PayPal
	bankDetail := i.config.Business.BankDetail
	contactName := i.config.Business.ContactName
	contactEmail := i.config.Business.ContactEmail
	contactPhone := i.config.Business.ContactPhone
	i.engine.RegisterFooter(func() {
		i.engine.Row(37, func() {
			i.engine.Col(4, func() {
				i.engine.Text("Payment Details:", props.Text{
					Top:   4,
					Align: consts.Left,
					Size:  12,
					Color: color.Color{
						Red:   124,
						Green: 124,
						Blue:  124,
					},
				})
				top := 10.0
				if paypal != "" {
					i.engine.Text("PAYPAL:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
				}
				if bankDetail != nil {
					i.engine.Text("BENEFICIARY NAME:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text("BENEFICIARY ACCOUNT NUMBER:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text("BANK NAME:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text("BANK ADDRESS:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text("BANK SWIFT CODE:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
				}
				if i.config.Business.Esewa != nil {
					i.engine.Text("ESEWA NAME:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text("ESEWA ID:", props.Text{
						Top:   top,
						Align: consts.Left,
						Size:  10,
						Color: *primaryColor,
					})
				}
			})

			i.engine.Col(4, func() {
				i.engine.Text("", props.Text{
					Top:   4,
					Align: consts.Left,
					Size:  12,
					Color: color.Color{
						Red:   124,
						Green: 124,
						Blue:  124,
					},
				})
				top := 10.0
				if paypal != "" {
					i.engine.Text(paypal, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
				}
				if bankDetail != nil {
					i.engine.Text(bankDetail.AccountName, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text(bankDetail.AccountNumber, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text(bankDetail.BankName, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text(bankDetail.BankAddress, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text(bankDetail.SwiftCode, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
				}
				if i.config.Business.Esewa != nil {
					i.engine.Text(i.config.Business.Esewa.ID, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
					top += 5
					i.engine.Text(i.config.Business.Esewa.Name, props.Text{
						Top:   top,
						Align: consts.Left,
						Style: consts.Bold,
						Size:  10,
						Color: *primaryColor,
					})
				}
			})
			i.engine.Col(4, func() {
				i.engine.Text("Contact Information:", props.Text{
					Top:   4,
					Align: consts.Right,
					Size:  12,
					Color: color.Color{
						Red:   124,
						Green: 124,
						Blue:  124,
					},
				})
				i.engine.Text(contactName, props.Text{
					Top:   10,
					Align: consts.Right,
					Size:  10,
					Color: *primaryColor,
				})
				i.engine.Text(contactEmail, props.Text{
					Top:   15,
					Align: consts.Right,
					Size:  10,
					Color: *primaryColor,
				})
				i.engine.Text(contactPhone, props.Text{
					Top:   20,
					Align: consts.Right,
					Size:  10,
					Color: *primaryColor,
				})
			})
		})
		/*i.engine.Line(1)*/
		i.engine.Row(3, func() {})
		i.engine.Row(6, func() {
			i.engine.Col(2, func() {
				if i.logo != nil {
					_ = i.engine.Base64Image(i.logo.base64Image, i.logo.ext, props.Rect{
						Percent: 100,
					})
				}
			})
			i.engine.Col(4, func() {
				i.engine.Text(businessDetails.Name, props.Text{
					Top:   0,
					Align: consts.Right,
					Size:  10,
					Style: consts.Bold,
					Color: *primaryColor,
				})
				i.engine.Text(businessDetails.Address1, props.Text{
					Top:   4,
					Align: consts.Right,
					Size:  8,
				})
			})
			i.engine.Col(4, func() {
				i.engine.Text("Tel: "+businessDetails.Telephone, props.Text{
					Top:   0,
					Align: consts.Right,
					Size:  8,
				})
				i.engine.Text(businessDetails.Email, props.Text{
					Top:   4,
					Align: consts.Right,
					Size:  8,
				})
			})
			i.engine.Col(2, func() {
				i.engine.Text(strconv.Itoa(i.engine.GetCurrentPage())+"/{nb}", props.Text{
					Align: consts.Right,
					Size:  8,
				})
			})
		})
		i.engine.Row(3, func() {
			i.engine.Col(12, func() {
				if i.config.Secret != "" {
					bt := i.Byte(detail)
					checksum.Default(str.ToByte(i.config.Secret))
					signature := checksum.Make(bt)
					i.engine.Text(signature, props.Text{
						Top:         4,
						Align:       consts.Center,
						Size:        1.5,
						Extrapolate: true,
					})
				}
			})
		})
	})
}

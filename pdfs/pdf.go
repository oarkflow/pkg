package pdfs

import (
	"bytes"
	"time"

	"github.com/oarkflow/pkg/gofpdf"
	"github.com/oarkflow/pkg/pdfs/color"
	"github.com/oarkflow/pkg/pdfs/consts"
	"github.com/oarkflow/pkg/pdfs/props"
)

const (
	defaultTopMargin   = 10
	defaultLeftMargin  = 10
	defaultRightMargin = 10
	defaultFontSize    = 16
)

// IPdf is the principal abstraction to create a PDF document.
type IPdf interface {
	// Grid System
	Row(height float64, closure func())
	Col(width uint, closure func())
	ColSpace(gridSize uint)
	// Registers
	RegisterHeader(closure func())
	RegisterFooter(closure func())
	// Outside Col/Row Components
	TableList(header []string, contents [][]string, prop ...props.TableList)
	Line(spaceHeight float64, prop ...props.Line)
	// Inside Col/Row Components
	Text(text string, prop ...props.Text)
	FileImage(filePathName string, prop ...props.Rect) (err error)
	Base64Image(base64 string, extension consts.Extension, prop ...props.Rect) (err error)
	Barcode(code string, prop ...props.Barcode) error
	QrCode(code string, prop ...props.Rect)
	DataMatrixCode(code string, prop ...props.Rect)
	Signature(label string, prop ...props.Font)

	// File System
	OutputFileAndClose(filePathName string) error
	Output() (bytes.Buffer, error)

	// Helpers
	AddPage()
	SetBorder(on bool)
	SetBackgroundColor(color color.Color)
	SetAliasNbPages(alias string)
	SetFirstPageNb(number int)
	GetBorder() bool
	GetPageSize() (width float64, height float64)
	GetCurrentPage() int
	GetCurrentOffset() float64
	SetPageMargins(left, top, right float64)
	GetPageMargins() (left float64, top float64, right float64, bottom float64)
	SetMaxGridSum(maxGridSum float64)

	// Fonts
	AddUTF8Font(familyStr string, styleStr consts.Style, fileStr string)
	AddUTF8FontFromBytes(familyStr string, styleStr consts.Style, utf8Bytes []byte)
	SetFontLocation(fontDirStr string)
	SetDefaultFontFamily(fontFamily string)
	GetDefaultFontFamily() string

	// Metadata
	SetCompression(compress bool)
	SetProtection(actionFlag byte, userPassStr, ownerPassStr string)
	SetAuthor(author string, isUTF8 bool)
	SetCreator(creator string, isUTF8 bool)
	SetSubject(subject string, isUTF8 bool)
	SetTitle(title string, isUTF8 bool)
	SetCreationDate(time time.Time)
}

// Pdf is the principal structure which implements IPdf abstraction.
type Pdf struct {
	// Gofpdf wrapper.
	Pdf Fpdf

	// Components.
	Math            Math
	Font            Font
	TextHelper      Text
	SignHelper      Signature
	Image           Image
	Code            Code
	TableListHelper TableList
	LineHelper      Line

	// Closures with IPdf Header and Footer logic.
	headerClosure func()
	footerClosure func()

	// Computed values.
	pageIndex                 int
	offsetY                   float64
	rowHeight                 float64
	firstPageNb               int
	xColOffset                float64
	colWidth                  float64
	footerHeight              float64
	headerFooterContextActive bool

	// Page configs.
	marginTop         float64
	maxGridSum        float64
	calculationMode   bool
	backgroundColor   color.Color
	debugMode         bool
	orientation       consts.Orientation
	pageSize          consts.PageSize
	defaultFontFamily string
}

// NewPdfCustomSize creates a IPdf instance returning a pointer to Pdf
// Receive an Orientation and a PageSize.
// Use if custom page size is needed. Otherwise use New() shorthand if using page sizes from consts.Pagesize.
// If using custom width and height, pageSize is just a string value for the format and takes no effect.
// Width and height inputs are measurements of the page in Portrait orientation.
func NewPdfCustomSize(orientation consts.Orientation, pageSize consts.PageSize, unitStr string, width, height float64) IPdf {
	fpdf := gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: string(orientation),
		UnitStr:        unitStr,
		SizeStr:        string(pageSize),
		Size: gofpdf.SizeType{
			Wd: width,
			Ht: height,
		},
		FontDirStr: "",
	})
	fpdf.SetMargins(defaultLeftMargin, defaultTopMargin, defaultRightMargin)

	math := NewMath(fpdf)
	font := NewFont(fpdf, defaultFontSize, consts.Arial, consts.Bold)
	text := NewText(fpdf, math, font)

	signature := NewSignature(fpdf, math, text)

	image := NewImage(fpdf, math)

	code := NewCode(fpdf, math)

	tableList := NewTableList(text, font, consts.DefaultMaxGridSum)

	lineHelper := NewLine(fpdf)

	pdf := &Pdf{
		Pdf:               fpdf,
		Math:              math,
		Font:              font,
		TextHelper:        text,
		SignHelper:        signature,
		Image:             image,
		Code:              code,
		TableListHelper:   tableList,
		LineHelper:        lineHelper,
		pageSize:          pageSize,
		orientation:       orientation,
		calculationMode:   false,
		backgroundColor:   color.NewWhite(),
		defaultFontFamily: consts.Arial,
		maxGridSum:        consts.DefaultMaxGridSum,
	}

	pdf.TableListHelper.BindGrid(pdf)

	pdf.Font.SetFamily(consts.Arial)
	pdf.Font.SetStyle(consts.Bold)
	pdf.Font.SetSize(defaultFontSize)
	pdf.debugMode = false

	pdf.Pdf.AddPage()

	return pdf
}

// New create a IPdf instance returning a pointer to Pdf
// Receive an Orientation and a PageSize.
// Shorthand when using a preset page size from consts.PageSize.
func New(orientation consts.Orientation, pageSize consts.PageSize) IPdf {
	return NewPdfCustomSize(orientation, pageSize, "mm", 0, 0)
}

// AddPage adds a new page in the PDF.
func (s *Pdf) AddPage() {
	_, pageHeight := s.Pdf.GetPageSize()
	_, top, _, bottom := s.Pdf.GetMargins()

	totalOffsetY := int(s.offsetY + s.footerHeight)
	maxOffsetPage := int(pageHeight - bottom - top)

	s.Row(float64(maxOffsetPage-totalOffsetY), func() {
		s.ColSpace(uint(s.maxGridSum))
	})
}

// AddUTF8FontFromBytes adds a custom UTF8 font from the provided bytes. familyStr is the name of the custom font
// registered in pdf. styleStr is the style of the font and fileStr is the path to the .ttf file.
func (s *Pdf) AddUTF8FontFromBytes(familyStr string, styleStr consts.Style, utf8Bytes []byte) {
	s.Pdf.AddUTF8FontFromBytes(familyStr, string(styleStr), utf8Bytes)
}

// SetMaxGridSum changes the max grid size of the page
func (s *Pdf) SetMaxGridSum(maxGridSum float64) {
	s.maxGridSum = maxGridSum
}

// RegisterHeader define a sequence of Rows, Lines ou TableLists
// which will be added in every new page.
func (s *Pdf) RegisterHeader(closure func()) {
	s.headerClosure = closure
}

// RegisterFooter define a sequence of Rows, Lines ou TableLists
// which will be added in every new page.
func (s *Pdf) RegisterFooter(closure func()) {
	s.footerClosure = closure

	// calculation mode execute all row flow but
	// only to calculate the sum of heights.
	s.calculationMode = true
	closure()
	s.calculationMode = false
}

// GetCurrentPage obtain the current page index
// this can be used inside a RegisterFooter/RegisterHeader
// to draw the current page, or to another purposes.
func (s *Pdf) GetCurrentPage() int {
	return s.pageIndex + s.firstPageNb
}

// GetCurrentOffset obtain the current offset in y axis.
func (s *Pdf) GetCurrentOffset() float64 {
	return s.offsetY
}

// SetPageMargins overrides default margins (10,10,10)
// the new page margin will affect all PDF pages.
func (s *Pdf) SetPageMargins(left, top, right float64) {
	if top > defaultTopMargin {
		s.marginTop = top - defaultTopMargin
	}

	s.Pdf.SetMargins(left, defaultTopMargin, right)
}

// GetPageMargins returns the set page margins. Comes in order of Left, Top, Right, Bottom
// Default page margins is left: 10, top: 10, right: 10.
func (s *Pdf) GetPageMargins() (left float64, top float64, right float64, bottom float64) {
	left, top, right, bottom = s.Pdf.GetMargins()
	top += s.marginTop

	return
}

// Signature add a space for a signature inside a cell,
// the space will have a line and a text below.
func (s *Pdf) Signature(label string, prop ...props.Font) {
	signProp := props.Font{
		Color: color.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
	}
	if len(prop) > 0 {
		signProp = prop[0]
	}

	signProp.MakeValid(s.defaultFontFamily)

	cell := Cell{
		X:      s.xColOffset,
		Y:      s.offsetY,
		Width:  s.colWidth,
		Height: s.rowHeight,
	}

	s.SignHelper.AddSpaceFor(label, cell, signProp.ToTextProp(consts.Center, 0.0, false, 0))
}

// TableList create a table with multiple rows and columns,
// so is not possible use this component inside a row or
// inside a column.
// Headers define the amount of columns from each row.
// Headers have bold style, and localized at the top of table.
// Contents are array of arrays. Each array is one line.
func (s *Pdf) TableList(header []string, contents [][]string, prop ...props.TableList) {
	s.TableListHelper.Create(header, contents, s.defaultFontFamily, prop...)
}

// SetBorder enable the draw of lines in every cell.
// Draw borders in all columns created.
func (s *Pdf) SetBorder(on bool) {
	s.debugMode = on
}

// SetBackgroundColor define the background color of the PDF.
// This method can be used to toggle background from rows.
func (s *Pdf) SetBackgroundColor(color color.Color) {
	s.backgroundColor = color
	s.Pdf.SetFillColor(s.backgroundColor.Red, s.backgroundColor.Green, s.backgroundColor.Blue)
}

// SetFirstPageNb define first page number
// Default: 0.
func (s *Pdf) SetFirstPageNb(number int) {
	s.firstPageNb = number
}

// SetAliasNbPages Defines an alias for the total number of pages.
// It will be substituted as the document is closed.
func (s *Pdf) SetAliasNbPages(alias string) {
	s.Pdf.AliasNbPages(alias)
}

// GetBorder return the actual border value.
func (s *Pdf) GetBorder() bool {
	return s.debugMode
}

// GetPageSize return the actual page size.
func (s *Pdf) GetPageSize() (width float64, height float64) {
	return s.Pdf.GetPageSize()
}

// Line draw a line from margin left to margin right
// in the current row.
func (s *Pdf) Line(spaceHeight float64, prop ...props.Line) {
	lineProp := props.Line{
		Color: color.NewBlack(),
	}
	if len(prop) > 0 {
		lineProp = prop[0]
	}
	lineProp.MakeValid(spaceHeight)

	s.Row(spaceHeight, func() {
		s.Col(0, func() {
			width, _ := s.Pdf.GetPageSize()
			left, top, right, _ := s.Pdf.GetMargins()

			const divisorToGetHalf = 2.0
			cell := Cell{
				X:      left,
				Y:      s.offsetY + top + (spaceHeight / divisorToGetHalf),
				Width:  width - right,
				Height: s.offsetY + top + (spaceHeight / divisorToGetHalf),
			}

			s.LineHelper.Draw(cell, lineProp)
		})
	})
}

// Row define a row and enable add columns inside the row.
// IPdf do not support recursive rows or rows inside columns.
func (s *Pdf) Row(height float64, closure func()) {
	// Used to calculate the height of the footer.
	if s.calculationMode {
		s.footerHeight += height
		return
	}

	_, pageHeight := s.Pdf.GetPageSize()
	_, top, _, bottom := s.Pdf.GetMargins()

	totalOffsetY := int(s.offsetY + height + s.footerHeight)
	maxOffsetPage := int(pageHeight - bottom - top)

	// Note: The headerFooterContextActive is needed to avoid recursive
	// calls without end, because footerClosure and headerClosure actually
	// have Row calls too.

	// If the new cell to be added pass the useful space counting the
	// height of the footer, add the footer.
	if totalOffsetY > maxOffsetPage {
		if !s.headerFooterContextActive {
			s.headerFooterContextActive = true
			s.footer()
			s.headerFooterContextActive = false
			s.offsetY = 0
			s.pageIndex++
		}
	}

	// If is a new page, add the header.
	if !s.headerFooterContextActive {
		if s.offsetY == 0 {
			s.headerFooterContextActive = true
			s.header()
			s.headerFooterContextActive = false
		}
	}

	s.rowHeight = height
	s.xColOffset = 0

	// This closure has the Cols to be executed.
	closure()

	s.offsetY += s.rowHeight
	s.Pdf.Ln(s.rowHeight)
}

// Col create a column inside a row and enable to add
// components inside. IPdf do not support recursive
// columns or rows inside columns.
func (s *Pdf) Col(width uint, closure func()) {
	if width == 0 {
		width = uint(s.maxGridSum)
	}

	percent := float64(width) / s.maxGridSum

	pageWidth, _ := s.Pdf.GetPageSize()
	left, _, right, _ := s.Pdf.GetMargins()
	widthPerCol := (pageWidth - right - left) * percent

	s.colWidth = widthPerCol
	s.createColSpace(widthPerCol)

	// This closure has the components to be executed.
	closure()

	s.xColOffset += s.colWidth
}

// ColSpace create an empty column inside a row.
func (s *Pdf) ColSpace(gridSize uint) {
	s.Col(gridSize, func() {})
}

// Text create a text inside a cell.
func (s *Pdf) Text(text string, prop ...props.Text) {
	textProp := props.Text{
		Color: color.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
	}

	if len(prop) > 0 {
		textProp = prop[0]
	}

	textProp.MakeValid(s.defaultFontFamily)

	if textProp.Top > s.rowHeight {
		textProp.Top = s.rowHeight
	}

	if textProp.Left > s.colWidth {
		textProp.Left = s.colWidth
	}

	if textProp.Right > s.colWidth {
		textProp.Right = s.colWidth
	}

	cellWidth := s.colWidth - textProp.Left - textProp.Right
	if cellWidth < 0 {
		cellWidth = 0
	}

	cell := Cell{
		X:      s.xColOffset + textProp.Left,
		Y:      s.offsetY + textProp.Top,
		Width:  cellWidth,
		Height: 0,
	}

	s.TextHelper.Add(text, cell, textProp)
}

// FileImage add an Image reading from disk inside a cell.
// Defining Image properties.
func (s *Pdf) FileImage(filePathName string, prop ...props.Rect) error {
	rectProp := props.Rect{}
	if len(prop) > 0 {
		rectProp = prop[0]
	}

	rectProp.MakeValid()

	cell := Cell{
		X:      s.xColOffset,
		Y:      s.offsetY + rectProp.Top,
		Width:  s.colWidth,
		Height: s.rowHeight,
	}

	return s.Image.AddFromFile(filePathName, cell, rectProp)
}

// Base64Image add an Image reading byte slices inside a cell.
// Defining Image properties.
func (s *Pdf) Base64Image(base64 string, extension consts.Extension, prop ...props.Rect) error {
	rectProp := props.Rect{}
	if len(prop) > 0 {
		rectProp = prop[0]
	}

	rectProp.MakeValid()

	cell := Cell{
		X:      s.xColOffset,
		Y:      s.offsetY + rectProp.Top,
		Width:  s.colWidth,
		Height: s.rowHeight,
	}

	return s.Image.AddFromBase64(base64, cell, rectProp, extension)
}

// Barcode create an barcode inside a cell.
func (s *Pdf) Barcode(code string, prop ...props.Barcode) (err error) {
	barcodeProp := props.Barcode{}
	if len(prop) > 0 {
		barcodeProp = prop[0]
	}

	barcodeProp.MakeValid()

	cell := Cell{
		X:      s.xColOffset,
		Y:      s.offsetY + barcodeProp.Top,
		Width:  s.colWidth,
		Height: s.rowHeight,
	}

	err = s.Code.AddBar(code, cell, barcodeProp)

	return
}

// DataMatrixCode creates an datamatrix code inside a cell.
func (s *Pdf) DataMatrixCode(code string, prop ...props.Rect) {
	rectProp := props.Rect{}
	if len(prop) > 0 {
		rectProp = prop[0]
	}
	rectProp.MakeValid()

	cell := Cell{
		X:      s.xColOffset,
		Y:      s.offsetY + rectProp.Top,
		Width:  s.colWidth,
		Height: s.rowHeight,
	}

	s.Code.AddDataMatrix(code, cell, rectProp)
}

// QrCode create a qrcode inside a cell.
func (s *Pdf) QrCode(code string, prop ...props.Rect) {
	rectProp := props.Rect{}
	if len(prop) > 0 {
		rectProp = prop[0]
	}

	rectProp.MakeValid()

	cell := Cell{
		X:      s.xColOffset,
		Y:      s.offsetY + rectProp.Top,
		Width:  s.colWidth,
		Height: s.rowHeight,
	}

	s.Code.AddQr(code, cell, rectProp)
}

// OutputFileAndClose save pdf in disk.
func (s *Pdf) OutputFileAndClose(filePathName string) (err error) {
	s.drawLastFooter()
	err = s.Pdf.OutputFileAndClose(filePathName)

	return
}

// Output extract PDF in byte slices.
func (s *Pdf) Output() (bytes.Buffer, error) {
	s.drawLastFooter()
	var buffer bytes.Buffer
	err := s.Pdf.Output(&buffer)
	return buffer, err
}

// AddUTF8Font add a custom utf8 font. familyStr is the name of the custom font registered in pdf.
// styleStr is the style of the font and fileStr is the path to the .ttf file.
func (s *Pdf) AddUTF8Font(familyStr string, styleStr consts.Style, fileStr string) {
	s.Pdf.AddUTF8Font(familyStr, string(styleStr), fileStr)
}

// SetFontLocation allows you to change the fonts lookup location.  fontDirStr is an absolute path where the fonts should be located.
func (s *Pdf) SetFontLocation(fontDirStr string) {
	s.Pdf.SetFontLocation(fontDirStr)
}

// SetCompression allows to set/unset compression for a page
// Compression is on by default.
func (s *Pdf) SetCompression(compress bool) {
	s.Pdf.SetCompression(compress)
}

// SetProtection define a password to open the pdf.
func (s *Pdf) SetProtection(actionFlag byte, userPassStr, ownerPassStr string) {
	s.Pdf.SetProtection(actionFlag, userPassStr, ownerPassStr)
}

// SetAuthor allows to set author name.
func (s *Pdf) SetAuthor(author string, isUTF8 bool) {
	s.Pdf.SetAuthor(author, isUTF8)
}

// SetCreator allows to set creator.
func (s *Pdf) SetCreator(creator string, isUTF8 bool) {
	s.Pdf.SetCreator(creator, isUTF8)
}

// SetSubject allows to set subject.
func (s *Pdf) SetSubject(subject string, isUTF8 bool) {
	s.Pdf.SetSubject(subject, isUTF8)
}

// SetTitle allows to set title.
func (s *Pdf) SetTitle(title string, isUTF8 bool) {
	s.Pdf.SetTitle(title, isUTF8)
}

// SetCreationDate allows to set creation date.
func (s *Pdf) SetCreationDate(time time.Time) {
	s.Pdf.SetCreationDate(time)
}

// SetDefaultFontFamily allows you to customize the default font. By default Arial is the original value.
func (s *Pdf) SetDefaultFontFamily(fontFamily string) {
	s.defaultFontFamily = fontFamily
}

// GetDefaultFontFamily allows you to get the current default font family.
func (s *Pdf) GetDefaultFontFamily() string {
	return s.defaultFontFamily
}

func (s *Pdf) createColSpace(actualWidthPerCol float64) {
	border := ""

	if s.debugMode {
		border = "1"
	}

	s.Pdf.CellFormat(actualWidthPerCol, s.rowHeight, "", border, 0, "C", !s.backgroundColor.IsWhite(), 0, "")
}

func (s *Pdf) drawLastFooter() {
	if s.footerClosure != nil {
		_, pageHeight := s.Pdf.GetPageSize()
		_, top, _, bottom := s.Pdf.GetMargins()

		if s.offsetY+s.footerHeight < pageHeight-bottom-top {
			totalOffsetY := int(s.offsetY + s.footerHeight)
			maxOffsetPage := int(pageHeight - bottom - top)

			s.Row(float64(maxOffsetPage-totalOffsetY), func() {
				s.ColSpace(12)
			})

			s.headerFooterContextActive = true
			s.footerClosure()
			s.headerFooterContextActive = false
		}
	}
}

func (s *Pdf) footer() {
	backgroundColor := s.backgroundColor
	s.SetBackgroundColor(color.NewWhite())

	_, pageHeight := s.Pdf.GetPageSize()
	_, top, _, bottom := s.Pdf.GetMargins()

	totalOffsetY := int(s.offsetY + s.footerHeight)
	maxOffsetPage := int(pageHeight - bottom - top)

	s.Row(float64(maxOffsetPage-totalOffsetY), func() {
		s.ColSpace(uint(s.maxGridSum))
	})

	if s.footerClosure != nil {
		s.footerClosure()
	}

	s.SetBackgroundColor(backgroundColor)
}

func (s *Pdf) header() {
	backgroundColor := s.backgroundColor
	s.SetBackgroundColor(color.NewWhite())

	s.Row(s.marginTop, func() {
		s.ColSpace(uint(s.maxGridSum))
	})

	if s.headerClosure != nil {
		s.headerClosure()
	}

	s.SetBackgroundColor(backgroundColor)
}

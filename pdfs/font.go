package pdfs

import (
	"github.com/oarkflow/pkg/pdfs/color"
	"github.com/oarkflow/pkg/pdfs/consts"
)

const (
	gofpdfFontScale1 = 72.0
	gofpdfFontScale2 = 25.4
)

// Font is the abstraction which deals of how to set font configurations.
type Font interface {
	SetFamily(family string)
	SetStyle(style consts.Style)
	SetSize(size float64)
	SetFont(family string, style consts.Style, size float64)
	GetFamily() string
	GetStyle() consts.Style
	GetSize() float64
	GetFont() (string, consts.Style, float64)
	GetScaleFactor() (scaleFactor float64)
	SetColor(color color.Color)
	GetColor() color.Color
	SetBackgroundColor(color color.Color)
	GetBackgroundColor() color.Color
}

type font struct {
	pdf             Fpdf
	size            float64
	family          string
	style           consts.Style
	scaleFactor     float64
	fontColor       color.Color
	backgroundColor color.Color
}

// NewFont create a Font.
func NewFont(pdf Fpdf, size float64, family string, style consts.Style) *font {
	return &font{
		pdf:             pdf,
		size:            size,
		family:          family,
		style:           style,
		scaleFactor:     gofpdfFontScale1 / gofpdfFontScale2, // Value defined inside gofpdf constructor,
		fontColor:       color.Color{Red: 0, Green: 0, Blue: 0},
		backgroundColor: color.Color{Red: 255, Green: 255, Blue: 255},
	}
}

// GetFamily return the currently Font family configured.
func (s *font) GetFamily() string {
	return s.family
}

// GetStyle return the currently Font style configured.
func (s *font) GetStyle() consts.Style {
	return s.style
}

// GetSize return the currently Font size configured.
func (s *font) GetSize() float64 {
	return s.size
}

// GetFont return all the currently Font properties configured.
func (s *font) GetFont() (string, consts.Style, float64) {
	return s.family, s.style, s.size
}

// SetFamily defines a new Font family.
func (s *font) SetFamily(family string) {
	s.family = family

	s.pdf.SetFont(s.family, string(s.style), s.size)
}

// SetStyle defines a new Font style.
func (s *font) SetStyle(style consts.Style) {
	s.style = style

	s.pdf.SetFontStyle(string(s.style))
}

// SetSize defines a new Font size.
func (s *font) SetSize(size float64) {
	s.size = size
	s.pdf.SetFontSize(s.size)
}

// SetFont defines all new Font properties.
func (s *font) SetFont(family string, style consts.Style, size float64) {
	s.family = family
	s.style = style
	s.size = size

	s.pdf.SetFont(s.family, string(s.style), s.size)
}

// GetScaleFactor retrieve the scale factor defined in the instantiation of gofpdf.
func (s *font) GetScaleFactor() (scaleFactor float64) {
	return s.scaleFactor
}

func (s *font) SetColor(color color.Color) {
	s.fontColor = color
	s.pdf.SetTextColor(color.Red, color.Green, color.Blue)
}

func (s *font) GetColor() color.Color {
	return s.fontColor
}

func (s *font) SetBackgroundColor(color color.Color) {
	s.backgroundColor = color
	s.pdf.SetFillColor(color.Red, color.Green, color.Blue)
}

func (s *font) GetBackgroundColor() color.Color {
	return s.backgroundColor
}

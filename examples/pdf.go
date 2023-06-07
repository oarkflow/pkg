package main

import (
	"strings"

	"github.com/oarkflow/pkg/gofpdf"
)

func loremList() []string {
	return []string{
		"Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		"Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		"Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
		"Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	}
}

func main() {
	const (
		colCount = 3
		colWd    = 60.0
		marginH  = 15.0
		lineHt   = 5.5
		cellGap  = 2.0
	)
	// var colStrList [colCount]string
	type cellType struct {
		str  string
		list [][]byte
		ht   float64
	}
	var (
		cellList [colCount]cellType
		cell     cellType
	)

	pdf := gofpdf.New("P", "mm", "A4", "") // 210 x 297
	header := [colCount]string{"Column A", "Column B", "Column C"}
	alignList := [colCount]string{"L", "C", "R"}
	strList := loremList()
	pdf.SetMargins(marginH, 15, marginH)
	pdf.SetFont("Arial", "", 14)
	pdf.AddPage()

	// Headers
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFillColor(64, 64, 64)
	for colJ := 0; colJ < colCount; colJ++ {
		pdf.CellFormat(colWd, 10, header[colJ], "", 0, "CM", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(24, 24, 24)
	pdf.SetFillColor(255, 255, 255)

	// Rows
	y := pdf.GetY()
	count := 0
	for rowJ := 0; rowJ < 2; rowJ++ {
		maxHt := lineHt
		// Cell height calculation loop
		for colJ := 0; colJ < colCount; colJ++ {
			count++
			if count > len(strList) {
				count = 1
			}
			cell.str = strings.Join(strList[0:count], " ")
			cell.list = pdf.SplitLines([]byte(cell.str), colWd-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList[colJ] = cell
		}
		// Cell render loop
		x := marginH
		for colJ := 0; colJ < colCount; colJ++ {
			// Apply border for each cell
			pdf.Rect(x, y, colWd, maxHt+cellGap+cellGap, "D")
			cell = cellList[colJ]
			cellY := y + cellGap + (maxHt-cell.ht)/2
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.CellFormat(colWd-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", 0, alignList[colJ], false, 0, "")
				cellY += lineHt
			}
			x += colWd
		}
		y += maxHt + cellGap + cellGap
	}

	pdf.OutputFileAndClose("test.pdf")

}

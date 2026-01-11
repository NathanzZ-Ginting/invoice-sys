package invoice

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// addItemsTable adds the invoice items table to the PDF
func addItemsTable(pdf *gofpdf.Fpdf, items []Item) {
	pdf.SetDrawColor(25, 103, 210)
	pdf.SetFillColor(25, 103, 210)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 10)

	// Table headers
	pdf.CellFormat(80, 8, "DESCRIPTION", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 8, "QTY", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "UNIT PRICE", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 8, "AMOUNT", "1", 0, "R", true, 0, "")
	pdf.Ln(-1)

	// Table rows
	pdf.SetTextColor(44, 62, 80)
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetFillColor(240, 245, 250)

	fill := false
	for _, item := range items {
		lineTotal := float64(item.Quantity) * item.UnitPrice

		pdf.CellFormat(80, 8, item.Description, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", item.UnitPrice), "1", 0, "R", fill, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", lineTotal), "1", 0, "R", fill, 0, "")
		pdf.Ln(-1)

		fill = !fill
	}
}

package invoice

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
)

// addTotalsSection adds the subtotal, tax, discount, and total to the PDF
func addTotalsSection(pdf *gofpdf.Fpdf, inv Invoice) {
	pdf.Ln(2)
	pdf.SetFont("Helvetica", "", 9)
	pdf.Cell(130, 6, "")

	// Calculate values
	subtotal := inv.Subtotal
	if subtotal == 0 {
		// Fallback calculation if subtotal not provided
		for _, item := range inv.Items {
			subtotal += float64(item.Quantity) * item.UnitPrice
		}
	}
	taxAmount := subtotal * (inv.Tax / 100)
	finalTotal := inv.Total
	if finalTotal == 0 {
		finalTotal = subtotal + taxAmount - inv.Discount
	}

	// Subtotal - Always display in USD for PDF compatibility
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(30, 6, "Subtotal:")
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(30, 6, fmt.Sprintf("$%.2f", subtotal))
	pdf.Ln(6)

	// Tax
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(130, 6, "")
	pdf.Cell(30, 6, fmt.Sprintf("Tax (%.1f%%):", inv.Tax))
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(30, 6, fmt.Sprintf("$%.2f", taxAmount))
	pdf.Ln(6)

	// Discount (only show if > 0)
	if inv.Discount > 0 {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.Cell(130, 6, "")
		pdf.Cell(30, 6, "Discount:")
		pdf.SetFont("Helvetica", "", 10)
		pdf.Cell(30, 6, fmt.Sprintf("-$%.2f", inv.Discount))
		pdf.Ln(6)
	}

	// Total
	pdf.SetFillColor(25, 103, 210)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 12)
	pdf.Cell(130, 8, "")
	pdf.CellFormat(30, 8, "TOTAL:", "0", 0, "R", true, 0, "")
	pdf.CellFormat(30, 8, fmt.Sprintf("$%.2f", finalTotal), "0", 0, "R", true, 0, "")
	pdf.Ln(12)
}

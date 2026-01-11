package invoice

import (
	"github.com/jung-kurt/gofpdf"
)

// addHeader adds a professional header to the PDF
func addHeader(pdf *gofpdf.Fpdf) {
	// Header background
	pdf.SetFillColor(25, 103, 210)
	pdf.Rect(0, 0, 210, 70, "F")

	// Company Logo placeholder (circle)
	pdf.SetFillColor(255, 255, 255)
	pdf.Circle(30, 25, 12, "F")

	// Company Name
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Helvetica", "B", 20)
	pdf.SetXY(48, 20)
	pdf.Cell(0, 10, "InvoicePro")

	// Tagline
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetXY(48, 30)
	pdf.Cell(0, 8, "Professional Invoicing Made Simple")

	// Contact info (right side)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetXY(140, 20)
	pdf.Cell(0, 5, "www.invoicepro.com")
	pdf.SetXY(140, 26)
	pdf.Cell(0, 5, "billing@invoicepro.com")
	pdf.SetXY(140, 32)
	pdf.Cell(0, 5, "+62 812 3456 7890")
}

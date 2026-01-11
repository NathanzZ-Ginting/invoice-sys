package pdf

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePDF creates a professional PDF invoice and returns it as bytes
func GeneratePDF(inv Invoice) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header section
	addHeader(pdf)

	// Invoice title and line
	pdf.SetDrawColor(25, 103, 210)
	pdf.SetLineWidth(0.5)
	pdf.Line(20, 80, 190, 80)
	pdf.Ln(5)

	pdf.SetFont("Helvetica", "B", 24)
	pdf.SetTextColor(25, 103, 210)
	pdf.Cell(0, 10, "INVOICE")
	pdf.Ln(15)

	// Invoice details
	addInvoiceDetails(pdf, inv)

	// From/To section
	addFromToSection(pdf, inv)

	// Items table - always uses USD in PDF
	addItemsTable(pdf, inv.Items)

	// Totals section - always uses USD in PDF
	addTotalsSection(pdf, inv)

	// Notes and footer
	addNotesAndFooter(pdf)

	// Generate PDF bytes
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// addInvoiceDetails adds invoice number and dates
func addInvoiceDetails(pdf *gofpdf.Fpdf, inv Invoice) {
	pdf.SetFont("Helvetica", "", 10)
	pdf.SetTextColor(44, 62, 80)

	// Invoice number
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(50, 6, "Invoice Number:")
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, inv.ID)
	pdf.Ln(6)

	// Issue date
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(50, 6, "Issue Date:")
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, time.Now().Format("02 January 2006"))
	pdf.Ln(6)

	// Due date
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(50, 6, "Due Date:")
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, time.Now().AddDate(0, 0, 30).Format("02 January 2006"))
	pdf.Ln(12)
}

// addFromToSection adds company and customer details
func addFromToSection(pdf *gofpdf.Fpdf, inv Invoice) {
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(25, 103, 210)
	pdf.Cell(95, 7, "FROM:")
	pdf.Cell(0, 7, "BILL TO:")
	pdf.Ln(8)

	// Company details (FROM)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(44, 62, 80)
	pdf.Cell(95, 6, "InvoicePro Systems")
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 100, 100)
	startY := pdf.GetY()

	pdf.SetX(20)
	pdf.MultiCell(90, 5, "Jl. Merdeka No. 123\nJakarta, Indonesia 12345\nPhone: +62 812 3456 7890\nEmail: billing@invoicepro.com", "", "L", false)

	// Customer details (BILL TO)
	customerY := startY
	pdf.SetXY(110, customerY)
	pdf.SetFont("Helvetica", "B", 11)
	pdf.SetTextColor(44, 62, 80)
	pdf.Cell(80, 6, inv.CustomerName)
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.SetX(110)

	addressLines := inv.CustomerAddress
	if inv.CustomerCity != "" {
		addressLines += "\n" + inv.CustomerCity
	}
	if inv.CustomerPostalCode != "" {
		addressLines += " " + inv.CustomerPostalCode
	}
	if inv.CustomerPhone != "" {
		addressLines += "\nPhone: " + inv.CustomerPhone
	}
	if inv.CustomerEmail != "" {
		addressLines += "\nEmail: " + inv.CustomerEmail
	}

	pdf.MultiCell(80, 5, addressLines, "", "L", false)
	pdf.Ln(8)
}

// addNotesAndFooter adds notes section and page footer
func addNotesAndFooter(pdf *gofpdf.Fpdf) {
	pdf.SetTextColor(44, 62, 80)
	pdf.SetFont("Helvetica", "B", 10)
	pdf.Cell(0, 6, "Notes & Terms:")
	pdf.Ln(6)

	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(100, 100, 100)
	pdf.MultiCell(0, 5, "Thank you for your business! Payment is due within 30 days of invoice date.\nPlease make payment to the bank details provided separately.\nFor inquiries, please contact us at billing@invoicepro.com", "", "L", false)

	// Bottom footer
	pdf.SetY(-30)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, pdf.GetY()-5, 190, pdf.GetY()-5)

	pdf.SetX(20)
	pdf.Cell(0, 5, "InvoicePro Systems | www.invoicepro.com | billing@invoicepro.com")
	pdf.Ln(5)
	pdf.Cell(0, 5, fmt.Sprintf("Generated on %s | Page %d", time.Now().Format("02-01-2006 15:04"), pdf.PageNo()))
}

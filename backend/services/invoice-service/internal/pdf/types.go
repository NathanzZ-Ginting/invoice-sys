package pdf

import (
	"time"
)

// Invoice represents invoice data for PDF generation
type Invoice struct {
	ID                 string
	CustomerName       string
	CustomerEmail      string
	CustomerAddress    string
	CustomerCity       string
	CustomerPostalCode string
	CustomerPhone      string
	Items              []Item
	Subtotal           float64
	Tax                float64 // Tax percentage
	Discount           float64 // Discount amount
	Total              float64
	Currency           string    // Currency code (USD, IDR, EUR, etc.)
	CreatedAt          time.Time
}

// Item represents a line item in an invoice
type Item struct {
	Description string
	Quantity    int
	UnitPrice   float64
}

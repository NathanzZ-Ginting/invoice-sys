package repository

import (
	"fmt"
	"time"

	"invoice-backend/services/shared/pkg/database"
	"invoice-backend/services/shared/pkg/types"
)

type InvoiceRepository struct {
	db *database.Client
}

func NewInvoiceRepository(db *database.Client) *InvoiceRepository {
	return &InvoiceRepository{db: db}
}

// GetAll returns all invoices with optional filters
func (r *InvoiceRepository) GetAll(status, currency string) ([]types.Invoice, error) {
	query := r.db.Supabase.From("invoices").
		Select("*, customers(name, email)", "", false).
		Order("created_at", nil)

	if status != "" && status != "all" {
		query = query.Eq("payment_status", status)
	}

	if currency != "" && currency != "all" {
		query = query.Eq("currency", currency)
	}

	var invoices []types.Invoice
	_, err := query.ExecuteTo(&invoices)
	return invoices, err
}

// GetByID returns an invoice by ID
func (r *InvoiceRepository) GetByID(id string) (*types.Invoice, error) {
	var invoices []types.Invoice
	_, err := r.db.Supabase.From("invoices").
		Select("*, customers(name, email)", "", false).
		Eq("id", id).
		ExecuteTo(&invoices)
	
	if err != nil {
		return nil, err
	}
	if len(invoices) == 0 {
		return nil, fmt.Errorf("invoice not found")
	}
	return &invoices[0], nil
}

// GetItems returns items for an invoice
func (r *InvoiceRepository) GetItems(invoiceID string) ([]types.Item, error) {
	var items []types.Item
	_, err := r.db.Supabase.From("invoice_items").
		Select("*", "", false).
		Eq("invoice_id", invoiceID).
		ExecuteTo(&items)
	return items, err
}

// Create creates a new invoice with items
func (r *InvoiceRepository) Create(customerID string, items []types.Item, tax, discount float64, status, notes, dueDate, currency string) (*types.Invoice, error) {
	// Calculate totals
	var subtotal float64
	for _, item := range items {
		item.Total = float64(item.Quantity) * item.UnitPrice
		subtotal += item.Total
	}

	// Calculate tax and total
	taxAmount := subtotal * (tax / 100)
	total := subtotal + taxAmount - discount

	// Generate invoice number
	invoiceNumber, err := r.generateInvoiceNumber()
	if err != nil {
		return nil, err
	}

	// Set defaults
	if status == "" {
		status = "pending"
	}
	if currency == "" {
		currency = "USD"
	}
	if dueDate == "" {
		dueDate = time.Now().AddDate(0, 0, 30).Format("2006-01-02")
	}

	// Create invoice
	invoiceData := map[string]interface{}{
		"customer_id":    customerID,
		"invoice_number": invoiceNumber,
		"due_date":       dueDate,
		"status":         status,
		"payment_status": "unpaid",
		"subtotal":       subtotal,
		"tax":            tax,
		"discount":       discount,
		"total":          total,
		"paid_amount":    0,
		"currency":       currency,
		"notes":          notes,
		"items":          items,
	}

	var result []types.Invoice
	_, err = r.db.Supabase.From("invoices").
		Insert(invoiceData, false, "", "", "").
		ExecuteTo(&result)
	
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no invoice created")
	}

	// Items are already stored as JSON in the invoices table
	// No need to insert into invoice_items table

	// Return the created invoice
	return &result[0], nil
}

// Update updates an invoice
func (r *InvoiceRepository) Update(id string, status, notes string) (*types.Invoice, error) {
	updateData := map[string]interface{}{
		"status": status,
		"notes":  notes,
	}

	var result []types.Invoice
	_, err := r.db.Supabase.From("invoices").
		Update(updateData, "", "").
		Eq("id", id).
		ExecuteTo(&result)
	
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("invoice not found")
	}
	return &result[0], nil
}

// Delete deletes an invoice
func (r *InvoiceRepository) Delete(id string) error {
	// Delete invoice items first
	_, _, err := r.db.Supabase.From("invoice_items").
		Delete("", "").
		Eq("invoice_id", id).
		Execute()
	
	if err != nil {
		return err
	}

	// Delete invoice
	_, _, err = r.db.Supabase.From("invoices").
		Delete("", "").
		Eq("id", id).
		Execute()
	
	return err
}

// generateInvoiceNumber generates a unique invoice number
func (r *InvoiceRepository) generateInvoiceNumber() (string, error) {
	year := time.Now().Year()
	
	// Get count of invoices this year
	var invoices []types.Invoice
	_, err := r.db.Supabase.From("invoices").
		Select("id", "", false).
		Like("invoice_number", fmt.Sprintf("INV-%d-%%", year)).
		ExecuteTo(&invoices)
	
	if err != nil {
		return "", err
	}

	nextNumber := len(invoices) + 1
	return fmt.Sprintf("INV-%d-%04d", year, nextNumber), nil
}

// GetCustomer returns a customer by ID
func (r *InvoiceRepository) GetCustomer(id string) (*types.Customer, error) {
	var customers []types.Customer
	_, err := r.db.Supabase.From("customers").
		Select("*", "", false).
		Eq("id", id).
		ExecuteTo(&customers)
	
	if err != nil {
		return nil, err
	}
	if len(customers) == 0 {
		return nil, fmt.Errorf("customer not found")
	}
	return &customers[0], nil
}

// GetCurrencyRates returns all currency exchange rates
func (r *InvoiceRepository) GetCurrencyRates() ([]types.CurrencyRate, error) {
	var rates []types.CurrencyRate
	_, err := r.db.Supabase.From("currency_rates").
		Select("*", "", false).
		Order("from_currency", nil).
		ExecuteTo(&rates)
	return rates, err
}

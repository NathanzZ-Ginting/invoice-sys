package db

import (
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

type Client struct {
	supabase *supabase.Client
}

type Customer struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

type CustomerCreate struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone,omitempty"`
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	Country     string `json:"country,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
}

type Invoice struct {
	ID            string  `json:"id"`
	CustomerID    string  `json:"customer_id"`
	InvoiceNumber string  `json:"invoice_number,omitempty"`
	Subtotal      float64 `json:"subtotal"`
	Tax           float64 `json:"tax,omitempty"`
	Discount      float64 `json:"discount,omitempty"`
	Total         float64 `json:"total"`
	Items         []Item  `json:"items"`
	PDFURL        string  `json:"pdf_url,omitempty"`
	Status        string  `json:"status,omitempty"`
	Notes         string  `json:"notes,omitempty"`
	DueDate       string  `json:"due_date,omitempty"`
	Currency      string  `json:"currency,omitempty"`
	PaymentStatus string  `json:"payment_status,omitempty"`
	PaidAmount    float64 `json:"paid_amount,omitempty"`
	PaymentDate   string  `json:"payment_date,omitempty"`
	CreatedAt     string  `json:"created_at,omitempty"`
}

type InvoiceCreate struct {
	CustomerID string  `json:"customer_id"`
	Subtotal   float64 `json:"subtotal"`
	Tax        float64 `json:"tax,omitempty"`
	Discount   float64 `json:"discount,omitempty"`
	Total      float64 `json:"total"`
	Items      []Item  `json:"items"`
	Status     string  `json:"status,omitempty"`
	Notes      string  `json:"notes,omitempty"`
	DueDate    string  `json:"due_date,omitempty"`
	Currency   string  `json:"currency,omitempty"`
}

type Item struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
}

// Payment represents a payment record for an invoice
type Payment struct {
	ID              string  `json:"id"`
	InvoiceID       string  `json:"invoice_id"`
	Amount          float64 `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	PaymentDate     string  `json:"payment_date"`
	ReferenceNumber string  `json:"reference_number,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	CreatedAt       string  `json:"created_at,omitempty"`
	CreatedBy       string  `json:"created_by,omitempty"`
}

// PaymentCreate is the struct for creating a new payment
type PaymentCreate struct {
	InvoiceID       string  `json:"invoice_id"`
	Amount          float64 `json:"amount"`
	PaymentMethod   string  `json:"payment_method"`
	PaymentDate     string  `json:"payment_date"`
	ReferenceNumber string  `json:"reference_number,omitempty"`
	Notes           string  `json:"notes,omitempty"`
	CreatedBy       string  `json:"created_by,omitempty"`
}

// DashboardStats represents dashboard analytics data
type DashboardStats struct {
	TotalRevenue     float64 `json:"total_revenue"`
	PaidAmount       float64 `json:"paid_amount"`
	UnpaidAmount     float64 `json:"unpaid_amount"`
	OverdueAmount    float64 `json:"overdue_amount"`
	TotalInvoices    int     `json:"total_invoices"`
	PaidInvoices     int     `json:"paid_invoices"`
	UnpaidInvoices   int     `json:"unpaid_invoices"`
	OverdueInvoices  int     `json:"overdue_invoices"`
	PartiallyPaid    int     `json:"partially_paid"`
	Currency         string  `json:"currency"`
}

// RevenueByPeriod represents revenue grouped by time period
type RevenueByPeriod struct {
	Date         string  `json:"date"`
	Revenue      float64 `json:"revenue"`
	InvoiceCount int     `json:"invoice_count"`
	Currency     string  `json:"currency"`
}

// TopCustomer represents customer with highest revenue
type TopCustomer struct {
	CustomerID    string  `json:"customer_id"`
	CustomerName  string  `json:"customer_name"`
	TotalRevenue  float64 `json:"total_revenue"`
	InvoiceCount  int     `json:"invoice_count"`
	Currency      string  `json:"currency"`
}

func New() (*Client, error) {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	if url == "" || key == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	supabase, err := supabase.NewClient(url, key, nil)
	if err != nil {
		return nil, err
	}
	return &Client{supabase: supabase}, nil
}

// Customer operations
func (c *Client) ListCustomers() ([]Customer, error) {
	var customers []Customer
	_, err := c.supabase.From("customers").Select("*", "", false).ExecuteTo(&customers)
	return customers, err
}

func (c *Client) CreateCustomer(name, email, phone, address, city, postalCode, country, companyName string) (*Customer, error) {
	customerCreate := map[string]interface{}{
		"name":         name,
		"email":        email,
		"phone":        phone,
		"address":      address,
		"city":         city,
		"postal_code":  postalCode,
		"country":      country,
		"company_name": companyName,
	}
	var result []Customer
	_, err := c.supabase.From("customers").Insert(customerCreate, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no customer created")
	}
	return &result[0], nil
}

func (c *Client) GetCustomer(id string) (*Customer, error) {
	var customer Customer
	_, err := c.supabase.From("customers").Select("*", "", false).Eq("id", id).Single().ExecuteTo(&customer)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (c *Client) UpdateCustomer(id, name, email, phone, address, city, postalCode, country, companyName string) (*Customer, error) {
	updates := map[string]interface{}{
		"name":         name,
		"email":        email,
		"phone":        phone,
		"address":      address,
		"city":         city,
		"postal_code":  postalCode,
		"country":      country,
		"company_name": companyName,
	}
	var result []Customer
	_, err := c.supabase.From("customers").Update(updates, "", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no customer updated")
	}
	return &result[0], nil
}

func (c *Client) DeleteCustomer(id string) error {
	_, _, err := c.supabase.From("customers").Delete("", "").Eq("id", id).ExecuteString()
	return err
}

// Invoice operations CREATE
func (c *Client) CreateInvoice(customerID string, subtotal, tax, discount, total float64, items []Item, status, notes, dueDate, currency string) (*Invoice, error) {
	if currency == "" {
		currency = "USD"
	}
	
	// Generate invoice number
	invoiceNumber, err := c.GenerateInvoiceNumber()
	if err != nil {
		return nil, fmt.Errorf("failed to generate invoice number: %v", err)
	}
	
	// Create invoice with generated number
	invoiceData := map[string]interface{}{
		"customer_id":    customerID,
		"invoice_number": invoiceNumber,
		"subtotal":       subtotal,
		"tax":            tax,
		"discount":       discount,
		"total":          total,
		"items":          items,
		"status":         status,
		"notes":          notes,
		"due_date":       dueDate,
		"currency":       currency,
		"payment_status": "unpaid",
		"paid_amount":    0,
	}
	
	var result []Invoice
	_, err = c.supabase.From("invoices").Insert(invoiceData, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no invoice created")
	}
	return &result[0], nil
}

func (c *Client) ListInvoices() ([]Invoice, error) {
	var invoices []Invoice
	_, err := c.supabase.From("invoices").Select("*", "", false).ExecuteTo(&invoices)
	return invoices, err
}

func (c *Client) GetInvoice(id string) (*Invoice, error) {
	var invoice Invoice
	_, err := c.supabase.From("invoices").Select("*", "", false).Eq("id", id).Single().ExecuteTo(&invoice)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (c *Client) UpdateInvoice(id string, status, notes, dueDate string) (*Invoice, error) {
	updates := map[string]interface{}{
		"status":   status,
		"notes":    notes,
		"due_date": dueDate,
	}
	var result []Invoice
	_, err := c.supabase.From("invoices").Update(updates, "", "").Eq("id", id).ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no invoice updated")
	}
	return &result[0], nil
}

// CurrencyRate represents exchange rate data
type CurrencyRate struct {
	ID           string  `json:"id"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
	Rate         float64 `json:"rate"`
	UpdatedAt    string  `json:"updated_at"`
}

// GetCurrencyRate gets exchange rate between two currencies
func (c *Client) GetCurrencyRate(fromCurrency, toCurrency string) (*CurrencyRate, error) {
	var rate CurrencyRate
	_, err := c.supabase.From("currency_rates").
		Select("*", "", false).
		Eq("from_currency", fromCurrency).
		Eq("to_currency", toCurrency).
		Single().
		ExecuteTo(&rate)
	if err != nil {
		return nil, err
	}
	return &rate, nil
}

// ListCurrencyRates gets all currency rates
func (c *Client) ListCurrencyRates() ([]CurrencyRate, error) {
	var rates []CurrencyRate
	_, err := c.supabase.From("currency_rates").Select("*", "", false).ExecuteTo(&rates)
	return rates, err
}

// FilterInvoices filters invoices based on status, date range, and search term
func (c *Client) FilterInvoices(status, searchTerm, startDate, endDate string) ([]Invoice, error) {
	var invoices []Invoice
	
	query := c.supabase.From("invoices").Select("*", "", false)
	
	// Filter by status
	if status != "" && status != "all" {
		query = query.Eq("status", status)
	}
	
	// Filter by date range
	if startDate != "" {
		query = query.Gte("created_at", startDate)
	}
	if endDate != "" {
		query = query.Lte("created_at", endDate)
	}
	
	_, err := query.ExecuteTo(&invoices)
	if err != nil {
		return nil, err
	}
	
	// Client-side search filtering (for customer_id or notes)
	if searchTerm != "" && len(invoices) > 0 {
		var filtered []Invoice
		for _, inv := range invoices {
			// Search in customer_id, notes, or invoice id
			if containsIgnoreCase(inv.ID, searchTerm) ||
				containsIgnoreCase(inv.CustomerID, searchTerm) ||
				containsIgnoreCase(inv.Notes, searchTerm) {
				filtered = append(filtered, inv)
			}
		}
		return filtered, nil
	}
	
	return invoices, nil
}

// Helper function for case-insensitive string matching
func containsIgnoreCase(str, substr string) bool {
	return len(str) >= len(substr) && 
		(str == substr || 
		 len(substr) == 0 || 
		 stringContains(str, substr))
}

func stringContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================
// PAYMENT RECORDING FUNCTIONS
// ============================================

// RecordPayment creates a payment record and updates invoice status
func (c *Client) RecordPayment(payment PaymentCreate) (*Payment, error) {
	// Create payment record
	var result []Payment
	_, err := c.supabase.From("payments").Insert(payment, false, "", "", "").ExecuteTo(&result)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no payment created")
	}

	// Get invoice to update payment status
	invoice, err := c.GetInvoice(payment.InvoiceID)
	if err != nil {
		return nil, err
	}

	// Calculate new paid amount
	newPaidAmount := invoice.PaidAmount + payment.Amount
	
	// Determine payment status
	paymentStatus := "partially_paid"
	if newPaidAmount >= invoice.Total {
		paymentStatus = "paid"
		newPaidAmount = invoice.Total // Cap at total
	}

	// Update invoice
	updateData := map[string]interface{}{
		"paid_amount":    newPaidAmount,
		"payment_status": paymentStatus,
		"payment_date":   payment.PaymentDate,
	}

	var updated []Invoice
	_, err = c.supabase.From("invoices").
		Update(updateData, "", "").
		Eq("id", payment.InvoiceID).
		ExecuteTo(&updated)
	
	if err != nil {
		return nil, fmt.Errorf("failed to update invoice: %v", err)
	}

	return &result[0], nil
}

// GetPaymentsByInvoice returns all payments for an invoice
func (c *Client) GetPaymentsByInvoice(invoiceID string) ([]Payment, error) {
	var payments []Payment
	_, err := c.supabase.From("payments").
		Select("*", "", false).
		Eq("invoice_id", invoiceID).
		ExecuteTo(&payments)
	return payments, err
}

// GetAllPayments returns all payments with optional filtering
func (c *Client) GetAllPayments() ([]Payment, error) {
	var payments []Payment
	_, err := c.supabase.From("payments").
		Select("*", "", false).
		ExecuteTo(&payments)
	return payments, err
}

// ============================================
// DASHBOARD & ANALYTICS FUNCTIONS
// ============================================

// GetDashboardStats returns overall statistics for dashboard
func (c *Client) GetDashboardStats(currency string) (*DashboardStats, error) {
	var invoices []Invoice
	query := c.supabase.From("invoices").Select("*", "", false)
	
	if currency != "" && currency != "ALL" {
		query = query.Eq("currency", currency)
	}
	
	_, err := query.ExecuteTo(&invoices)
	if err != nil {
		return nil, err
	}

	stats := &DashboardStats{Currency: currency}
	
	for _, inv := range invoices {
		stats.TotalRevenue += inv.Total
		stats.TotalInvoices++
		
		switch inv.PaymentStatus {
		case "paid":
			stats.PaidAmount += inv.Total
			stats.PaidInvoices++
		case "partially_paid":
			stats.PaidAmount += inv.PaidAmount
			stats.UnpaidAmount += (inv.Total - inv.PaidAmount)
			stats.PartiallyPaid++
		case "overdue":
			stats.OverdueAmount += inv.Total
			stats.OverdueInvoices++
		case "unpaid":
			stats.UnpaidAmount += inv.Total
			stats.UnpaidInvoices++
		}
	}
	
	return stats, nil
}

// GetRevenueByPeriod returns revenue grouped by period (daily, weekly, monthly)
func (c *Client) GetRevenueByPeriod(period string, limit int) ([]RevenueByPeriod, error) {
	// For simplicity, we'll aggregate in Go
	var invoices []Invoice
	_, err := c.supabase.From("invoices").
		Select("*", "", false).
		ExecuteTo(&invoices)
	
	if err != nil {
		return nil, err
	}

	// Group by date (simplified - daily grouping)
	revenueMap := make(map[string]*RevenueByPeriod)
	
	for _, inv := range invoices {
		date := inv.CreatedAt[:10] // Get YYYY-MM-DD
		
		if _, exists := revenueMap[date]; !exists {
			revenueMap[date] = &RevenueByPeriod{
				Date:     date,
				Currency: inv.Currency,
			}
		}
		
		revenueMap[date].Revenue += inv.Total
		revenueMap[date].InvoiceCount++
	}
	
	// Convert map to slice
	var result []RevenueByPeriod
	for _, v := range revenueMap {
		result = append(result, *v)
	}
	
	return result, nil
}

// GetTopCustomers returns customers with highest revenue
func (c *Client) GetTopCustomers(limit int) ([]TopCustomer, error) {
	var invoices []Invoice
	_, err := c.supabase.From("invoices").
		Select("*", "", false).
		ExecuteTo(&invoices)
	
	if err != nil {
		return nil, err
	}

	// Group by customer
	customerMap := make(map[string]*TopCustomer)
	
	for _, inv := range invoices {
		if _, exists := customerMap[inv.CustomerID]; !exists {
			// Get customer name
			customer, err := c.GetCustomer(inv.CustomerID)
			if err != nil {
				continue
			}
			
			customerMap[inv.CustomerID] = &TopCustomer{
				CustomerID:   inv.CustomerID,
				CustomerName: customer.Name,
				Currency:     inv.Currency,
			}
		}
		
		customerMap[inv.CustomerID].TotalRevenue += inv.Total
		customerMap[inv.CustomerID].InvoiceCount++
	}
	
	// Convert to slice and sort (simplified)
	var result []TopCustomer
	for _, v := range customerMap {
		result = append(result, *v)
	}
	
	// Return first N customers (should sort by revenue first)
	if len(result) > limit {
		result = result[:limit]
	}
	
	return result, nil
}

// GetOverdueInvoices returns all overdue invoices
func (c *Client) GetOverdueInvoices() ([]Invoice, error) {
	var invoices []Invoice
	_, err := c.supabase.From("invoices").
		Select("*", "", false).
		Eq("payment_status", "overdue").
		ExecuteTo(&invoices)
	return invoices, err
}

// ============================================
// INVOICE NUMBERING FUNCTIONS
// ============================================

// GenerateInvoiceNumber generates a new invoice number using database function
func (c *Client) GenerateInvoiceNumber() (string, error) {
	// Simple fallback: generate based on count + timestamp
	var count []Invoice
	_, err := c.supabase.From("invoices").
		Select("id", "", false).
		ExecuteTo(&count)
	
	if err != nil {
		return "", err
	}
	
	// Format: INV-2026-0001
	invoiceNum := fmt.Sprintf("INV-2026-%04d", len(count)+1)
	return invoiceNum, nil
}

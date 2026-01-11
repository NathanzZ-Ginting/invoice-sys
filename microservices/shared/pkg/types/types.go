package types

// Invoice represents an invoice record
type Invoice struct {
	ID              string  `json:"id"`
	CustomerID      string  `json:"customer_id"`
	InvoiceNumber   string  `json:"invoice_number"`
	Date            string  `json:"date"`
	DueDate         string  `json:"due_date"`
	Status          string  `json:"status"`
	PaymentStatus   string  `json:"payment_status"`
	Total           float64 `json:"total"`
	PaidAmount      float64 `json:"paid_amount"`
	Currency        string  `json:"currency"`
	Notes           string  `json:"notes,omitempty"`
	Items           []Item  `json:"items,omitempty"`
	CustomerName    string  `json:"customer_name,omitempty"`
	CustomerEmail   string  `json:"customer_email,omitempty"`
	PaymentDate     string  `json:"payment_date,omitempty"`
	CreatedAt       string  `json:"created_at,omitempty"`
	CreatedBy       string  `json:"created_by,omitempty"`
}

// Item represents an invoice line item
type Item struct {
	Description string  `json:"description"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Total       float64 `json:"total"`
}

// Customer represents a customer record
type Customer struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Address   string `json:"address,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
}

// Payment represents a payment record
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
	TotalRevenue    float64 `json:"total_revenue"`
	PaidAmount      float64 `json:"paid_amount"`
	UnpaidAmount    float64 `json:"unpaid_amount"`
	OverdueAmount   float64 `json:"overdue_amount"`
	TotalInvoices   int     `json:"total_invoices"`
	PaidInvoices    int     `json:"paid_invoices"`
	UnpaidInvoices  int     `json:"unpaid_invoices"`
	OverdueInvoices int     `json:"overdue_invoices"`
}

// RevenueData represents revenue data for charts
type RevenueData struct {
	Period string  `json:"period"`
	Amount float64 `json:"amount"`
}

// TopCustomer represents top customer by revenue
type TopCustomer struct {
	CustomerID    string  `json:"customer_id"`
	CustomerName  string  `json:"customer_name"`
	CustomerEmail string  `json:"customer_email"`
	TotalRevenue  float64 `json:"total_revenue"`
	InvoiceCount  int     `json:"invoice_count"`
}

// APIResponse is a generic API response wrapper
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// ServiceConfig holds common service configuration
type ServiceConfig struct {
	ServiceName  string
	Port         string
	DatabaseURL  string
	SupabaseURL  string
	SupabaseKey  string
}

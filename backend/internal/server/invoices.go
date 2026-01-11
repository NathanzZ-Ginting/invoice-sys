package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"invoice-backend/internal/db"
	"invoice-backend/internal/invoice"

	"github.com/gorilla/mux"
)

// listInvoices handles GET /invoices
func (s *Server) listInvoices(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	invoices, err := s.db.ListInvoices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

// createInvoice handles POST /invoices
func (s *Server) createInvoice(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	var req struct {
		CustomerID string    `json:"customer_id"`
		Items      []db.Item `json:"items"`
		Tax        float64   `json:"tax,omitempty"`
		Discount   float64   `json:"discount,omitempty"`
		Status     string    `json:"status,omitempty"`
		Notes      string    `json:"notes,omitempty"`
		DueDate    string    `json:"due_date,omitempty"`
		Currency   string    `json:"currency,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.CustomerID == "" || len(req.Items) == 0 {
		http.Error(w, "customer_id and items required", http.StatusBadRequest)
		return
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "pending"
	}

	// Set default currency if not provided
	if req.Currency == "" {
		req.Currency = "USD"
	}

	// Calculate subtotal
	var subtotal float64
	for _, item := range req.Items {
		subtotal += float64(item.Quantity) * item.UnitPrice
	}

	// Calculate total with tax and discount
	taxAmount := subtotal * (req.Tax / 100)
	total := subtotal + taxAmount - req.Discount

	// Get customer for PDF
	customer, err := s.db.GetCustomer(req.CustomerID)
	if err != nil {
		http.Error(w, "Customer not found", http.StatusBadRequest)
		return
	}

	// Create invoice record in database
	invRecord, err := s.db.CreateInvoice(req.CustomerID, subtotal, req.Tax, req.Discount, total, req.Items, req.Status, req.Notes, req.DueDate, req.Currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate PDF
	pdfBytes, err := generateInvoicePDF(invRecord, customer, req.Items, subtotal, req.Tax, req.Discount, total)
	if err != nil {
		log.Printf("Failed to generate PDF: %v", err)
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Return invoice with PDF data
	response := map[string]interface{}{
		"id":          invRecord.ID,
		"customer_id": invRecord.CustomerID,
		"total":       invRecord.Total,
		"items":       invRecord.Items,
		"currency":    invRecord.Currency,
		"pdf_data":    pdfBytes,
		"created_at":  invRecord.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// getInvoice handles GET /invoices/{id}
func (s *Server) getInvoice(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	invoice, err := s.db.GetInvoice(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

// updateInvoice handles PUT /invoices/{id}
func (s *Server) updateInvoice(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Status  string `json:"status"`
		Notes   string `json:"notes,omitempty"`
		DueDate string `json:"due_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Status == "" {
		http.Error(w, "status is required", http.StatusBadRequest)
		return
	}

	invoice, err := s.db.UpdateInvoice(id, req.Status, req.Notes, req.DueDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoice)
}

// generateInvoicePDF is a helper function to generate PDF from invoice data
func generateInvoicePDF(invRecord *db.Invoice, customer *db.Customer, items []db.Item, subtotal, tax, discount, total float64) ([]byte, error) {
	invItems := make([]invoice.Item, len(items))
	for i, item := range items {
		invItems[i] = invoice.Item{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}

	inv := invoice.Invoice{
		ID:                 invRecord.ID,
		CustomerName:       customer.Name,
		CustomerEmail:      customer.Email,
		CustomerAddress:    customer.Address,
		CustomerCity:       customer.City,
		CustomerPostalCode: customer.PostalCode,
		CustomerPhone:      customer.Phone,
		Items:              invItems,
		Subtotal:           subtotal,
		Tax:                tax,
		Discount:           discount,
		Total:              total,
		Currency:           invRecord.Currency, // Add currency from invoice record
	}

	return invoice.GeneratePDF(inv)
}

// filterInvoices handles GET /invoices/filter with query parameters
func (s *Server) filterInvoices(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	// Get query parameters
	status := r.URL.Query().Get("status")
	searchTerm := r.URL.Query().Get("search")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	invoices, err := s.db.FilterInvoices(status, searchTerm, startDate, endDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

// getCurrencyRates handles GET /currency-rates
func (s *Server) getCurrencyRates(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	rates, err := s.db.ListCurrencyRates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}

// convertCurrency handles GET /convert?amount=100&from=USD&to=IDR
func (s *Server) convertCurrency(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	// Get query parameters
	amountStr := r.URL.Query().Get("amount")
	fromCurrency := r.URL.Query().Get("from")
	toCurrency := r.URL.Query().Get("to")

	if amountStr == "" || fromCurrency == "" || toCurrency == "" {
		http.Error(w, "amount, from, and to parameters required", http.StatusBadRequest)
		return
	}

	var amount float64
	if _, err := fmt.Sscanf(amountStr, "%f", &amount); err != nil {
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	// If same currency, no conversion needed
	if fromCurrency == toCurrency {
		response := map[string]interface{}{
			"from":             fromCurrency,
			"to":               toCurrency,
			"amount":           amount,
			"converted_amount": amount,
			"rate":             1.0,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get exchange rate
	rate, err := s.db.GetCurrencyRate(fromCurrency, toCurrency)
	if err != nil {
		http.Error(w, "Currency rate not found", http.StatusNotFound)
		return
	}

	convertedAmount := amount * rate.Rate

	response := map[string]interface{}{
		"from":             fromCurrency,
		"to":               toCurrency,
		"amount":           amount,
		"converted_amount": convertedAmount,
		"rate":             rate.Rate,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

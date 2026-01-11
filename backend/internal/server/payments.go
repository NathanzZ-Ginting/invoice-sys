package server

import (
	"encoding/json"
	"net/http"

	"invoice-backend/internal/db"

	"github.com/gorilla/mux"
)

// recordPayment handles POST /payments
func (s *Server) recordPayment(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	var req db.PaymentCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.InvoiceID == "" || req.Amount <= 0 || req.PaymentMethod == "" {
		http.Error(w, "invoice_id, amount, and payment_method required", http.StatusBadRequest)
		return
	}

	payment, err := s.db.RecordPayment(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// getInvoicePayments handles GET /invoices/{id}/payments
func (s *Server) getInvoicePayments(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	invoiceID := vars["id"]

	payments, err := s.db.GetPaymentsByInvoice(invoiceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

// getAllPayments handles GET /payments
func (s *Server) getAllPayments(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	payments, err := s.db.GetAllPayments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)
}

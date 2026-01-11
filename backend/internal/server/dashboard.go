package server

import (
	"encoding/json"
	"net/http"
)

// getDashboardStats handles GET /dashboard/stats
func (s *Server) getDashboardStats(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "ALL"
	}

	stats, err := s.db.GetDashboardStats(currency)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// getRevenueByPeriod handles GET /dashboard/revenue
func (s *Server) getRevenueByPeriod(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	limit := 30 // Default to 30 days
	
	revenue, err := s.db.GetRevenueByPeriod(period, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(revenue)
}

// getTopCustomers handles GET /dashboard/top-customers
func (s *Server) getTopCustomers(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	limit := 10 // Top 10 customers
	
	customers, err := s.db.GetTopCustomers(limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// getOverdueInvoices handles GET /dashboard/overdue
func (s *Server) getOverdueInvoices(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	invoices, err := s.db.GetOverdueInvoices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(invoices)
}

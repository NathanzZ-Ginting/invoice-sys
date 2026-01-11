package server

import (
	"log"
	"net/http"

	"invoice-backend/internal/db"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	mux *mux.Router
	db  *db.Client
}

// New creates a new HTTP server with all routes configured
func New() http.Handler {
	d, err := db.New()
	if err != nil {
		log.Printf("Warning: DB not configured: %v", err)
		d = nil
	}

	r := mux.NewRouter()

	srv := &Server{
		mux: r,
		db:  d,
	}

	// Health check
	r.HandleFunc("/health", srv.health).Methods("GET")

	// Customer endpoints
	r.HandleFunc("/customers", srv.listCustomers).Methods("GET")
	r.HandleFunc("/customers", srv.createCustomer).Methods("POST")
	r.HandleFunc("/customers/{id}", srv.getCustomer).Methods("GET")
	r.HandleFunc("/customers/{id}", srv.updateCustomer).Methods("PUT")
	r.HandleFunc("/customers/{id}", srv.deleteCustomer).Methods("DELETE")

	// Invoice endpoints
	r.HandleFunc("/invoices", srv.listInvoices).Methods("GET")
	r.HandleFunc("/invoices/filter", srv.filterInvoices).Methods("GET")
	r.HandleFunc("/invoices", srv.createInvoice).Methods("POST")
	r.HandleFunc("/invoices/{id}", srv.getInvoice).Methods("GET")
	r.HandleFunc("/invoices/{id}", srv.updateInvoice).Methods("PUT")
	r.HandleFunc("/invoices/{id}/payments", srv.getInvoicePayments).Methods("GET")

	// Payment endpoints
	r.HandleFunc("/payments", srv.recordPayment).Methods("POST")
	r.HandleFunc("/payments", srv.getAllPayments).Methods("GET")

	// Dashboard & Analytics endpoints
	r.HandleFunc("/dashboard/stats", srv.getDashboardStats).Methods("GET")
	r.HandleFunc("/dashboard/revenue", srv.getRevenueByPeriod).Methods("GET")
	r.HandleFunc("/dashboard/top-customers", srv.getTopCustomers).Methods("GET")
	r.HandleFunc("/dashboard/overdue", srv.getOverdueInvoices).Methods("GET")

	// Currency endpoints
	r.HandleFunc("/currency-rates", srv.getCurrencyRates).Methods("GET")
	r.HandleFunc("/convert", srv.convertCurrency).Methods("GET")

	// Enable CORS for React frontend
	return handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)(r)
}

// health is a simple health check endpoint
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
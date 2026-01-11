package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// listCustomers handles GET /customers
func (s *Server) listCustomers(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	customers, err := s.db.ListCustomers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}

// createCustomer handles POST /customers
func (s *Server) createCustomer(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Phone       string `json:"phone,omitempty"`
		Address     string `json:"address,omitempty"`
		City        string `json:"city,omitempty"`
		PostalCode  string `json:"postal_code,omitempty"`
		Country     string `json:"country,omitempty"`
		CompanyName string `json:"company_name,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	customer, err := s.db.CreateCustomer(req.Name, req.Email, req.Phone, req.Address, req.City, req.PostalCode, req.Country, req.CompanyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(customer)
}

// getCustomer handles GET /customers/{id}
func (s *Server) getCustomer(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	customer, err := s.db.GetCustomer(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// updateCustomer handles PUT /customers/{id}
func (s *Server) updateCustomer(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Name        string `json:"name"`
		Email       string `json:"email"`
		Phone       string `json:"phone,omitempty"`
		Address     string `json:"address,omitempty"`
		City        string `json:"city,omitempty"`
		PostalCode  string `json:"postal_code,omitempty"`
		Country     string `json:"country,omitempty"`
		CompanyName string `json:"company_name,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" {
		http.Error(w, "Name and email are required", http.StatusBadRequest)
		return
	}

	customer, err := s.db.UpdateCustomer(id, req.Name, req.Email, req.Phone, req.Address, req.City, req.PostalCode, req.Country, req.CompanyName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customer)
}

// deleteCustomer handles DELETE /customers/{id}
func (s *Server) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	if s.db == nil {
		http.Error(w, "Database not configured", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	if err := s.db.DeleteCustomer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package handler

import (
	"encoding/json"
	"net/http"

	"invoice-backend/services/customer-service/internal/repository"
	"invoice-backend/services/shared/pkg/types"
	"invoice-backend/services/shared/pkg/utils"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	repo *repository.CustomerRepository
}

func NewCustomerHandler(repo *repository.CustomerRepository) *CustomerHandler {
	return &CustomerHandler{repo: repo}
}

// GetAll handles GET /customers
func (h *CustomerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	customers, err := h.repo.GetAll()
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Success(w, customers)
}

// GetByID handles GET /customers/{id}
func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	customer, err := h.repo.GetByID(id)
	if err != nil {
		utils.NotFound(w, err.Error())
		return
	}
	utils.Success(w, customer)
}

// Create handles POST /customers
func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var customerCreate types.CustomerCreate
	if err := json.NewDecoder(r.Body).Decode(&customerCreate); err != nil {
		utils.BadRequest(w, "Invalid JSON")
		return
	}

	if customerCreate.Name == "" || customerCreate.Email == "" {
		utils.BadRequest(w, "name and email are required")
		return
	}

	// Convert to Customer struct for database insertion
	customer := types.Customer{
		Name:    customerCreate.Name,
		Email:   customerCreate.Email,
		Phone:   customerCreate.Phone,
		Address: customerCreate.Address,
	}

	result, err := h.repo.Create(customer)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Created(w, result)
}

// Update handles PUT /customers/{id}
func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var customer types.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.BadRequest(w, "Invalid JSON")
		return
	}

	result, err := h.repo.Update(id, customer)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Success(w, result)
}

// Delete handles DELETE /customers/{id}
func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Success(w, map[string]string{"message": "Customer deleted successfully"})
}

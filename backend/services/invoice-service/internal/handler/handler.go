package handler

import (
	"encoding/json"
	"net/http"

	"invoice-backend/services/invoice-service/internal/pdf"
	"invoice-backend/services/invoice-service/internal/repository"
	"invoice-backend/services/shared/pkg/types"
	"invoice-backend/services/shared/pkg/utils"

	"github.com/gorilla/mux"
)

type InvoiceHandler struct {
	repo *repository.InvoiceRepository
}

func NewInvoiceHandler(repo *repository.InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{repo: repo}
}

// GetAll handles GET /invoices
func (h *InvoiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	currency := r.URL.Query().Get("currency")

	invoices, err := h.repo.GetAll(status, currency)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Success(w, invoices)
}

// GetByID handles GET /invoices/{id}
func (h *InvoiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	invoice, err := h.repo.GetByID(id)
	if err != nil {
		utils.NotFound(w, err.Error())
		return
	}

	// Get items
	invoice.Items, _ = h.repo.GetItems(invoice.ID)
	
	utils.Success(w, invoice)
}

// Create handles POST /invoices
func (h *InvoiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string       `json:"customer_id"`
		Items      []types.Item `json:"items"`
		Tax        float64      `json:"tax,omitempty"`      // Tax percentage
		Discount   float64      `json:"discount,omitempty"` // Discount amount
		Status     string       `json:"status,omitempty"`
		Notes      string       `json:"notes,omitempty"`
		DueDate    string       `json:"due_date,omitempty"`
		Currency   string       `json:"currency,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON")
		return
	}

	if req.CustomerID == "" || len(req.Items) == 0 {
		utils.BadRequest(w, "customer_id and items are required")
		return
	}

	invoice, err := h.repo.Create(req.CustomerID, req.Items, req.Tax, req.Discount, req.Status, req.Notes, req.DueDate, req.Currency)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Created(w, invoice)
}

// Update handles PUT /invoices/{id}
func (h *InvoiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var req struct {
		Status string `json:"status"`
		Notes  string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON")
		return
	}

	invoice, err := h.repo.Update(id, req.Status, req.Notes)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Success(w, invoice)
}

// Delete handles DELETE /invoices/{id}
func (h *InvoiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Success(w, map[string]string{"message": "Invoice deleted successfully"})
}

// GeneratePDF handles GET /invoices/{id}/pdf
func (h *InvoiceHandler) GeneratePDF(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Get invoice
	invoice, err := h.repo.GetByID(id)
	if err != nil {
		utils.NotFound(w, "Invoice not found")
		return
	}

	// Get items
	items, err := h.repo.GetItems(id)
	if err != nil {
		utils.InternalError(w, "Failed to get invoice items")
		return
	}

	// Get customer
	customer, err := h.repo.GetCustomer(invoice.CustomerID)
	if err != nil {
		utils.InternalError(w, "Failed to get customer")
		return
	}

	// Convert to PDF invoice type
	pdfItems := make([]pdf.Item, len(items))
	for i, item := range items {
		pdfItems[i] = pdf.Item{
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
		}
	}
	
	pdfInvoice := pdf.Invoice{
		ID:           invoice.InvoiceNumber,
		CustomerName: customer.Name,
		CustomerAddress: customer.Address,
		CustomerEmail: customer.Email,
		CustomerPhone: customer.Phone,
		Items:        pdfItems,
		Total:        invoice.Total,
		Currency:     invoice.Currency,
	}

	// Generate PDF
	pdfBytes, err := pdf.GeneratePDF(pdfInvoice)
	if err != nil {
		utils.InternalError(w, "Failed to generate PDF: "+err.Error())
		return
	}

	// Send PDF
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=invoice-"+invoice.InvoiceNumber+".pdf")
	w.Write(pdfBytes)
}

// GetCurrencyRates handles GET /currency-rates
func (h *InvoiceHandler) GetCurrencyRates(w http.ResponseWriter, r *http.Request) {
	rates, err := h.repo.GetCurrencyRates()
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Success(w, rates)
}

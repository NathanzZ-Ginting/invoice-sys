package handler

import (
	"encoding/json"
	"net/http"

	"invoice-backend/services/payment-service/internal/repository"
	"invoice-backend/services/shared/pkg/types"
	"invoice-backend/services/shared/pkg/utils"

	"github.com/gorilla/mux"
)

type PaymentHandler struct {
	repo *repository.PaymentRepository
}

func NewPaymentHandler(repo *repository.PaymentRepository) *PaymentHandler {
	return &PaymentHandler{repo: repo}
}

// Create handles POST /payments
func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req types.PaymentCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON")
		return
	}

	if req.InvoiceID == "" || req.Amount <= 0 || req.PaymentMethod == "" {
		utils.BadRequest(w, "invoice_id, amount, and payment_method required")
		return
	}

	// Get invoice to update payment status
	invoice, err := h.repo.GetInvoice(req.InvoiceID)
	if err != nil {
		utils.NotFound(w, "Invoice not found")
		return
	}

	// Create payment record
	payment, err := h.repo.Create(req)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	// Calculate new paid amount
	newPaidAmount := invoice.PaidAmount + req.Amount
	
	// Determine payment status
	paymentStatus := "partially_paid"
	if newPaidAmount >= invoice.Total {
		paymentStatus = "paid"
		newPaidAmount = invoice.Total // Cap at total
	}

	// Update invoice
	if err := h.repo.UpdateInvoice(req.InvoiceID, newPaidAmount, paymentStatus, req.PaymentDate); err != nil {
		utils.InternalError(w, "Failed to update invoice: "+err.Error())
		return
	}

	utils.Created(w, payment)
}

// GetByInvoiceID handles GET /payments/invoice/{id}
func (h *PaymentHandler) GetByInvoiceID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["id"]

	payments, err := h.repo.GetByInvoiceID(invoiceID)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Success(w, payments)
}

// GetAll handles GET /payments
func (h *PaymentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	payments, err := h.repo.GetAll()
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}
	utils.Success(w, payments)
}

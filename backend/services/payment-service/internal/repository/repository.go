package repository

import (
	"fmt"

	"invoice-backend/services/shared/pkg/database"
	"invoice-backend/services/shared/pkg/types"
)

type PaymentRepository struct {
	db *database.Client
}

func NewPaymentRepository(db *database.Client) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create records a new payment
func (r *PaymentRepository) Create(payment types.PaymentCreate) (*types.Payment, error) {
	var result []types.Payment
	_, err := r.db.Supabase.From("payments").
		Insert(payment, false, "", "", "").
		ExecuteTo(&result)
	
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no payment created")
	}
	return &result[0], nil
}

// GetByInvoiceID returns all payments for an invoice
func (r *PaymentRepository) GetByInvoiceID(invoiceID string) ([]types.Payment, error) {
	var payments []types.Payment
	_, err := r.db.Supabase.From("payments").
		Select("*", "", false).
		Eq("invoice_id", invoiceID).
		Order("payment_date", nil).
		ExecuteTo(&payments)
	return payments, err
}

// GetAll returns all payments
func (r *PaymentRepository) GetAll() ([]types.Payment, error) {
	var payments []types.Payment
	_, err := r.db.Supabase.From("payments").
		Select("*", "", false).
		Order("payment_date", nil).
		ExecuteTo(&payments)
	return payments, err
}

// GetInvoice returns an invoice by ID
func (r *PaymentRepository) GetInvoice(id string) (*types.Invoice, error) {
	var invoices []types.Invoice
	_, err := r.db.Supabase.From("invoices").
		Select("*", "", false).
		Eq("id", id).
		ExecuteTo(&invoices)
	
	if err != nil {
		return nil, err
	}
	if len(invoices) == 0 {
		return nil, fmt.Errorf("invoice not found")
	}
	return &invoices[0], nil
}

// UpdateInvoice updates invoice payment status
func (r *PaymentRepository) UpdateInvoice(id string, paidAmount float64, paymentStatus, paymentDate string) error {
	updateData := map[string]interface{}{
		"paid_amount":    paidAmount,
		"payment_status": paymentStatus,
		"payment_date":   paymentDate,
	}

	var updated []types.Invoice
	_, err := r.db.Supabase.From("invoices").
		Update(updateData, "", "").
		Eq("id", id).
		ExecuteTo(&updated)
	
	return err
}

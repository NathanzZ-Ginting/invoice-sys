package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"invoice-backend/services/notification-service/internal/email"
	"invoice-backend/services/shared/pkg/utils"
)

type NotificationHandler struct {
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// SendEmail handles POST /notifications/send
func (h *NotificationHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Body    string `json:"body"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON")
		return
	}

	if req.To == "" || req.Subject == "" || req.Body == "" {
		utils.BadRequest(w, "to, subject, and body are required")
		return
	}

	if err := email.SendEmail(req.To, req.Subject, req.Body, "", nil); err != nil {
		utils.InternalError(w, "Failed to send email: "+err.Error())
		return
	}

	utils.Success(w, map[string]string{
		"message": "Email sent successfully",
		"to":      req.To,
	})
}

// SendPaymentReminder handles POST /notifications/reminder
func (h *NotificationHandler) SendPaymentReminder(w http.ResponseWriter, r *http.Request) {
	var req struct {
		InvoiceID     string  `json:"invoice_id"`
		CustomerEmail string  `json:"customer_email"`
		CustomerName  string  `json:"customer_name"`
		InvoiceNumber string  `json:"invoice_number"`
		Amount        float64 `json:"amount"`
		DueDate       string  `json:"due_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.BadRequest(w, "Invalid JSON")
		return
	}

	if req.CustomerEmail == "" || req.InvoiceNumber == "" {
		utils.BadRequest(w, "customer_email and invoice_number are required")
		return
	}

	// Send payment reminder email
	subject := "Payment Reminder - Invoice " + req.InvoiceNumber
	body := h.generateReminderBody(req.CustomerName, req.InvoiceNumber, req.Amount, req.DueDate)

	if err := email.SendEmail(req.CustomerEmail, subject, body, "", nil); err != nil {
		utils.InternalError(w, "Failed to send reminder: "+err.Error())
		return
	}

	utils.Success(w, map[string]string{
		"message": "Payment reminder sent successfully",
		"to":      req.CustomerEmail,
	})
}

// generateReminderBody generates email body for payment reminder
func (h *NotificationHandler) generateReminderBody(customerName, invoiceNumber string, amount float64, dueDate string) string {
	return `Dear ` + customerName + `,

This is a friendly reminder that payment for Invoice ` + invoiceNumber + ` is due.

Invoice Details:
- Invoice Number: ` + invoiceNumber + `
- Amount: $` + fmt.Sprintf("%.2f", amount) + `
- Due Date: ` + dueDate + `

Please process the payment at your earliest convenience.

If you have already made the payment, please disregard this message.

Thank you for your business!

Best regards,
InvoicePro Systems
`
}

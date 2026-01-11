package handler

import (
	"net/http"
	"strconv"

	"invoice-backend/services/analytics-service/internal/repository"
	"invoice-backend/services/shared/pkg/utils"
)

type AnalyticsHandler struct {
	repo *repository.AnalyticsRepository
}

func NewAnalyticsHandler(repo *repository.AnalyticsRepository) *AnalyticsHandler {
	return &AnalyticsHandler{repo: repo}
}

// GetDashboardStats handles GET /dashboard/stats
func (h *AnalyticsHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	currency := r.URL.Query().Get("currency")
	if currency == "" {
		currency = "ALL"
	}

	stats, err := h.repo.GetDashboardStats(currency)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Success(w, stats)
}

// GetRevenueByPeriod handles GET /dashboard/revenue
func (h *AnalyticsHandler) GetRevenueByPeriod(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "daily"
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 30 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	revenue, err := h.repo.GetRevenueByPeriod(period, limit)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Success(w, revenue)
}

// GetTopCustomers handles GET /dashboard/top-customers
func (h *AnalyticsHandler) GetTopCustomers(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	limit := 10 // default top 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	customers, err := h.repo.GetTopCustomers(limit)
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Success(w, customers)
}

// GetOverdueInvoices handles GET /dashboard/overdue
func (h *AnalyticsHandler) GetOverdueInvoices(w http.ResponseWriter, r *http.Request) {
	invoices, err := h.repo.GetOverdueInvoices()
	if err != nil {
		utils.InternalError(w, err.Error())
		return
	}

	utils.Success(w, invoices)
}

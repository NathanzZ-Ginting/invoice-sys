package repository

import (
	"fmt"
	"time"

	"invoice-backend/services/shared/pkg/database"
	"invoice-backend/services/shared/pkg/types"
)

type AnalyticsRepository struct {
	db *database.Client
}

func NewAnalyticsRepository(db *database.Client) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

// GetDashboardStats returns overall statistics
func (r *AnalyticsRepository) GetDashboardStats(currency string) (*types.DashboardStats, error) {
	var invoices []types.Invoice
	query := r.db.Supabase.From("invoices").Select("*", "", false)

	if currency != "" && currency != "ALL" {
		query = query.Eq("currency", currency)
	}

	_, err := query.ExecuteTo(&invoices)
	if err != nil {
		return nil, err
	}

	stats := &types.DashboardStats{}

	for _, inv := range invoices {
		stats.TotalRevenue += inv.Total
		stats.TotalInvoices++

		switch inv.PaymentStatus {
		case "paid":
			stats.PaidAmount += inv.Total
			stats.PaidInvoices++
		case "unpaid":
			stats.UnpaidAmount += inv.Total
			stats.UnpaidInvoices++
		case "partially_paid":
			stats.PaidAmount += inv.PaidAmount
			stats.UnpaidAmount += (inv.Total - inv.PaidAmount)
		}

		// Check if overdue
		if inv.DueDate != "" && inv.PaymentStatus != "paid" {
			dueDate, err := time.Parse("2006-01-02", inv.DueDate)
			if err == nil && dueDate.Before(time.Now()) {
				stats.OverdueAmount += (inv.Total - inv.PaidAmount)
				stats.OverdueInvoices++
			}
		}
	}

	return stats, nil
}

// GetRevenueByPeriod returns revenue data grouped by period
func (r *AnalyticsRepository) GetRevenueByPeriod(period string, limit int) ([]types.RevenueData, error) {
	var invoices []types.Invoice
	_, err := r.db.Supabase.From("invoices").
		Select("created_at, total, payment_status", "", false).
		Eq("payment_status", "paid").
		Order("created_at", nil).
		Limit(limit, "").
		ExecuteTo(&invoices)

	if err != nil {
		return nil, err
	}

	// Group by period
	revenueMap := make(map[string]float64)

	for _, inv := range invoices {
		// Parse CreatedAt timestamp
		if inv.CreatedAt == "" {
			continue
		}
		
		createdAt, err := time.Parse(time.RFC3339, inv.CreatedAt)
		if err != nil {
			continue
		}

		var periodKey string
		switch period {
		case "daily":
			periodKey = createdAt.Format("2006-01-02")
		case "weekly":
			year, week := createdAt.ISOWeek()
			periodKey = fmt.Sprintf("%d-W%02d", year, week)
		case "monthly":
			periodKey = createdAt.Format("2006-01")
		default:
			periodKey = createdAt.Format("2006-01-02")
		}

		revenueMap[periodKey] += inv.Total
	}

	// Convert map to slice
	var result []types.RevenueData
	for period, amount := range revenueMap {
		result = append(result, types.RevenueData{
			Period: period,
			Amount: amount,
		})
	}

	return result, nil
}

// GetTopCustomers returns top customers by revenue
func (r *AnalyticsRepository) GetTopCustomers(limit int) ([]types.TopCustomer, error) {
	var invoices []types.Invoice
	_, err := r.db.Supabase.From("invoices").
		Select("customer_id, total", "", false).
		ExecuteTo(&invoices)

	if err != nil {
		return nil, err
	}

	// Group by customer
	customerMap := make(map[string]*types.TopCustomer)

	for _, inv := range invoices {
		if _, exists := customerMap[inv.CustomerID]; !exists {
			customerMap[inv.CustomerID] = &types.TopCustomer{
				CustomerID:    inv.CustomerID,
				TotalRevenue:  0,
				InvoiceCount:  0,
			}
		}
		customerMap[inv.CustomerID].TotalRevenue += inv.Total
		customerMap[inv.CustomerID].InvoiceCount++
	}

	// Fetch customer details for each customer
	for customerID, topCustomer := range customerMap {
		var customers []types.Customer
		_, err := r.db.Supabase.From("customers").
			Select("name, email", "", false).
			Eq("id", customerID).
			ExecuteTo(&customers)
		
		if err == nil && len(customers) > 0 {
			topCustomer.CustomerName = customers[0].Name
			topCustomer.CustomerEmail = customers[0].Email
		}
	}

	// Convert to slice and sort
	var result []types.TopCustomer
	for _, customer := range customerMap {
		result = append(result, *customer)
	}

	// Sort by revenue (descending)
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].TotalRevenue > result[i].TotalRevenue {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	// Limit results
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

// GetOverdueInvoices returns invoices that are past due date
func (r *AnalyticsRepository) GetOverdueInvoices() ([]types.Invoice, error) {
	var invoices []types.Invoice
	today := time.Now().Format("2006-01-02")

	_, err := r.db.Supabase.From("invoices").
		Select("*, customers(name, email)", "", false).
		Lt("due_date", today).
		Neq("payment_status", "paid").
		Order("due_date", nil).
		ExecuteTo(&invoices)

	return invoices, err
}

// API Configuration
export const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

// API Endpoints
export const API_ENDPOINTS = {
  // Customers
  customers: `${API_BASE_URL}/customers`,
  customer: (id) => `${API_BASE_URL}/customers/${id}`,
  
  // Invoices
  invoices: `${API_BASE_URL}/invoices`,
  invoice: (id) => `${API_BASE_URL}/invoices/${id}`,
  invoicePDF: (id) => `${API_BASE_URL}/invoices/${id}/pdf`,
  
  // Payments
  payments: `${API_BASE_URL}/payments`,
  invoicePayments: (id) => `${API_BASE_URL}/invoices/${id}/payments`,
  
  // Dashboard & Analytics
  dashboardStats: (currency = 'ALL') => `${API_BASE_URL}/dashboard/stats?currency=${currency}`,
  dashboardRevenue: (period = 'daily') => `${API_BASE_URL}/dashboard/revenue?period=${period}`,
  dashboardTopCustomers: `${API_BASE_URL}/dashboard/top-customers`,
  dashboardOverdue: `${API_BASE_URL}/dashboard/overdue`,
  
  // Notifications
  sendEmail: `${API_BASE_URL}/notifications/send`,
  sendReminder: `${API_BASE_URL}/notifications/reminder`,
};

export default API_BASE_URL;

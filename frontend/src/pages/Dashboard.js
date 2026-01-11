import React, { useState, useEffect } from 'react';

const API_BASE = 'http://localhost:8080';

function Dashboard() {
  const [stats, setStats] = useState(null);
  const [topCustomers, setTopCustomers] = useState([]);
  const [overdueInvoices, setOverdueInvoices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedCurrency, setSelectedCurrency] = useState('USD');

  useEffect(() => {
    fetchDashboardData();
  }, [selectedCurrency]);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      
      // Fetch dashboard data in parallel
      const [statsRes, topCustomersRes, overdueRes] = await Promise.all([
        fetch(`${API_BASE}/dashboard/stats?currency=${selectedCurrency}`),
        fetch(`${API_BASE}/dashboard/top-customers`),
        fetch(`${API_BASE}/dashboard/overdue`)
      ]);

      if (!statsRes.ok || !topCustomersRes.ok || !overdueRes.ok) {
        throw new Error('Failed to fetch dashboard data');
      }

      const [statsData, topCustomersData, overdueData] = await Promise.all([
        statsRes.json(),
        topCustomersRes.json(),
        overdueRes.json()
      ]);

      setStats(statsData);
      setTopCustomers(topCustomersData || []);
      setOverdueInvoices(overdueData || []);
    } catch (err) {
      console.error('Dashboard error:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const getCurrencySymbol = (currency) => {
    const symbols = {
      'USD': '$', 'EUR': '€', 'GBP': '£', 'JPY': '¥', 'AUD': 'A$',
      'CAD': 'C$', 'CHF': 'Fr', 'CNY': '¥', 'INR': '₹', 'MYR': 'RM',
      'SGD': 'S$', 'IDR': 'Rp', 'THB': '฿', 'PHP': '₱'
    };
    return symbols[currency] || currency + ' ';
  };

  const formatCurrencyAmount = (amount, currency) => {
    if (!amount) return '0';
    if (currency === 'IDR' || currency === 'JPY') {
      return Math.round(amount).toLocaleString('en-US', { minimumFractionDigits: 0, maximumFractionDigits: 0 });
    }
    return amount.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <svg className="animate-spin h-12 w-12 text-indigo-600" fill="none" viewBox="0 0 24 24">
          <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
          <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg">
        Error loading dashboard: {error}
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">📊 Dashboard</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">Overview of your invoice analytics</p>
        </div>
        <div className="flex items-center gap-3">
          <label className="text-sm font-medium text-gray-700 dark:text-gray-300">Currency:</label>
          <select
            value={selectedCurrency}
            onChange={(e) => setSelectedCurrency(e.target.value)}
            className="px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
          >
            <option value="USD">🇺🇸 USD</option>
            <option value="EUR">🇪🇺 EUR</option>
            <option value="GBP">🇬🇧 GBP</option>
            <option value="JPY">🇯🇵 JPY</option>
            <option value="IDR">🇮🇩 IDR</option>
            <option value="MYR">🇲🇾 MYR</option>
            <option value="SGD">🇸🇬 SGD</option>
          </select>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {/* Total Revenue */}
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl shadow-lg p-6 text-white">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium opacity-90">Total Revenue</h3>
            <span className="text-2xl">💰</span>
          </div>
          <div className="text-3xl font-bold mb-1">
            {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats?.total_revenue || 0, selectedCurrency)}
          </div>
          <p className="text-xs opacity-75">{stats?.total_invoices || 0} total invoices</p>
        </div>

        {/* Paid Amount */}
        <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-xl shadow-lg p-6 text-white">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium opacity-90">Paid Amount</h3>
            <span className="text-2xl">✅</span>
          </div>
          <div className="text-3xl font-bold mb-1">
            {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats?.paid_amount || 0, selectedCurrency)}
          </div>
          <p className="text-xs opacity-75">{stats?.paid_invoices || 0} paid invoices</p>
        </div>

        {/* Unpaid Amount */}
        <div className="bg-gradient-to-br from-yellow-500 to-yellow-600 rounded-xl shadow-lg p-6 text-white">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium opacity-90">Unpaid Amount</h3>
            <span className="text-2xl">⏳</span>
          </div>
          <div className="text-3xl font-bold mb-1">
            {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats?.unpaid_amount || 0, selectedCurrency)}
          </div>
          <p className="text-xs opacity-75">{stats?.unpaid_invoices || 0} unpaid invoices</p>
        </div>

        {/* Overdue Amount */}
        <div className="bg-gradient-to-br from-red-500 to-red-600 rounded-xl shadow-lg p-6 text-white">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium opacity-90">Overdue Amount</h3>
            <span className="text-2xl">🚨</span>
          </div>
          <div className="text-3xl font-bold mb-1">
            {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats?.overdue_amount || 0, selectedCurrency)}
          </div>
          <p className="text-xs opacity-75">{stats?.overdue_invoices || 0} overdue invoices</p>
        </div>
      </div>

      {/* Two Column Layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Customers */}
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-6">👥 Top Customers</h2>
          {topCustomers.length > 0 ? (
            <div className="space-y-4">
              {topCustomers.map((customer, index) => (
                <div key={customer.customer_id} className="flex items-center gap-4 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
                  <div className="flex-shrink-0 w-8 h-8 flex items-center justify-center bg-indigo-100 dark:bg-indigo-900/50 text-indigo-600 dark:text-indigo-400 font-bold rounded-full">
                    {index + 1}
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="font-semibold text-gray-900 dark:text-gray-100 truncate">
                      {customer.customer_name}
                    </div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      {customer.invoice_count} {customer.invoice_count === 1 ? 'invoice' : 'invoices'}
                    </div>
                  </div>
                  <div className="text-right">
                    <div className="font-bold text-gray-900 dark:text-gray-100">
                      {getCurrencySymbol(customer.currency)}{formatCurrencyAmount(customer.total_amount, customer.currency)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
              No customer data available
            </div>
          )}
        </div>

        {/* Overdue Invoices */}
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-6">🚨 Overdue Invoices</h2>
          {overdueInvoices.length > 0 ? (
            <div className="space-y-3 max-h-96 overflow-y-auto">
              {overdueInvoices.map((invoice) => (
                <div key={invoice.id} className="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
                  <div className="flex items-center justify-between mb-2">
                    <span className="font-mono text-sm font-semibold text-red-600 dark:text-red-400">
                      {invoice.invoice_number || `INV-${invoice.id.substring(0, 8)}`}
                    </span>
                    <span className="text-sm font-bold text-red-600 dark:text-red-400">
                      {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(invoice.total, invoice.currency)}
                    </span>
                  </div>
                  <div className="text-xs text-gray-600 dark:text-gray-400">
                    Customer ID: {invoice.customer_id.substring(0, 8)}...
                  </div>
                  {invoice.due_date && (
                    <div className="text-xs text-red-600 dark:text-red-400 mt-1">
                      Due: {new Date(invoice.due_date).toLocaleDateString('en-US', { year: 'numeric', month: 'short', day: 'numeric' })}
                    </div>
                  )}
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500 dark:text-gray-400">
              🎉 No overdue invoices!
            </div>
          )}
        </div>
      </div>

      {/* Payment Status Breakdown */}
      {stats && (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-6">📊 Payment Status Breakdown</h2>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            <div className="text-center p-4 bg-green-50 dark:bg-green-900/20 rounded-lg">
              <div className="text-2xl font-bold text-green-600 dark:text-green-400">{stats.paid_invoices}</div>
              <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">Paid</div>
            </div>
            <div className="text-center p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
              <div className="text-2xl font-bold text-blue-600 dark:text-blue-400">{stats.partially_paid_invoices}</div>
              <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">Partially Paid</div>
            </div>
            <div className="text-center p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-lg">
              <div className="text-2xl font-bold text-yellow-600 dark:text-yellow-400">{stats.unpaid_invoices}</div>
              <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">Unpaid</div>
            </div>
            <div className="text-center p-4 bg-red-50 dark:bg-red-900/20 rounded-lg">
              <div className="text-2xl font-bold text-red-600 dark:text-red-400">{stats.overdue_invoices}</div>
              <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">Overdue</div>
            </div>
            <div className="text-center p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg">
              <div className="text-2xl font-bold text-gray-900 dark:text-gray-100">{stats.total_invoices}</div>
              <div className="text-sm text-gray-600 dark:text-gray-400 mt-1">Total</div>
            </div>
          </div>
        </div>
      )}

      {/* Quick Actions */}
      <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-6">
        <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <a href="/customers" className="flex items-center p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors border border-gray-200 dark:border-gray-600">
            <div className="w-10 h-10 bg-blue-500 rounded-lg flex items-center justify-center mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
              </svg>
            </div>
            <div>
              <p className="font-semibold text-gray-900 dark:text-gray-100">Add Customer</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">Create new customer</p>
            </div>
          </a>

          <a href="/invoices" className="flex items-center p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors border border-gray-200 dark:border-gray-600">
            <div className="w-10 h-10 bg-blue-500 rounded-lg flex items-center justify-center mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div>
              <p className="font-semibold text-gray-900 dark:text-gray-100">Create Invoice</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">Generate new invoice</p>
            </div>
          </a>

          <a href="/analytics" className="flex items-center p-4 bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 rounded-lg hover:from-purple-100 hover:to-indigo-100 dark:hover:from-purple-900/30 dark:hover:to-indigo-900/30 transition-colors border border-purple-200 dark:border-purple-800">
            <div className="w-10 h-10 bg-gradient-to-br from-purple-500 to-indigo-600 rounded-lg flex items-center justify-center mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div>
              <p className="font-semibold text-gray-900 dark:text-gray-100">📊 Analytics</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">Revenue & Charts</p>
            </div>
          </a>

          <a href="/settings" className="flex items-center p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors border border-gray-200 dark:border-gray-600">
            <div className="w-10 h-10 bg-gray-500 rounded-lg flex items-center justify-center mr-4">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div>
              <p className="font-semibold text-gray-900 dark:text-gray-100">Settings</p>
              <p className="text-sm text-gray-600 dark:text-gray-400">Company settings</p>
            </div>
          </a>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;

import React, { useState, useEffect } from 'react';

const Analytics = () => {
  const [revenue, setRevenue] = useState([]);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [revenuePeriod, setRevenuePeriod] = useState('month');
  const [selectedCurrency, setSelectedCurrency] = useState('USD');

  useEffect(() => {
    fetchAnalyticsData();
  }, [revenuePeriod, selectedCurrency]);

  const fetchAnalyticsData = async () => {
    try {
      setLoading(true);
      
      const [revenueRes, statsRes] = await Promise.all([
        fetch(`http://localhost:8080/dashboard/revenue?period=${revenuePeriod}`),
        fetch(`http://localhost:8080/dashboard/stats?currency=${selectedCurrency}`)
      ]);

      if (!revenueRes.ok || !statsRes.ok) {
        throw new Error('Failed to fetch analytics data');
      }

      const [revenueData, statsData] = await Promise.all([
        revenueRes.json(),
        statsRes.json()
      ]);

      setRevenue(revenueData || []);
      setStats(statsData);
    } catch (err) {
      console.error('Analytics error:', err);
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
        Error loading analytics: {error}
      </div>
    );
  }

  // Calculate chart data
  const maxRevenue = revenue.length > 0 ? Math.max(...revenue.map(r => r.total_amount)) : 0;
  const totalRevenue = revenue.reduce((sum, item) => sum + item.total_amount, 0);
  const avgRevenue = revenue.length > 0 ? totalRevenue / revenue.length : 0;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">📊 Analytics</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">Detailed revenue analysis and trends</p>
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

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-gradient-to-br from-blue-500 to-blue-600 rounded-xl shadow-lg p-6 text-white">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium opacity-90">Total Revenue ({revenuePeriod})</h3>
            <span className="text-2xl">💰</span>
          </div>
          <div className="text-3xl font-bold mb-1">
            {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(totalRevenue, selectedCurrency)}
          </div>
          <p className="text-xs opacity-75">{revenue.length} periods</p>
        </div>

        <div className="bg-gradient-to-br from-purple-500 to-purple-600 rounded-xl shadow-lg p-6 text-white">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium opacity-90">Average per Period</h3>
            <span className="text-2xl">📈</span>
          </div>
          <div className="text-3xl font-bold mb-1">
            {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(avgRevenue, selectedCurrency)}
          </div>
          <p className="text-xs opacity-75">Per {revenuePeriod}</p>
        </div>

        <div className="bg-gradient-to-br from-green-500 to-green-600 rounded-xl shadow-lg p-6 text-white">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium opacity-90">Peak Revenue</h3>
            <span className="text-2xl">🚀</span>
          </div>
          <div className="text-3xl font-bold mb-1">
            {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(maxRevenue, selectedCurrency)}
          </div>
          <p className="text-xs opacity-75">Highest period</p>
        </div>
      </div>

      {/* Revenue Bar Chart */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">📈 Revenue Trends</h2>
          <div className="flex items-center gap-2">
            <button
              onClick={() => setRevenuePeriod('day')}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                revenuePeriod === 'day'
                  ? 'bg-indigo-600 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              }`}
            >
              Daily
            </button>
            <button
              onClick={() => setRevenuePeriod('week')}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                revenuePeriod === 'week'
                  ? 'bg-indigo-600 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              }`}
            >
              Weekly
            </button>
            <button
              onClick={() => setRevenuePeriod('month')}
              className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                revenuePeriod === 'month'
                  ? 'bg-indigo-600 text-white'
                  : 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-600'
              }`}
            >
              Monthly
            </button>
          </div>
        </div>
        
        {revenue.length > 0 ? (
          <div className="space-y-3">
            {revenue.map((item, index) => {
              const widthPercent = (item.total_amount / maxRevenue) * 100;
              
              return (
                <div key={index} className="flex items-center gap-4">
                  <div className="w-32 text-sm font-medium text-gray-700 dark:text-gray-300 flex-shrink-0">
                    {item.period}
                  </div>
                  <div className="flex-1">
                    <div 
                      className="h-12 bg-gradient-to-r from-indigo-500 to-indigo-600 hover:from-indigo-600 hover:to-indigo-700 rounded-lg relative transition-all group cursor-pointer" 
                      style={{ width: `${widthPercent}%`, minWidth: '80px' }}
                    >
                      <div className="absolute inset-0 flex items-center justify-end pr-3">
                        <span className="text-white font-semibold text-sm">
                          {getCurrencySymbol(item.currency)}{formatCurrencyAmount(item.total_amount, item.currency)}
                        </span>
                      </div>
                      {/* Tooltip on hover */}
                      <div className="absolute bottom-full left-1/2 -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap">
                        {item.invoice_count} {item.invoice_count === 1 ? 'invoice' : 'invoices'}
                      </div>
                    </div>
                  </div>
                  <div className="w-24 text-right text-sm text-gray-500 dark:text-gray-400 flex-shrink-0">
                    {item.invoice_count} inv
                  </div>
                </div>
              );
            })}
          </div>
        ) : (
          <div className="text-center py-12 text-gray-500 dark:text-gray-400">
            <svg className="w-16 h-16 mx-auto mb-4 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <p className="text-lg font-medium">No revenue data available</p>
            <p className="text-sm mt-1">Create some invoices to see analytics</p>
          </div>
        )}
      </div>

      {/* Payment Status Pie Chart Representation */}
      {stats && (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-6">💳 Payment Status Distribution</h2>
          <div className="grid grid-cols-2 md:grid-cols-5 gap-4">
            <div className="text-center p-6 bg-green-50 dark:bg-green-900/20 rounded-xl border-2 border-green-200 dark:border-green-800 hover:shadow-lg transition-shadow">
              <div className="text-4xl font-bold text-green-600 dark:text-green-400 mb-2">{stats.paid_invoices}</div>
              <div className="text-sm font-medium text-gray-700 dark:text-gray-300">Paid</div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats.paid_amount, selectedCurrency)}
              </div>
            </div>
            <div className="text-center p-6 bg-blue-50 dark:bg-blue-900/20 rounded-xl border-2 border-blue-200 dark:border-blue-800 hover:shadow-lg transition-shadow">
              <div className="text-4xl font-bold text-blue-600 dark:text-blue-400 mb-2">{stats.partially_paid_invoices}</div>
              <div className="text-sm font-medium text-gray-700 dark:text-gray-300">Partially Paid</div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats.partially_paid_amount || 0, selectedCurrency)}
              </div>
            </div>
            <div className="text-center p-6 bg-yellow-50 dark:bg-yellow-900/20 rounded-xl border-2 border-yellow-200 dark:border-yellow-800 hover:shadow-lg transition-shadow">
              <div className="text-4xl font-bold text-yellow-600 dark:text-yellow-400 mb-2">{stats.unpaid_invoices}</div>
              <div className="text-sm font-medium text-gray-700 dark:text-gray-300">Unpaid</div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats.unpaid_amount, selectedCurrency)}
              </div>
            </div>
            <div className="text-center p-6 bg-red-50 dark:bg-red-900/20 rounded-xl border-2 border-red-200 dark:border-red-800 hover:shadow-lg transition-shadow">
              <div className="text-4xl font-bold text-red-600 dark:text-red-400 mb-2">{stats.overdue_invoices}</div>
              <div className="text-sm font-medium text-gray-700 dark:text-gray-300">Overdue</div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats.overdue_amount, selectedCurrency)}
              </div>
            </div>
            <div className="text-center p-6 bg-gray-50 dark:bg-gray-700/50 rounded-xl border-2 border-gray-300 dark:border-gray-600 hover:shadow-lg transition-shadow">
              <div className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-2">{stats.total_invoices}</div>
              <div className="text-sm font-medium text-gray-700 dark:text-gray-300">Total</div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(stats.total_revenue, selectedCurrency)}
              </div>
            </div>
          </div>

          {/* Visual Progress Bars */}
          <div className="mt-8 space-y-4">
            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="font-medium text-gray-700 dark:text-gray-300">Collection Rate</span>
                <span className="text-gray-500 dark:text-gray-400">
                  {stats.total_invoices > 0 ? Math.round((stats.paid_amount / stats.total_revenue) * 100) : 0}%
                </span>
              </div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-gradient-to-r from-green-500 to-green-600 rounded-full transition-all duration-500"
                  style={{ width: `${stats.total_invoices > 0 ? (stats.paid_amount / stats.total_revenue) * 100 : 0}%` }}
                ></div>
              </div>
            </div>

            <div>
              <div className="flex justify-between text-sm mb-2">
                <span className="font-medium text-gray-700 dark:text-gray-300">Outstanding Amount</span>
                <span className="text-gray-500 dark:text-gray-400">
                  {stats.total_invoices > 0 ? Math.round(((stats.unpaid_amount + stats.overdue_amount) / stats.total_revenue) * 100) : 0}%
                </span>
              </div>
              <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
                <div 
                  className="h-full bg-gradient-to-r from-yellow-500 to-red-600 rounded-full transition-all duration-500"
                  style={{ width: `${stats.total_invoices > 0 ? ((stats.unpaid_amount + stats.overdue_amount) / stats.total_revenue) * 100 : 0}%` }}
                ></div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default Analytics;

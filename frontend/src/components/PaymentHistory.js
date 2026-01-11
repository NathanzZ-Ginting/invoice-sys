import React, { useState, useEffect } from 'react';

const PaymentHistory = ({ invoice }) => {
  const [payments, setPayments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchPayments();
  }, [invoice.id]);

  const fetchPayments = async () => {
    try {
      setLoading(true);
      const response = await fetch(`http://localhost:8080/invoices/${invoice.id}/payments`);
      
      if (!response.ok) {
        throw new Error('Failed to fetch payments');
      }

      const data = await response.json();
      setPayments(data || []);
    } catch (err) {
      console.error('Error fetching payments:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const getCurrencySymbol = (currency) => {
    const symbols = {
      'USD': '$',
      'EUR': '€',
      'GBP': '£',
      'JPY': '¥',
      'AUD': 'A$',
      'CAD': 'C$',
      'CHF': 'Fr',
      'CNY': '¥',
      'INR': '₹',
      'MYR': 'RM',
      'SGD': 'S$',
      'IDR': 'Rp',
      'THB': '฿',
      'PHP': '₱'
    };
    return symbols[currency] || currency + ' ';
  };

  const formatCurrencyAmount = (amount, currency) => {
    if (currency === 'IDR' || currency === 'JPY') {
      return Math.round(amount).toLocaleString('en-US', { minimumFractionDigits: 0, maximumFractionDigits: 0 });
    }
    return amount.toLocaleString('en-US', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
  };

  const getPaymentMethodIcon = (method) => {
    const icons = {
      'bank_transfer': '🏦',
      'credit_card': '💳',
      'cash': '💵',
      'check': '📝',
      'paypal': '🅿️',
      'other': '🔄'
    };
    return icons[method] || '💰';
  };

  const getPaymentMethodLabel = (method) => {
    const labels = {
      'bank_transfer': 'Bank Transfer',
      'credit_card': 'Credit Card',
      'cash': 'Cash',
      'check': 'Check',
      'paypal': 'PayPal',
      'other': 'Other'
    };
    return labels[method] || method;
  };

  if (loading) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
        <div className="flex items-center justify-center py-8">
          <svg className="animate-spin h-8 w-8 text-indigo-600" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
        <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
          Error loading payments: {error}
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700">
      <div className="p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100">
            💳 Payment History
          </h3>
          <span className="px-3 py-1 bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300 rounded-full text-sm font-semibold">
            {payments.length} {payments.length === 1 ? 'Payment' : 'Payments'}
          </span>
        </div>
      </div>

      <div className="p-6">
        {payments.length === 0 ? (
          <div className="text-center py-8">
            <div className="text-4xl mb-2">💸</div>
            <p className="text-gray-500 dark:text-gray-400 text-sm">No payments recorded yet</p>
          </div>
        ) : (
          <div className="space-y-4">
            {payments.map((payment, index) => (
              <div
                key={payment.id}
                className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors"
              >
                <div className="flex items-start justify-between">
                  {/* Left side - Payment info */}
                  <div className="flex-1">
                    <div className="flex items-center gap-3 mb-2">
                      <span className="text-2xl">{getPaymentMethodIcon(payment.payment_method)}</span>
                      <div>
                        <div className="font-semibold text-gray-900 dark:text-gray-100">
                          {getPaymentMethodLabel(payment.payment_method)}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          {new Date(payment.payment_date).toLocaleString('en-US', {
                            year: 'numeric',
                            month: 'short',
                            day: 'numeric',
                            hour: '2-digit',
                            minute: '2-digit'
                          })}
                        </div>
                      </div>
                    </div>

                    {/* Reference Number */}
                    {payment.reference_number && (
                      <div className="flex items-center gap-2 mb-1">
                        <span className="text-xs text-gray-500 dark:text-gray-400">Ref:</span>
                        <span className="text-xs font-mono text-gray-700 dark:text-gray-300 bg-gray-100 dark:bg-gray-700 px-2 py-0.5 rounded">
                          {payment.reference_number}
                        </span>
                      </div>
                    )}

                    {/* Notes */}
                    {payment.notes && (
                      <div className="mt-2 text-sm text-gray-600 dark:text-gray-400 italic">
                        "{payment.notes}"
                      </div>
                    )}
                  </div>

                  {/* Right side - Amount */}
                  <div className="text-right ml-4">
                    <div className="text-2xl font-bold text-green-600 dark:text-green-400">
                      +{getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(payment.amount, invoice.currency)}
                    </div>
                    <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                      Payment #{payments.length - index}
                    </div>
                  </div>
                </div>
              </div>
            ))}

            {/* Total Summary */}
            <div className="border-t-2 border-gray-300 dark:border-gray-600 pt-4 mt-4">
              <div className="flex justify-between items-center">
                <span className="text-sm font-semibold text-gray-700 dark:text-gray-300">Total Paid:</span>
                <span className="text-xl font-bold text-green-600 dark:text-green-400">
                  {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(invoice.paid_amount || 0, invoice.currency)}
                </span>
              </div>
              <div className="flex justify-between items-center mt-2">
                <span className="text-sm font-semibold text-gray-700 dark:text-gray-300">Remaining Balance:</span>
                <span className="text-xl font-bold text-red-600 dark:text-red-400">
                  {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(invoice.total - (invoice.paid_amount || 0), invoice.currency)}
                </span>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default PaymentHistory;

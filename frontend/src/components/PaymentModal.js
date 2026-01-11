import React, { useState } from 'react';

const PaymentModal = ({ invoice, onClose, onPaymentRecorded }) => {
  const [formData, setFormData] = useState({
    amount: '',
    payment_method: 'bank_transfer',
    reference_number: '',
    notes: ''
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const remainingAmount = invoice.total - (invoice.paid_amount || 0);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    // Validation - Clean the amount string (remove commas, spaces)
    const cleanedAmount = formData.amount.toString().replace(/,/g, '').replace(/\s/g, '');
    const amount = parseFloat(cleanedAmount);
    
    console.log('💰 Payment Debug:', {
      raw: formData.amount,
      cleaned: cleanedAmount,
      parsed: amount,
      isValid: !isNaN(amount) && amount > 0
    });
    
    if (isNaN(amount) || amount <= 0) {
      setError('Please enter a valid amount');
      setLoading(false);
      return;
    }

    if (amount > remainingAmount) {
      setError(`Amount cannot exceed remaining balance: ${getCurrencySymbol(invoice.currency)}${formatCurrencyAmount(remainingAmount, invoice.currency)}`);
      setLoading(false);
      return;
    }

    const payload = {
      invoice_id: invoice.id,
      amount: amount,
      payment_method: formData.payment_method,
      payment_date: new Date().toISOString(), // Add current timestamp
      reference_number: formData.reference_number.trim(),
      notes: formData.notes.trim()
    };
    
    console.log('📤 Sending payment:', payload);

    try {
      const response = await fetch('http://localhost:8080/payments', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload),
      });

      console.log('📥 Response status:', response.status);

      if (!response.ok) {
        const errorText = await response.text();
        console.error('❌ Payment error response:', errorText);
        try {
          const errorData = JSON.parse(errorText);
          throw new Error(errorData.error || 'Failed to record payment');
        } catch (parseErr) {
          throw new Error(errorText || 'Failed to record payment');
        }
      }

      const result = await response.json();
      console.log('✅ Payment success:', result);
      
      // Call callback to refresh invoice data
      if (onPaymentRecorded) {
        onPaymentRecorded(result);
      }
      
      // Close modal
      onClose();
    } catch (err) {
      console.error('Payment recording error:', err);
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

  const handleQuickFill = (percentage) => {
    const amount = (remainingAmount * percentage / 100).toFixed(2);
    setFormData({ ...formData, amount });
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-md w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <div className="flex items-center justify-between">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              💰 Record Payment
            </h2>
            <button
              onClick={onClose}
              className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
            >
              <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        {/* Invoice Info */}
        <div className="p-6 bg-gray-50 dark:bg-gray-700/50 border-b border-gray-200 dark:border-gray-700">
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="text-sm text-gray-600 dark:text-gray-400">Invoice Number:</span>
              <span className="text-sm font-mono font-semibold text-indigo-600 dark:text-indigo-400">
                {invoice.invoice_number || `INV-${invoice.id.substring(0, 8)}`}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm text-gray-600 dark:text-gray-400">Total Amount:</span>
              <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(invoice.total, invoice.currency)}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="text-sm text-gray-600 dark:text-gray-400">Paid Amount:</span>
              <span className="text-sm font-semibold text-green-600 dark:text-green-400">
                {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(invoice.paid_amount || 0, invoice.currency)}
              </span>
            </div>
            <div className="flex justify-between pt-2 border-t border-gray-300 dark:border-gray-600">
              <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">Remaining Balance:</span>
              <span className="text-lg font-bold text-red-600 dark:text-red-400">
                {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(remainingAmount, invoice.currency)}
              </span>
            </div>
          </div>
        </div>

        {/* Form */}
        <form onSubmit={handleSubmit} noValidate className="p-6 space-y-4">
          {/* Error Message */}
          {error && (
            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded-lg text-sm">
              {error}
            </div>
          )}

          {/* Payment Amount */}
          <div>
            <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
              Payment Amount *
            </label>
            <div className="relative">
              <span className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500 dark:text-gray-400 font-semibold">
                {getCurrencySymbol(invoice.currency)}
              </span>
              <input
                type="text"
                inputMode="decimal"
                value={formData.amount}
                onChange={(e) => {
                  // Allow only numbers and decimal point
                  const value = e.target.value.replace(/[^\d.]/g, '');
                  // Ensure only one decimal point
                  const parts = value.split('.');
                  const cleaned = parts.length > 2 ? parts[0] + '.' + parts.slice(1).join('') : value;
                  setFormData({ ...formData, amount: cleaned });
                }}
                className="w-full pl-10 pr-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
                placeholder="0.00"
              />
            </div>
            {/* Quick Fill Buttons */}
            <div className="flex gap-2 mt-2">
              <button
                type="button"
                onClick={() => handleQuickFill(25)}
                className="flex-1 px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded transition-colors"
              >
                25%
              </button>
              <button
                type="button"
                onClick={() => handleQuickFill(50)}
                className="flex-1 px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded transition-colors"
              >
                50%
              </button>
              <button
                type="button"
                onClick={() => handleQuickFill(75)}
                className="flex-1 px-2 py-1 text-xs bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded transition-colors"
              >
                75%
              </button>
              <button
                type="button"
                onClick={() => setFormData({ ...formData, amount: remainingAmount.toFixed(2) })}
                className="flex-1 px-2 py-1 text-xs bg-indigo-100 dark:bg-indigo-900/50 hover:bg-indigo-200 dark:hover:bg-indigo-900 text-indigo-700 dark:text-indigo-300 rounded transition-colors font-semibold"
              >
                Full
              </button>
            </div>
          </div>

          {/* Payment Method */}
          <div>
            <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
              Payment Method *
            </label>
            <select
              value={formData.payment_method}
              onChange={(e) => setFormData({ ...formData, payment_method: e.target.value })}
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
            >
              <option value="bank_transfer">🏦 Bank Transfer</option>
              <option value="credit_card">💳 Credit Card</option>
              <option value="cash">💵 Cash</option>
              <option value="check">📝 Check</option>
              <option value="paypal">🅿️ PayPal</option>
              <option value="other">🔄 Other</option>
            </select>
          </div>

          {/* Reference Number */}
          <div>
            <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
              Reference Number
            </label>
            <input
              type="text"
              value={formData.reference_number}
              onChange={(e) => setFormData({ ...formData, reference_number: e.target.value })}
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              placeholder="e.g., TXN123456, Check #789"
            />
          </div>

          {/* Notes */}
          <div>
            <label className="block text-sm font-semibold text-gray-700 dark:text-gray-300 mb-2">
              Notes
            </label>
            <textarea
              value={formData.notes}
              onChange={(e) => setFormData({ ...formData, notes: e.target.value })}
              rows="3"
              className="w-full px-4 py-3 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100"
              placeholder="Additional notes about this payment..."
            />
          </div>

          {/* Action Buttons */}
          <div className="flex gap-3 pt-4">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-3 bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-300 rounded-lg font-semibold transition-colors"
              disabled={loading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className="flex-1 px-4 py-3 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg font-semibold transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={loading}
            >
              {loading ? (
                <span className="flex items-center justify-center">
                  <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Recording...
                </span>
              ) : (
                '💰 Record Payment'
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default PaymentModal;

import React, { useState, useEffect } from 'react';
import PaymentModal from '../components/PaymentModal';
import PaymentHistory from '../components/PaymentHistory';

const API_BASE = 'http://localhost:8080';

function Invoices() {
  const [customers, setCustomers] = useState([]);
  const [invoices, setInvoices] = useState([]);
  const [filteredInvoices, setFilteredInvoices] = useState([]);
  const [showForm, setShowForm] = useState(false);
  const [showEditModal, setShowEditModal] = useState(false);
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [showPaymentHistory, setShowPaymentHistory] = useState(false);
  const [selectedInvoice, setSelectedInvoice] = useState(null);
  const [invoiceData, setInvoiceData] = useState({
    customer_id: '',
    items: [{ description: '', quantity: 1, unit_price: 0 }],
    tax: 0,
    discount: 0,
    status: 'pending',
    notes: '',
    due_date: '',
    currency: 'USD'
  });
  const [previousCurrency, setPreviousCurrency] = useState('USD'); // Track previous currency for conversion
  const [editData, setEditData] = useState({
    status: '',
    notes: '',
    due_date: ''
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState(null);

  // Filter states
  const [filters, setFilters] = useState({
    status: 'all',
    search: '',
    startDate: '',
    endDate: ''
  });

  // Currency states
  const [currencies, setCurrencies] = useState(['USD', 'IDR', 'EUR', 'GBP', 'JPY', 'CNY', 'SGD', 'MYR']);
  const [selectedCurrency, setSelectedCurrency] = useState('USD');
  const [currencyRates, setCurrencyRates] = useState({});

  const subtotal = invoiceData.items.reduce((sum, item) => sum + (item.quantity * item.unit_price), 0);
  const taxAmount = subtotal * (invoiceData.tax / 100);
  const total = subtotal + taxAmount - invoiceData.discount;

  // Currency symbol helper
  const getCurrencySymbol = (currencyCode) => {
    const symbols = {
      'USD': '$',
      'IDR': 'Rp',
      'EUR': '€',
      'GBP': '£',
      'JPY': '¥',
      'CNY': '¥',
      'SGD': 'S$',
      'MYR': 'RM'
    };
    return symbols[currencyCode] || currencyCode;
  };

  // Format currency amount based on currency type
  const formatCurrencyAmount = (amount, currencyCode) => {
    // Currencies without decimals
    const noDecimalCurrencies = ['IDR', 'JPY', 'CNY'];
    
    if (noDecimalCurrencies.includes(currencyCode)) {
      // Format without decimals and with thousand separator
      return Math.round(amount).toLocaleString('en-US');
    } else {
      // Format with 2 decimals and thousand separator
      return amount.toLocaleString('en-US', {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2
      });
    }
  };

  useEffect(() => {
    loadCustomers();
    loadInvoices();
    loadCurrencyRates();
  }, []);

  useEffect(() => {
    applyFilters();
  }, [invoices, filters]);

  const loadCustomers = async () => {
    try {
      const response = await fetch(`${API_BASE}/customers`);
      const data = await response.json();
      setCustomers(data || []);
    } catch (error) {
      console.error('Error loading customers:', error);
    }
  };

  const loadInvoices = async () => {
    try {
      const response = await fetch(`${API_BASE}/invoices`);
      const data = await response.json();
      setInvoices(data || []);
    } catch (error) {
      console.error('Error loading invoices:', error);
    }
  };

  const loadCurrencyRates = async () => {
    try {
      const response = await fetch(`${API_BASE}/currency-rates`);
      const data = await response.json();
      const ratesMap = {};
      data.forEach(rate => {
        const key = `${rate.from_currency}_${rate.to_currency}`;
        ratesMap[key] = rate.rate;
      });
      setCurrencyRates(ratesMap);
    } catch (error) {
      console.error('Error loading currency rates:', error);
    }
  };

  const applyFilters = () => {
    let filtered = [...invoices];

    // Filter by status
    if (filters.status && filters.status !== 'all') {
      filtered = filtered.filter(inv => inv.status === filters.status);
    }

    // Filter by search term (customer name or invoice ID)
    if (filters.search) {
      const searchLower = filters.search.toLowerCase();
      filtered = filtered.filter(inv => {
        const customer = customers.find(c => c.id === inv.customer_id);
        const customerName = customer?.name?.toLowerCase() || '';
        const invoiceId = inv.id?.toLowerCase() || '';
        return customerName.includes(searchLower) || invoiceId.includes(searchLower);
      });
    }

    // Filter by date range
    if (filters.startDate) {
      filtered = filtered.filter(inv => new Date(inv.created_at) >= new Date(filters.startDate));
    }
    if (filters.endDate) {
      filtered = filtered.filter(inv => new Date(inv.created_at) <= new Date(filters.endDate + 'T23:59:59'));
    }

    setFilteredInvoices(filtered);
  };

  const convertAmount = (amount, fromCurrency, toCurrency) => {
    if (fromCurrency === toCurrency) return amount;
    const key = `${fromCurrency}_${toCurrency}`;
    const rate = currencyRates[key] || 1;
    return amount * rate;
  };

  // Handle currency change with auto-conversion of prices
  const handleCurrencyChange = (newCurrency) => {
    const oldCurrency = invoiceData.currency;
    
    // Convert all item prices
    const convertedItems = invoiceData.items.map(item => ({
      ...item,
      unit_price: convertAmount(item.unit_price, oldCurrency, newCurrency)
    }));

    // Convert discount
    const convertedDiscount = convertAmount(invoiceData.discount, oldCurrency, newCurrency);

    // Update invoice data with converted values
    setInvoiceData({
      ...invoiceData,
      currency: newCurrency,
      items: convertedItems,
      discount: convertedDiscount
    });
    
    setPreviousCurrency(newCurrency);
  };

  const addItem = () => {
    setInvoiceData({
      ...invoiceData,
      items: [...invoiceData.items, { description: '', quantity: 1, unit_price: 0 }]
    });
  };

  const updateItem = (index, field, value) => {
    const newItems = [...invoiceData.items];
    newItems[index][field] = field === 'quantity' || field === 'unit_price' ? parseFloat(value) || 0 : value;
    setInvoiceData({ ...invoiceData, items: newItems });
  };

  const removeItem = (index) => {
    if (invoiceData.items.length > 1) {
      const newItems = invoiceData.items.filter((_, i) => i !== index);
      setInvoiceData({ ...invoiceData, items: newItems });
    }
  };

  const createInvoice = async (e) => {
    e.preventDefault();

    if (invoiceData.customer_id === "") {
      setMessage({ type: 'error', text: '✗ Please select a customer' });
      return;
    }
    if (invoiceData.items.length === 0) {
      setMessage({ type: 'error', text: '✗ Please add at least one item' });
      return;
    }
    for (let item of invoiceData.items) {
      if (!item.description.trim()) {
        setMessage({ type: 'error', text: '✗ Please fill in all item descriptions' });
        return;
      }
      if (item.quantity <= 0) {
        setMessage({ type: 'error', text: '✗ Quantity must be greater than 0' });
        return;
      }
      if (item.unit_price <= 0) {
        setMessage({ type: 'error', text: '✗ Unit price must be greater than 0' });
        return;
      }
    }

    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/invoices`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(invoiceData)
      });
      if (response.ok) {
        const data = await response.json();
        setInvoiceData({
          customer_id: '',
          items: [{ description: '', quantity: 1, unit_price: 0 }],
          tax: 0,
          discount: 0,
          status: 'pending',
          notes: '',
          due_date: '',
          currency: 'USD'
        });
        setPreviousCurrency('USD'); // Reset previous currency
        setShowForm(false);

        if (data.pdf_data) {
          const pdfBlob = new Blob([new Uint8Array(atob(data.pdf_data).split('').map(c => c.charCodeAt(0)))], { type: 'application/pdf' });
          const url = URL.createObjectURL(pdfBlob);
          const a = document.createElement('a');
          a.href = url;
          a.download = `invoice_${data.id}.pdf`;
          document.body.appendChild(a);
          a.click();
          document.body.removeChild(a);
          URL.revokeObjectURL(url);
        }

        setMessage({ type: 'success', text: '✓ Invoice created and PDF downloaded!' });
        setTimeout(() => setMessage(null), 3000);
        loadInvoices(); // Reload invoices list
      } else {
        const errorText = await response.text();
        setMessage({ type: 'error', text: `✗ ${errorText}` });
      }
    } catch (error) {
      setMessage({ type: 'error', text: '✗ Error creating invoice' });
    }
    setLoading(false);
  };

  const openEditModal = (invoice) => {
    setSelectedInvoice(invoice);
    setEditData({
      status: invoice.status || 'pending',
      notes: invoice.notes || '',
      due_date: invoice.due_date ? invoice.due_date.split('T')[0] : ''
    });
    setShowEditModal(true);
  };

  const updateInvoice = async (e) => {
    e.preventDefault();
    setLoading(true);
    try {
      const response = await fetch(`${API_BASE}/invoices/${selectedInvoice.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(editData)
      });
      if (response.ok) {
        setShowEditModal(false);
        setMessage({ type: 'success', text: '✓ Invoice updated successfully!' });
        setTimeout(() => setMessage(null), 3000);
        loadInvoices();
      } else {
        const errorText = await response.text();
        setMessage({ type: 'error', text: `✗ ${errorText}` });
      }
    } catch (error) {
      setMessage({ type: 'error', text: '✗ Error updating invoice' });
    }
    setLoading(false);
  };

  const openPaymentModal = (invoice) => {
    setSelectedInvoice(invoice);
    setShowPaymentModal(true);
  };

  const openPaymentHistory = (invoice) => {
    setSelectedInvoice(invoice);
    setShowPaymentHistory(true);
  };

  const handlePaymentRecorded = async () => {
    setMessage({ type: 'success', text: '✓ Payment recorded successfully!' });
    setTimeout(() => setMessage(null), 3000);
    await loadInvoices();
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Invoices</h1>
          <p className="text-gray-600 dark:text-gray-400 mt-2">Create and manage your invoices</p>
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="px-6 py-3 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors flex items-center shadow-md"
        >
          <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
            <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
          </svg>
          Create Invoice
        </button>
      </div>

      {/* Alert Messages */}
      {message && (
        <div className={`p-4 rounded-lg font-semibold flex items-start ${
          message.type === 'success'
            ? 'bg-green-50 text-green-800 border border-green-200 dark:bg-green-900/20 dark:text-green-300 dark:border-green-500/30'
            : 'bg-red-50 text-red-800 border border-red-200 dark:bg-red-900/20 dark:text-red-300 dark:border-red-500/30'
        }`}>
          <div className={`mr-3 mt-0.5 ${message.type === 'success' ? 'text-green-500' : 'text-red-500'}`}>
            {message.type === 'success' ? (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            ) : (
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clipRule="evenodd" />
              </svg>
            )}
          </div>
          <div>{message.text}</div>
        </div>
      )}

      {/* Create Invoice Form */}
      {showForm && (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-6">Create New Invoice</h2>
          <form onSubmit={createInvoice} className="space-y-5">
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Select Customer *</label>
              <select
                value={invoiceData.customer_id}
                onChange={(e) => setInvoiceData({...invoiceData, customer_id: e.target.value})}
                required
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
              >
                <option value="">Choose a customer</option>
                {customers && customers.map(customer => (
                  <option key={customer.id} value={customer.id}>{customer.name}</option>
                ))}
              </select>
            </div>

            <div>
              <div className="flex justify-between items-center mb-3">
                <h3 className="font-semibold text-gray-900 dark:text-gray-100">Invoice Items</h3>
                <button
                  type="button"
                  onClick={addItem}
                  className="px-3 py-1.5 border border-transparent text-xs font-medium rounded text-blue-700 bg-purple-100 hover:bg-purple-200 transition-colors"
                >
                  + Add Item
                </button>
              </div>
              <div className="space-y-3">
                {invoiceData.items.map((item, index) => (
                  <div key={index} className="flex gap-3 p-3 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-700">
                    <div className="flex-1">
                      <input
                        type="text"
                        placeholder="Item description"
                        value={item.description}
                        onChange={(e) => updateItem(index, 'description', e.target.value)}
                        required
                        className="w-full px-3 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                      />
                    </div>
                    <div className="w-20">
                      <input
                        type="number"
                        placeholder="Qty"
                        value={item.quantity}
                        onChange={(e) => updateItem(index, 'quantity', e.target.value)}
                        min="1"
                        required
                        className="w-full px-2 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                      />
                    </div>
                    <div className="w-24">
                      <input
                        type="number"
                        placeholder="Price"
                        value={item.unit_price}
                        onChange={(e) => updateItem(index, 'unit_price', e.target.value)}
                        step="0.01"
                        min="0"
                        required
                        className="w-full px-2 py-2 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md text-sm focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                      />
                    </div>
                    {invoiceData.items.length > 1 && (
                      <button
                        type="button"
                        onClick={() => removeItem(index)}
                        className="px-2.5 py-2 bg-red-500 hover:bg-red-600 text-white rounded-md text-sm font-medium transition-colors"
                      >
                        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                          <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                        </svg>
                      </button>
                    )}
                  </div>
                ))}
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Status *</label>
                <select
                  value={invoiceData.status}
                  onChange={(e) => setInvoiceData({...invoiceData, status: e.target.value})}
                  required
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                >
                  <option value="pending">Pending</option>
                  <option value="paid">Paid</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Currency *</label>
                <select
                  value={invoiceData.currency}
                  onChange={(e) => handleCurrencyChange(e.target.value)}
                  required
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                >
                  {currencies.map(currency => (
                    <option key={currency} value={currency}>{currency}</option>
                  ))}
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Due Date</label>
                <input
                  type="date"
                  value={invoiceData.due_date}
                  onChange={(e) => setInvoiceData({...invoiceData, due_date: e.target.value})}
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                />
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tax (%)</label>
                <input
                  type="number"
                  placeholder="0"
                  value={invoiceData.tax}
                  onChange={(e) => setInvoiceData({...invoiceData, tax: parseFloat(e.target.value) || 0})}
                  min="0"
                  max="100"
                  step="0.1"
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Discount ({getCurrencySymbol(invoiceData.currency)})</label>
                <input
                  type="number"
                  placeholder="0.00"
                  value={invoiceData.discount}
                  onChange={(e) => setInvoiceData({...invoiceData, discount: parseFloat(e.target.value) || 0})}
                  min="0"
                  step="0.01"
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                />
              </div>
            </div>

            <div className="p-4 bg-gray-50 dark:bg-gray-700/50 rounded-lg border border-gray-200 dark:border-gray-700 space-y-2">
              <div className="flex justify-between items-center text-sm">
                <span className="text-gray-600 dark:text-gray-400">Subtotal</span>
                <span className="text-gray-900 dark:text-gray-100 font-medium">{getCurrencySymbol(invoiceData.currency)}{formatCurrencyAmount(subtotal, invoiceData.currency)}</span>
              </div>
              <div className="flex justify-between items-center text-sm">
                <span className="text-gray-600 dark:text-gray-400">Tax ({invoiceData.tax}%)</span>
                <span className="text-gray-900 dark:text-gray-100 font-medium">{getCurrencySymbol(invoiceData.currency)}{formatCurrencyAmount(taxAmount, invoiceData.currency)}</span>
              </div>
              <div className="flex justify-between items-center text-sm">
                <span className="text-gray-600 dark:text-gray-400">Discount</span>
                <span className="text-gray-900 dark:text-gray-100 font-medium">-{getCurrencySymbol(invoiceData.currency)}{formatCurrencyAmount(invoiceData.discount, invoiceData.currency)}</span>
              </div>
              <div className="border-t border-gray-300 dark:border-gray-600 pt-2 mt-2"></div>
              <div className="flex justify-between items-center">
                <span className="text-gray-700 dark:text-gray-300 font-semibold">Total Amount</span>
                <span className="text-2xl font-bold text-gray-900 dark:text-gray-100">{getCurrencySymbol(invoiceData.currency)}{formatCurrencyAmount(total, invoiceData.currency)}</span>
              </div>
            </div>

            <div className="flex gap-3 pt-2">
              <button
                type="submit"
                disabled={loading}
                className="flex-1 py-3 px-4 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors disabled:opacity-50 flex items-center justify-center"
              >
                {loading ? 'Creating...' : (
                  <>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                      <path fillRule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd" />
                    </svg>
                    Generate & Download PDF
                  </>
                )}
              </button>
              <button
                type="button"
                onClick={() => {
                  setShowForm(false);
                  setPreviousCurrency('USD'); // Reset previous currency on cancel
                  setInvoiceData({
                    customer_id: '',
                    items: [{ description: '', quantity: 1, unit_price: 0 }],
                    tax: 0,
                    discount: 0,
                    status: 'pending',
                    notes: '',
                    due_date: '',
                    currency: 'USD'
                  });
                }}
                className="px-6 py-3 bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium rounded-lg transition-colors dark:bg-gray-600 dark:text-gray-300 dark:hover:bg-gray-500"
              >
                Cancel
              </button>
            </div>
          </form>
        </div>
      )}

      {/* Filter & Search Section - Hidden when creating invoice */}
      {!showForm && (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6 mb-6">
          <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">Filter & Search</h2>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Status</label>
            <select
              value={filters.status}
              onChange={(e) => setFilters({...filters, status: e.target.value})}
              className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:text-white"
            >
              <option value="all">All Status</option>
              <option value="pending">Pending</option>
              <option value="paid">Paid</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Search</label>
            <input
              type="text"
              placeholder="Search by customer or ID..."
              value={filters.search}
              onChange={(e) => setFilters({...filters, search: e.target.value})}
              className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:text-white dark:placeholder-gray-400"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">From Date</label>
            <input
              type="date"
              value={filters.startDate}
              onChange={(e) => setFilters({...filters, startDate: e.target.value})}
              className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:text-white"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">To Date</label>
            <input
              type="date"
              value={filters.endDate}
              onChange={(e) => setFilters({...filters, endDate: e.target.value})}
              className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:text-white"
            />
          </div>
        </div>
        <div className="mt-4 flex items-center justify-between">
          <div className="text-sm text-gray-600 dark:text-gray-400">
            Showing {filteredInvoices.length} of {invoices.length} invoices
          </div>
          <div className="flex items-center gap-2">
            <label className="text-sm font-medium text-gray-700 dark:text-gray-300">View in:</label>
            <select
              value={selectedCurrency}
              onChange={(e) => setSelectedCurrency(e.target.value)}
              className="px-3 py-1.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg text-sm focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:text-white"
            >
              {currencies.map(currency => (
                <option key={currency} value={currency}>{currency}</option>
              ))}
            </select>
          </div>
        </div>
      </div>
      )}

      {/* Invoices List - Hidden when creating invoice */}
      {!showForm && (
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700">
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Recent Invoices</h2>
        </div>
        <div className="p-6">
          {filteredInvoices && filteredInvoices.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead className="bg-gray-50 dark:bg-gray-700/50">
                  <tr>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Invoice #</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Customer</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Total</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Payment Status</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Items</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Created</th>
                    <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody className="bg-white dark:bg-gray-800 divide-y divide-gray-200 dark:divide-gray-700">
                  {filteredInvoices.map((invoice) => {
                    const invoiceCurrency = invoice.currency || 'USD';
                    const convertedTotal = convertAmount(invoice.total, invoiceCurrency, selectedCurrency);
                    
                    return (
                    <tr key={invoice.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-mono font-semibold text-indigo-600 dark:text-indigo-400">
                        {invoice.invoice_number || `INV-${invoice.id.substring(0, 8)}`}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 dark:text-gray-300">
                        {customers.find(c => c.id === invoice.customer_id)?.name || invoice.customer_id.substring(0, 8)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm font-semibold text-gray-900 dark:text-gray-100">
                        {getCurrencySymbol(selectedCurrency)}{formatCurrencyAmount(convertedTotal, selectedCurrency)}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap">
                        <span className={`px-2 py-1 text-xs font-semibold rounded-full ${
                          invoice.payment_status === 'paid' ? 'bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300' :
                          invoice.payment_status === 'overdue' ? 'bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300' :
                          invoice.payment_status === 'partially_paid' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/50 dark:text-blue-300' :
                          'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/50 dark:text-yellow-300'
                        }`}>
                          {invoice.payment_status ? 
                            invoice.payment_status.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ') 
                            : 'Unpaid'}
                        </span>
                        {invoice.payment_status === 'partially_paid' && (
                          <div className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                            {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(invoice.paid_amount || 0, invoice.currency)} / {getCurrencySymbol(invoice.currency)}{formatCurrencyAmount(invoice.total, invoice.currency)}
                          </div>
                        )}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        {Array.isArray(invoice.items) ? invoice.items.length : 0} item(s)
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                        {new Date(invoice.created_at).toLocaleDateString('en-US', { 
                          year: 'numeric', 
                          month: 'short', 
                          day: 'numeric',
                          hour: '2-digit',
                          minute: '2-digit'
                        })}
                      </td>
                      <td className="px-6 py-4 whitespace-nowrap text-sm">
                        <div className="flex items-center gap-2">
                          {/* Record Payment Button - Only for unpaid/partially_paid/overdue */}
                          {(invoice.payment_status === 'unpaid' || invoice.payment_status === 'partially_paid' || invoice.payment_status === 'overdue') && (
                            <button
                              onClick={() => openPaymentModal(invoice)}
                              className="px-3 py-1.5 bg-green-600 hover:bg-green-700 text-white rounded-lg font-medium transition-colors text-xs"
                              title="Record Payment"
                            >
                              💰 Pay
                            </button>
                          )}
                          
                          {/* Payment History Button - Always show if there are payments */}
                          {(invoice.paid_amount > 0 || invoice.payment_status === 'paid') && (
                            <button
                              onClick={() => openPaymentHistory(invoice)}
                              className="px-3 py-1.5 bg-blue-100 hover:bg-blue-200 dark:bg-blue-900/50 dark:hover:bg-blue-900 text-blue-700 dark:text-blue-300 rounded-lg font-medium transition-colors text-xs"
                              title="View Payment History"
                            >
                              📜 History
                            </button>
                          )}
                          
                          {/* Edit Button */}
                          <button
                            onClick={() => openEditModal(invoice)}
                            className="text-indigo-600 hover:text-indigo-900 dark:text-indigo-400 dark:hover:text-indigo-300 font-medium"
                          >
                            ✏️ Edit
                          </button>
                        </div>
                      </td>
                    </tr>
                  );
                  })}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="text-center py-12">
              <svg xmlns="http://www.w3.org/2000/svg" className="h-16 w-16 mx-auto text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p className="mt-4 text-gray-500 text-lg">No invoices yet</p>
              <p className="text-sm text-gray-400">Create your first invoice to get started</p>
            </div>
          )}
        </div>
      </div>
      )}

      {/* Edit Invoice Modal */}
      {showEditModal && selectedInvoice && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-md w-full p-6">
            <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">Edit Invoice</h2>
            <form onSubmit={updateInvoice} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Invoice ID</label>
                <input
                  type="text"
                  value={selectedInvoice.id.substring(0, 8) + '...'}
                  disabled
                  className="w-full px-4 py-2.5 bg-gray-100 dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg text-gray-600 dark:text-gray-400"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Status *</label>
                <select
                  value={editData.status}
                  onChange={(e) => setEditData({...editData, status: e.target.value})}
                  required
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                >
                  <option value="pending">Pending</option>
                  <option value="paid">Paid</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Due Date</label>
                <input
                  type="date"
                  value={editData.due_date}
                  onChange={(e) => setEditData({...editData, due_date: e.target.value})}
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                />
              </div>
              <div>
                <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Notes</label>
                <textarea
                  value={editData.notes}
                  onChange={(e) => setEditData({...editData, notes: e.target.value})}
                  rows="3"
                  className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 dark:placeholder-gray-400 dark:text-white"
                  placeholder="Add notes here..."
                />
              </div>
              <div className="flex gap-3 pt-2">
                <button
                  type="submit"
                  disabled={loading}
                  className="flex-1 py-2.5 px-4 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors disabled:opacity-50"
                >
                  {loading ? 'Updating...' : 'Update Invoice'}
                </button>
                <button
                  type="button"
                  onClick={() => setShowEditModal(false)}
                  className="px-6 py-2.5 bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium rounded-lg transition-colors dark:bg-gray-600 dark:text-gray-300 dark:hover:bg-gray-500"
                >
                  Cancel
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Payment Modal */}
      {showPaymentModal && selectedInvoice && (
        <PaymentModal
          invoice={selectedInvoice}
          onClose={() => setShowPaymentModal(false)}
          onPaymentRecorded={handlePaymentRecorded}
        />
      )}

      {/* Payment History Modal */}
      {showPaymentHistory && selectedInvoice && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="p-6 border-b border-gray-200 dark:border-gray-700 flex items-center justify-between sticky top-0 bg-white dark:bg-gray-800">
              <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100">
                Payment History - {selectedInvoice.invoice_number || `INV-${selectedInvoice.id.substring(0, 8)}`}
              </h2>
              <button
                onClick={() => setShowPaymentHistory(false)}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <div className="p-6">
              <PaymentHistory invoice={selectedInvoice} />
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default Invoices;

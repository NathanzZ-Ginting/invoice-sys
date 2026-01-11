import React, { useState } from 'react';

function Settings() {
  const [companyInfo, setCompanyInfo] = useState({
    name: 'InvoicePro Company',
    email: 'contact@invoicepro.com',
    phone: '+1 (555) 123-4567',
    address: '123 Business Street',
    city: 'New York',
    postal_code: '10001',
    country: 'United States',
    website: 'www.invoicepro.com',
    tax_id: ''
  });
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState(null);

  const handleChange = (e) => {
    setCompanyInfo({
      ...companyInfo,
      [e.target.name]: e.target.value
    });
  };

  const saveSettings = async (e) => {
    e.preventDefault();
    setLoading(true);
    
    // Simulate API call since backend endpoint doesn't exist yet
    setTimeout(() => {
      setMessage({ type: 'success', text: '✓ Settings saved successfully!' });
      setLoading(false);
      setTimeout(() => setMessage(null), 3000);
    }, 500);
  };

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100">Settings</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-2">Configure your company information for invoices</p>
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

      {/* Company Information Form */}
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-md border border-gray-200 dark:border-gray-700">
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Company Information</h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">This information will appear on all your invoices</p>
        </div>
        
        <form onSubmit={saveSettings} className="p-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            {/* Company Name */}
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Company Name *
              </label>
              <input
                type="text"
                name="name"
                value={companyInfo.name}
                onChange={handleChange}
                required
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="Your Company Name"
              />
            </div>

            {/* Email */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Email Address *
              </label>
              <input
                type="email"
                name="email"
                value={companyInfo.email}
                onChange={handleChange}
                required
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="company@example.com"
              />
            </div>

            {/* Phone */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Phone Number
              </label>
              <input
                type="text"
                name="phone"
                value={companyInfo.phone}
                onChange={handleChange}
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="+1 (555) 123-4567"
              />
            </div>

            {/* Address */}
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Street Address
              </label>
              <input
                type="text"
                name="address"
                value={companyInfo.address}
                onChange={handleChange}
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="123 Business Street"
              />
            </div>

            {/* City */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                City
              </label>
              <input
                type="text"
                name="city"
                value={companyInfo.city}
                onChange={handleChange}
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="New York"
              />
            </div>

            {/* Postal Code */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Postal Code
              </label>
              <input
                type="text"
                name="postal_code"
                value={companyInfo.postal_code}
                onChange={handleChange}
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="10001"
              />
            </div>

            {/* Country */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Country
              </label>
              <input
                type="text"
                name="country"
                value={companyInfo.country}
                onChange={handleChange}
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="United States"
              />
            </div>

            {/* Website */}
            <div>
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Website
              </label>
              <input
                type="text"
                name="website"
                value={companyInfo.website}
                onChange={handleChange}
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="www.yourcompany.com"
              />
            </div>

            {/* Tax ID */}
            <div className="md:col-span-2">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                Tax ID / VAT Number
              </label>
              <input
                type="text"
                name="tax_id"
                value={companyInfo.tax_id}
                onChange={handleChange}
                className="w-full px-4 py-2.5 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:placeholder-gray-400 dark:text-white"
                placeholder="Enter your tax identification number"
              />
            </div>
          </div>

          {/* Save Button */}
          <div className="mt-6 flex justify-end">
            <button
              type="submit"
              disabled={loading}
              className="px-6 py-3 bg-blue-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors shadow-md disabled:opacity-50 flex items-center"
            >
              {loading ? 'Saving...' : (
                <>
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M7.707 10.293a1 1 0 10-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 11.586V6h5a2 2 0 012 2v7a2 2 0 01-2 2H4a2 2 0 01-2-2V8a2 2 0 012-2h5v5.586l-1.293-1.293zM9 4a1 1 0 012 0v2H9V4z" />
                  </svg>
                  Save Settings
                </>
              )}
            </button>
          </div>
        </form>
      </div>

      {/* Preview Section */}
      <div className="bg-white dark:bg-gray-800 rounded-xl shadow-md border border-gray-200 dark:border-gray-700">
        <div className="p-6 border-b border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100">Invoice Preview</h2>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">How your information will appear on invoices</p>
        </div>
        
        <div className="p-6">
          <div className="border-2 border-dashed border-gray-300 dark:border-gray-600 rounded-lg p-6 bg-gray-50 dark:bg-gray-700/50">
            <div className="flex items-start justify-between">
              <div>
                <div className="w-16 h-16 bg-blue-500 rounded-lg flex items-center justify-center mb-3">
                  <span className="text-white text-2xl font-bold">
                    {companyInfo.name.substring(0, 2).toUpperCase()}
                  </span>
                </div>
                <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100">{companyInfo.name}</h3>
              </div>
              <div className="text-right text-sm text-gray-600 dark:text-gray-400 space-y-1">
                <p>{companyInfo.email}</p>
                <p>{companyInfo.phone}</p>
                <p>{companyInfo.address}</p>
                <p>{companyInfo.city}, {companyInfo.postal_code}</p>
                <p>{companyInfo.country}</p>
                {companyInfo.website && <p>{companyInfo.website}</p>}
                {companyInfo.tax_id && <p>Tax ID: {companyInfo.tax_id}</p>}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Settings;

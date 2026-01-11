-- =====================================================
-- ADD FILTER & MULTI-CURRENCY SUPPORT
-- =====================================================

-- Add currency field to invoices table
ALTER TABLE invoices 
ADD COLUMN IF NOT EXISTS currency VARCHAR(3) DEFAULT 'USD';

-- Create currency_rates table for exchange rates
CREATE TABLE IF NOT EXISTS currency_rates (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    from_currency VARCHAR(3) NOT NULL,
    to_currency VARCHAR(3) NOT NULL,
    rate DECIMAL(20,6) NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(from_currency, to_currency)
);

-- Create index for faster queries
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_currency ON invoices(currency);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);
CREATE INDEX IF NOT EXISTS idx_currency_rates_pair ON currency_rates(from_currency, to_currency);

-- Enable RLS for currency_rates
ALTER TABLE currency_rates ENABLE ROW LEVEL SECURITY;

-- Policy for currency_rates
DROP POLICY IF EXISTS "Allow all operations on currency_rates" ON currency_rates;
CREATE POLICY "Allow all operations on currency_rates" 
ON currency_rates FOR ALL USING (true);

-- Insert sample currency rates (base currency: USD)
INSERT INTO currency_rates (from_currency, to_currency, rate) VALUES
    ('USD', 'USD', 1.000000),
    ('USD', 'IDR', 15800.000000),
    ('USD', 'EUR', 0.920000),
    ('USD', 'GBP', 0.790000),
    ('USD', 'JPY', 148.500000),
    ('USD', 'CNY', 7.240000),
    ('USD', 'SGD', 1.340000),
    ('USD', 'MYR', 4.720000),
    
    ('IDR', 'USD', 0.000063),
    ('IDR', 'IDR', 1.000000),
    ('IDR', 'EUR', 0.000058),
    ('IDR', 'GBP', 0.000050),
    
    ('EUR', 'USD', 1.087000),
    ('EUR', 'IDR', 17174.000000),
    ('EUR', 'EUR', 1.000000),
    ('EUR', 'GBP', 0.859000),
    
    ('GBP', 'USD', 1.266000),
    ('GBP', 'IDR', 20003.000000),
    ('GBP', 'EUR', 1.164000),
    ('GBP', 'GBP', 1.000000)
ON CONFLICT (from_currency, to_currency) 
DO UPDATE SET rate = EXCLUDED.rate, updated_at = NOW();

-- Add comment to explain the feature
COMMENT ON COLUMN invoices.currency IS 'Currency code (ISO 4217) for the invoice';
COMMENT ON TABLE currency_rates IS 'Exchange rates for currency conversion';

-- =====================================================
-- INVOICE SYSTEM - COMPLETE DATABASE SETUP
-- Run this entire file to set up the database from scratch
-- =====================================================

-- Drop existing tables if they exist (use with caution!)
DROP TABLE IF EXISTS invoices CASCADE;
DROP TABLE IF EXISTS customers CASCADE;
DROP TABLE IF EXISTS company_info CASCADE;

-- =====================================================
-- TABLE: customers
-- =====================================================
CREATE TABLE customers (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    phone TEXT,
    address TEXT,
    city TEXT,
    postal_code TEXT,
    country TEXT DEFAULT 'Indonesia',
    company_name TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =====================================================
-- TABLE: invoices
-- =====================================================
CREATE TABLE invoices (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    subtotal DECIMAL(10,2) NOT NULL,
    tax DECIMAL(10,2) DEFAULT 0,
    discount DECIMAL(10,2) DEFAULT 0,
    total DECIMAL(10,2) NOT NULL,
    items JSONB NOT NULL,
    pdf_url TEXT,
    status TEXT DEFAULT 'pending',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    due_date TIMESTAMP WITH TIME ZONE
);

-- =====================================================
-- TABLE: company_info
-- =====================================================
CREATE TABLE company_info (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT,
    phone TEXT,
    address TEXT,
    city TEXT,
    postal_code TEXT,
    country TEXT,
    website TEXT,
    tax_id TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- =====================================================
-- INDEXES FOR PERFORMANCE
-- =====================================================
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_customers_created_at ON customers(created_at);
CREATE INDEX idx_invoices_customer_id ON invoices(customer_id);
CREATE INDEX idx_invoices_created_at ON invoices(created_at);
CREATE INDEX idx_invoices_status ON invoices(status);

-- =====================================================
-- ROW LEVEL SECURITY (RLS)
-- =====================================================
ALTER TABLE customers ENABLE ROW LEVEL SECURITY;
ALTER TABLE invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE company_info ENABLE ROW LEVEL SECURITY;

-- =====================================================
-- POLICIES (Allow all operations for development)
-- =====================================================
CREATE POLICY "Allow all operations on customers" 
ON customers FOR ALL USING (true);

CREATE POLICY "Allow all operations on invoices" 
ON invoices FOR ALL USING (true);

CREATE POLICY "Allow all operations on company_info" 
ON company_info FOR ALL USING (true);

-- =====================================================
-- INSERT SAMPLE DATA (Optional)
-- =====================================================

-- Sample customer
INSERT INTO customers (name, email, phone, address, city, postal_code, country, company_name)
VALUES 
    ('John Doe', 'john@example.com', '081234567890', 'Jl. Contoh No. 123', 'Jakarta', '12345', 'Indonesia', 'PT Example'),
    ('Jane Smith', 'jane@example.com', '081234567891', 'Jl. Sample No. 456', 'Bandung', '54321', 'Indonesia', 'CV Sample');

-- Sample company info
INSERT INTO company_info (name, email, phone, address, city, postal_code, country, website, tax_id)
VALUES 
    ('Your Company Name', 'info@yourcompany.com', '021-12345678', 'Jl. Company No. 789', 'Jakarta', '10110', 'Indonesia', 'www.yourcompany.com', '01.234.567.8-901.000');

-- =====================================================
-- VERIFICATION QUERIES
-- =====================================================

-- Check all tables
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
ORDER BY table_name;

-- Check customers table structure
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns 
WHERE table_name = 'customers' 
ORDER BY ordinal_position;

-- Check invoices table structure
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns 
WHERE table_name = 'invoices' 
ORDER BY ordinal_position;

-- Check company_info table structure
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns 
WHERE table_name = 'company_info' 
ORDER BY ordinal_position;

-- Count records
SELECT 
    (SELECT COUNT(*) FROM customers) as customers_count,
    (SELECT COUNT(*) FROM invoices) as invoices_count,
    (SELECT COUNT(*) FROM company_info) as company_info_count;

-- =====================================================
-- SETUP COMPLETE!
-- =====================================================
-- Next steps:
-- 1. Update your .env file with database credentials
-- 2. Run backend: cd backend && go run main.go
-- 3. Run frontend: cd frontend && npm start
-- =====================================================

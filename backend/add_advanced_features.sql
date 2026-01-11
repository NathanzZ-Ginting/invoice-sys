-- ============================================
-- Advanced Features Migration
-- Features: Invoice Numbering, Payment Recording, Payment Reminders
-- ============================================

-- 1. Add invoice_number column to invoices table
ALTER TABLE invoices 
ADD COLUMN IF NOT EXISTS invoice_number VARCHAR(50) UNIQUE;

-- 2. Add payment tracking columns to invoices table
ALTER TABLE invoices
ADD COLUMN IF NOT EXISTS payment_status VARCHAR(20) DEFAULT 'unpaid',
ADD COLUMN IF NOT EXISTS paid_amount DECIMAL(20,2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS payment_date TIMESTAMP,
ADD COLUMN IF NOT EXISTS last_reminder_sent TIMESTAMP;

-- Add check constraint for payment_status
ALTER TABLE invoices
DROP CONSTRAINT IF EXISTS check_payment_status;

ALTER TABLE invoices
ADD CONSTRAINT check_payment_status 
CHECK (payment_status IN ('unpaid', 'partially_paid', 'paid', 'overdue'));

-- 3. Create payments table for payment history
CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    amount DECIMAL(20,2) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_date TIMESTAMP NOT NULL DEFAULT NOW(),
    reference_number VARCHAR(100),
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    created_by VARCHAR(100)
);

-- 4. Create invoice_settings table for numbering configuration
CREATE TABLE IF NOT EXISTS invoice_settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    setting_key VARCHAR(100) UNIQUE NOT NULL,
    setting_value TEXT NOT NULL,
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Insert default invoice numbering settings
INSERT INTO invoice_settings (setting_key, setting_value) VALUES
    ('invoice_prefix', 'INV'),
    ('invoice_number_format', '{PREFIX}-{YEAR}-{NUMBER}'),
    ('invoice_counter', '1'),
    ('invoice_counter_reset', 'yearly'),
    ('invoice_counter_padding', '4')
ON CONFLICT (setting_key) DO NOTHING;

-- 5. Create payment_reminders table for tracking reminders
CREATE TABLE IF NOT EXISTS payment_reminders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID REFERENCES invoices(id) ON DELETE CASCADE,
    reminder_type VARCHAR(50) NOT NULL, -- 'before_due', 'on_due', 'after_due'
    days_offset INTEGER NOT NULL, -- -3, 0, 3, 7, etc.
    sent_at TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending', -- 'pending', 'sent', 'failed'
    email_to VARCHAR(255),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 6. Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_invoices_invoice_number ON invoices(invoice_number);
CREATE INDEX IF NOT EXISTS idx_invoices_payment_status ON invoices(payment_status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);
CREATE INDEX IF NOT EXISTS idx_payments_invoice_id ON payments(invoice_id);
CREATE INDEX IF NOT EXISTS idx_payments_payment_date ON payments(payment_date);
CREATE INDEX IF NOT EXISTS idx_payment_reminders_invoice_id ON payment_reminders(invoice_id);
CREATE INDEX IF NOT EXISTS idx_payment_reminders_status ON payment_reminders(status);

-- 7. Create view for invoice analytics
CREATE OR REPLACE VIEW invoice_analytics AS
SELECT 
    DATE_TRUNC('day', created_at) as date,
    COUNT(*) as invoice_count,
    SUM(total) as total_amount,
    SUM(CASE WHEN payment_status = 'paid' THEN total ELSE 0 END) as paid_amount,
    SUM(CASE WHEN payment_status = 'unpaid' THEN total ELSE 0 END) as unpaid_amount,
    SUM(CASE WHEN payment_status = 'overdue' THEN total ELSE 0 END) as overdue_amount,
    currency
FROM invoices
GROUP BY DATE_TRUNC('day', created_at), currency;

-- 8. Create view for payment summary
CREATE OR REPLACE VIEW payment_summary AS
SELECT 
    i.id as invoice_id,
    i.invoice_number,
    i.total,
    i.paid_amount,
    i.total - i.paid_amount as outstanding_amount,
    i.payment_status,
    i.due_date,
    c.name as customer_name,
    c.email as customer_email,
    CASE 
        WHEN i.due_date < NOW() AND i.payment_status != 'paid' THEN 'overdue'
        ELSE i.payment_status
    END as current_status
FROM invoices i
LEFT JOIN customers c ON i.customer_id = c.id;

-- 9. Update existing invoices with default payment status
UPDATE invoices 
SET payment_status = CASE 
    WHEN status = 'paid' THEN 'paid'
    WHEN due_date < NOW() AND status != 'paid' THEN 'overdue'
    ELSE 'unpaid'
END
WHERE payment_status IS NULL;

-- 10. Function to auto-generate invoice number
CREATE OR REPLACE FUNCTION generate_invoice_number()
RETURNS TEXT AS $$
DECLARE
    v_prefix TEXT;
    v_format TEXT;
    v_counter INTEGER;
    v_padding INTEGER;
    v_year TEXT;
    v_number TEXT;
    v_invoice_number TEXT;
BEGIN
    -- Get settings
    SELECT setting_value INTO v_prefix FROM invoice_settings WHERE setting_key = 'invoice_prefix';
    SELECT setting_value INTO v_format FROM invoice_settings WHERE setting_key = 'invoice_number_format';
    SELECT setting_value::INTEGER INTO v_counter FROM invoice_settings WHERE setting_key = 'invoice_counter';
    SELECT setting_value::INTEGER INTO v_padding FROM invoice_settings WHERE setting_key = 'invoice_counter_padding';
    
    -- Get current year
    v_year := TO_CHAR(NOW(), 'YYYY');
    
    -- Pad the counter
    v_number := LPAD(v_counter::TEXT, v_padding, '0');
    
    -- Build invoice number
    v_invoice_number := v_format;
    v_invoice_number := REPLACE(v_invoice_number, '{PREFIX}', v_prefix);
    v_invoice_number := REPLACE(v_invoice_number, '{YEAR}', v_year);
    v_invoice_number := REPLACE(v_invoice_number, '{NUMBER}', v_number);
    
    -- Increment counter
    UPDATE invoice_settings 
    SET setting_value = (v_counter + 1)::TEXT 
    WHERE setting_key = 'invoice_counter';
    
    RETURN v_invoice_number;
END;
$$ LANGUAGE plpgsql;

-- 11. Trigger to auto-update payment_status based on due_date
CREATE OR REPLACE FUNCTION update_payment_status()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.due_date < NOW() AND NEW.payment_status NOT IN ('paid', 'partially_paid') THEN
        NEW.payment_status := 'overdue';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_payment_status
BEFORE INSERT OR UPDATE ON invoices
FOR EACH ROW
EXECUTE FUNCTION update_payment_status();

-- 12. Comments for documentation
COMMENT ON TABLE payments IS 'Stores payment history for invoices - supports partial payments';
COMMENT ON TABLE payment_reminders IS 'Tracks payment reminder emails sent to customers';
COMMENT ON TABLE invoice_settings IS 'Configurable settings for invoice numbering and other features';
COMMENT ON COLUMN invoices.invoice_number IS 'Auto-generated unique invoice number (e.g., INV-2026-0001)';
COMMENT ON COLUMN invoices.payment_status IS 'Current payment status: unpaid, partially_paid, paid, overdue';
COMMENT ON COLUMN invoices.paid_amount IS 'Total amount paid so far (for partial payments)';

-- Success message
DO $$
BEGIN
    RAISE NOTICE '✅ Advanced features migration completed successfully!';
    RAISE NOTICE '📊 New tables: payments, payment_reminders, invoice_settings';
    RAISE NOTICE '🔢 Invoice numbering: Auto-enabled with format INV-YYYY-####';
    RAISE NOTICE '💳 Payment tracking: Ready for payment recording';
    RAISE NOTICE '🔔 Reminder system: Infrastructure ready';
END $$;

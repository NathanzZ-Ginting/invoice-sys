-- Migration: Add tax, discount, and subtotal columns to invoices table
-- Run this in Supabase SQL Editor if you already have existing invoices table

-- Add new columns
ALTER TABLE invoices 
  ADD COLUMN IF NOT EXISTS subtotal DECIMAL(10,2),
  ADD COLUMN IF NOT EXISTS tax DECIMAL(10,2) DEFAULT 0,
  ADD COLUMN IF NOT EXISTS discount DECIMAL(10,2) DEFAULT 0;

-- Update existing records: set subtotal = total, keep total as is
UPDATE invoices 
SET subtotal = total 
WHERE subtotal IS NULL;

-- Make subtotal NOT NULL after setting values
ALTER TABLE invoices 
  ALTER COLUMN subtotal SET NOT NULL;

-- Verify the changes
SELECT column_name, data_type, is_nullable, column_default
FROM information_schema.columns 
WHERE table_name = 'invoices' 
ORDER BY ordinal_position;

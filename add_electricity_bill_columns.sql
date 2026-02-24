-- Add electricity bill columns to payment table
ALTER TABLE payment ADD COLUMN electricity_bill DECIMAL(10,2) NULL;
ALTER TABLE payment ADD COLUMN paid_bill DECIMAL(10,2) NULL; 
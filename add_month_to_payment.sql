-- Add month field to payment table
ALTER TABLE payment ADD COLUMN month VARCHAR(20) DEFAULT NULL AFTER due_rent; 
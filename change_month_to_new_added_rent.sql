-- Migration: Change month column to new_added_rent with integer type
-- This script modifies the payment table to replace the month column with new_added_rent

-- Step 1: Add the new column
ALTER TABLE payment ADD COLUMN new_added_rent INT NULL AFTER rent;

-- Step 2: Update existing records to set new_added_rent based on rent value
-- For existing records, we'll set new_added_rent = rent since that represents the amount added
UPDATE payment SET new_added_rent = rent WHERE new_added_rent IS NULL;

-- Step 3: Drop the old month column
ALTER TABLE payment DROP COLUMN month;

-- Step 4: Verify the changes
DESCRIBE payment; 
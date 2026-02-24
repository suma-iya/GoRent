-- Update advance table structure
-- Remove month column and ensure id is unique with random numbers

-- First, let's check if the month column exists and remove it
ALTER TABLE advance DROP COLUMN IF EXISTS month;

-- Update the table structure to ensure id is unique
-- The id column should already be set as PRIMARY KEY which makes it unique
-- If you need to regenerate existing IDs with random numbers, you can do:

-- Note: This is a destructive operation that will change existing IDs
-- Only run this if you want to regenerate all existing advance payment IDs
-- UPDATE advance SET id = FLOOR(RAND() * 900000000) + 100000000 WHERE id IS NOT NULL;

-- Verify the table structure
DESCRIBE advance; 
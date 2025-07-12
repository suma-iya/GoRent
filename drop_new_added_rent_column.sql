-- Drop new_added_rent column from payment table
-- Since we're using the rent field for new added rent values

ALTER TABLE payment DROP COLUMN new_added_rent; 
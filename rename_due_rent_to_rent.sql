-- Rename due_rent column to rent in payment table
ALTER TABLE payment CHANGE COLUMN due_rent rent int(11) NULL; 
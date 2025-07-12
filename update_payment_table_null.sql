-- Update payment table to allow NULL values for due_rent and received_money
ALTER TABLE payment MODIFY COLUMN due_rent int(11) NULL;
ALTER TABLE payment MODIFY COLUMN recieved_money int(11) NULL; 
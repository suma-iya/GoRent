-- Add photo column to property table
ALTER TABLE property ADD COLUMN photo VARCHAR(255) NULL;

-- Add comment to explain the column
COMMENT ON COLUMN property.photo IS 'URL or file path to property photo'; 
-- Add FCM token column to user table
ALTER TABLE user ADD COLUMN fcm_token VARCHAR(255) NULL;

-- Add index for better performance
CREATE INDEX idx_user_fcm_token ON user(fcm_token); 
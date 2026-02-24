-- MySQL Setup Script
-- Run this as root user: mysql -u root -p < setup_mysql_user.sql

-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS rent;

-- Create or update the user 'suma' with the password from config
CREATE USER IF NOT EXISTS 'suma'@'localhost' IDENTIFIED BY 'tMyc6mApj]wgzHl7';

-- Grant all privileges on the rent database
GRANT ALL PRIVILEGES ON rent.* TO 'suma'@'localhost';

-- If user already exists, update the password
ALTER USER 'suma'@'localhost' IDENTIFIED BY 'tMyc6mApj]wgzHl7';

-- Flush privileges to apply changes
FLUSH PRIVILEGES;

-- Show the user to confirm
SELECT User, Host FROM mysql.user WHERE User='suma';




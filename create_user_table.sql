-- Create user table for rent database
USE rent;

CREATE TABLE IF NOT EXISTS user (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(255) NULL,
    NID VARCHAR(20) NULL,
    password VARCHAR(255) NOT NULL,
    manager BOOLEAN NULL DEFAULT FALSE,
    fcm_token VARCHAR(255) NULL,
    created_at DATE NOT NULL,
    created_by BIGINT NOT NULL,
    updated_at DATE NOT NULL,
    updated_by BIGINT NOT NULL,
    INDEX idx_phone_number (phone_number),
    INDEX idx_fcm_token (fcm_token)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;



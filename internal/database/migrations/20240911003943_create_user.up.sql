-- Create the extension for UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create the enum type for role
CREATE TYPE user_role AS ENUM ('user', 'admin');

-- Create the users table with role column using the enum type
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    verified_at TIMESTAMP,
    otp VARCHAR(6),
    otp_expired_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
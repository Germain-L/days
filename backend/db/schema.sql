-- Database schema for the days project

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Days table - representing individual days/entries
CREATE TABLE days (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    title VARCHAR(255),
    description TEXT,
    mood INTEGER CHECK (mood >= 1 AND mood <= 10), -- 1-10 mood scale
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, date) -- One entry per user per day
);

-- Create indexes for better performance
CREATE INDEX idx_days_user_id ON days(user_id);
CREATE INDEX idx_days_date ON days(date);
CREATE INDEX idx_users_email ON users(email);

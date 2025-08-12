-- Required for gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Users table for authentication
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Calendars table - each user can have multiple calendars
CREATE TABLE calendars (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, name)
);

-- Color meanings for each calendar (e.g., red = stressed, green = relaxed)
CREATE TABLE color_meanings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    color_hex VARCHAR(7) NOT NULL, -- e.g., "#FF0000"
    meaning VARCHAR(50) NOT NULL,   -- e.g., "stressed", "relaxed", "good", "bad"
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Day entries - tracks a specific day for a specific calendar
CREATE TABLE day_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calendar_id UUID NOT NULL REFERENCES calendars(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    color_meaning_id UUID NOT NULL REFERENCES color_meanings(id) ON DELETE CASCADE,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(calendar_id, date) -- One entry per day per calendar
);

-- Indexes for better performance
CREATE INDEX idx_calendars_user_id ON calendars(user_id);
CREATE INDEX idx_color_meanings_calendar_id ON color_meanings(calendar_id);
CREATE INDEX idx_day_entries_calendar_id ON day_entries(calendar_id);
CREATE INDEX idx_day_entries_date ON day_entries(date);
CREATE INDEX idx_day_entries_calendar_date ON day_entries(calendar_id, date);

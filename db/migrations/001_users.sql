-- Registered users. Each gets a 6-digit login code at registration; only its bcrypt hash is stored.
CREATE TABLE IF NOT EXISTS users (
    id               SERIAL PRIMARY KEY,
    email            VARCHAR(254) NOT NULL UNIQUE,   -- 254 = max length of a valid email address
    first_name       VARCHAR(100) NOT NULL,
    last_name        VARCHAR(100) NOT NULL,
    login_code_hash  TEXT NOT NULL,                  -- bcrypt hash of the 6-digit code, never the code itself
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT users_email_lowercase CHECK (email = LOWER(email)),
    CONSTRAINT users_first_name_not_blank CHECK (BTRIM(first_name) <> ''),
    CONSTRAINT users_last_name_not_blank CHECK (BTRIM(last_name) <> '')
);

-- Row Level Security with no policies: blocks Supabase's public REST API (anon/authenticated keys)
-- from reading this table. The Go API connects as the table owner, which RLS does not restrict.
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

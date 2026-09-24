-- Submitted checkout forms. No payment is processed; the data is just recorded.
CREATE TABLE IF NOT EXISTS checkouts (
    id                SERIAL PRIMARY KEY,
    user_id           INTEGER REFERENCES users(id) ON DELETE SET NULL,  -- NULL for guest checkouts
    email             VARCHAR(254) NOT NULL,
    phone             CHAR(10) NOT NULL,                 -- Indian mobile number, 10 digits
    address_label     VARCHAR(30) NOT NULL,              -- name for the address: 'Home', 'Work' or a custom one
    address_line1     VARCHAR(200) NOT NULL,             -- house / flat, street
    address_line2     VARCHAR(200) NOT NULL DEFAULT '',  -- optional: area, landmark
    city              VARCHAR(100) NOT NULL,
    state             VARCHAR(100) NOT NULL,
    pincode           CHAR(6) NOT NULL,                  -- Indian PIN code
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT checkouts_phone_format CHECK (phone ~ '^[6-9][0-9]{9}$'),
    CONSTRAINT checkouts_address_label_not_blank CHECK (BTRIM(address_label) <> ''),
    CONSTRAINT checkouts_address_line1_not_blank CHECK (BTRIM(address_line1) <> ''),
    CONSTRAINT checkouts_city_not_blank CHECK (BTRIM(city) <> ''),
    CONSTRAINT checkouts_state_not_blank CHECK (BTRIM(state) <> ''),
    CONSTRAINT checkouts_pincode_format CHECK (pincode ~ '^[1-9][0-9]{5}$')
);

-- A user's recent orders, newest first: used to offer their saved phone and address at checkout.
CREATE INDEX IF NOT EXISTS idx_checkouts_user_id_created_at ON checkouts(user_id, created_at DESC);

-- See 001_users.sql: only the Go API (table owner) may access checkouts.
ALTER TABLE checkouts ENABLE ROW LEVEL SECURITY;

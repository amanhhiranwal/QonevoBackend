CREATE TABLE IF NOT EXISTS contacts (

    id BIGSERIAL PRIMARY KEY,

    full_name VARCHAR(255) NOT NULL,

    email VARCHAR(255) NOT NULL,

    phone_number VARCHAR(20) NOT NULL,

    company_name VARCHAR(255),

    website_url TEXT,

    help_message TEXT NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- Migration: create_external_links_table
-- Created: 2026-05-09T03:10:54Z
-- Description: Add description here

BEGIN;

-- Add your migration SQL here
CREATE TABLE IF NOT EXISTS external_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    company_id UUID NOT NULL,
    platform VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW(),
        CONSTRAINT fk_external_links_company FOREIGN KEY (company_id) REFERENCES company_profiles (id) ON DELETE CASCADE
);

CREATE INDEX idx_external_links_company_id ON external_links (company_id);

COMMIT;
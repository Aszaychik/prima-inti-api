-- Migration: create_series_table
-- Created: 2026-05-11T08:17:08Z
-- Description: Add description here

BEGIN;

-- Add your migration SQL here
CREATE TABLE IF NOT EXISTS series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    brand_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW(),
        CONSTRAINT fk_series_brand FOREIGN KEY (brand_id) REFERENCES brands (id) ON DELETE CASCADE,
        UNIQUE (brand_id, name)
);

COMMIT;
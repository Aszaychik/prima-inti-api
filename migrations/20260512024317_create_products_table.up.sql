-- Migration: create_products_table
-- Created: 2026-05-12T02:43:17Z
-- Description: Add description here

BEGIN;

-- Add your migration SQL here
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    brand_id UUID NOT NULL,
    category_id UUID NOT NULL,
    series_id UUID,
    model VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2),
    stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW(),
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT NOW(),
        CONSTRAINT fk_products_brand FOREIGN KEY (brand_id) REFERENCES brands (id) ON DELETE RESTRICT,
        CONSTRAINT fk_products_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE RESTRICT,
        CONSTRAINT fk_products_series FOREIGN KEY (series_id) REFERENCES series (id) ON DELETE SET NULL
);

CREATE INDEX idx_products_brand ON products (brand_id);

CREATE INDEX idx_products_category ON products (category_id);

CREATE INDEX idx_products_series ON products (series_id);

COMMIT;
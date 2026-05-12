-- Migration: create_products_table (rollback)
-- Created: 2026-05-12T02:43:17Z

BEGIN;

-- Add your rollback SQL here
DROP TABLE IF EXISTS products;

COMMIT;
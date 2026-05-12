-- Migration: create_brands_table (rollback)
-- Created: 2026-05-11T04:31:53Z

BEGIN;

-- Add your rollback SQL here
DROP TABLE IF EXISTS brands;

COMMIT;
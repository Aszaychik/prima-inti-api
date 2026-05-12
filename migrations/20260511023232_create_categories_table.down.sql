-- Migration: create_categories_table (rollback)
-- Created: 2026-05-11T02:32:32Z

BEGIN;

-- Add your rollback SQL here
DROP TABLE IF EXISTS categories;

COMMIT;
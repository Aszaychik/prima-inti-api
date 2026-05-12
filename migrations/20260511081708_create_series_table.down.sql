-- Migration: create_series_table (rollback)
-- Created: 2026-05-11T08:17:08Z

BEGIN;

-- Add your rollback SQL here
DROP TABLE IF EXISTS series;

COMMIT;
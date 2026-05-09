-- Migration: create_external_links_table (rollback)
-- Created: 2026-05-09T03:10:54Z

BEGIN;

-- Add your rollback SQL here
DROP TABLE IF EXISTS external_links;

COMMIT;
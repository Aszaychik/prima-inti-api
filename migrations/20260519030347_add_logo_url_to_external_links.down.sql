-- Migration: add_logo_url_to_external_links (rollback)
-- Created: 2026-05-19T03:03:47Z

BEGIN;

-- Add your rollback SQL here
ALTER TABLE external_links ADD COLUMN logo_url TEXT;

COMMIT;
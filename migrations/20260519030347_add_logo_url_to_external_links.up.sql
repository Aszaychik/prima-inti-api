-- Migration: add_logo_url_to_external_links
-- Created: 2026-05-19T03:03:47Z
-- Description: Add description here

BEGIN;

-- Add your migration SQL here
ALTER TABLE external_links ADD COLUMN logo_url TEXT;

COMMIT;
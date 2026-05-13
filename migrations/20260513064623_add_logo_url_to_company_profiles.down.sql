-- Migration: add_logo_url_to_company_profiles (rollback)
-- Created: 2026-05-13T06:46:23Z

BEGIN;

-- Add your rollback SQL here
ALTER TABLE company_profiles DROP COLUMN logo_url;

COMMIT;
-- Migration: add_logo_url_to_company_profiles
-- Created: 2026-05-13T06:46:23Z
-- Description: Add description here

BEGIN;

-- Add your migration SQL here
ALTER TABLE company_profiles ADD COLUMN logo_url TEXT;

COMMIT;
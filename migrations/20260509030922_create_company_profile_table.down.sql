-- Migration: create_company_profile_table (rollback)
-- Created: 2026-05-09T03:09:22Z

BEGIN;

-- Add your rollback SQL here
DROP TABLE IF EXISTS company_profile CASCADE;

COMMIT;
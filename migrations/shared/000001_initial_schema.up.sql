-- Shared schema for cross-cutting infrastructure
-- This schema contains only infrastructure-level tables

CREATE SCHEMA IF NOT EXISTS shared;

-- Migration tracking table (managed by golang-migrate)
-- This is automatically created by golang-migrate

-- Version info table
CREATE TABLE IF NOT EXISTS shared.version_info (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert initial version info
INSERT INTO shared.version_info (key, value) VALUES
    ('schema_version', '1'),
    ('migrated_at', NOW()::TEXT)
ON CONFLICT (key) DO UPDATE SET
    value = EXCLUDED.value,
    updated_at = NOW();
-- Audit schema
CREATE SCHEMA IF NOT EXISTS audit;

-- Enable pgcrypto for UUID generation
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Audit events table
CREATE TABLE IF NOT EXISTS audit.audit_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    event_type TEXT NOT NULL,
    actor_type TEXT NOT NULL, -- user, service, agent
    actor_id TEXT NOT NULL,
    actor_tenant_id TEXT,
    resource_type TEXT NOT NULL, -- workflow, execution, sandbox, capability, tool
    resource_id TEXT NOT NULL,
    resource_tenant_id TEXT,
    action TEXT NOT NULL,
    outcome TEXT NOT NULL, -- success, failure, denied
    details JSONB,
    security_context JSONB,
    prev_hash TEXT,
    hash TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_audit_events_timestamp ON audit.audit_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_events_tenant ON audit.audit_events(resource_tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_type ON audit.audit_events(event_type);
CREATE INDEX IF NOT EXISTS idx_audit_events_actor ON audit.audit_events(actor_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_resource ON audit.audit_events(resource_type, resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_outcome ON audit.audit_events(outcome);
CREATE INDEX IF NOT EXISTS idx_audit_events_hash ON audit.audit_events(hash);

-- Integrity chain table for audit log verification
CREATE TABLE IF NOT EXISTS audit.integrity_chain (
    event_id UUID PRIMARY KEY REFERENCES audit.audit_events(event_id),
    prev_hash TEXT,
    hash TEXT NOT NULL,
    verified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    verified_by TEXT
);

CREATE INDEX IF NOT EXISTS idx_integrity_chain_hash ON audit.integrity_chain(hash);
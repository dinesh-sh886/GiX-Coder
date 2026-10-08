-- Gateway schema
CREATE SCHEMA IF NOT EXISTS gateway;

-- Executions table
CREATE TABLE IF NOT EXISTS gateway.executions (
    execution_id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL,
    workflow_type TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    input JSONB,
    output JSONB,
    error_code TEXT,
    error_message TEXT,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    correlation_id UUID,
    tenant_id TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_gateway_executions_tenant ON gateway.executions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_gateway_executions_status ON gateway.executions(status);
CREATE INDEX IF NOT EXISTS idx_gateway_executions_correlation ON gateway.executions(correlation_id);
CREATE INDEX IF NOT EXISTS idx_gateway_executions_started ON gateway.executions(started_at);

-- Idempotency keys table
DROP TABLE IF EXISTS gateway.idempotency_keys CASCADE;
CREATE TABLE gateway.idempotency_keys (
    idempotency_key TEXT PRIMARY KEY,
    request_hash VARCHAR(64) NOT NULL,
    response_status INT NOT NULL,
    response_body BYTEA,
    response_content_type VARCHAR(128),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_gateway_idempotency_expires ON gateway.idempotency_keys(expires_at);

-- Rate limit tracking (optional, for distributed rate limiting)
CREATE TABLE IF NOT EXISTS gateway.rate_limits (
    key TEXT PRIMARY KEY,
    count INTEGER NOT NULL DEFAULT 0,
    window_start TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_gateway_rate_limits_expires ON gateway.rate_limits(expires_at);
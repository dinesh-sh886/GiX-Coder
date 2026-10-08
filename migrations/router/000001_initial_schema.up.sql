-- Router schema
CREATE SCHEMA IF NOT EXISTS router;

-- Provider registry table
CREATE TABLE IF NOT EXISTS router.provider_registry (
    model_id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    name TEXT NOT NULL,
    max_tokens INTEGER NOT NULL DEFAULT 4096,
    max_output_tokens INTEGER NOT NULL DEFAULT 4096,
    input_cost_per_1k NUMERIC(10, 6) NOT NULL DEFAULT 0,
    output_cost_per_1k NUMERIC(10, 6) NOT NULL DEFAULT 0,
    latency_p50_ms INTEGER NOT NULL DEFAULT 0,
    latency_p99_ms INTEGER NOT NULL DEFAULT 0,
    capabilities TEXT[] NOT NULL DEFAULT '{}',
    available BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_router_providers_provider ON router.provider_registry(provider);
CREATE INDEX IF NOT EXISTS idx_router_providers_available ON router.provider_registry(available);

-- Model metadata table
CREATE TABLE IF NOT EXISTS router.model_metadata (
    model_id TEXT PRIMARY KEY,
    provider TEXT NOT NULL,
    name TEXT NOT NULL,
    max_tokens INTEGER NOT NULL DEFAULT 4096,
    max_output_tokens INTEGER NOT NULL DEFAULT 4096,
    context_window INTEGER NOT NULL DEFAULT 4096,
    input_cost_per_1k NUMERIC(10, 6) NOT NULL DEFAULT 0,
    output_cost_per_1k NUMERIC(10, 6) NOT NULL DEFAULT 0,
    latency_p50_ms INTEGER NOT NULL DEFAULT 0,
    latency_p99_ms INTEGER NOT NULL DEFAULT 0,
    capabilities TEXT[] NOT NULL DEFAULT '{}',
    available BOOLEAN NOT NULL DEFAULT TRUE,
    metadata JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_router_model_metadata_provider ON router.model_metadata(provider);

-- Budget tracking table
CREATE TABLE IF NOT EXISTS router.budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    workflow_id UUID,
    request_id UUID,
    model_id TEXT,
    tokens_reserved INTEGER NOT NULL DEFAULT 0,
    tokens_used INTEGER NOT NULL DEFAULT 0,
    cost_estimate NUMERIC(10, 6) NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'reserved',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_router_budgets_tenant ON router.budgets(tenant_id);
CREATE INDEX IF NOT EXISTS idx_router_budgets_workflow ON router.budgets(workflow_id);
CREATE INDEX IF NOT EXISTS idx_router_budgets_status ON router.budgets(status);
CREATE INDEX IF NOT EXISTS idx_router_budgets_expires ON router.budgets(expires_at);
-- Policy schema
CREATE SCHEMA IF NOT EXISTS policy;

-- Capability policies table
CREATE TABLE IF NOT EXISTS policy.capability_policies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    workflow_id UUID,
    workflow_type TEXT,
    capability_type TEXT NOT NULL,
    resource_scope TEXT,
    allowed BOOLEAN NOT NULL DEFAULT TRUE,
    constraints JSONB,
    priority INTEGER NOT NULL DEFAULT 100,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_policy_policies_tenant ON policy.capability_policies(tenant_id);
CREATE INDEX IF NOT EXISTS idx_policy_policies_workflow ON policy.capability_policies(workflow_id);
CREATE INDEX IF NOT EXISTS idx_policy_policies_capability ON policy.capability_policies(capability_type);

-- Tenant allowlists table
CREATE TABLE IF NOT EXISTS policy.tenant_allowlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    model_id TEXT NOT NULL,
    allowed BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, model_id)
);

CREATE INDEX IF NOT EXISTS idx_policy_allowlists_tenant ON policy.tenant_allowlists(tenant_id);
CREATE INDEX IF NOT EXISTS idx_policy_allowlists_model ON policy.tenant_allowlists(model_id);
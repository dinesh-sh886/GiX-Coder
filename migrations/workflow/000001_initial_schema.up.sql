-- Workflow schema
CREATE SCHEMA IF NOT EXISTS workflow;

-- Workflow definitions table
CREATE TABLE IF NOT EXISTS workflow.workflow_definitions (
    workflow_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    version TEXT NOT NULL,
    definition JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_workflow_definitions_name ON workflow.workflow_definitions(name);

-- Execution records table
CREATE TABLE IF NOT EXISTS workflow.execution_records (
    execution_id UUID PRIMARY KEY,
    workflow_id UUID NOT NULL REFERENCES workflow.workflow_definitions(workflow_id),
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

CREATE INDEX IF NOT EXISTS idx_workflow_executions_tenant ON workflow.execution_records(tenant_id);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_workflow ON workflow.execution_records(workflow_id);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_status ON workflow.execution_records(status);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_correlation ON workflow.execution_records(correlation_id);
CREATE INDEX IF NOT EXISTS idx_workflow_executions_started ON workflow.execution_records(started_at);
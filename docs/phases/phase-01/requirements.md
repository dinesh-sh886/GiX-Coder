# Phase 01 Requirements

## Overview

This document decomposes the Phase 01 specification into implementation-level requirements with full traceability. Each requirement defines its module ownership, inputs, outputs, dependencies, security constraints, observability requirements, test strategy, and acceptance evidence.

**Source**: [Phase 01 Specification](specification.md)
**Baseline**: [Phase 00 Architecture Baseline](../../architecture/baseline.md)

---

## Functional Requirements

### Gateway Module (`gateway/`)

| Field                          | Value                                                                                                                             |
| ------------------------------ | --------------------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-01                                                                                                                             |
| **Description**                | Submit workflow via REST API (`POST /workflows/execute`)                                                                          |
| **Priority**                   | P0                                                                                                                                |
| **Module**                     | gateway                                                                                                                           |
| **Inputs**                     | `WorkflowExecutionRequest` DTO (workflow definition, input parameters, capability requests)                                       |
| **Outputs**                    | `WorkflowExecutionResponse` DTO (execution_id, status, result or error)                                                           |
| **Dependencies**               | FR-02 (auth), FR-03 (authz), FR-05 (rate limit), FR-04 (routing)                                                                  |
| **Security Constraints**       | JWT/OIDC validation, request size limit (10MB), input sanitization, audit log on submit                                           |
| **Observability Requirements** | Structured log (correlation_id, tenant_id, workflow_id), RED metrics, trace span, audit event `workflow.submitted`                |
| **Test Strategy**              | Unit: DTO validation, auth middleware; Integration: full request flow with valid/invalid JWT; Contract: OpenAPI schema validation |
| **Acceptance Evidence**        | AC-01: Integration test passes with valid request; AC-02: Invalid JWT rejected with 401                                           |

| Field                          | Value                                                                                                                                |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------ |
| **ID**                         | FR-02                                                                                                                                |
| **Description**                | Authenticate requests via JWT/OIDC                                                                                                   |
| **Priority**                   | P0                                                                                                                                   |
| **Module**                     | gateway                                                                                                                              |
| **Inputs**                     | Authorization header (Bearer token), JWKS endpoint, issuer, audience config                                                          |
| **Outputs**                    | Authenticated `Identity` (subject, tenant_id, scopes, expiry) or error                                                               |
| **Dependencies**               | Configuration (JWKS, issuer, audience), shared/identity package                                                                      |
| **Security Constraints**       | RS256 validation, issuer validation, audience validation, clock skew tolerance (30s), key rotation support, token expiry enforcement |
| **Observability Requirements** | Structured log on auth success/failure, metric `auth_total{result="success                                                           | failure"}`, trace attribute `auth.method=jwt` |
| **Test Strategy**              | Unit: token validation logic (valid, expired, wrong issuer, wrong audience, malformed); Integration: end-to-end with real JWKS       |
| **Acceptance Evidence**        | AC-02: Invalid JWT rejected; valid JWT extracts identity correctly                                                                   |

| Field                          | Value                                                                                                 |
| ------------------------------ | ----------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-03                                                                                                 |
| **Description**                | Authorize via policy engine (capability grant evaluation)                                             |
| **Priority**                   | P0                                                                                                    |
| **Module**                     | gateway, policy                                                                                       |
| **Inputs**                     | `AuthorizationRequest` (identity, resource, action, context)                                          |
| **Outputs**                    | `AuthorizationDecision` (allow/deny, capability grants, obligations)                                  |
| **Dependencies**               | FR-02 (identity), policy engine stub, FR-07-FR-11 (capability definitions)                            |
| **Security Constraints**       | Default deny, explicit allow only, time-limited grants, audit every decision                          |
| **Observability Requirements** | Structured log (decision, grants), metric `authz_total{decision="allow                                | deny"}`, audit event `capability.grant | deny` |
| **Test Strategy**              | Unit: policy evaluation logic; Integration: gateway → policy engine; Contract: policy decision schema |
| **Acceptance Evidence**        | AC-03: Rate limit enforced; capability grants match policy                                            |

| Field                          | Value                                                                                                 |
| ------------------------------ | ----------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-04                                                                                                 |
| **Description**                | Route to workflow engine (direct execution adapter for Phase 01)                                      |
| **Priority**                   | P0                                                                                                    |
| **Module**                     | gateway, workflow                                                                                     |
| **Inputs**                     | Validated `WorkflowExecutionRequest`, authenticated identity, capability grants                       |
| **Outputs**                    | `ExecutionHandle` (execution_id, status channel)                                                      |
| **Dependencies**               | FR-01, FR-02, FR-03, workflow adapter interface                                                       |
| **Security Constraints**       | No customer code execution in gateway, mTLS to workflow adapter (future)                              |
| **Observability Requirements** | Trace span `gateway.route_to_workflow`, metric `workflow_routed_total`, audit event `workflow.routed` |
| **Test Strategy**              | Unit: routing logic; Integration: gateway → workflow adapter stub                                     |
| **Acceptance Evidence**        | AC-01: Request routed and execution handle returned                                                   |

| Field                          | Value                                                                                               |
| ------------------------------ | --------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-05                                                                                               |
| **Description**                | Rate limit per tenant/workflow                                                                      |
| **Priority**                   | P0                                                                                                  |
| **Module**                     | gateway                                                                                             |
| **Inputs**                     | Identity (tenant_id), workflow_id, configured limits (requests/min, burst)                          |
| **Outputs**                    | Allow/deny with retry-after header                                                                  |
| **Dependencies**               | FR-02 (identity), Redis for distributed rate limiting, configuration                                |
| **Security Constraints**       | No rate limit bypass, per-tenant isolation, burst protection                                        |
| **Observability Requirements** | Metric `rate_limit_total{result="allow                                                              | deny", tenant}`, structured log on denial, audit event `rate_limit.exceeded` |
| **Test Strategy**              | Unit: token bucket algorithm; Integration: load test with concurrent requests; Chaos: Redis failure |
| **Acceptance Evidence**        | AC-03: Load test shows rate limiting at configured threshold                                        |

### Sandbox Module (`sandbox/`)

| Field                          | Value                                                                                                                                                                     |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-06                                                                                                                                                                     |
| **Description**                | Execute code in isolated sandbox (gVisor primary, Firecracker fallback)                                                                                                   |
| **Priority**                   | P0                                                                                                                                                                        |
| **Module**                     | sandbox                                                                                                                                                                   |
| **Inputs**                     | `SandboxExecutionRequest` (command, args, env, capabilities, resource limits, workspace)                                                                                  |
| **Outputs**                    | `SandboxExecutionResult` (exit_code, stdout, stderr, duration, resource_usage)                                                                                            |
| **Dependencies**               | gVisor/Firecracker runtime, capability framework (FR-07-FR-11), resource enforcement (FR-12)                                                                              |
| **Security Constraints**       | No host access, no cross-tenant access, seccomp, capability dropping, read-only rootfs, network namespace isolation                                                       |
| **Observability Requirements** | Structured log (execution_id, tool, duration, resource_usage), metrics `sandbox_execution_total`, `sandbox_duration_seconds`, trace span, audit event `sandbox.execution` |
| **Test Strategy**              | Unit: request/response validation; Integration: gVisor sandbox execution; Security: escape attempts, resource exhaustion                                                  |
| **Acceptance Evidence**        | AC-04: File read/write works; AC-08: CPU/memory limits enforced                                                                                                           |

| Field                          | Value                                                                                                                                                |
| ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-07                                                                                                                                                |
| **Description**                | Grant file capabilities (path-scoped read/write/list)                                                                                                |
| **Priority**                   | P0                                                                                                                                                   |
| **Module**                     | sandbox, policy                                                                                                                                      |
| **Inputs**                     | `CapabilityGrantRequest` (type=filesystem, paths, operations, execution_id)                                                                          |
| **Outputs**                    | `CapabilityGrant` (grant_id, paths, operations, expires_at, constraints)                                                                             |
| **Dependencies**               | FR-03 (policy evaluation), capability schema, FR-06 (enforcement)                                                                                    |
| **Security Constraints**       | Path canonicalization, symlink resolution, traversal prevention (`../`), absolute path rejection, workspace root enforcement, file size/count limits |
| **Observability Requirements** | Audit event `capability.grant.filesystem`, metric `capability_grant_total{type="filesystem"}`, log grant details                                     |
| **Test Strategy**              | Unit: path validation; Security: traversal, symlink escape, absolute path, oversized file, special files                                             |
| **Acceptance Evidence**        | AC-04: Authorized paths work; AC-05: Unauthorized paths blocked                                                                                      |

| Field                          | Value                                                                                                                                                         |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-08                                                                                                                                                         |
| **Description**                | Grant shell capabilities (command allowlist, timeout, resource limits)                                                                                        |
| **Priority**                   | P0                                                                                                                                                            |
| **Module**                     | sandbox, policy                                                                                                                                               |
| **Inputs**                     | `CapabilityGrantRequest` (type=shell, commands, args_pattern, timeout, resource_limits, execution_id)                                                         |
| **Outputs**                    | `CapabilityGrant` (grant_id, allowed_commands, constraints, expires_at)                                                                                       |
| **Dependencies**               | FR-03, FR-06, shell tool implementation                                                                                                                       |
| **Security Constraints**       | Executable allowlist only, argument validation (no injection), environment restrictions, working directory confinement, stdout/stderr limits, signal handling |
| **Observability Requirements** | Audit event `capability.grant.shell`, metric `capability_grant_total{type="shell"}`, log command execution                                                    |
| **Test Strategy**              | Unit: command validation; Security: command injection, unauthorized command, resource exhaustion, env var leakage                                             |
| **Acceptance Evidence**        | AC-06: Allowed commands execute; AC-07: Disallowed commands blocked                                                                                           |

| Field                          | Value                                                                                                                |
| ------------------------------ | -------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-09                                                                                                                |
| **Description**                | Grant Git capabilities (repo-scoped, branch-scoped operations)                                                       |
| **Priority**                   | P0                                                                                                                   |
| **Module**                     | sandbox, policy                                                                                                      |
| **Inputs**                     | `CapabilityGrantRequest` (type=git, repo_url, ref, operations, credentials_ref, execution_id)                        |
| **Outputs**                    | `CapabilityGrant` (grant_id, repo_scope, branch_scope, operations, expires_at)                                       |
| **Dependencies**               | FR-03, FR-06, git tool implementation, Vault for credential injection                                                |
| **Security Constraints**       | Repo-scoped credentials, branch-scoped operations, no credential leakage in logs, credential cleanup on grant expiry |
| **Observability Requirements** | Audit event `capability.grant.git`, metric `capability_grant_total{type="git"}`, log git operations                  |
| **Test Strategy**              | Unit: repo/branch scope validation; Integration: clone, fetch, checkout, diff, status; Security: credential exposure |
| **Acceptance Evidence**        | Git operations within scope succeed; out-of-scope operations denied                                                  |

| Field                          | Value                                                                                                                   |
| ------------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-10                                                                                                                   |
| **Description**                | Grant test execution capabilities (framework-agnostic)                                                                  |
| **Priority**                   | P0                                                                                                                      |
| **Module**                     | sandbox, policy                                                                                                         |
| **Inputs**                     | `CapabilityGrantRequest` (type=test, framework, args, timeout, execution_id)                                            |
| **Outputs**                    | `CapabilityGrant` (grant_id, framework, constraints, expires_at)                                                        |
| **Dependencies**               | FR-03, FR-06, test tool implementation                                                                                  |
| **Security Constraints**       | Framework allowlist, argument validation, timeout enforcement, resource limits, no network access during test execution |
| **Observability Requirements** | Audit event `capability.grant.test`, metric `capability_grant_total{type="test"}`, log test execution                   |
| **Acceptance Evidence**        | Test execution works for allowed frameworks; resource limits enforced                                                   |

| Field                          | Value                                                                                                  |
| ------------------------------ | ------------------------------------------------------------------------------------------------------ |
| **ID**                         | FR-11                                                                                                  |
| **Description**                | Grant MCP tool capabilities (tool-scoped, approved registry)                                           |
| **Priority**                   | P0                                                                                                     |
| **Module**                     | harness, sandbox, policy                                                                               |
| **Inputs**                     | `CapabilityGrantRequest` (type=mcp, server_id, tool_name, args_schema, execution_id)                   |
| **Outputs**                    | `CapabilityGrant` (grant_id, server_id, tool_name, constraints, expires_at)                            |
| **Dependencies**               | FR-03, FR-06, FR-19 (MCP client), MCP registry                                                         |
| **Security Constraints**       | Approved registry only, tool-scoped grants, argument validation per schema, output validation, timeout |
| **Observability Requirements** | Audit event `capability.grant.mcp`, metric `capability_grant_total{type="mcp"}`, log MCP invocation    |
| **Acceptance Evidence**        | AC-14: MCP tool executes with valid grant; unauthorized tools denied                                   |

| Field                          | Value                                                                                                                    |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------ |
| **ID**                         | FR-12                                                                                                                    |
| **Description**                | Enforce resource limits (CPU, memory, pids, I/O, wall-time)                                                              |
| **Priority**                   | P0                                                                                                                       |
| **Module**                     | sandbox                                                                                                                  |
| **Inputs**                     | `ResourceLimits` (cpu_cores, memory_bytes, pids_limit, io_bps, wall_time_seconds)                                        |
| **Outputs**                    | Enforcement via cgroups v2; termination on limit exceed                                                                  |
| **Dependencies**               | FR-06 (sandbox execution), Linux cgroups v2, gVisor/Firecracker integration                                              |
| **Security Constraints**       | Hard limits (OOM kill, CPU throttle, pid limit), no limit bypass, cleanup on termination                                 |
| **Observability Requirements** | Metrics `sandbox_cpu_usage`, `sandbox_memory_usage`, `sandbox_pids`, `sandbox_io`, audit event `resource_limit.exceeded` |
| **Test Strategy**              | Unit: limit parsing; Integration: stress test with limit exceed; Chaos: OOM, CPU throttle                                |
| **Acceptance Evidence**        | AC-08: Stress test shows limits enforced with termination                                                                |

### Router Module (`router/`)

| Field                          | Value                                                                                                   |
| ------------------------------ | ------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-13                                                                                                   |
| **Description**                | Classify task for model routing (complexity, domain, requirements)                                      |
| **Priority**                   | P0                                                                                                      |
| **Module**                     | router                                                                                                  |
| **Inputs**                     | `TaskClassificationRequest` (prompt, context, files, repo_state, workflow_type)                         |
| **Outputs**                    | `TaskClassification` (complexity, domain, required_capabilities, estimated_tokens, latency_requirement) |
| **Dependencies**               | Classifier model/heuristics, context management, configuration                                          |
| **Security Constraints**       | Input sanitization, no PII in classification, bounded compute                                           |
| **Observability Requirements** | Metric `task_classification_total{complexity,domain}`, trace span, structured log                       |
| **Test Strategy**              | Unit: classifier logic with fixtures; Integration: router → classifier                                  |
| **Acceptance Evidence**        | AC-10: Unit tests with fixtures show correct classification                                             |

| Field                          | Value                                                                                  |
| ------------------------------ | -------------------------------------------------------------------------------------- |
| **ID**                         | FR-14                                                                                  |
| **Description**                | Select model per policy (cost, latency, capability, policy)                            |
| **Priority**                   | P0                                                                                     |
| **Module**                     | router                                                                                 |
| **Inputs**                     | `ModelSelectionRequest` (task_classification, tenant_policy, model_allowlist, budgets) |
| **Outputs**                    | `ModelSelection` (model_id, endpoint, token_budget, parameters)                        |
| **Dependencies**               | FR-13, policy engine, model provider adapters, budget manager                          |
| **Security Constraints**       | Model allowlist enforcement, tenant policy compliance, no unauthorized model access    |
| **Observability Requirements** | Metric `model_selection_total{model,tenant}`, trace span, audit event `model.selected` |
| **Test Strategy**              | Unit: selection logic with policies; Integration: router → provider adapters           |
| **Acceptance Evidence**        | AC-11: Integration test shows model selected per policy                                |

| Field                          | Value                                                                             |
| ------------------------------ | --------------------------------------------------------------------------------- |
| **ID**                         | FR-15                                                                             |
| **Description**                | Enforce token budgets (per request, per workflow, per tenant)                     |
| **Priority**                   | P0                                                                                |
| **Module**                     | router                                                                            |
| **Inputs**                     | `BudgetCheckRequest` (tenant_id, workflow_id, estimated_tokens, model_pricing)    |
| **Outputs**                    | `BudgetDecision` (allow/deny, remaining_budget, estimated_cost)                   |
| **Dependencies**               | FR-14, budget storage (Redis), configuration                                      |
| **Security Constraints**       | Hard budget limits, real-time tracking, overage prevention, audit on denial       |
| **Observability Requirements** | Metric `budget_check_total{result="allow                                          | deny"}`, `token_usage_total{tenant,model}`, audit event `budget.exceeded` |
| **Test Strategy**              | Unit: budget calculation; Integration: concurrent requests with budget exhaustion |
| **Acceptance Evidence**        | AC-12: Integration test shows budget enforcement                                  |

| Field                          | Value                                                                                        |
| ------------------------------ | -------------------------------------------------------------------------------------------- |
| **ID**                         | FR-16                                                                                        |
| **Description**                | Validate model outputs (schema, safety, PII)                                                 |
| **Priority**                   | P0                                                                                           |
| **Module**                     | router                                                                                       |
| **Inputs**                     | `ModelOutput` (raw response, expected_schema, safety_rules, pii_patterns)                    |
| **Outputs**                    | `ValidatedOutput` (validated data, warnings, rejections)                                     |
| **Dependencies**               | Schema validator, safety classifier, PII detector                                            |
| **Security Constraints**       | PII detection/redaction, schema validation, safety policy enforcement, no raw output leakage |
| **Observability Requirements** | Metric `output_validation_total{result="valid                                                | invalid | warning"}`, trace span, audit event `output.validated` |
| **Test Strategy**              | Unit: validator logic; Integration: model response → validator                               |
| **Acceptance Evidence**        | Schema validation works; PII detected; safety violations rejected                            |

### Harness Module (`harness/`)

| Field                          | Value                                                                                                                              |
| ------------------------------ | ---------------------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-17                                                                                                                              |
| **Description**                | Spawn/manage agent lifecycle (CREATED → INITIALIZING → READY → RUNNING → CHECKPOINTING → PAUSED → TERMINATING → TERMINATED/FAILED) |
| **Priority**                   | P0                                                                                                                                 |
| **Module**                     | harness                                                                                                                            |
| **Inputs**                     | `AgentSpawnRequest` (agent_definition, capability_grants, initial_state, workflow_context)                                         |
| **Outputs**                    | `AgentHandle` (agent_id, state_channel, execution_channel)                                                                         |
| **Dependencies**               | FR-07-FR-11 (capability grants), FR-18 (tool registry), FR-19 (MCP), FR-20 (checkpointing), shared/state machine                   |
| **Security Constraints**       | State transitions validated, capability grants bound to lifecycle, no privilege escalation, crash recovery                         |
| **Observability Requirements** | Metrics `agent_lifecycle_total{state,transition}`, trace spans per state, structured log, audit events for state changes           |
| **Test Strategy**              | Unit: state machine; Integration: full lifecycle with capabilities; Chaos: crash recovery                                          |
| **Acceptance Evidence**        | AC-13: Integration test spawns agent with grants; lifecycle transitions verified                                                   |

| Field                          | Value                                                                                                                                         |
| ------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-18                                                                                                                                         |
| **Description**                | Register/execute tools (built-in + MCP)                                                                                                       |
| **Priority**                   | P0                                                                                                                                            |
| **Module**                     | harness                                                                                                                                       |
| **Inputs**                     | `ToolRegistration` (id, version, schema, required_capabilities, timeout, resource_limits), `ToolInvocation` (tool_id, args, capability_grant) |
| **Outputs**                    | `ToolResult` (output, duration, resource_usage, audit_context)                                                                                |
| **Dependencies**               | FR-07-FR-11 (capabilities), FR-19 (MCP), sandbox executor, tool registry                                                                      |
| **Security Constraints**       | Tool schema validation, capability verification before execution, timeout/resource enforcement, audit every invocation                        |
| **Observability Requirements** | Metrics `tool_execution_total{tool,result}`, `tool_duration_seconds{tool}`, trace span, audit event `tool.invoked`                            |
| **Test Strategy**              | Unit: registry operations; Integration: tool execution with capabilities; Contract: tool schema                                               |
| **Acceptance Evidence**        | Tool registration works; execution enforces capabilities; results returned                                                                    |

| Field                          | Value                                                                                                                     |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------- |
| **ID**                         | FR-19                                                                                                                     |
| **Description**                | Integrate MCP tools (client implementation)                                                                               |
| **Priority**                   | P0                                                                                                                        |
| **Module**                     | harness                                                                                                                   |
| **Inputs**                     | `MCPServerConfig` (url, transport, auth), `MCPToolInvocation` (server_id, tool_name, args)                                |
| **Outputs**                    | `MCPToolResult` (output, metadata)                                                                                        |
| **Dependencies**               | FR-11 (MCP capability), FR-18 (tool registry), MCP SDK, capability framework                                              |
| **Security Constraints**       | Approved registry only, capability grant required per invocation, argument/output validation, timeout, circuit breaker    |
| **Observability Requirements** | Metrics `mcp_invocation_total{server,tool,result}`, `mcp_duration_seconds{server}`, trace span, audit event `mcp.invoked` |
| **Test Strategy**              | Unit: MCP client; Integration: approved server tool invocation; Security: unapproved server denied                        |
| **Acceptance Evidence**        | AC-14: MCP tool executes with valid grant; unapproved server denied                                                       |

| Field                          | Value                                                                                                                          |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------ |
| **ID**                         | FR-20                                                                                                                          |
| **Description**                | Checkpoint/restore agent state (serialized to object storage)                                                                  |
| **Priority**                   | P1                                                                                                                             |
| **Module**                     | harness                                                                                                                        |
| **Inputs**                     | `CheckpointRequest` (agent_id, state, execution_id), `RestoreRequest` (checkpoint_id)                                          |
| **Outputs**                    | `CheckpointResult` (checkpoint_id, size, timestamp), `RestoredAgentState`                                                      |
| **Dependencies**               | Object storage abstraction (S3/GCS/local), FR-17 (lifecycle), FR-18 (tool state), serialization format                         |
| **Security Constraints**       | Encrypted at rest, tenant-scoped access, integrity verification (hash), no secrets in checkpoint, cleanup policy               |
| **Observability Requirements** | Metrics `checkpoint_total{result}`, `checkpoint_size_bytes`, `restore_total{result}`, `restore_duration_seconds`, audit events |
| **Test Strategy**              | Unit: serialization/deserialization; Integration: checkpoint → restore → continue execution; Chaos: partial restore            |
| **Acceptance Evidence**        | AC-15: Integration test checkpoint/restore works                                                                               |

---

## Non-Functional Requirements

| ID     | Requirement                   | Target                      | Module                    | Measurement Method                                |
| ------ | ----------------------------- | --------------------------- | ------------------------- | ------------------------------------------------- |
| NFR-01 | Gateway p99 latency           | < 100ms                     | gateway                   | Load test (k6) at 1000 RPS                        |
| NFR-02 | Sandbox cold start            | < 10s                       | sandbox                   | Repeated cold starts (n=100), p99                 |
| NFR-03 | Sandbox warm start            | < 2s                        | sandbox                   | Warm pool, repeated starts (n=100), p99           |
| NFR-04 | Concurrent executions         | 1000                        | gateway, sandbox, harness | Load test with 1000 concurrent workflows          |
| NFR-05 | Availability                  | 99.9%                       | all                       | SLO measurement over 30 days                      |
| NFR-06 | Audit log durability          | 11 nines                    | audit                     | Immutable storage verification, DR test           |
| NFR-07 | Zero critical vulnerabilities | Continuous                  | all                       | SAST/SCA/Container scans in CI                    |
| NFR-08 | Sandbox escape resistance     | gVisor/Firecracker verified | sandbox                   | gVisor test suite, escape attempts                |
| NFR-09 | Configuration immutability    | Runtime                     | all                       | Config validation at startup, no runtime mutation |
| NFR-10 | Observability coverage        | 100% endpoints              | all                       | Trace span on every endpoint, RED metrics         |

---

## Requirement Dependencies

```
FR-01 (Submit workflow)
  ├─ FR-02 (Auth) ──────────────────┐
  ├─ FR-03 (Authz) ─────────────────┤
  ├─ FR-05 (Rate limit) ────────────┤
  └─ FR-04 (Route to workflow) ◄────┘
       │
       ▼
FR-04 (Route to workflow)
       │
       ▼
WorkflowCoordinator → DirectExecutionAdapter (Phase 01 stub for Temporal)
       │
       ▼
FR-17 (Agent lifecycle) ◄──────────────────┐
       │                                    │
       ├─ FR-07..FR-11 (Capabilities) ◄────┤
       ├─ FR-18 (Tool registry) ◄──────────┤
       ├─ FR-19 (MCP) ◄────────────────────┤
       └─ FR-20 (Checkpoint) ◄─────────────┘
              │
              ▼
       FR-06 (Sandbox execution)
              │
       ├─ FR-07 (Files) ──► Path validation, traversal prevention
       ├─ FR-08 (Shell) ──► Command allowlist, injection prevention
       ├─ FR-09 (Git) ────► Repo/branch scoping, credential handling
       ├─ FR-10 (Test) ───► Framework allowlist, resource limits
       └─ FR-11 (MCP) ────► Tool-scoped, registry approval
              │
              ▼
       FR-12 (Resource limits) ──► cgroups v2 enforcement

FR-13 (Task classification) ──┐
FR-14 (Model selection) ◄─────┤
FR-15 (Budget enforcement) ◄──┤
FR-16 (Output validation) ◄───┘
```

**Note**: Phase 01 uses `WorkflowCoordinator` + `DirectExecutionAdapter` (synchronous, no Temporal). Phase 02 replaces `DirectExecutionAdapter` with `TemporalWorkflowAdapter`.

---

## Security Constraints Summary

| Constraint                                   | Requirements Affected        | Enforcement Point                          |
| -------------------------------------------- | ---------------------------- | ------------------------------------------ |
| Control plane never executes customer code   | FR-06, FR-17, FR-18          | Architecture (ADR-0001), module boundaries |
| Default deny, explicit capability grants     | FR-07..FR-11, FR-03          | Policy engine, sandbox supervisor          |
| Path canonicalization & traversal prevention | FR-07                        | Sandbox filesystem tool                    |
| Command allowlist & injection prevention     | FR-08                        | Sandbox shell tool                         |
| Repo/branch-scoped Git credentials           | FR-09                        | Git tool, Vault injection                  |
| Network egress allowlist only                | FR-11 (MCP), sandbox network | CNI policy, sandbox supervisor             |
| Token budget hard limits                     | FR-15                        | Router budget manager                      |
| PII detection in model I/O                   | FR-16                        | Router output validator                    |
| Encrypted checkpoints, no secrets            | FR-20                        | Checkpoint store, serialization            |
| mTLS between control plane services          | All inter-service            | Infrastructure, service mesh               |
| Immutable audit log for all security events  | All                          | Audit service                              |

---

## Observability Requirements Summary

| Signal  | Requirements                                         | Implementation                                         |
| ------- | ---------------------------------------------------- | ------------------------------------------------------ |
| Logs    | Structured JSON, correlation_id, tenant_id, trace_id | zerolog (Go), pino (TS)                                |
| Metrics | RED + business + system per service                  | Prometheus exposition, histogram buckets per standards |
| Traces  | W3C TraceContext, span per logical operation         | OpenTelemetry, 100% errors, 10% success sampling       |
| Audit   | Immutable, tamper-evident, separate store            | Merkle tree/hash chain, 7yr retention, strict RBAC     |

---

## Test Strategy Summary

| Layer        | Target                 | Tools                      | Phase 01 Specific                                          |
| ------------ | ---------------------- | -------------------------- | ---------------------------------------------------------- |
| Unit         | 80%+ coverage          | Go testing, Vitest         | All validation, state machines, classifiers                |
| Integration  | Critical paths         | testcontainers, localstack | Gateway → WorkflowCoordinator → Harness → Sandbox → Router |
| Contract     | All external APIs      | Pact/Schemathesis          | Protobuf + OpenAPI for all services                        |
| Architecture | All modules            | Custom (go list, eslint)   | Import boundaries, cyclic deps, layer violations           |
| Security     | Sandbox attack surface | gVisor test suite, custom  | Traversal, injection, escape, resource exhaustion          |
| Performance  | NFR targets            | k6                         | Load (1000 concurrent), stress, soak                       |
| Chaos        | System resilience      | Chaos Mesh                 | Sandbox kill, network partition, latency                   |

---

## Acceptance Criteria Traceability

| AC ID | Requirement(s)                    | Test Method          |
| ----- | --------------------------------- | -------------------- |
| AC-01 | FR-01, FR-02, FR-03, FR-04, FR-05 | Integration test     |
| AC-02 | FR-02                             | Unit + Integration   |
| AC-03 | FR-05                             | Load test            |
| AC-04 | FR-06, FR-07                      | Integration test     |
| AC-05 | FR-07                             | Security test        |
| AC-06 | FR-06, FR-08                      | Integration test     |
| AC-07 | FR-08                             | Security test        |
| AC-08 | FR-12                             | Stress test          |
| AC-09 | FR-11 (network)                   | Network test         |
| AC-10 | FR-13                             | Unit test (fixtures) |
| AC-11 | FR-14                             | Integration test     |
| AC-12 | FR-15                             | Integration test     |
| AC-13 | FR-17, FR-07..FR-11               | Integration test     |
| AC-14 | FR-19, FR-11                      | Integration test     |
| AC-15 | FR-20                             | Integration test     |
| AC-16 | All FRs                           | E2E test             |
| AC-17 | All audit events                  | Verification         |
| AC-18 | All quality gates                 | CI/CD pipeline       |
| AC-19 | NFR-07                            | Security scan        |
| AC-20 | Documentation                     | Review               |

---

## Metadata

---

title: Phase 01 Requirements
type: phase
phase: 01
status: Planning - Awaiting Plan Approval
author: Planner
date: 2024-10-07
reviewers: Product Owner, Architect, Security Lead
approved_by: N/A - Awaiting Plan Approval
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---

# Phase 01 Architecture

## Overview

This document defines the detailed architecture for Phase 01: Agent Gateway & Execution Sandbox Foundation. It extends the Phase 00 Architecture Baseline with concrete module designs, interfaces, data flows, and security boundaries.

**Source**: [Phase 01 Requirements](requirements.md), [Phase 00 Baseline](../../architecture/baseline.md)
**ADRs**: ADR-0001, ADR-0002, ADR-0003, ADR-0004

---

## Module Structure (Modular Monolith)

```
gix-coder/
├── gateway/          # Agent Gateway - API, auth, routing, rate limiting
├── sandbox/          # Execution Sandbox - isolation, tool execution
├── router/           # Context/Policy/Model Router - classification, selection, budgets
├── harness/          # Agent Harness - lifecycle, capabilities, tools, MCP, checkpointing
├── workflow/         # Workflow Adapter - direct execution (stub for Temporal Phase 02)
├── policy/           # Policy Engine - capability evaluation (stub for OPA/Cedar Phase 02)
├── audit/            # Audit Logging - immutable event log
├── shared/           # Shared Kernel - DTOs, errors, logging, config, metrics, tracing
├── api/              # API Contracts - Protobuf + OpenAPI for all services
└── deploy/           # Deployment - K8s manifests, Dockerfiles, Helm, docker-compose
```

### Module Independence Rules (ADR-0002)

1. **No circular dependencies** between modules
2. **Shared kernel only** - modules import from `shared/`, not from each other
3. **DTOs at boundaries** - all external communication uses DTOs
4. **Internal packages private** - `internal/` not importable by other modules
5. **Explicit contracts** - API definitions in `api/` are the contract

---

## Control Plane / Data Plane Separation (ADR-0001)

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                          CONTROL PLANE (Trusted)                            │
│  ┌──────────┐  ┌────────────┐  ┌──────────┐  ┌────────────┐  ┌────────┐  │
│  │ Gateway  │  │  Workflow  │  │  Router  │  │   Policy   │  │ Audit  │  │
│  │          │  │  Adapter   │  │          │  │   Engine   │  │        │  │
│  └──────────┘  └────────────┘  └──────────┘  └────────────┘  └────────┘  │
│  ┌──────────┐  ┌────────────┐                                           │
│  │ Harness  │  │   Shared   │                                           │
│  └──────────┘  └────────────┘                                           │
│       │              │              │              │                     │
│       ▼              ▼              ▼              ▼                     │
│    mTLS           mTLS           mTLS           mTLS                     │
└───────│──────────────│──────────────│──────────────│────────────────────┘
        │              │              │              │
        ▼              ▼              ▼              ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           DATA PLANE (Untrusted)                            │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                        SANDBOX (gVisor / Firecracker)                │  │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌───────────────┐  │  │
│  │  │ Files   │ │ Shell   │ │ Git     │ │ Test    │ │  Network      │  │  │
│  │  │ Tool    │ │ Tool    │ │ Tool    │ │ Tool    │ │  Policy       │  │  │
│  │  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └───────────────┘  │  │
│  │  Capability Grants (scoped, time-limited, audited, enforced)         │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────────────────┘
```

**Boundaries:**

- All control plane → data plane communication via gRPC with mTLS
- Data plane has NO direct access to control plane services
- Capability grants are explicit, scoped, time-limited, audited
- No shared databases, no shared memory, no shared filesystem
- Control plane services are stateless (except audit)

---

## Module Architectures

### 1. Gateway Module (`gateway/`)

#### Responsibility

API entry point for workflow submission. Handles authentication, authorization, validation, rate limiting, routing, and observability.

#### Owned State

- None (stateless). Session/rate limit state in Redis.

#### Public Interfaces

**REST (OpenAPI 3.1)**

```
POST   /workflows/execute        Submit workflow for execution
GET    /workflows/{execution_id} Get execution status/result
GET    /health/live              Liveness probe
GET    /health/ready             Readiness probe
GET    /metrics                  Prometheus metrics
```

**gRPC (Protobuf)**

```protobuf
service GatewayService {
  rpc ExecuteWorkflow(ExecuteWorkflowRequest) returns (ExecuteWorkflowResponse);
  rpc GetExecution(GetExecutionRequest) returns (GetExecutionResponse);
}
```

#### Key Components

| Component        | Responsibility                                              |
| ---------------- | ----------------------------------------------------------- |
| `AuthMiddleware` | JWT/OIDC validation, identity extraction, token caching     |
| `RateLimiter`    | Token bucket per tenant/workflow, Redis-backed, distributed |
| `Validator`      | DTO validation (request/response), schema enforcement       |
| `Authorizer`     | Policy engine integration, capability grant evaluation      |
| `Router`         | Request routing to workflow adapter                         |
| `AuditLogger`    | Structured audit events for all security-relevant actions   |

#### Data Flow

```
Request
  → AuthMiddleware (JWT validation, identity)
  → Request ID / Trace Context generation
  → Validator (DTO schema validation)
  → RateLimiter (tenant/workflow limits)
  → Authorizer (policy evaluation, capability grants)
  → Router (forward to Workflow Adapter)
  → Response (with correlation_id, audit event)
```

#### Error Taxonomy

| Code                 | HTTP | Description                                        |
| -------------------- | ---- | -------------------------------------------------- |
| `AUTH_INVALID`       | 401  | JWT missing, malformed, expired, invalid signature |
| `AUTHZ_DENIED`       | 403  | Policy denial, insufficient capabilities           |
| `VALIDATION_ERROR`   | 400  | Request DTO validation failure                     |
| `RATE_LIMITED`       | 429  | Tenant/workflow rate limit exceeded                |
| `WORKFLOW_NOT_FOUND` | 404  | Execution ID not found                             |
| `INTERNAL_ERROR`     | 500  | Unexpected error (never leaks internal details)    |

#### Idempotency

- `POST /workflows/execute` supports `Idempotency-Key` header
- Key scoped to tenant+workflow, TTL 24h
- Duplicate key returns original response

#### Dependencies

- **Allowed**: `shared/`, `api/`, `workflow/` (via interface), `policy/` (via interface), `audit/` (via interface), Redis, Vault
- **Forbidden**: Direct `sandbox/`, `harness/`, `router/` imports

#### Failure Modes

| Scenario                     | Behavior                           |
| ---------------------------- | ---------------------------------- |
| JWT validation failure       | 401, audit log, no downstream call |
| Rate limit exceeded          | 429 with `Retry-After`, audit log  |
| Policy engine unavailable    | 503 (fail closed), audit log       |
| Workflow adapter unavailable | 503, audit log, circuit breaker    |

---

### 2. Sandbox Module (`sandbox/`)

#### Responsibility

Isolated execution environment for untrusted code. Provides capability-enforced tool execution with resource limits.

#### Owned State

- Sandbox instances (ephemeral, per-execution)
- Capability grants (validated per request)

#### Public Interfaces

**gRPC (Protobuf)**

```protobuf
service SandboxService {
  rpc Execute(ExecuteRequest) returns (ExecuteResponse);
  rpc CreateSession(CreateSessionRequest) returns (CreateSessionResponse);
  rpc DestroySession(DestroySessionRequest) returns (DestroySessionResponse);
}
```

#### Architecture

```
Sandbox Supervisor (Control Plane)
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│                    ISOLATION RUNTIME                            │
│  ┌─────────────────────────────────────────────────────────┐   │
│  │ gVisor (runsc) / Firecracker (firecracker-containerd)  │   │
│  │  • User namespace (rootless)                            │   │
│  │  • seccomp profile (per capability)                     │   │
│  │  • Linux capabilities (drop all, add minimal)           │   │
│  │  • Read-only rootfs, tmpfs overlays                     │   │
│  │  • cgroups v2 (CPU, memory, pids, I/O)                  │   │
│  │  • Network namespace (none by default)                  │   │
│  │  • Wall-clock & CPU time limits                         │   │
│  │  • Audit: all syscalls logged via seccomp notify        │   │
│  └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
    │
    ▼
┌─────────────────────────────────────────────────────────────────┐
│                  CAPABILITY ENFORCEMENT                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐           │
│  │ Files    │ │ Shell    │ │ Git      │ │ Test     │ ...       │
│  │ Tool     │ │ Tool     │ │ Tool     │ │ Tool     │           │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘           │
│  Each tool: validates capability grant → enforces constraints  │
└─────────────────────────────────────────────────────────────────┘
```

#### Tool Interfaces

**Filesystem Tool**

```protobuf
message FileOperation {
  enum Op { READ = 0; WRITE = 1; LIST = 2; }
  Op op = 1;
  string path = 2;
  bytes content = 3;        // for WRITE
  FileOptions options = 4;  // encoding, max_size, etc.
}

message FileResult {
  bytes content = 1;
  repeated FileEntry entries = 2;
  int64 size = 3;
}
```

- **Capabilities**: `filesystem.read`, `filesystem.write`, `filesystem.list`
- **Constraints**: Path-scoped, canonicalized, traversal prevention, symlink resolution, max file size, max file count
- **Security**: Reject absolute paths, `../`, symlinks outside workspace, special files

**Shell Tool**

```protobuf
message ShellCommand {
  string command = 1;           // from allowlist
  repeated string args = 2;     // validated per command schema
  map<string, string> env = 3;  // restricted env vars
  int32 timeout_seconds = 4;
  ResourceLimits limits = 5;
}
```

- **Capabilities**: `shell.execute`
- **Constraints**: Executable allowlist, argument validation (no injection), env restrictions, working dir confinement, stdout/stderr limits, signal handling
- **Security**: No shell metacharacters, no command substitution, no pipeline

**Git Tool**

```protobuf
message GitOperation {
  enum Op { CLONE = 0; FETCH = 1; CHECKOUT = 2; DIFF = 3; STATUS = 4; BRANCH = 5; COMMIT = 6; PUSH = 7; }
  Op op = 1;
  string repo_url = 2;
  string ref = 3;              // branch/tag/commit
  GitOptions options = 4;
}
```

- **Capabilities**: `git.read`, `git.write`
- **Constraints**: Repo-scoped, branch-scoped, credential injection via Vault, no credential leakage
- **Operations**: Clone, fetch, checkout, diff, status, branch (read); commit, push (write)

**Test Tool**

```protobuf
message TestExecution {
  string framework = 1;         // go, jest, pytest, cargo, etc.
  repeated string args = 2;
  int32 timeout_seconds = 3;
  ResourceLimits limits = 4;
}
```

- **Capabilities**: `test.execute`
- **Constraints**: Framework allowlist, argument validation, timeout, resource limits, no network

**Network Tool**

```protobuf
message NetworkRequest {
  string url = 1;
  string method = 2;
  map<string, string> headers = 3;
  bytes body = 4;
  int32 timeout_seconds = 5;
}
```

- **Capabilities**: `network.egress`
- **Constraints**: Destination allowlist, port allowlist, DNS policy, no private network access, metadata service protection

#### Capability Model

```protobuf
message CapabilityGrant {
  string grant_id = 1;
  string execution_id = 2;
  CapabilityType type = 3;
  ResourceScope scope = 4;
  Constraints constraints = 5;
  google.protobuf.Timestamp issued_at = 6;
  google.protobuf.Timestamp expires_at = 7;
  string policy_version = 8;
  AuditContext audit_context = 9;
}

enum CapabilityType {
  FILESYSTEM_READ = 0;
  FILESYSTEM_WRITE = 1;
  FILESYSTEM_LIST = 2;
  SHELL_EXECUTE = 3;
  GIT_READ = 4;
  GIT_WRITE = 5;
  TEST_EXECUTE = 6;
  NETWORK_EGRESS = 7;
  MCP_INVOKE = 8;
}
```

**Default Behavior: DENY**

Every capability must be:

- Explicit (requested and granted)
- Scoped (resource, path, repo, command)
- Time-limited (expires_at)
- Auditable (grant_id, audit_context)

#### Dependencies

- **Allowed**: `shared/`, `api/`, gVisor/Firecracker runtime, Linux kernel (cgroups, namespaces, seccomp)
- **Forbidden**: Direct `gateway/`, `harness/`, `router/`, `policy/` imports

#### Failure Modes

| Scenario                         | Behavior                                         |
| -------------------------------- | ------------------------------------------------ |
| Capability grant missing/invalid | 403, audit `capability.deny`, no execution       |
| Resource limit exceeded          | SIGKILL/SIGTERM, audit `resource_limit.exceeded` |
| Sandbox runtime failure          | 500, audit, automatic cleanup, circuit breaker   |
| Network policy violation         | Connection refused, audit `network.denied`       |
| Credential injection failure     | 500, audit, no execution                         |

---

### 3. Router Module (`router/`)

#### Responsibility

Task classification, policy evaluation, model selection, budget enforcement, and output validation.

#### Owned State

- Model provider registry (cached)
- Budget state (Redis)
- Classification models (in-memory/embedded)

#### Public Interfaces

**gRPC (Protobuf)**

```protobuf
service RouterService {
  rpc ClassifyTask(ClassifyTaskRequest) returns (ClassifyTaskResponse);
  rpc SelectModel(SelectModelRequest) returns (SelectModelResponse);
  rpc ValidateOutput(ValidateOutputRequest) returns (ValidateOutputResponse);
  rpc CheckBudget(CheckBudgetRequest) returns (CheckBudgetResponse);
}
```

#### Component Diagram

```
Task Request
    │
    ▼
┌─────────────────────┐
│  Task Classifier    │  → complexity, domain, capabilities, tokens, latency
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Policy Evaluator   │  → tenant policy, model allowlist, capability requirements
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Model Selector     │  → cost, latency, capability match, provider health
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Budget Manager     │  → token budget check, reservation, enforcement
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Provider Adapter   │  → OpenAI, Anthropic, Ollama (pluggable)
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│  Output Validator   │  → schema, safety, PII
└─────────┬───────────┘
```

#### Provider Adapter Interface

```go
type ModelProvider interface {
    // Metadata
    Name() string
    Models() []ModelMetadata  // capability, context_limit, pricing, latency
    HealthCheck(ctx context.Context) error

    // Inference
    Complete(ctx context.Context, req CompletionRequest) (CompletionResponse, error)
    StreamComplete(ctx context.Context, req CompletionRequest) (<-chan CompletionChunk, error)

    // Budget
    EstimateTokens(req CompletionRequest) (int, error)
    EstimateCost(tokens int) (decimal.Decimal, error)
}
```

**Supported Providers (Phase 01):**

- OpenAI (GPT-4, GPT-4o, GPT-3.5-turbo)
- Anthropic (Claude 3 Opus, Sonnet, Haiku)
- Ollama/local (Llama, Mistral, CodeLlama)

#### Budget Enforcement

```
Budget Hierarchy:
Tenant (monthly)
  └─ Workflow (per execution)
       └─ Request (per model call)
```

- Hard limits at all levels
- Real-time reservation on selection
- Release on completion/error
- Audit on denial/exhaustion

#### Output Validation

| Validation | Implementation                                 |
| ---------- | ---------------------------------------------- |
| Schema     | JSON Schema / Protobuf validation              |
| Safety     | Classifier (prompt injection, harmful content) |
| PII        | Regex + ML detector (email, SSN, keys, tokens) |
| Size       | Max response size limit                        |

#### Dependencies

- **Allowed**: `shared/`, `api/`, `policy/` (via interface), Redis, Vault, model provider SDKs
- **Forbidden**: Direct `gateway/`, `sandbox/`, `harness/`, `workflow/` imports

---

### 4. Harness Module (`harness/`)

#### Responsibility

Agent lifecycle management, capability grant management, tool registry, MCP client, sandbox communication, checkpointing.

#### Owned State

- Agent instances (state machine)
- Tool registry
- MCP server connections
- Checkpoint store (object storage)

#### Public Interfaces

**gRPC (Protobuf)**

```protobuf
service HarnessService {
  rpc SpawnAgent(SpawnAgentRequest) returns (SpawnAgentResponse);
  rpc GetAgentState(GetAgentStateRequest) returns (AgentState);
  rpc ExecuteStep(ExecuteStepRequest) returns (ExecuteStepResponse);
  rpc Checkpoint(CheckpointRequest) returns (CheckpointResponse);
  rpc Restore(RestoreRequest) returns (RestoreResponse);
  rpc TerminateAgent(TerminateAgentRequest) returns (TerminateAgentResponse);
}
```

#### Agent Lifecycle State Machine

```
CREATED
    │
    ▼
INITIALIZING  ──(capability grants received)──► READY
    │                                            │
    │         (execution request)                │
    ▼                                            ▼
RUNNING ◄─────────────────────────────────────────┘
    │
    ├─(checkpoint)──► CHECKPOINTING ──► RUNNING
    │
    ├─(pause)────────► PAUSED ◄──────────────────┤
    │                    │                       │
    │         (resume)   │                       │
    │                    ▼                       │
    └──────────────────► RUNNING                 │
    │
    ├─(complete)──────► TERMINATING ──► TERMINATED
    │
    └─(error/crash)───► TERMINATING ──► FAILED
```

**Invalid Transitions**: Any not shown above (enforced by state machine)

#### Tool Registry

```go
type ToolDefinition struct {
    ID              string
    Version         string
    Schema          ToolSchema        // JSON Schema for args/result
    RequiredCapabilities []CapabilityType
    Timeout         time.Duration
    ResourceLimits  ResourceLimits
    AuditClassification string        // e.g., "filesystem.read", "shell.execute"
}

type ToolExecutor interface {
    Execute(ctx context.Context, invocation ToolInvocation) (ToolResult, error)
    Validate(args json.RawMessage) error
}
```

**Built-in Tools**: Filesystem, Shell, Git, Test
**MCP Tools**: Registered via MCP client, capability-mapped

#### Checkpoint Store Abstraction

```go
type CheckpointStore interface {
    Save(ctx context.Context, checkpoint Checkpoint) (CheckpointID, error)
    Load(ctx context.Context, id CheckpointID) (Checkpoint, error)
    Delete(ctx context.Context, id CheckpointID) error
    List(ctx context.Context, filter CheckpointFilter) ([]CheckpointMetadata, error)
}

// Phase 01 Implementation
type DirectCheckpointStore struct {
    // Serialized agent state → Object Storage (S3/GCS/local)
    // Encrypted at rest, tenant-scoped, integrity verified
}

// Phase 02 Implementation
type TemporalCheckpointStore struct {
    // Temporal-native: workflow history, replay, signals
}
```

**Checkpoint Contents:**

- Agent state (serialized)
- Tool states (MCP connections, file handles)
- Execution context (workflow, step, variables)
- Capability grants (active)
- **Excluded**: Secrets, credentials, raw model outputs

#### Dependencies

- **Allowed**: `shared/`, `api/`, `sandbox/` (via interface), `router/` (via interface), `policy/` (via interface), `audit/` (via interface), object storage, Redis, MCP SDK
- **Forbidden**: Direct `gateway/`, `workflow/` imports

---

### 5. Workflow Adapter (`workflow/`)

#### Responsibility

Direct execution adapter for Phase 01. Stub for Temporal integration (Phase 02).

#### Owned State

- Execution records (in-memory/PostgreSQL)
- No durable workflow state

#### Public Interfaces

**gRPC (Protobuf)**

```protobuf
service WorkflowService {
  rpc StartExecution(StartExecutionRequest) returns (StartExecutionResponse);
  rpc GetExecution(GetExecutionRequest) returns (ExecutionRecord);
  rpc CancelExecution(CancelExecutionRequest) returns (CancelExecutionResponse);
}
```

#### Phase 01 Behavior

- Synchronous execution coordination
- Calls Harness → Router → Sandbox directly
- No retry, no saga, no human-in-the-loop signals
- Audit log for execution start/complete/fail

#### Phase 02 Migration Path

- Replace with Temporal workflow definitions
- Activities → Sandbox execution
- Signals → Human approval
- History → Temporal event store
- Checkpointing → `TemporalCheckpointStore`

#### Dependencies

- **Allowed**: `shared/`, `api/`, `harness/` (via interface), `router/` (via interface), `audit/` (via interface), PostgreSQL
- **Forbidden**: Direct `gateway/`, `sandbox/`, `router/`, `policy/` imports

---

### 6. Policy Engine (`policy/`)

#### Responsibility

Capability grant evaluation, authorization decisions. Phase 01: stub with allowlist-based evaluation. Phase 02: OPA/Cedar integration.

#### Owned State

- Policy rules (loaded from config)
- Tenant/workflow capability policies

#### Public Interfaces

**gRPC (Protobuf)**

```protobuf
service PolicyService {
  rpc Evaluate(EvaluateRequest) returns (EvaluateResponse);
  rpc GetCapabilities(GetCapabilitiesRequest) returns (CapabilitiesResponse);
}
```

#### Phase 01 Behavior

- Static policy rules (YAML/JSON config)
- Tenant → capability allowlist
- Workflow → capability requirements
- Evaluation: intersection of allowlist and requirements
- Audit every decision

#### Phase 02 Migration

- OPA/Cedar policy engine
- Dynamic policy updates
- ABAC with context
- Policy testing/validation

#### Dependencies

- **Allowed**: `shared/`, `api/`, config, Vault
- **Forbidden**: Direct `gateway/`, `sandbox/`, `harness/`, `router/`, `workflow/` imports

---

### 7. Audit Module (`audit/`)

#### Responsibility

Immutable, tamper-evident audit logging for all security-relevant events.

#### Owned State

- Audit event log (append-only)
- Merkle tree / hash chain for integrity
- Separate storage (immutable object storage)

#### Public Interfaces

**gRPC (Protobuf)**

```protobuf
service AuditService {
  rpc LogEvent(LogEventRequest) returns (LogEventResponse);
  rpc QueryEvents(QueryEventsRequest) returns (QueryEventsResponse);
  rpc VerifyIntegrity(VerifyIntegrityRequest) returns (VerifyIntegrityResponse);
}
```

#### Audit Event Structure

```protobuf
message AuditEvent {
  string event_id = 1;                    // UUID v7
  google.protobuf.Timestamp timestamp = 2;
  string event_type = 3;                  // e.g., "workflow.submitted", "capability.grant"
  Actor actor = 4;                        // user|service|agent
  Resource resource = 5;                  // workflow, execution, sandbox, etc.
  string action = 6;                      // execute, grant, deny, read, write
  string outcome = 7;                     // success, failure, denied
  EventDetails details = 8;               // structured, type-specific
  SecurityContext security_context = 9;   // source_ip, user_agent, permissions
  Integrity integrity = 10;               // prev_hash, hash (Merkle tree)
}

message Actor {
  string type = 1;      // user, service, agent
  string id = 2;
  string tenant_id = 3;
}

message Resource {
  string type = 1;      // workflow, execution, sandbox, capability, tool
  string id = 2;
  string tenant_id = 3;
}
```

#### Event Types (Phase 01)

| Category       | Events                                                                                                                     |
| -------------- | -------------------------------------------------------------------------------------------------------------------------- |
| Authentication | `auth.success`, `auth.failure`                                                                                             |
| Authorization  | `authz.allow`, `authz.deny`                                                                                                |
| Workflow       | `workflow.submitted`, `workflow.routed`, `workflow.started`, `workflow.completed`, `workflow.failed`, `workflow.cancelled` |
| Capability     | `capability.grant`, `capability.deny`, `capability.revoked`, `capability.expired`                                          |
| Sandbox        | `sandbox.created`, `sandbox.execution`, `sandbox.terminated`, `sandbox.resource_exceeded`                                  |
| Tool           | `tool.invoked`, `tool.completed`, `tool.failed`                                                                            |
| MCP            | `mcp.invoked`, `mcp.completed`, `mcp.failed`                                                                               |
| Model          | `model.selected`, `model.request`, `model.response`, `model.validated`                                                     |
| Budget         | `budget.allow`, `budget.deny`, `budget.exhausted`                                                                          |
| Checkpoint     | `checkpoint.created`, `checkpoint.restored`, `checkpoint.failed`                                                           |
| Agent          | `agent.spawned`, `agent.state_change`, `agent.terminated`                                                                  |

#### Security

- **Never log**: Secrets, tokens, credentials, raw model outputs, PII
- **Immutability**: Append-only, Merkle tree, hash chain, object storage (WORM)
- **Retention**: 7 years minimum
- **Query**: Separate API, strict RBAC, audit trail on query access

#### Dependencies

- **Allowed**: `shared/`, `api/`, immutable object storage, Vault (signing keys)
- **Forbidden**: Direct imports from any other module

---

### 8. Shared Kernel (`shared/`)

#### Responsibility

Common utilities, types, and contracts used by all modules.

#### Structure

```
shared/
├── errors/           # Typed error definitions with codes
├── logging/          # Structured logging (zerolog/pino wrappers)
├── config/           # Configuration loading (Viper), validation
├── metrics/          # Prometheus metrics helpers, standard buckets
├── tracing/          # OpenTelemetry tracing utilities, propagation
├── validation/       # DTO validation, sanitization
├── dto/              # Shared DTOs (Identity, CapabilityGrant, etc.)
├── constants/        # Shared constants (error codes, event types)
├── security/         # Crypto helpers, sanitization, PII detection
└── version/          # Version info, build metadata
```

#### Key Types

```go
// Typed errors with codes
type ErrorCode string

const (
    ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
    ErrCodeAuthFailed     ErrorCode = "AUTH_FAILED"
    ErrCodeAuthzDenied    ErrorCode = "AUTHZ_DENIED"
    ErrCodeRateLimited    ErrorCode = "RATE_LIMITED"
    ErrCodeNotFound       ErrorCode = "NOT_FOUND"
    ErrCodeInternal       ErrorCode = "INTERNAL_ERROR"
    ErrCodeCapabilityDenied ErrorCode = "CAPABILITY_DENIED"
    ErrCodeResourceExceeded ErrorCode = "RESOURCE_EXCEEDED"
    ErrCodeSandboxError   ErrorCode = "SANDBOX_ERROR"
    ErrCodeBudgetExceeded ErrorCode = "BUDGET_EXCEEDED"
)

type AppError struct {
    Code    ErrorCode
    Message string
    Details map[string]any
    Cause   error
}

// Identity extracted from JWT
type Identity struct {
    Subject     string
    TenantID    string
    Scopes      []string
    ExpiresAt   time.Time
    IssuedAt    time.Time
}

// Correlation IDs for tracing
type CorrelationContext struct {
    TraceID      string
    SpanID       string
    CorrelationID string
    TenantID     string
    UserID       string
}
```

#### Dependencies

- **Allowed**: Standard library, zerolog, pino, Viper, OpenTelemetry, Prometheus client
- **Forbidden**: Any other module imports (this is the root of the dependency tree)

---

### 9. API Contracts (`api/`)

#### Protobuf Services

| Service         | Package           | Methods                                                                     |
| --------------- | ----------------- | --------------------------------------------------------------------------- |
| GatewayService  | `gix.gateway.v1`  | ExecuteWorkflow, GetExecution                                               |
| SandboxService  | `gix.sandbox.v1`  | Execute, CreateSession, DestroySession                                      |
| RouterService   | `gix.router.v1`   | ClassifyTask, SelectModel, ValidateOutput, CheckBudget                      |
| HarnessService  | `gix.harness.v1`  | SpawnAgent, GetAgentState, ExecuteStep, Checkpoint, Restore, TerminateAgent |
| WorkflowService | `gix.workflow.v1` | StartExecution, GetExecution, CancelExecution                               |
| PolicyService   | `gix.policy.v1`   | Evaluate, GetCapabilities                                                   |
| AuditService    | `gix.audit.v1`    | LogEvent, QueryEvents, VerifyIntegrity                                      |

#### Versioning Strategy

- Package version in proto path: `gix.gateway.v1`
- Breaking changes → new package version (`v2`), ADR required
- Field additions: optional, backward compatible
- Field removal: deprecated annotation, migration period

#### OpenAPI 3.1 (REST)

- Generated from protobuf via `buf` + `protoc-gen-openapiv2`
- Served at `/openapi.json` per service
- Includes examples, error schemas, security schemes (Bearer JWT)

#### Generated Clients

- Go: `buf generate` → `github.com/gix-coder/api/go/...`
- TypeScript: `buf generate` → `@gix-coder/api`
- Committed to repository, versioned with API

---

## Data Flow (Complete — Phase 01 Direct Execution)

```
1. Client → POST /workflows/execute {workflow, input, capabilities}
2. Gateway → AuthMiddleware: Validate JWT, extract Identity
3. Gateway → Validator: Validate DTO (request schema)
4. Gateway → RateLimiter: Check tenant/workflow limits
5. Gateway → Authorizer: Evaluate policy, get CapabilityGrants
6. Gateway → Audit: LogEvent(workflow.submitted)
7. Gateway → WorkflowCoordinator: StartExecution(request, grants)
8. WorkflowCoordinator → DirectExecutionAdapter: StartExecution(request, grants)
9. DirectExecutionAdapter → Harness: SpawnAgent(definition, grants)
10. Harness → Policy: GetCapabilities(tenant, workflow)
11. Harness → Router: ClassifyTask(prompt, context)
12. Harness → Sandbox: CreateSession(capabilities)
13. Agent → Router: SelectModel(classification, policy, budget)
14. Router → Provider: Complete(request)
15. Router → Validator: ValidateOutput(response)
16. Agent → Sandbox: ExecuteStep(tool, args, grant)
17. Sandbox → Tool: Validate grant, enforce constraints, execute
18. Sandbox → Audit: LogEvent(tool.invoked / capability.grant/deny)
19. Sandbox → Result: Return to Agent
20. Agent → Harness: Checkpoint(state) [if long-running]
21. Harness → CheckpointStore: Save(checkpoint)
22. Harness → DirectExecutionAdapter: CompleteExecution(result)
23. DirectExecutionAdapter → WorkflowCoordinator: ExecutionResult
24. WorkflowCoordinator → Gateway: ExecutionResult
25. Gateway → Audit: LogEvent(workflow.completed)
26. Gateway → Response: Return to Client
```

**Note**: Phase 01 uses `WorkflowCoordinator` + `DirectExecutionAdapter` (synchronous, no Temporal). Phase 02 replaces `DirectExecutionAdapter` with `TemporalWorkflowAdapter` and `DirectCheckpointStore` with `TemporalCheckpointStore`.

---

## Security Boundaries Detail

### Sandbox Isolation (ADR-0004)

| Layer                | Mechanism                                       | Purpose                                 |
| -------------------- | ----------------------------------------------- | --------------------------------------- |
| 1. Kernel Boundary   | gVisor (runsc) / Firecracker                    | Process isolation, syscall interception |
| 2. Syscall Filtering | seccomp-bpf (per capability)                    | Allow only required syscalls            |
| 3. Capabilities      | Linux capabilities (drop all, add minimal)      | Limit privileged operations             |
| 4. User Namespace    | Rootless (uid/gid mapping)                      | No root in sandbox                      |
| 5. Filesystem        | Read-only root, tmpfs overlays, path scoping    | Prevent host fs access                  |
| 6. Network           | Network namespace, CNI policy, egress allowlist | No ingress, controlled egress           |
| 7. Resources         | cgroups v2 (CPU, memory, pids, I/O)             | Prevent resource exhaustion             |
| 8. Time Limits       | cgroups + sandbox supervisor                    | Wall-clock and CPU time                 |
| 9. Audit             | seccomp notify, gVisor logging                  | Full syscall audit trail                |

### Capability Enforcement Flow

```
Tool Invocation
    │
    ▼
┌─────────────────────┐
│ Capability Validator │  → Check grant_id exists, not expired, matches tool
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Constraint Checker  │  → Path scope, command allowlist, repo scope, etc.
└─────────┬───────────┘
          │
          ▼
┌─────────────────────┐
│ Resource Enforcer   │  → cgroups limits, timeout, pids
└─────────┬───────────┘
          │
          ▼
    Execute Tool
          │
          ▼
┌─────────────────────┐
│ Audit Logger        │  → LogEvent(tool.invoked, grant_id, result)
└─────────────────────┘
```

### Secret Handling

| Secret Type             | Storage | Injection                             | Rotation |
| ----------------------- | ------- | ------------------------------------- | -------- |
| JWT signing keys        | Vault   | Runtime (CSI)                         | 90 days  |
| Model provider API keys | Vault   | Runtime (env)                         | 90 days  |
| Git credentials         | Vault   | Per-execution (injected into sandbox) | Per-use  |
| MCP server tokens       | Vault   | Runtime (MCP client)                  | 90 days  |
| Database passwords      | Vault   | Runtime (sidecar)                     | 30 days  |
| TLS certificates        | Vault   | Runtime (file)                        | 90 days  |

**Never in**: Code, config files, env vars (prod), logs, checkpoints, audit events

---

## Persistence Boundaries

### Schema Ownership (Single PostgreSQL Instance with Per-Module Schemas)

| Module   | PostgreSQL Schema | Tables                                      | Ownership                |
| -------- | ----------------- | ------------------------------------------- | ------------------------ |
| Gateway  | `gateway`         | `executions`, `idempotency_keys`            | Exclusive                |
| Router   | `router`          | `provider_registry`, `model_metadata`       | Exclusive                |
| Workflow | `workflow`        | `execution_records`, `workflow_definitions` | Exclusive                |
| Policy   | `policy`          | `capability_policies`, `tenant_allowlists`  | Exclusive                |
| Audit    | `audit`           | `audit_events`, `integrity_chain`           | Exclusive                |
| Shared   | `shared`          | `migrations`, `version_info`                | Shared (migrations only) |

### Access Rules

| Module   | Allowed Database Access            | Forbidden              |
| -------- | ---------------------------------- | ---------------------- |
| Gateway  | `gateway.*`, `shared.migrations`   | Any other schema       |
| Router   | `router.*`, `shared.migrations`    | Any other schema       |
| Workflow | `workflow.*`, `shared.migrations`  | Any other schema       |
| Policy   | `policy.*`, `shared.migrations`    | Any other schema       |
| Audit    | `audit.*`, `shared.migrations`     | Any other schema       |
| Harness  | None (uses object storage + Redis) | All PostgreSQL schemas |
| Sandbox  | None                               | All PostgreSQL schemas |

### Transaction Boundaries

- Each module owns its transactions. No distributed transactions across modules.
- Cross-module operations use eventual consistency via gRPC + audit events.
- Sagas/compensations are application-level (Phase 02: Temporal).

### Migration Ownership

- Each module owns its schema migrations in `migrations/<module>/`.
- `shared/migrations` contains only cross-cutting infrastructure (extensions, users).
- Migrations applied via CI/CD per module; no shared migration runner.

### Database Credentials

- Per-module PostgreSQL users with schema-scoped permissions (GRANT ON SCHEMA).
- Connection pooling per module (pgBouncer or native).
- Credentials injected via Vault at runtime; never in config.

---

## Redis State: Authoritative vs Ephemeral

### Authoritative State (PostgreSQL / Object Storage — Redis Loss = Safe)

| State               | Primary Store                             | Redis Role                |
| ------------------- | ----------------------------------------- | ------------------------- |
| Execution metadata  | PostgreSQL (`gateway.executions`)         | Cache / read acceleration |
| Agent checkpoints   | Object Storage (S3/GCS)                   | None                      |
| Workflow state      | PostgreSQL (`workflow.execution_records`) | None                      |
| Durable policy      | PostgreSQL (`policy.*`)                   | Cache                     |
| Audit records       | Immutable Object Storage                  | None                      |
| Capability policies | PostgreSQL (`policy.capability_policies`) | Cache                     |

### Ephemeral State (Redis Loss = Recoverable, No Data Loss)

| State                 | Redis Key Pattern                        | Recovery Behavior                                |
| --------------------- | ---------------------------------------- | ------------------------------------------------ |
| Rate-limit counters   | `ratelimit:{tenant}:{workflow}:{window}` | Rebuilt from request logs; brief burst tolerance |
| Budget reservations   | `budget:{tenant}:{workflow}:{request}`   | Re-evaluated on next request; no over-spend      |
| Ephemeral agent state | `agent:{agent_id}:state`                 | Reconstructed from checkpoint + execution log    |
| Distributed locks     | `lock:{resource}:{holder}`               | TTL-based expiry; auto-release                   |
| Session/token caches  | `session:{token_hash}`, `jwks:{issuer}`  | Refetched from source (Vault, JWKS)              |

### Redis Unavailability Behavior

| Scenario              | Behavior                                                                                                                                                        |
| --------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Redis unavailable     | Gateway: rate limiting fails open with warning; Router: budget checks use PostgreSQL fallback; Harness: agent state reconstructed from checkpoint; no data loss |
| Redis partial failure | Circuit breaker per Redis client; failover to direct PostgreSQL reads where applicable                                                                          |
| Redis data corruption | TTL-based keys auto-expire; caches rebuilt on next access; no authoritative data in Redis                                                                       |

**Invariant**: Loss of Redis MUST NOT cause loss of authoritative agent state. All authoritative state resides in PostgreSQL or immutable object storage.

---

## Workflow Boundary: Phase 01 vs Phase 02

### Phase 01 Execution Path (No Temporal Required)

```
Gateway
    ↓
WorkflowCoordinator (direct execution coordinator)
    ↓
DirectExecutionAdapter
    ↓
Harness
    ↓
Router / Sandbox
```

### Interface for Phase 02 Migration

```go
// Phase 01: DirectExecutionAdapter
type ExecutionAdapter interface {
    StartExecution(ctx context.Context, req StartExecutionRequest) (ExecutionHandle, error)
    GetExecution(ctx context.Context, id string) (ExecutionRecord, error)
    CancelExecution(ctx context.Context, id string) error
}

// Phase 02: TemporalWorkflowAdapter implements ExecutionAdapter
type TemporalWorkflowAdapter struct {
    // Temporal client, workflow definitions, signal handlers
}
```

### Checkpoint Abstraction (Phase 01 Ready for Phase 02)

```go
// Phase 01: DirectCheckpointStore
type CheckpointStore interface {
    Save(ctx context.Context, checkpoint Checkpoint) (CheckpointID, error)
    Load(ctx context.Context, id CheckpointID) (Checkpoint, error)
    Delete(ctx context.Context, id CheckpointID) error
}

// Phase 01 implementation: DirectCheckpointStore (object storage)
// Phase 02 implementation: TemporalCheckpointStore (Temporal history + replay)
```

### Temporal Dependency in Phase 01

- **NOT required** for normal Phase 01 execution.
- **Optional** for local development only (embedded Temporal in docker-compose for API compatibility testing).
- **Production Phase 01** runs without Temporal.
- Phase 02 introduces Temporal as the `TemporalWorkflowAdapter` and `TemporalCheckpointStore`.

## Observability Architecture

### Per-Module Requirements

| Module   | Logs | Metrics                    | Traces                        | Audit          |
| -------- | ---- | -------------------------- | ----------------------------- | -------------- |
| Gateway  | ✓    | ✓ (RED)                    | ✓ (inbound/outbound)          | ✓              |
| Sandbox  | ✓    | ✓ (executions, resources)  | ✓ (tool execution)            | ✓              |
| Router   | ✓    | ✓ (selections, budgets)    | ✓ (classification, selection) | ✓              |
| Harness  | ✓    | ✓ (lifecycle, checkpoints) | ✓ (agent lifecycle)           | ✓              |
| Workflow | ✓    | ✓ (executions)             | ✓ (execution flow)            | ✓              |
| Policy   | ✓    | ✓ (evaluations)            | ✓                             | ✓              |
| Audit    | N/A  | ✓ (event rate)             | N/A                           | N/A (is audit) |

### Correlation & Propagation

```
Headers: X-Correlation-ID, X-Request-ID, traceparent
Generation: UUID v7 (correlation), UUID v4 (request), W3C traceparent (trace)
Propagation: All gRPC/HTTP calls forward all three headers
```

### Health Endpoints (Per Service)

| Endpoint          | Checks                                              |
| ----------------- | --------------------------------------------------- |
| `/health/live`    | Process alive, no deadlock                          |
| `/health/ready`   | Dependencies reachable, config loaded, not draining |
| `/health/startup` | Initialization complete                             |

---

## Deployment Architecture

### Local (Docker Compose)

```yaml
services:
  gateway: # + shared, api
  harness: # + shared, api
  sandbox: # + shared, api (gVisor runtime)
  router: # + shared, api
  workflow: # + shared, api (DirectExecutionAdapter stub)
  policy: # + shared, api (allowlist stub)
  audit: # + shared, api
  postgres: # shared DB (per-module schemas)
  redis: # shared cache (ephemeral state only)
  vault: # dev mode
  otel-collector: # logs, metrics, traces
  prometheus: # metrics
  grafana: # dashboards
  loki: # logs
  tempo: # traces
```

**Note**: Temporal is NOT included in Phase 01 local development. The `workflow` service runs the `DirectExecutionAdapter`. Temporal is introduced in Phase 02.

### DEV/STAGE/PROD (Kubernetes)

- **Namespace per environment**
- **Multi-stage Docker builds** (builder → distroless runtime)
- **Resource requests/limits** per service (from config)
- **Network policies** (deny by default, explicit allow)
- **Service mesh**: mTLS (Istio/Linkerd), authorization policies
- **Secret injection**: Vault CSI driver / sidecar
- **Health probes**: `/health/live`, `/health/ready`, `/health/startup`
- **Autoscaling**: HPA (CPU, memory, custom metrics), VPA, Cluster Autoscaler
- **Rollout**: Rolling (DEV), Blue/Green (STAGE), Canary (PROD)

---

## API Contract Design Details

### Gateway REST API

```yaml
paths:
  /workflows/execute:
    post:
      summary: Submit workflow for execution
      security:
        - BearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/ExecuteWorkflowRequest"
      responses:
        "200":
          description: Execution started
          content:
            application/json:
              schema:
                $ref: "#/components/schemas/ExecuteWorkflowResponse"
        "401": { $ref: "#/components/responses/Unauthorized" }
        "403": { $ref: "#/components/responses/Forbidden" }
        "422": { $ref: "#/components/responses/ValidationError" }
        "429": { $ref: "#/components/responses/RateLimited" }
```

### Error Schema (Standardized)

```yaml
components:
  schemas:
    ErrorResponse:
      type: object
      required: [error]
      properties:
        error:
          type: object
          required: [code, message]
          properties:
            code:
              type: string
              enum:
                [
                  AUTH_INVALID,
                  AUTHZ_DENIED,
                  VALIDATION_ERROR,
                  RATE_LIMITED,
                  NOT_FOUND,
                  INTERNAL_ERROR,
                ]
            message:
              type: string
            details:
              type: object
            correlation_id:
              type: string
              format: uuid
```

### Idempotency

```yaml
headers:
  Idempotency-Key:
    type: string
    description: Unique key for idempotent execution (tenant+workflow scoped, 24h TTL)
    x-go-name: IdempotencyKey
```

---

## ADR Review for Phase 01

| Potential ADR                              | Decision Required                         | Covered by Existing ADR?          |
| ------------------------------------------ | ----------------------------------------- | --------------------------------- |
| API contract strategy (protobuf + OpenAPI) | Versioning, breaking changes, client gen  | No - **NEW ADR NEEDED**           |
| Capability model schema                    | Grant structure, types, constraints       | Partially (ADR-0004) - **EXTEND** |
| Sandbox runtime architecture               | gVisor vs Firecracker selection, fallback | ADR-0004 covers - **SUFFICIENT**  |
| Checkpoint abstraction                     | Direct vs Temporal store interface        | No - **NEW ADR NEEDED**           |
| Model provider abstraction                 | Interface, provider selection, budget     | No - **NEW ADR NEEDED**           |
| Policy engine choice                       | OPA vs Cedar vs custom                    | No - **NEW ADR NEEDED**           |
| MCP security boundary                      | Capability mapping, registry approval     | Partially (ADR-0004) - **EXTEND** |
| Persistence ownership                      | Per-module DB, no cross-access            | ADR-0002 covers - **SUFFICIENT**  |

**New ADRs Required for Phase 01 Plan Approval:**

1. **ADR-0005**: API Contract Strategy (protobuf versioning, OpenAPI generation, client generation)
2. **ADR-0006**: Capability Model Schema (grant structure, types, constraints, lifecycle)
3. **ADR-0007**: Checkpoint Store Abstraction (Direct vs Temporal, interface, migration)
4. **ADR-0008**: Model Provider Abstraction (interface, selection, budget, validation)
5. **ADR-0009**: Policy Engine Selection (OPA vs Cedar, integration timeline)

---

## Cross-Module Dependency Graph

```
shared/ (ROOT - no deps)
    │
    ├──► api/ (depends on shared)
    │
    ├──► audit/ (depends on shared, api)
    │
    ├──► policy/ (depends on shared, api)
    │
    ├──► gateway/ (depends on shared, api, workflow, policy, audit)
    │
    ├──► workflow/ (depends on shared, api, harness, router, audit)
    │
    ├──► router/ (depends on shared, api, policy)
    │
    ├──► sandbox/ (depends on shared, api)
    │
    └──► harness/ (depends on shared, api, sandbox, router, policy, audit)

FORBIDDEN: Any reverse or lateral dependencies not shown above
```

---

## Configuration per Module

| Module   | Config File                   | Key Sections                                                                   |
| -------- | ----------------------------- | ------------------------------------------------------------------------------ |
| Gateway  | `configs/{env}/gateway.yaml`  | server, auth, rate_limit, routing, workflow_adapter, observability             |
| Sandbox  | `configs/{env}/sandbox.yaml`  | runtime (gvisor/firecracker), seccomp, capabilities, tools, resources, network |
| Router   | `configs/{env}/router.yaml`   | classification, providers, budgets, validation, policy                         |
| Harness  | `configs/{env}/harness.yaml`  | lifecycle, tools, mcp, checkpoint_store, sandbox_client                        |
| Workflow | `configs/{env}/workflow.yaml` | adapter_mode (direct/temporal), execution, harness_client                      |
| Policy   | `configs/{env}/policy.yaml`   | rules, evaluation, storage                                                     |
| Audit    | `configs/{env}/audit.yaml`    | storage, integrity, retention, query                                           |

**Layered Config Priority**: Defaults → Config File → Env Vars → Vault → Feature Flags

---

## Metadata

---

title: Phase 01 Architecture
type: phase
phase: 01
status: Planning - Awaiting Plan Approval
author: Architect
date: 2024-10-07
reviewers: Architect, Platform Lead, Security Lead
approved_by: N/A - Awaiting Plan Approval
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---

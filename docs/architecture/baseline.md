# GiX-Coder System Architecture Baseline

## Architectural Vision

```
Developer
    ↓
CLI / IDE / Web / API
    ↓
Agent Gateway (Control Plane)
    ↓
Durable Workflow Engine
    ↓
Agent Harness
    ↓
Context + Policy + Model Router
    ↓
Isolated Execution Sandbox (Data Plane)
    ↓
Files / Shell / Git / Tests / MCP
    ↓
GitHub / GitLab / etc.
```

## Core Architectural Decisions

### 1. Control Plane / Data Plane Separation

**Decision**: Strict separation between control plane services and execution/data plane.

**Rationale**:

- Customer-controlled code must never execute inside control-plane services
- Enables independent scaling, security boundaries, and failure domains
- Supports zero-trust between execution boundaries

**Boundaries**:

- Control Plane: Agent Gateway, Workflow Engine, API Gateway, AuthZ/AuthN, Policy Engine
- Data Plane: Execution Sandboxes, Agent Harness, Tool Adapters, Model Router

### 2. Modular Monolith First

**Decision**: Start with a modular monolith with clear domain boundaries, not microservices.

**Rationale**:

- Avoids premature distributed system complexity
- Clear module boundaries enable future extraction
- Simpler operations, debugging, and testing
- Transactional consistency within domains

**Module Boundaries**:

```
gix-coder/
├── gateway/          # Agent Gateway - API, auth, routing
├── workflow/         # Durable Workflow Engine
├── harness/          # Agent Harness - agent lifecycle, tools
├── sandbox/          # Execution Sandbox - isolation, execution
├── router/           # Context/Policy/Model Router
├── policy/           # Policy Engine - authorization, quotas
├── audit/            # Audit Logging - immutable event log
└── shared/           # Shared kernel - DTOs, errors, utilities
```

### 3. Agent Execution Model

**Decision**: Agents execute in isolated sandboxes with explicit capability grants.

**Capabilities** (scoped, auditable, time-limited):

- File system access (read/write/list) - path-scoped
- Shell command execution - command allowlist, timeout
- Git operations - repo-scoped, branch-scoped
- Test execution - framework-agnostic
- MCP tool invocation - tool-scoped, approved registry
- Network access - egress allowlist only

### 4. Durable Workflow Foundation

**Decision**: All long-running operations use durable workflow engine.

**Characteristics**:

- Event-sourced state machine
- Automatic checkpointing
- Retry with exponential backoff
- Human-in-the-loop support
- Compensation/rollback for failures

### 5. Model Router & Policy

**Decision**: Centralized model routing with policy enforcement.

**Routing Factors**:

- Task complexity classification
- Cost optimization
- Latency requirements
- Tenant/model preferences
- Capability requirements

**Policy Enforcement**:

- Model allowlist per tenant
- Token budget limits
- PII/data classification guards
- Output validation

### 6. Configuration Strategy

**Decision**: Immutable configuration at runtime, environment-specific via layered config.

**Layers** (highest priority last):

1. Defaults (code)
2. Config files (versioned)
3. Environment variables
4. Secret manager (runtime)
5. Feature flags (runtime)

### 7. Observability Baseline

**Decision**: OpenTelemetry-native from day one.

**Signals**:

- Structured logs (JSON, correlation IDs)
- Metrics (RED: Rate, Errors, Duration)
- Traces (W3C TraceContext)
- Audit events (immutable, tamper-evident)

## Technology Choices (Baseline)

| Layer         | Technology                                              | Rationale                           |
| ------------- | ------------------------------------------------------- | ----------------------------------- |
| Language      | Go 1.25+ / TypeScript 5+                                | Performance, type safety, ecosystem |
| API           | gRPC + REST (OpenAPI 3.1)                               | Internal gRPC, external REST        |
| Workflow      | Temporal.io                                             | Durable execution, visibility       |
| Sandbox       | gVisor / Firecracker                                    | Strong isolation                    |
| Config        | Viper / dotenv                                          | Layered config                      |
| Secrets       | HashiCorp Vault / AWS Secrets Manager                   | Secret management                   |
| Observability | OpenTelemetry Collector → Prometheus/Grafana/Loki/Tempo | Vendor-neutral                      |
| CI/CD         | GitHub Actions                                          | Native GitHub integration           |
| Container     | Docker / containerd                                     | Standard runtime                    |
| Orchestration | Kubernetes (EKS/GKE)                                    | Production-grade                    |

## Domain Model (High-Level)

```
Tenant
    ├── Projects
    │   ├── Workflows
    │   │   ├── Executions
    │   │   └── Checkpoints
    │   ├── Agents
    │   │   ├── Capabilities
    │   │   └── Policies
    │   └── AuditLog
    ├── Users
    │   ├── Roles
    │   └── Permissions
    └── Models
        ├── Allowlist
        └── Quotas
```

## Security Architecture

### Threat Model Coverage

- Prompt injection → Input sanitization, output validation, capability scoping
- Malicious repositories → Sandbox isolation, network egress control
- Malicious dependencies → Dependency scanning, SBOM, allowlist
- Secret exfiltration → Secret detection, no secrets in control plane
- Command execution → Command allowlist, sandbox escape prevention
- Network abuse → Egress allowlist, rate limiting
- Cross-tenant access → Tenant isolation at data plane
- MCP tool abuse → Tool registry, capability grants
- Supply-chain attacks → SLSA, reproducible builds, sigstore
- Sandbox escape → gVisor/Firecracker, seccomp, capability dropping

### Zero Trust Boundaries

- Every service-to-service call authenticated (mTLS)
- Every agent capability explicitly granted
- Every data access audited
- No implicit trust between execution boundaries

## Data Flow

```
User Request
    → API Gateway (auth, rate limit, routing)
    → Agent Gateway (validation, tenant resolution)
    → Workflow Engine (durable execution start)
    → Agent Harness (agent instantiation)
    → Model Router (model selection, policy check)
    → Execution Sandbox (capability grant)
    → Tool Execution (file/shell/git/test/MCP)
    → Result Aggregation
    → Audit Log (immutable)
    → Response
```

## Non-Functional Requirements

| Requirement              | Target                                  |
| ------------------------ | --------------------------------------- |
| Availability             | 99.9%                                   |
| Latency (p99)            | < 500ms (gateway), < 5s (sandbox start) |
| Throughput               | 1000 concurrent executions              |
| Sandbox startup          | < 2s (warm), < 10s (cold)               |
| Audit log durability     | 99.9999999% (11 nines)                  |
| Recovery time objective  | < 10 min                                |
| Recovery point objective | < 1 min                                 |

## Evolution Path

```
Phase 00: Governance ✓
Phase 01: Gateway + Sandbox + Router + Harness (Modular Monolith)
Phase 02: Durable Workflows + Checkpointing
Phase 03: CLI/IDE/Web/API Interfaces
Phase 04: Multi-tenancy + RBAC + Compliance
Future:   Service extraction where justified by scale/team boundaries
```

---

## References

- [Product Charter](../product/charter.md)
- [Repository Structure](repository-structure.md)
- [Engineering Principles](engineering-principles.md)
- ADR-0001: Control Plane / Data Plane Separation
- ADR-0002: Modular Monolith
- ADR-0003: Durable Workflow Engine
- ADR-0004: Sandbox Isolation

---

## Metadata

---

title: GiX-Coder System Architecture Baseline
type: architecture
phase: 00
status: Verified
author: Architect
date: 2024-10-07
reviewers: Platform Lead, Security Lead, Product Owner
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---

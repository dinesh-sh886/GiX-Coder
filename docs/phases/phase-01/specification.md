# GiX-Coder Phase 01 Specification

## Phase Overview

**Phase 01**: Agent Gateway & Execution Sandbox Foundation
**Duration**: 4 weeks (target)
**Dependencies**: Phase 00 complete

## Goals

Deliver the core execution platform:
1. **Agent Gateway** - API entry point, auth, routing, rate limiting
2. **Execution Sandbox** - Isolated, capability-controlled execution environment
3. **Context/Policy/Model Router** - Intelligent model routing with policy enforcement
4. **Agent Harness** - Agent lifecycle, tool registry, MCP integration

## Scope

### In Scope

#### Agent Gateway (`gateway/`)
- REST/gRPC API for workflow submission
- Authentication (JWT, OIDC)
- Authorization (policy engine integration)
- Request routing to workflow engine
- Rate limiting (per tenant, per workflow)
- Request/response validation (DTOs)
- Observability (logs, metrics, traces, audit)

#### Execution Sandbox (`sandbox/`)
- gVisor-based isolation (primary) / Firecracker (alternative)
- Capability-based access control:
  - Filesystem (read/write/list, path-scoped)
  - Shell (command allowlist, timeout, resource limits)
  - Git (repo-scoped, branch-scoped operations)
  - Test execution (framework-agnostic)
  - Network (egress allowlist only)
- Resource limits (CPU, memory, pids, I/O, wall-time)
- Deterministic snapshot/restore for checkpointing
- Audit logging of all operations

#### Context/Policy/Model Router (`router/`)
- Task classification (complexity, domain, requirements)
- Model selection (cost, latency, capability, policy)
- Policy evaluation (OPA/Cedar integration)
- Context management (conversation, files, repo state)
- Token budget enforcement
- Output validation (schema, safety)

#### Agent Harness (`harness/`)
- Agent lifecycle (spawn, execute, checkpoint, terminate)
- Capability grant management
- Tool registry (built-in + MCP)
- MCP client implementation
- Agent-to-sandbox communication
- State persistence (checkpointing)

#### Shared Infrastructure (`shared/`)
- Typed error definitions
- Structured logging (zerolog/pino)
- Configuration management (Viper)
- Metrics utilities (Prometheus)
- Tracing utilities (OpenTelemetry)
- Validation utilities
- DTOs for all service boundaries

#### API Definitions (`api/`)
- Protobuf definitions for all services
- OpenAPI 3.1 for REST endpoints
- Generated clients (Go, TypeScript)
- Contract tests (Pact)

#### Deployment (`deploy/`)
- Kubernetes manifests (base + overlays)
- Dockerfiles (multi-stage, distroless)
- Docker Compose for local development
- Helm charts for production

### Out of Scope
- Durable Workflow Engine (Phase 02)
- CLI/IDE/Web interfaces (Phase 03)
- Multi-tenancy/RBAC (Phase 04)
- Advanced scheduling/orchestration
- Custom model fine-tuning
- Audit log UI

## Requirements

### Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| FR-01 | Submit workflow via REST API | P0 |
| FR-02 | Authenticate requests via JWT/OIDC | P0 |
| FR-03 | Authorize via policy engine | P0 |
| FR-04 | Route to workflow engine | P0 |
| FR-05 | Rate limit per tenant/workflow | P0 |
| FR-06 | Execute code in isolated sandbox | P0 |
| FR-07 | Grant file capabilities (path-scoped) | P0 |
| FR-08 | Grant shell capabilities (allowlist) | P0 |
| FR-09 | Grant Git capabilities (repo-scoped) | P0 |
| FR-10 | Grant test execution capabilities | P0 |
| FR-11 | Grant MCP tool capabilities | P0 |
| FR-12 | Enforce resource limits | P0 |
| FR-13 | Classify task for model routing | P0 |
| FR-14 | Select model per policy | P0 |
| FR-15 | Enforce token budgets | P0 |
| FR-16 | Validate model outputs | P0 |
| FR-17 | Spawn/manage agent lifecycle | P0 |
| FR-18 | Register/execute tools | P0 |
| FR-19 | Integrate MCP tools | P0 |
| FR-20 | Checkpoint/restore agent state | P1 |

### Non-Functional Requirements

| ID | Requirement | Target |
|----|-------------|--------|
| NFR-01 | Gateway p99 latency | < 100ms |
| NFR-02 | Sandbox cold start | < 10s |
| NFR-03 | Sandbox warm start | < 2s |
| NFR-04 | Concurrent executions | 1000 |
| NFR-05 | Availability | 99.9% |
| NFR-06 | Audit log durability | 11 nines |
| NFR-07 | Zero critical vulnerabilities | Continuous |
| NFR-08 | Sandbox escape resistance | gVisor/Firecracker verified |
| NFR-09 | Configuration immutability | Runtime |
| NFR-10 | Observability coverage | 100% endpoints |

## Architecture

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
                        CLIENT (CLI/IDE/API)
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
                    AGENT GATEWAY (gateway/)
  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
  │  Auth   │ │ Rate    │ │ Validate│ │  Route  │ │ Audit   │
  │  (JWT)  │ │ Limit   │ │  (DTO)  │ │(Workflow)│ │  Log    │
  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
                 DURABLE WORKFLOW ENGINE (workflow/)
  (Phase 02 - stubbed for Phase 01 with direct execution)
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
                      AGENT HARNESS (harness/)
  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
  │ Lifecycle│ │Capabilities│ │ Tool    │ │  MCP    │ │ Check-  │
  │ Manager │ │ Manager │ │ Registry│ │ Client  │ │ point   │
  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
            CONTEXT/POLICY/MODEL ROUTER (router/)
  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
  │ Classify│ │  Policy │ │  Model  │ │ Context │ │ Budget  │
  │ Task    │ │ Engine  │ │ Selector│ │ Manager │ │ Manager │
  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
                    EXECUTION SANDBOX (sandbox/)
  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
  │ Files   │ │ Shell   │ │ Git     │ │ Test    │ │ Network │
  │ Tool    │ │ Tool    │ │ Tool    │ │ Tool    │ │ Policy  │
  └─────────┘ └─────────┘ └─────────┘ └─────────┘ └─────────┘
  ┌─────────────────────────────────────────────────────────┐
  │              gVisor / Firecracker Isolation             │
  └─────────────────────────────────────────────────────────┘
└─────────────────────────────────────────────────────────────┘
```

### Data Flow

```
1. Client → POST /workflows/execute {workflow, input, capabilities}
2. Gateway → Validate JWT, check rate limit, validate DTO
3. Gateway → Audit log: workflow.submitted
4. Gateway → Forward to Workflow Engine (direct call for Phase 01)
5. Workflow → Create execution record
6. Workflow → Spawn Agent via Harness
7. Harness → Request capabilities from Policy Engine
8. Policy → Evaluate tenant/workflow caps, return grants
9. Harness → Initialize Agent with grants
10. Agent → Request model via Router
11. Router → Classify task, check policy, select model
12. Router → Return model endpoint + token budget
13. Agent → Execute steps via Sandbox tools
14. Sandbox → Enforce capabilities, resource limits
15. Sandbox → Return results
16. Agent → Checkpoint state (if long-running)
17. Harness → Aggregate results
18. Workflow → Complete execution record
19. Gateway → Return execution result to client
20. Audit → Log: workflow.completed (with full trace)
```

### Security Boundaries

```
┌─────────────────────────────────────────────────────────────┐
                     CONTROL PLANE (Trusted)
  ┌──────────┐ ┌────────────┐ ┌──────────┐ ┌────────────┐
  │ Gateway  │ │ Workflow   │ │ Router   │ │ Policy     │
  └──────────┘ └────────────┘ └──────────┘ └────────────┘
  ┌──────────┐ ┌────────────┐
  │ Harness  │ │ Audit      │
  └──────────┘ └────────────┘
  ▼                    ▼                    ▼
  mTLS              mTLS                 mTLS
  ▼                    ▼                    ▼
┌─────────────────────────────────────────────────────────────┐
                    DATA PLANE (Untrusted)
  ┌────────────────────────────────────────────────────────┐
  │  SANDBOX (gVisor)                                      │
  │  ┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐ ┌────────────┐  │
  │  │Files │ │Shell │ │ Git  │ │ Test │ │  Network   │  │
  │  └──────┘ └──────┘ └──────┘ └──────┘ └────────────┘  │
  │  Capability Grants (scoped, time-limited, audited)     │
  └────────────────────────────────────────────────────────┘
└─────────────────────────────────────────────────────────────┘
```

## Implementation Plan

### Week 1: Foundation & Gateway
- [ ] Shared kernel (errors, logging, config, metrics, tracing, validation, DTOs)
- [ ] API definitions (protobuf + OpenAPI)
- [ ] Gateway: Auth, rate limit, validation, routing stub
- [ ] Gateway: Observability (logs, metrics, traces, health)
- [ ] CI/CD pipeline for all modules
- [ ] Local dev environment (docker-compose)

### Week 2: Sandbox & Router
- [ ] Sandbox: gVisor integration, capability framework
- [ ] Sandbox: Filesystem tool (read/write/list)
- [ ] Sandbox: Shell tool (allowlist, limits)
- [ ] Sandbox: Resource enforcement (cgroups v2)
- [ ] Router: Task classification, policy integration
- [ ] Router: Model selection, budget enforcement
- [ ] Router: Output validation

### Week 3: Harness & Integration
- [ ] Harness: Agent lifecycle, capability grants
- [ ] Harness: Tool registry, built-in tools
- [ ] Harness: MCP client implementation
- [ ] Harness: Checkpoint/restore
- [ ] Integration: Gateway → Harness → Sandbox → Router
- [ ] End-to-end test: Simple workflow execution

### Week 4: Hardening & Validation
- [ ] Security hardening (seccomp, capabilities, network policies)
- [ ] Performance optimization (warm pools, connection reuse)
- [ ] Chaos testing (sandbox kill, network partition)
- [ ] Load testing (1000 concurrent)
- [ ] Documentation (API, runbooks, architecture)
- [ ] Phase verification & approval

## Acceptance Criteria

| ID | Criterion | Test Method |
|----|-----------|-------------|
| AC-01 | Gateway accepts valid workflow request | Integration test |
| AC-02 | Gateway rejects invalid JWT | Unit + Integration |
| AC-03 | Gateway enforces rate limits | Load test |
| AC-04 | Sandbox executes file read/write | Integration test |
| AC-05 | Sandbox blocks unauthorized paths | Security test |
| AC-06 | Sandbox executes allowed shell commands | Integration test |
| AC-07 | Sandbox blocks disallowed commands | Security test |
| AC-08 | Sandbox enforces CPU/memory limits | Stress test |
| AC-09 | Sandbox prevents network egress | Network test |
| AC-10 | Router classifies task correctly | Unit test (fixtures) |
| AC-11 | Router selects model per policy | Integration test |
| AC-12 | Router enforces token budget | Integration test |
| AC-13 | Harness spawns agent with grants | Integration test |
| AC-14 | Harness executes MCP tool | Integration test |
| AC-15 | Checkpoint/restore works | Integration test |
| AC-16 | Full workflow executes end-to-end | E2E test |
| AC-17 | Audit log captures all events | Verification |
| AC-18 | All quality gates pass | CI/CD pipeline |
| AC-19 | No critical/high vulnerabilities | Security scan |
| AC-20 | Documentation complete | Review |

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|------------|--------|------------|
| gVisor compatibility issues | Medium | High | Prototype early, Firecracker fallback |
| MCP protocol changes | Low | Medium | Pin version, adapter pattern |
| Model router complexity | Medium | Medium | Start simple, iterate |
| Sandbox performance | Medium | High | Warm pool, benchmark early |
| Policy engine integration | Low | High | Define interface early, mock |
| Cross-module integration | Medium | High | Contract tests, daily integration |

## Quality Gates (Phase 01 Specific)

Additional to standard gates:
- [ ] Sandbox escape tests pass (gVisor test suite)
- [ ] Capability enforcement tests pass (100% coverage)
- [ ] Resource limit enforcement verified (stress test)
- [ ] Contract tests for all service boundaries pass
- [ ] Load test: 1000 concurrent, p99 < 500ms gateway
- [ ] Chaos test: Sandbox kill → graceful degradation
- [ ] Security review completed (threat model + code)

## Deliverables

### Code
- `gateway/` - Complete, tested, documented
- `sandbox/` - Complete, tested, documented
- `router/` - Complete, tested, documented
- `harness/` - Complete, tested, documented
- `shared/` - Complete, tested, documented
- `workflow/` - Stub for integration
- `policy/` - Stub for integration
- `audit/` - Complete, tested, documented
- `api/` - Complete protobuf + OpenAPI

### Documentation
- `docs/phases/phase-01/README.md`
- `docs/phases/phase-01/requirements.md`
- `docs/phases/phase-01/architecture.md`
- `docs/phases/phase-01/plan.md`
- `docs/phases/phase-01/acceptance-criteria.md`
- `docs/phases/phase-01/implementation.md`
- `docs/phases/phase-01/verification.md`
- `docs/phases/phase-01/remediation.md` (if needed)
- `docs/phases/phase-01/approval.md`
- Module READMEs
- API documentation
- Runbooks (deploy, debug, scale)

### Infrastructure
- Kubernetes manifests (dev, stage, prod overlays)
- Docker images (signed, SBOM attested)
- Helm charts
- Docker Compose (local)
- CI/CD pipelines

## Dependencies

### External Services (Stubs for Phase 01)
- **Temporal** - Workflow engine (embedded for dev, external for stage/prod)
- **PostgreSQL** - State persistence
- **Redis** - Caching, rate limiting
- **Vault** - Secrets
- **Model Providers** - OpenAI, Anthropic, local (Ollama)
- **MCP Registry** - Tool registry

### Internal Dependencies
- Phase 00 governance baseline (complete)
- Shared kernel (built first)

## Success Metrics

- All acceptance criteria PASS
- Zero critical/high security findings
- Performance targets met
- Documentation completeness > 90%
- Team confidence > 4/5 (survey)

---

## Metadata

---
title: Phase 01 Specification
type: phase
phase: 01
status: Specified - Awaiting Phase 00 Completion
author: Architect
date: 2024-10-07
reviewers: Architect, Platform Lead, Security Lead, Product Owner
approved_by: N/A - Phase 00 Pending
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---
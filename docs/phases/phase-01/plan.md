# Phase 01 Implementation Plan

## Overview

Dependency-aware implementation plan for Phase 01: Agent Gateway & Execution Sandbox Foundation.

**Duration**: 8 weeks (target)
**Team**: 5-8 engineers (Architect, Platform, Security, Backend, QA)
**Methodology**: Milestone-based, PR-driven, CI-gated

---

## Milestone Summary

9 milestones over 8 calendar weeks. M0 is sequential; M1-M5 are sequential on the critical path; M6, M7, M8 run in parallel after M5 completes.

| Milestone | Name                    | Calendar Weeks | Duration          | Key Deliverables                                                                         |
| --------- | ----------------------- | -------------- | ----------------- | ---------------------------------------------------------------------------------------- |
| M0        | Contracts & Foundations | Week 1         | 1 week            | API contracts, shared kernel, CI/CD, local dev env                                       |
| M1        | Gateway Core            | Week 2         | 1 week            | Auth, rate limit, validation, routing stub, observability                                |
| M2        | Sandbox Core            | Weeks 3-4      | 2 weeks           | gVisor integration, capability framework, filesystem/shell tools, resource enforcement   |
| M3        | Router & Policy         | Week 5         | 1 week            | Task classification, model selection, budget enforcement, output validation, policy stub |
| M4        | Harness & Integration   | Week 6         | 1 week            | Agent lifecycle, tool registry, MCP client, checkpoint store, sandbox integration        |
| M5        | End-to-End Integration  | Week 7         | 1 week            | Full workflow execution, contract tests, cross-module integration                        |
| M6        | Security Hardening      | Week 8         | 1 week (parallel) | Seccomp, capabilities, network policies, escape tests, audit completeness                |
| M7        | Performance & Chaos     | Week 8         | 1 week (parallel) | Load testing, chaos experiments, warm pools, benchmark baselines                         |
| M8        | Release Readiness       | Week 8         | 1 week (parallel) | Documentation, runbooks, deployment validation, rollback testing                         |

---

## Detailed Work Breakdown

### Milestone 0: Contracts & Foundations (Week 1)

| ID    | Task                                                                                                          | Module   | Dependencies | Owner     | Effort | Test Requirement              | Acceptance Criterion                                     |
| ----- | ------------------------------------------------------------------------------------------------------------- | -------- | ------------ | --------- | ------ | ----------------------------- | -------------------------------------------------------- |
| M0.1  | Define Protobuf schemas for all 7 services                                                                    | api      | -            | Architect | 3d     | `buf lint`, `buf breaking`    | All `.proto` files compile, no breaking changes vs empty |
| M0.2  | Generate Go/TS clients from protobuf                                                                          | api      | M0.1         | Platform  | 2d     | Client generation test        | Generated code commits clean, imports work               |
| M0.3  | Generate OpenAPI 3.1 from protobuf                                                                            | api      | M0.1         | Platform  | 1d     | Schema validation             | OpenAPI spec valid, served at `/openapi.json`            |
| M0.4  | Implement shared/errors with typed error codes                                                                | shared   | -            | Backend   | 1d     | Unit tests 90%+               | All error codes defined, wrapping works                  |
| M0.5  | Implement shared/logging (zerolog/pino wrappers)                                                              | shared   | -            | Backend   | 1d     | Unit tests                    | Structured JSON output, correlation ID propagation       |
| M0.6  | Implement shared/config (Viper, layered, validation)                                                          | shared   | -            | Backend   | 2d     | Unit tests, config validation | All layers work, schema validation passes                |
| M0.7  | Implement shared/metrics (Prometheus helpers, standard buckets)                                               | shared   | -            | Platform  | 1d     | Unit tests                    | RED metrics helpers, histogram buckets per standards     |
| M0.8  | Implement shared/tracing (OTel, W3C propagation)                                                              | shared   | -            | Platform  | 2d     | Unit tests                    | Trace context propagation, sampling config               |
| M0.9  | Implement shared/validation (DTO validation, sanitization)                                                    | shared   | -            | Backend   | 1d     | Unit tests                    | Input validation at boundaries, PII detection            |
| M0.10 | Implement shared/dto (Identity, CapabilityGrant, etc.)                                                        | shared   | M0.1         | Backend   | 1d     | Unit tests                    | DTOs compile, serialization works                        |
| M0.11 | Implement shared/security (crypto, sanitization, PII)                                                         | shared   | -            | Security  | 2d     | Unit tests, security tests    | PII detection works, sanitization prevents injection     |
| M0.12 | Setup CI/CD pipeline for all modules                                                                          | deploy   | -            | Platform  | 3d     | Pipeline runs                 | All stages pass (validate, test, security, build)        |
| M0.13 | Create docker-compose for local development                                                                   | deploy   | M0.1         | Platform  | 2d     | `make dev-up` works           | All services start, health checks pass                   |
| M0.14 | Create Kubernetes base manifests + overlays                                                                   | deploy   | M0.1         | Platform  | 3d     | `kubectl apply --dry-run`     | Manifests valid, kustomize builds                        |
| M0.15 | Define ADR-0005 through ADR-0009 (API contracts, capability model, checkpoint, model provider, policy engine) | docs/adr | -            | Architect | 2d     | ADR review process            | ADRs proposed, reviewed, accepted                        |

**M0 Exit Criteria**: All protobuf schemas defined, shared kernel complete, CI/CD runs, local dev environment works, ADRs accepted.

---

### Milestone 1: Gateway Core (Week 2)

| ID    | Task                                                      | Module  | Dependencies            | Owner    | Effort | Test Requirement                                                        | Acceptance Criterion                                             |
| ----- | --------------------------------------------------------- | ------- | ----------------------- | -------- | ------ | ----------------------------------------------------------------------- | ---------------------------------------------------------------- |
| M1.1  | Implement JWT/OIDC authentication middleware              | gateway | M0.4, M0.5, M0.6, M0.10 | Backend  | 2d     | Unit: valid/expired/wrong issuer/audience; Integration: real JWKS       | AC-02: Invalid JWT rejected, valid extracts identity             |
| M1.2  | Implement rate limiter (token bucket, Redis-backed)       | gateway | M0.6, M0.7              | Backend  | 2d     | Unit: algorithm; Integration: concurrent requests; Chaos: Redis failure | AC-03: Rate limit enforced at threshold                          |
| M1.3  | Implement request/response DTO validation                 | gateway | M0.9, M0.10             | Backend  | 1d     | Unit: schema validation; Integration: malformed requests                | AC-01: Valid request accepted, invalid returns 422               |
| M1.4  | Implement policy authorizer (capability grant evaluation) | gateway | M0.10, policy stub      | Backend  | 2d     | Unit: grant evaluation; Integration: gateway→policy                     | Capability grants returned, audit logged                         |
| M1.5  | Implement workflow routing (direct adapter stub)          | gateway | workflow stub           | Backend  | 1d     | Integration: gateway→workflow                                           | Execution handle returned                                        |
| M1.6  | Implement observability (logs, metrics, traces, health)   | gateway | M0.5, M0.7, M0.8        | Platform | 1d     | Unit: metric emission; Integration: health endpoints                    | `/health/live`, `/health/ready`, RED metrics exposed             |
| M1.7  | Implement audit logging for gateway events                | gateway | audit stub              | Backend  | 1d     | Integration: audit events emitted                                       | `workflow.submitted`, `auth.success/failure`, `authz.allow/deny` |
| M1.8  | Implement idempotency key support                         | gateway | Redis                   | Backend  | 1d     | Unit: key storage/retrieval; Integration: duplicate requests            | Duplicate key returns original response                          |
| M1.9  | Contract tests for Gateway REST/gRPC API                  | gateway | M0.3                    | QA       | 2d     | Pact/Schemathesis                                                       | Provider verification passes                                     |
| M1.10 | Architecture tests (import boundaries, no cycles)         | gateway | M0.14                   | QA       | 1d     | CI arch test passes                                                     | No forbidden imports, no cycles                                  |

**M1 Exit Criteria**: Gateway accepts valid workflow request, rejects invalid JWT, enforces rate limits, routes to workflow adapter, all quality gates pass.

---

### Milestone 2: Sandbox Core (Weeks 3-4)

| ID    | Task                                                                  | Module  | Dependencies | Owner             | Effort | Test Requirement                                                                                        | Acceptance Criterion                                                    |
| ----- | --------------------------------------------------------------------- | ------- | ------------ | ----------------- | ------ | ------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| M2.1  | gVisor integration (runsc runtime, session management)                | sandbox | -            | Security/Backend  | 3d     | Integration: sandbox create/execute/destroy                                                             | Sandbox creates, executes command, destroys cleanly                     |
| M2.2  | Firecracker fallback implementation                                   | sandbox | M2.1         | Security/Backend  | 2d     | Integration: runtime selection via feature flag                                                         | Feature flag switches runtime, both pass tests                          |
| M2.3  | seccomp profiles per capability                                       | sandbox | M2.1         | Security          | 2d     | Security: syscall audit, escape attempts                                                                | Only allowed syscalls execute, audit logs show syscalls                 |
| M2.4  | Capability framework (grant validation, constraints, expiry)          | sandbox | M0.10        | Security/Backend  | 2d     | Unit: grant parsing; Security: forged grants                                                            | Invalid/expired grants denied, audit logged                             |
| M2.5  | Filesystem tool (read/write/list, path scoping, traversal prevention) | sandbox | M2.4         | Security/Backend  | 3d     | Unit: path canonicalization; Security: traversal, symlink, absolute path, oversized file, special files | AC-04: Authorized paths work; AC-05: Unauthorized blocked               |
| M2.6  | Shell tool (allowlist, argument validation, injection prevention)     | sandbox | M2.4         | Security/Backend  | 3d     | Unit: command validation; Security: injection, unauthorized cmd, resource exhaustion                    | AC-06: Allowed commands work; AC-07: Disallowed blocked                 |
| M2.7  | Git tool (repo/branch scope, credential injection)                    | sandbox | M2.4, Vault  | Backend           | 2d     | Integration: clone/fetch/checkout/diff/status; Security: credential leakage                             | Git operations within scope succeed                                     |
| M2.8  | Test tool (framework allowlist, resource limits)                      | sandbox | M2.4         | Backend           | 1d     | Integration: go/jest/pytest execution                                                                   | Test execution works, limits enforced                                   |
| M2.9  | Network tool (egress allowlist, DNS policy, metadata protection)      | sandbox | M2.4         | Security          | 2d     | Security: SSRF, private network, metadata service                                                       | AC-09: Egress to allowed destinations only                              |
| M2.10 | Resource enforcement (cgroups v2: CPU, memory, pids, I/O, wall-time)  | sandbox | M2.1         | Security/Platform | 2d     | Stress: OOM, CPU throttle, pid limit, I/O limit                                                         | AC-08: Limits enforced with termination                                 |
| M2.11 | Sandbox gRPC API implementation                                       | sandbox | M2.1         | Backend           | 1d     | Contract tests                                                                                          | `Execute`, `CreateSession`, `DestroySession` work                       |
| M2.12 | Sandbox audit logging (all operations)                                | sandbox | M2.4, audit  | Backend           | 1d     | Integration: audit events for all tools                                                                 | `sandbox.execution`, `capability.grant/deny`, `resource_limit.exceeded` |
| M2.13 | Sandbox contract tests                                                | sandbox | M0.3         | QA                | 1d     | Pact/Schemathesis                                                                                       | Provider verification passes                                            |
| M2.14 | Architecture tests (import boundaries)                                | sandbox | M0.14        | QA                | 1d     | CI arch test passes                                                                                     | No forbidden imports                                                    |

**M2 Exit Criteria**: Sandbox executes code in gVisor, capability grants enforced, filesystem/shell/git/test tools work, resource limits enforced, network egress controlled, all security tests pass.

---

### Milestone 3: Router & Policy (Week 5)

| ID    | Task                                                                | Module | Dependencies | Owner    | Effort | Test Requirement                                                   | Acceptance Criterion                                                |
| ----- | ------------------------------------------------------------------- | ------ | ------------ | -------- | ------ | ------------------------------------------------------------------ | ------------------------------------------------------------------- |
| M3.1  | Task classifier (complexity, domain, capabilities, tokens, latency) | router | M0.11        | Backend  | 2d     | Unit: fixtures for each type; Integration: classifier accuracy     | AC-10: Correct classification on test fixtures                      |
| M3.2  | Policy evaluator stub (tenant allowlist, workflow requirements)     | policy | M0.10        | Backend  | 2d     | Unit: allowlist intersection; Integration: policy→grants           | Capability grants match policy                                      |
| M3.3  | Model selector (cost, latency, capability match, provider health)   | router | M3.1, M0.10  | Backend  | 2d     | Unit: selection logic with policies; Integration: router→providers | AC-11: Model selected per policy                                    |
| M3.4  | Budget manager (tenant/workflow/request hierarchy, Redis)           | router | M0.6         | Backend  | 2d     | Unit: budget calc; Integration: concurrent exhaustion              | AC-12: Budget enforced at all levels                                |
| M3.5  | Provider adapters (OpenAI, Anthropic, Ollama)                       | router | M0.1         | Backend  | 2d     | Integration: real provider calls (mocked in CI)                    | All three adapters work, health checks pass                         |
| M3.5  | Output validator (schema, safety, PII)                              | router | M0.11        | Security | 2d     | Unit: validators; Integration: model response→validator            | Schema validation, PII detection, safety rejection                  |
| M3.7  | Router gRPC API implementation                                      | router | M3.1-M3.6    | Backend  | 1d     | Contract tests                                                     | `ClassifyTask`, `SelectModel`, `ValidateOutput`, `CheckBudget` work |
| M3.8  | Router audit logging                                                | router | audit        | Backend  | 1d     | Integration: audit events                                          | `model.selected`, `budget.allow/deny`, `output.validated`           |
| M3.9  | Router contract tests                                               | router | M0.3         | QA       | 1d     | Pact/Schemathesis                                                  | Provider verification passes                                        |
| M3.10 | Architecture tests                                                  | router | M0.14        | QA       | 1d     | CI arch test passes                                                | No forbidden imports                                                |

**M3 Exit Criteria**: Task classification works, model selection per policy, budget enforcement at all levels, output validation (schema/safety/PII), provider adapters functional.

---

### Milestone 4: Harness & Integration (Week 6)

| ID    | Task                                                                        | Module  | Dependencies         | Owner   | Effort | Test Requirement                                                          | Acceptance Criterion                                                                         |
| ----- | --------------------------------------------------------------------------- | ------- | -------------------- | ------- | ------ | ------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------- |
| M4.1  | Agent lifecycle state machine (CREATED→TERMINATED/FAILED)                   | harness | M0.10                | Backend | 2d     | Unit: state transitions; Integration: full lifecycle                      | AC-13: All valid transitions work, invalid rejected                                          |
| M4.2  | Capability grant management (request, bind, revoke)                         | harness | M2.4                 | Backend | 1d     | Integration: harness→policy→sandbox                                       | Grants bound to agent lifecycle                                                              |
| M4.3  | Tool registry (registration, schema, capability mapping)                    | harness | M0.10, M2.5-M2.9     | Backend | 2d     | Unit: registry ops; Integration: tool execution via registry              | Tools registered, capabilities verified before execution                                     |
| M4.4  | Built-in tool executors (filesystem, shell, git, test via sandbox)          | harness | M2.5-M2.9            | Backend | 2d     | Integration: harness→sandbox for each tool                                | AC-04, AC-06, AC-08 via harness                                                              |
| M4.5  | MCP client implementation (server registration, tool discovery, invocation) | harness | M2.11                | Backend | 2d     | Integration: approved server tool invocation; Security: unapproved denied | AC-14: MCP tool executes with grant                                                          |
| M4.6  | DirectCheckpointStore (serialized state → object storage, encrypted)        | harness | M0.6, object storage | Backend | 2d     | Unit: serialize/deserialize; Integration: checkpoint→restore→continue     | AC-15: Checkpoint/restore works                                                              |
| M4.7  | Harness gRPC API implementation                                             | harness | M4.1-M4.6            | Backend | 1d     | Contract tests                                                            | `SpawnAgent`, `GetAgentState`, `ExecuteStep`, `Checkpoint`, `Restore`, `TerminateAgent` work |
| M4.8  | Harness audit logging                                                       | harness | audit                | Backend | 1d     | Integration: audit events                                                 | `agent.spawned`, `agent.state_change`, `tool.invoked`, `checkpoint.created/restored`         |
| M4.9  | Harness contract tests                                                      | harness | M0.3                 | QA      | 1d     | Pact/Schemathesis                                                         | Provider verification passes                                                                 |
| M4.10 | Architecture tests                                                          | harness | M0.14                | QA      | 1d     | CI arch test passes                                                       | No forbidden imports                                                                         |

**M4 Exit Criteria**: Agent lifecycle complete, tool registry with capabilities, MCP client works, checkpoint/restore functional, all integrations verified.

---

### Milestone 5: End-to-End Integration (Week 7)

| ID   | Task                                                                     | Module   | Dependencies                  | Owner    | Effort | Test Requirement                         | Acceptance Criterion                |
| ---- | ------------------------------------------------------------------------ | -------- | ----------------------------- | -------- | ------ | ---------------------------------------- | ----------------------------------- |
| M5.1 | Gateway → Workflow → Harness → Sandbox → Router full flow                | all      | M1, M2, M3, M4                | Backend  | 2d     | Integration: full workflow execution     | AC-16: E2E test passes              |
| M5.2 | Contract tests for all service boundaries                                | all      | M0.3, M1.9, M2.13, M3.9, M4.9 | QA       | 2d     | Pact/Schemathesis all pass               | All provider verifications pass     |
| M5.3 | Cross-module integration tests (gateway→workflow→harness→sandbox→router) | all      | M5.1                          | QA       | 2d     | Integration test suite                   | All critical paths covered          |
| M5.4 | Workflow adapter direct execution coordination                           | workflow | M1.5, M4.7                    | Backend  | 1d     | Integration: workflow→harness            | Execution record created, completed |
| M5.5 | Audit log completeness verification                                      | audit    | M1.7, M2.12, M3.8, M4.8       | Security | 1d     | Verification: all 24 event types present | AC-17: All audit events captured    |
| M5.6 | Quality gates validation (all CI stages)                                 | all      | -                             | Platform | 1d     | Full CI pipeline                         | AC-18: All gates pass               |
| M5.7 | Security scan (SAST, SCA, container, secrets)                            | all      | -                             | Security | 1d     | CI security stage                        | AC-19: Zero critical/high findings  |

**M5 Exit Criteria**: Full workflow executes end-to-end, all contract tests pass, audit log complete, all quality gates pass, security scans clean.

---

### Milestone 6: Security Hardening (Week 8)

| ID   | Task                                                               | Module  | Dependencies | Owner             | Effort | Test Requirement                           | Acceptance Criterion               |
| ---- | ------------------------------------------------------------------ | ------- | ------------ | ----------------- | ------ | ------------------------------------------ | ---------------------------------- |
| M6.1 | Seccomp profile hardening (minimal syscalls per tool)              | sandbox | M2.3         | Security          | 1d     | gVisor test suite, escape attempts         | gVisor test suite passes           |
| M6.2 | Linux capability dropping (minimal set per tool)                   | sandbox | M2.1         | Security          | 1d     | `capsh --print` verification               | Only required capabilities present |
| M6.3 | Network policy enforcement (CNI, egress allowlist, no ingress)     | sandbox | M2.9         | Security/Platform | 1d     | Network tests: SSRF, private net, metadata | No unauthorized network access     |
| M6.4 | Sandbox escape test suite (gVisor tests + custom)                  | sandbox | M6.1-M6.3    | Security          | 2d     | Automated escape attempts                  | Zero successful escapes            |
| M6.5 | Capability enforcement 100% coverage test                          | sandbox | M2.4         | Security          | 1d     | All capability types tested                | 100% capability paths covered      |
| M6.6 | Resource limit stress test (CPU, memory, pids, I/O)                | sandbox | M2.10        | Security/Platform | 1d     | Stress: exceed limits                      | Limits enforced, termination clean |
| M6.7 | Credential handling audit (no leakage in logs, checkpoints, audit) | all     | M2.7, M4.6   | Security          | 1d     | Log scan, checkpoint scan, audit scan      | Zero secrets in any output         |
| M6.8 | mTLS verification (all control plane communication)                | all     | deploy       | Platform          | 1d     | `openssl s_client` verification            | All gRPC connections use mTLS      |
| M6.9 | Security review sign-off                                           | docs    | M6.1-M6.8    | Security          | 1d     | Review checklist                           | Security Lead signs approval.md    |

**M6 Exit Criteria**: gVisor test suite passes, capability enforcement 100% coverage, resource limits verified, escape tests pass, zero secrets leaked, mTLS verified, security review complete.

---

### Milestone 7: Performance & Chaos (Week 8, Parallel)

| ID   | Task                                                             | Module  | Dependencies | Owner       | Effort | Test Requirement               | Acceptance Criterion                        |
| ---- | ---------------------------------------------------------------- | ------- | ------------ | ----------- | ------ | ------------------------------ | ------------------------------------------- |
| M7.1 | Warm pool implementation (sandbox pre-warming)                   | sandbox | M2.1         | Platform    | 2d     | Benchmark: cold vs warm start  | NFR-02: cold < 10s, NFR-03: warm < 2s       |
| M7.2 | Connection pooling (gRPC, HTTP, Redis, DB)                       | all     | M0.6, M0.7   | Platform    | 1d     | Load test: connection reuse    | No connection exhaustion at 1000 concurrent |
| M7.3 | Load test: 1000 concurrent workflows, p99 gateway < 100ms        | gateway | M5.1         | Platform/QA | 2d     | k6 script, STAGE env           | NFR-01, NFR-04 met                          |
| M7.4 | Sandbox stress test (concurrent executions, resource contention) | sandbox | M6.6         | Platform/QA | 1d     | k6 script                      | 1000 concurrent sandbox executions          |
| M7.5 | Chaos test: sandbox kill → graceful degradation                  | sandbox | M5.1         | Platform/QA | 1d     | Chaos Mesh: pod kill           | Agent recovers or fails cleanly             |
| M7.6 | Chaos test: network partition (sandbox ↔ control plane)          | all     | M5.1         | Platform/QA | 1d     | Chaos Mesh: network partition  | Circuit breakers work, no cascade           |
| M7.7 | Chaos test: Redis/PostgreSQL/Vault failure                       | all     | M5.1         | Platform/QA | 1d     | Chaos Mesh: dependency failure | Graceful degradation, proper errors         |
| M7.8 | Benchmark baseline establishment                                 | all     | M7.1-M7.7    | Platform    | 1d     | Documented baselines           | Performance baselines recorded              |

**M7 Exit Criteria**: NFR-01 through NFR-04 met, warm pools working, chaos tests pass, baselines established.

---

### Milestone 8: Release Readiness (Week 8, Final)

| ID   | Task                                                        | Module | Dependencies | Owner    | Effort | Test Requirement    | Acceptance Criterion                      |
| ---- | ----------------------------------------------------------- | ------ | ------------ | -------- | ------ | ------------------- | ----------------------------------------- |
| M8.1 | Module READMEs (purpose, API, config, deps, run/deploy)     | all    | M1-M4        | Backend  | 1d     | Review              | All 8 modules have README                 |
| M8.2 | API documentation (OpenAPI served, examples)                | api    | M0.3         | Platform | 1d     | Review              | `/openapi.json` valid, examples render    |
| M8.3 | Runbooks (deploy, debug, scale, rollback)                   | deploy | M0.12, M0.14 | Platform | 2d     | Review              | Runbooks complete, tested                 |
| M8.4 | Deployment validation (DEV auto, STAGE manual, PROD canary) | deploy | M0.14        | Platform | 1d     | Deploy to DEV/STAGE | DEV auto-deploy works, STAGE manual works |
| M8.5 | Rollback test (revert merge, verify health)                 | deploy | M8.4         | Platform | 1d     | Manual rollback     | Rollback < 2 min, health verified         |
| M8.6 | Phase 01 documentation completeness review                  | docs   | M1-M8        | All      | 1d     | Checklist > 90%     | AC-20: Documentation complete             |
| M8.7 | Final CI pipeline validation (all stages)                   | all    | -            | Platform | 1d     | Full CI run         | All stages green                          |
| M8.8 | Phase verification preparation                              | docs   | M8.6         | Verifier | 1d     | Verification plan   | Verification.md template ready            |

**M8 Exit Criteria**: All documentation complete, deployment validated, rollback tested, CI green, ready for independent verification.

---

## Dependency Graph

```
M0 (Contracts & Foundations, Week 1)
    │
    ▼
M1 (Gateway Core, Week 2)
    │
    ▼
M2 (Sandbox Core, Weeks 3-4)
    │
    ▼
M3 (Router & Policy, Week 5)
    │
    ▼
M4 (Harness & Integration, Week 6)
    │
    ▼
M5 (End-to-End Integration, Week 7)
    │
    ├──► M6 (Security Hardening, Week 8, parallel)
    ├──► M7 (Performance & Chaos, Week 8, parallel)
    └──► M8 (Release Readiness, Week 8, parallel)
```

## Critical Path

```
M0 (Week 1) → M1 (Week 2) → M2 (Weeks 3-4) → M3 (Week 5) → M4 (Week 6) → M5 (Week 7) → M6/M7/M8 (Week 8 parallel)
```

**Critical Path Duration**: 8 calendar weeks (no slack on M0-M5 sequential chain)

**Parallel Opportunities**:

- M6, M7, M8 run in parallel during Week 8 after M5 completes
- M2.1 and M2.2 can be parallelized (gVisor + Firecracker)
- M3.1-M3.6 can be parallelized (classifier, policy, selector, budget, providers, validator)
- M4.1-M4.6 can be parallelized (lifecycle, grants, registry, tools, MCP, checkpoint)
- M1.1-M1.4 can be parallelized (auth, rate limit, validation, authorizer)

## Calendar Schedule

| Calendar Week | Milestones Active       | Notes                                                       |
| ------------- | ----------------------- | ----------------------------------------------------------- |
| Week 1        | M0                      | Contracts, shared kernel, CI/CD, ADRs                       |
| Week 2        | M1                      | Gateway auth, rate limit, routing stub                      |
| Week 3        | M2 (part 1)             | gVisor integration, capability framework                    |
| Week 4        | M2 (part 2)             | Filesystem/shell/git/test tools, resource enforcement       |
| Week 5        | M3                      | Router classification, model selection, budget, policy stub |
| Week 6        | M4                      | Harness lifecycle, tool registry, MCP, checkpoint           |
| Week 7        | M5                      | E2E integration, contract tests, audit verification         |
| Week 8        | M6 + M7 + M8 (parallel) | Security hardening, perf/chaos, release readiness           |

---

## Risk Register

| Risk                                     | Likelihood | Impact   | Mitigation                                          | Owner     |
| ---------------------------------------- | ---------- | -------- | --------------------------------------------------- | --------- |
| gVisor compatibility issues              | Medium     | High     | Prototype early (M2.1), Firecracker fallback (M2.2) | Security  |
| gVisor/Firecracker performance           | Medium     | High     | Warm pools (M7.1), benchmark early                  | Platform  |
| MCP protocol changes                     | Low        | Medium   | Pin version, adapter pattern                        | Backend   |
| Model router complexity                  | Medium     | Medium   | Start simple (M3.1-M3.3), iterate                   | Backend   |
| Policy engine integration                | Low        | High     | Define interface early (ADR-0009), mock in M1/M3    | Backend   |
| Cross-module integration failures        | Medium     | High     | Contract tests daily (M5.2), shared kernel first    | QA        |
| Sandbox escape vulnerability             | Low        | Critical | Defense in depth (M6.1-M6.4), security review       | Security  |
| Checkpoint/restore data corruption       | Low        | High     | Integrity verification (M4.6), encryption           | Backend   |
| Temporal migration complexity (Phase 02) | Medium     | Medium   | Clear abstraction (ADR-0007), stub in M4            | Architect |

---

## Resource Requirements

| Role              | Week 1 | Week 2 | Week 3-4 | Week 5 | Week 6 | Week 7 | Week 8 |
| ----------------- | ------ | ------ | -------- | ------ | ------ | ------ | ------ |
| Architect         | 1.0    | 0.5    | 0.5      | 0.5    | 0.5    | 0.5    | 0.5    |
| Platform Lead     | 1.0    | 0.5    | 0.5      | 0.5    | 0.5    | 1.0    | 1.0    |
| Security Lead     | 0.5    | 0.5    | 1.0      | 0.5    | 0.5    | 0.5    | 1.0    |
| Backend Engineers | 3.0    | 3.0    | 3.0      | 3.0    | 3.0    | 2.0    | 1.0    |
| QA Engineer       | 0.5    | 1.0    | 1.0      | 1.0    | 1.0    | 2.0    | 1.0    |

---

## Phase 01 Definition of Done

The phase is **DONE** only when ALL of the following are verified:

- [ ] All approved FRs implemented (FR-01 through FR-20)
- [ ] All approved NFRs measured and met (NFR-01 through NFR-10)
- [ ] All ACs tested and passing (AC-01 through AC-20)
- [ ] Security tests pass (gVisor suite, capability enforcement 100%, escape tests)
- [ ] Capability enforcement verified (all 9 types, all constraints)
- [ ] Contract tests pass for all 7 service boundaries
- [ ] All CI quality gates pass (validate, test, security, build)
- [ ] Observability verified (logs, metrics, traces, audit on all endpoints)
- [ ] Documentation complete (> 90%): module READMEs, API docs, runbooks
- [ ] Deployment tested (DEV auto, STAGE manual, PROD canary dry-run)
- [ ] Rollback tested (< 2 min, health verified)
- [ ] Independent verification PASS
- [ ] Human approval recorded (approval.md signed by all roles)

**No "works locally" definition of done.**

---

## Metadata

---

title: Phase 01 Plan
type: phase
phase: 01
status: Planning - Awaiting Plan Approval
author: Planner
date: 2024-10-07
reviewers: Architect, Platform Lead, Security Lead, Product Owner
approved_by: N/A - Awaiting Plan Approval
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---

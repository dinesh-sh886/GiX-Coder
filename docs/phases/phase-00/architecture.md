# Phase 00 Architecture

## Architecture Decisions Made

### ADR-0001: Control Plane / Data Plane Separation
**Decision**: Strict separation between control plane services and execution/data plane.
**Rationale**: Customer code must never execute in control plane; enables independent scaling and security boundaries.

### ADR-0002: Modular Monolith
**Decision**: Start with modular monolith with clear domain boundaries, not microservices.
**Rationale**: Avoids premature distributed complexity; clear boundaries enable future extraction.

### ADR-0003: Durable Workflow Engine
**Decision**: All long-running operations use durable workflow engine (Temporal).
**Rationale**: Event-sourced state machine, automatic checkpointing, retry, human-in-the-loop.

### ADR-0004: Sandbox Isolation
**Decision**: gVisor/Firecracker for strong isolation with capability-based access control.
**Rationale**: Defense in depth against sandbox escape, command execution, network abuse.

## Architecture Baseline Summary

See [Architecture Baseline](docs/architecture/baseline.md) for full details.

### High-Level Structure
```
Developer → CLI/IDE/Web/API → Agent Gateway → Workflow Engine → Agent Harness
    → Router → Execution Sandbox → Files/Shell/Git/Tests/MCP → GitHub/GitLab
```

### Module Boundaries (Modular Monolith)
```
gix-coder/
├── gateway/    # Agent Gateway - API, auth, routing
├── workflow/   # Durable Workflow Engine
├── harness/    # Agent Harness - lifecycle, tools
├── sandbox/    # Execution Sandbox - isolation
├── router/     # Context/Policy/Model Router
├── policy/     # Policy Engine
├── audit/      # Audit Logging
└── shared/     # Shared Kernel
```

### Technology Choices
| Layer | Technology |
|-------|------------|
| Language | Go 1.22+ / TypeScript 5+ |
| API | gRPC + REST (OpenAPI 3.1) |
| Workflow | Temporal.io |
| Sandbox | gVisor / Firecracker |
| Config | Viper / dotenv |
| Secrets | HashiCorp Vault / AWS Secrets Manager |
| Observability | OpenTelemetry → Prometheus/Grafana/Loki/Tempo |
| CI/CD | GitHub Actions |
| Container | Docker / containerd |
| Orchestration | Kubernetes (EKS/GKE) |

## Phase 00 Scope
- **Governance only** - No application implementation
- **Documentation** - All 20 objectives documented
- **Infrastructure** - CI/CD pipeline, GitHub configuration
- **Standards** - All engineering standards established
## Evolution Path

```
Phase 00: Governance ✓
Phase 01: Gateway + Sandbox + Router + Harness (Modular Monolith)
Phase 02: Durable Workflows + Checkpointing
Phase 03: CLI/IDE/Web/API Interfaces
Phase 04: Multi-tenancy + RBAC + Compliance
Future:   Service extraction where justified
```

---

## Metadata

---
title: Phase 00 Architecture
type: phase
phase: 00
status: Verified
author: Architect
date: 2024-10-07
reviewers: Architect, Platform Lead
approved_by: N/A - Awaiting Re-verification
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---
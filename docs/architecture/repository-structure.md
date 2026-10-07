# GiX-Coder Repository Structure

## Current Repository Structure (Phase 00)

```
GiX-Coder/
├── .github/                    # GitHub configuration
│   ├── workflows/              # CI/CD pipelines
│   ├── ISSUE_TEMPLATE/         # Issue templates
│   ├── CODEOWNERS              # Code ownership
│   └── pull_request_template.md
├── docs/                       # Documentation hierarchy
│   ├── product/                # Product-level docs
│   ├── architecture/           # Architecture decisions & baseline
│   ├── adr/                    # Architecture Decision Records
│   ├── security/               # Security documentation
│   ├── qa/                     # Quality assurance
│   └── phases/                 # Phase-specific documentation
│       └── phase-00/           # Phase 00 artifacts
├── .gitignore
├── .env.example                # Environment variable template
├── README.md
├── CONTRIBUTING.md
├── SECURITY.md
├── LICENSE
└── Makefile                    # Governance Makefile (Phase 00)
```

## Planned Repository Structure (Phase 01+)

The following directories and files are **PLANNED** for Phase 01 and beyond. They do not exist in Phase 00.

### Application Structure (Modular Monolith) - PLANNED

```
GiX-Coder/
├── gateway/                    # Agent Gateway Service (PLANNED Phase 01)
│   ├── cmd/
│   │   └── gateway/
│   ├── internal/
│   │   ├── api/               # REST/gRPC handlers
│   │   ├── auth/              # Authentication/Authorization
│   │   ├── routing/           # Request routing
│   │   └── middleware/        # HTTP/gRPC middleware
│   ├── pkg/
│   │   └── dto/               # Data Transfer Objects
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── workflow/                   # Durable Workflow Engine (PLANNED Phase 01/02)
│   ├── cmd/
│   │   └── workflow/
│   ├── internal/
│   │   ├── engine/            # Workflow execution engine
│   │   ├── activities/        # Temporal activities
│   │   ├── workflows/         # Workflow definitions
│   │   └── persistence/       # State persistence
│   ├── pkg/
│   │   └── dto/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── harness/                    # Agent Harness (PLANNED Phase 01)
│   ├── cmd/
│   │   └── harness/
│   ├── internal/
│   │   ├── agent/             # Agent lifecycle
│   │   ├── capabilities/      # Capability management
│   │   ├── tools/             # Tool registry & execution
│   │   └── mcp/               # MCP integration
│   ├── pkg/
│   │   └── dto/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── sandbox/                    # Execution Sandbox (PLANNED Phase 01)
│   ├── cmd/
│   │   └── sandbox/
│   ├── internal/
│   │   ├── executor/          # Command execution
│   │   ├── isolation/         # gVisor/Firecracker integration
│   │   ├── filesystem/        # Virtual filesystem
│   │   └── network/           # Network policy enforcement
│   ├── pkg/
│   │   └── dto/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── router/                     # Context/Policy/Model Router (PLANNED Phase 01)
│   ├── cmd/
│   │   └── router/
│   ├── internal/
│   │   ├── routing/           # Model routing logic
│   │   ├── policy/            # Policy evaluation
│   │   ├── context/           # Context management
│   │   └── classifiers/       # Task classification
│   ├── pkg/
│   │   └── dto/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── policy/                     # Policy Engine (PLANNED Phase 01/02)
│   ├── cmd/
│   │   └── policy/
│   ├── internal/
│   │   ├── evaluation/        # Policy evaluation engine
│   │   ├── storage/           # Policy storage
│   │   └── cache/             # Policy cache
│   ├── pkg/
│   │   └── dto/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── audit/                      # Audit Logging Service (PLANNED Phase 01)
│   ├── cmd/
│   │   └── audit/
│   ├── internal/
│   │   ├── logger/            # Structured audit logger
│   │   ├── storage/           # Immutable storage
│   │   ├── query/             # Audit query API
│   │   └── integrity/         # Tamper evidence
│   ├── pkg/
│   │   └── dto/
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   └── Makefile
├── shared/                     # Shared Kernel (PLANNED Phase 01)
│   ├── errors/                 # Typed error definitions
│   ├── logging/                # Structured logging
│   ├── config/                 # Configuration utilities
│   ├── metrics/                # Metrics utilities
│   ├── tracing/                # Tracing utilities
│   ├── validation/             # Validation utilities
│   ├── dto/                    # Shared DTOs
│   └── constants/              # Shared constants
├── api/                        # API Definitions (Protobuf/OpenAPI) (PLANNED Phase 01)
│   ├── proto/
│   │   ├── gateway/
│   │   ├── workflow/
│   │   ├── harness/
│   │   ├── sandbox/
│   │   ├── router/
│   │   ├── policy/
│   │   └── audit/
│   └── openapi/
├── deploy/                     # Deployment manifests (PLANNED Phase 01)
│   ├── kubernetes/
│   │   ├── base/
│   │   ├── overlays/
│   │   │   ├── dev/
│   │   │   ├── stage/
│   │   │   └── prod/
│   │   └── charts/
│   └── docker-compose/
│       ├── dev.yml
│       ├── stage.yml
│       └── prod.yml
└── test/                       # Cross-cutting tests (PLANNED Phase 01)
    ├── integration/
    ├── contract/
    ├── architecture/
    ├── security/
    └── performance/
```

## Module Independence Rules

1. **No circular dependencies** between modules
2. **Shared kernel only** - modules import from `shared/`, not from each other
3. **DTOs at boundaries** - all external communication uses DTOs
4. **Internal packages private** - `internal/` not importable by other modules
5. **Explicit contracts** - API definitions in `api/` are the contract

## Configuration Structure (PLANNED Phase 01+)

```
configs/
├── dev/
│   ├── gateway.yaml
│   ├── workflow.yaml
│   ├── harness.yaml
│   ├── sandbox.yaml
│   ├── router.yaml
│   ├── policy.yaml
│   └── audit.yaml
├── stage/
│   └── (same structure)
└── prod/
    └── (structure only - secrets via secret manager)
```

## Documentation Structure

### Current (Phase 00)

```
docs/
├── product/
│   └── charter.md
├── architecture/
│   ├── baseline.md
│   ├── ci-cd-strategy.md
│   ├── coding-standards.md
│   ├── configuration-strategy.md
│   ├── documentation-standards.md
│   ├── engineering-principles.md
│   ├── environment-strategy.md
│   ├── git-strategy.md
│   ├── nemotron-rules.md
│   ├── observability-strategy.md
│   ├── opencode-rules.md
│   ├── phase-lifecycle.md
│   ├── pr-strategy.md
│   └── repository-structure.md
├── adr/
│   ├── 0001-control-plane-data-plane-separation.md
│   ├── 0002-modular-monolith.md
│   ├── 0003-durable-workflow-engine.md
│   ├── 0004-sandbox-isolation.md
│   └── framework.md
├── security/
│   └── baseline.md
├── qa/
│   ├── testing-strategy.md
│   └── quality-gates.md
└── phases/
    └── phase-00/
        ├── README.md
        ├── requirements.md
        ├── architecture.md
        ├── plan.md
        ├── acceptance-criteria.md
        ├── implementation.md
        ├── verification.md
        ├── remediation.md
        └── approval.md
```

### Planned Documentation Structure (Phase 01+)

```
docs/
├── security/
│   ├── threat-model.md              # PLANNED
│   ├── incident-response.md         # PLANNED
├── qa/
│   ├── test-standards.md            # PLANNED
├── operations/
│   ├── runbooks/                    # PLANNED
│   ├── deployment.md                # PLANNED
│   ├── monitoring.md                # PLANNED
│   └── disaster-recovery.md         # PLANNED
└── phases/
    └── phase-01/                    # PLANNED
        ├── README.md
        ├── requirements.md
        ├── architecture.md
        ├── plan.md
        ├── acceptance-criteria.md
        ├── implementation.md
        ├── verification.md
        ├── remediation.md
        └── approval.md
```

## Naming Conventions

| Artifact        | Convention                       |
| --------------- | -------------------------------- |
| Directories     | kebab-case                       |
| Go packages     | lowercase, single word preferred |
| Go files        | snake_case.go                    |
| Proto files     | snake_case.proto                 |
| Config files    | service-name.yaml                |
| Dockerfiles     | Dockerfile (no extension)        |
| Makefiles       | Makefile                         |
| Scripts         | verb-noun.sh                     |
| ADRs            | NNNN-short-description.md        |
| Issue templates | type-description.yml             |

---

## References

- [Architecture Baseline](baseline.md)
- [Engineering Principles](engineering-principles.md)
- [Phase 00 Artifacts](../phases/phase-00/)
- ADR-0002: Modular Monolith

---

## Metadata

---

title: GiX-Coder Repository Structure
type: architecture
phase: 00
status: Verified
author: Architect
date: 2024-10-07
reviewers: Platform Lead, Architect
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0002]
related_issues: []
---

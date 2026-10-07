# GiX-Coder Repository Structure

## Root Structure

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
│   ├── operations/             # Operational runbooks
│   └── phases/                 # Phase-specific documentation
│       └── phase-00/           # Phase 00 artifacts
├── configs/                    # Configuration templates
│   ├── dev/                    # DEV environment configs
│   ├── stage/                  # STAGE environment configs
│   └── prod/                   # PROD environment configs (structure only)
├── scripts/                    # Operational scripts
├── tools/                      # Development tools
├── .gitignore
├── .env.example                # Environment variable template
├── README.md
├── CONTRIBUTING.md
├── SECURITY.md
├── LICENSE
└── go.mod / package.json       # Root module (if applicable)
```

## Application Structure (Modular Monolith)

```
GiX-Coder/
├── gateway/                    # Agent Gateway Service
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
├── workflow/                   # Durable Workflow Engine
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
├── harness/                    # Agent Harness
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
├── sandbox/                    # Execution Sandbox
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
├── router/                     # Context/Policy/Model Router
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
├── policy/                     # Policy Engine
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
├── audit/                      # Audit Logging Service
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
├── shared/                     # Shared Kernel
│   ├── errors/                 # Typed error definitions
│   ├── logging/                # Structured logging
│   ├── config/                 # Configuration utilities
│   ├── metrics/                # Metrics utilities
│   ├── tracing/                # Tracing utilities
│   ├── validation/             # Validation utilities
│   ├── dto/                    # Shared DTOs
│   └── constants/              # Shared constants
├── api/                        # API Definitions (Protobuf/OpenAPI)
│   ├── proto/
│   │   ├── gateway/
│   │   ├── workflow/
│   │   ├── harness/
│   │   ├── sandbox/
│   │   ├── router/
│   │   ├── policy/
│   │   └── audit/
│   └── openapi/
├── deploy/                     # Deployment manifests
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
└── test/                       # Cross-cutting tests
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

## Configuration Structure

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

```
docs/
├── product/
│   └── charter.md
├── architecture/
│   ├── baseline.md
│   └── repository-structure.md
├── adr/
│   ├── 0001-control-plane-data-plane-separation.md
│   ├── 0002-modular-monolith.md
│   ├── 0003-durable-workflow-engine.md
│   └── 0004-sandbox-isolation.md
├── security/
│   ├── threat-model.md
│   ├── baseline.md
│   └── incident-response.md
├── qa/
│   ├── testing-strategy.md
│   ├── quality-gates.md
│   └── test-standards.md
├── operations/
│   ├── runbooks/
│   ├── deployment.md
│   ├── monitoring.md
│   └── disaster-recovery.md
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

## Naming Conventions

| Artifact | Convention |
|----------|------------|
| Directories | kebab-case |
| Go packages | lowercase, single word preferred |
| Go files | snake_case.go |
| Proto files | snake_case.proto |
| Config files | service-name.yaml |
| Dockerfiles | Dockerfile (no extension) |
| Makefiles | Makefile |
| Scripts | verb-noun.sh |
| ADRs | NNNN-short-description.md |
| Issue templates | type-description.yml |
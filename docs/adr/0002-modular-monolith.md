# ADR-0002: Modular Monolith

## Status

Accepted

## Context

GiX-Coder needs clear module boundaries for maintainability and future extraction, but premature microservices create operational complexity, distributed system failures, and development friction.

Key forces:

- Team size: small to medium (5-20 engineers initially)
- Domain complexity: high (multiple bounded contexts)
- Deployment simplicity: single deployable unit preferred
- Evolution path: clear extraction boundaries needed
- Testing: need fast, reliable integration tests
- Debugging: single process easier than distributed traces

## Decision

Start with a **modular monolith** with explicit domain boundaries:

### Module Structure

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

### Rules

1. **No circular dependencies** between modules
2. **Shared kernel only** - modules import from `shared/`, not from each other
3. **DTOs at boundaries** - all external communication uses DTOs
4. **Internal packages private** - `internal/` not importable by other modules
5. **Explicit contracts** - API definitions in `api/` are the contract

### Communication

- **Within module**: Direct function calls
- **Between modules**: gRPC (async) or direct calls via interfaces (sync)
- **External**: REST/gRPC via Gateway only

### Extraction Criteria

A module may be extracted as a service when:

- Team ownership boundary aligns
- Independent scaling required
- Different technology stack needed
- Independent deployment velocity needed
- Clear SLO/SLA separation

## Consequences

### Positive

- Simple deployment: single binary/container
- Fast tests: no network in integration tests
- Easy debugging: single process, shared memory
- Transactional consistency: within module
- Refactoring safety: compiler catches cross-module breaks
- Clear evolution path: extraction when justified

### Negative

- Single failure domain: bug in one module can crash all
- Technology lock-in: all modules same language (Go)
- Scaling granularity: all modules scale together
- Deployment coupling: all modules deploy together

### Risks

- **Boundary erosion**: Mitigated by architecture tests (import rules)
- **God module emergence**: Mitigated by code review, architecture review
- **Premature extraction**: Mitigated by explicit extraction criteria

## Alternatives Considered

### Alternative 1: Microservices from Start

- **Pros**: Maximum isolation, independent scaling/deployment, polyglot
- **Cons**: Operational complexity, distributed debugging, network latency, eventual consistency, testing difficulty
- **Why rejected**: Team too small, domain not fully understood, premature optimization

### Alternative 2: Monolith with No Boundaries

- **Pros**: Simplest initially
- **Cons**: Spaghetti code, no extraction path, untestable, coupling
- **Why rejected**: Violates clean architecture, no evolution path

### Alternative 3: Modular Monolith with Shared Database

- **Pros**: Simpler persistence
- **Cons**: Schema coupling, transaction boundary violations, no extraction
- **Why rejected**: Each module owns its data (database per module)

## Implementation Plan

- [ ] Define module structure in repository
- [ ] Create architecture tests for import boundaries
- [ ] Implement shared kernel
- [ ] Build each module with internal/ public separation
- [ ] Define API contracts in `api/`
- [ ] CI enforcement of boundaries

## Related

- ADR-0001: Control Plane / Data Plane Separation
- ADR-0003: Durable Workflow Engine
- ADR-0004: Sandbox Isolation

## Metadata

- **Author**: Architect
- **Date**: 2024-01-15
- **Reviewers**: Platform Lead, Tech Leads
- **Approved By**: Architect (Principal Architect)

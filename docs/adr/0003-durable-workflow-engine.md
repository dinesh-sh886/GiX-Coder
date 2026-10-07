# ADR-0003: Durable Workflow Engine

## Status

Accepted

## Context

GiX-Coder executes long-running, multi-step engineering workflows that may take minutes to hours. These workflows must be:

- Resilient to failures (process crashes, network partitions, deployments)
- Auditable (full execution trace)
- Support human-in-the-loop (approvals, reviews)
- Support checkpointing and recovery
- Support compensation/rollback on failure

Key forces:

- Workflows are business-critical
- Infrastructure failures are inevitable
- Human approval gates required
- Cost optimization (pause/resume)
- Debugging requires full history

## Decision

Use **Temporal.io** as the durable workflow engine for all long-running operations.

### Characteristics

- Event-sourced state machine
- Automatic checkpointing at activity boundaries
- Retry with exponential backoff (configurable)
- Human-in-the-loop via signals/queries
- Compensation/rollback via saga pattern
- Visibility: full history, stack traces, replay

### Workflow Model

```
Workflow Definition (deterministic)
    ↓
Workflow Execution (durable)
    ├── Activity 1 (sandbox exec)
    ├── Activity 2 (model call)
    ├── Human Approval (signal wait)
    ├── Activity 3 (sandbox exec)
    └── Completion / Compensation
```

### Integration

- Control plane services call Temporal via gRPC
- Activities execute in data plane (sandbox)
- Workers run in control plane (harness)
- Sandbox execution via activity heartbeats

## Consequences

### Positive

- Battle-tested: used by Uber, Airbnb, Coinbase, etc.
- Full history: every decision, input, output recorded
- Replay: debug production issues by replaying locally
- Scalability: horizontal worker scaling
- Language support: Go, TypeScript, Java, Python, etc.
- Visibility: Web UI for debugging, monitoring

### Negative

- Operational complexity: Temporal cluster required
- Learning curve: event-sourcing paradigm
- Latency: activity heartbeats, history persistence
- Cost: Temporal cluster (self-hosted or cloud)
- Vendor dependency: though open source (MIT)

### Risks

- **Temporal cluster failure**: Mitigated by multi-AZ deployment, backup
- **History growth**: Mitigated by retention policies, archival
- **Non-determinism bugs**: Mitigated by Temporal's determinism checks, testing
- **Worker scaling**: Mitigated by autoscaling, sticky execution

## Alternatives Considered

### Alternative 1: Custom State Machine + Database

- **Pros**: Full control, no external dependency
- **Cons**: Reinventing durable execution, months of work, subtle bugs
- **Why rejected**: Not core competency, high risk

### Alternative 2: Message Queue + Saga Orchestrator

- **Pros**: Simpler infrastructure (Kafka/RabbitMQ)
- **Cons**: No built-in retry, no visibility, manual checkpointing, complex compensation
- **Why rejected**: Temporal provides all this out of box

### Alternative 3: AWS Step Functions / Azure Logic Apps

- **Pros**: Managed, serverless
- **Cons**: Vendor lock-in, limited expressiveness, cost at scale, no local replay
- **Why rejected**: Must run on any cloud/on-prem

### Alternative 4: Cadence (Temporal predecessor)

- **Pros**: Similar architecture
- **Cons**: Less active development, fewer features, smaller community
- **Why rejected**: Temporal is the active fork

## Implementation Plan

- [ ] Deploy Temporal (dev: embedded, stage/prod: cluster)
- [ ] Define workflow/activity interfaces
- [ ] Implement harness as Temporal worker
- [ ] Sandbox execution as activity
- [ ] Human approval via signals
- [ ] Compensation activities for rollback
- [ ] Integration tests with failure injection

## Related

- ADR-0001: Control Plane / Data Plane Separation
- ADR-0002: Modular Monolith
- ADR-0004: Sandbox Isolation

## Metadata

- **Author**: Platform Lead
- **Date**: 2024-01-16
- **Reviewers**: Architect, Security Lead
- **Approved By**: Architect (Principal Architect)

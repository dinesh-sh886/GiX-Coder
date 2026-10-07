# Phase 01: Agent Gateway & Execution Sandbox Foundation

## Overview

**Phase 01** delivers the core execution platform for GiX-Coder: the Agent Gateway, Execution Sandbox, Context/Policy/Model Router, and Agent Harness. This phase implements the modular monolith foundation with strict control plane / data plane separation as defined in ADR-0001 and ADR-0002.

**Duration**: 8 weeks (target)
**Dependencies**: Phase 00 complete (verified, approved)
**Status**: PLANNING — Awaiting Plan Approval

## Scope Summary

| Module      | Responsibility                                                               | Phase 01 Deliverable         |
| ----------- | ---------------------------------------------------------------------------- | ---------------------------- |
| `gateway/`  | API entry point, auth, routing, rate limiting                                | Complete, tested, documented |
| `sandbox/`  | Isolated, capability-controlled execution                                    | Complete, tested, documented |
| `router/`   | Task classification, policy evaluation, model selection, budget enforcement  | Complete, tested, documented |
| `harness/`  | Agent lifecycle, capability grants, tool registry, MCP client, checkpointing | Complete, tested, documented |
| `shared/`   | Typed errors, logging, config, metrics, tracing, validation, DTOs            | Complete, tested, documented |
| `api/`      | Protobuf + OpenAPI contracts for all services                                | Complete, tested, documented |
| `workflow/` | WorkflowCoordinator + DirectExecutionAdapter (synchronous, no Temporal)      | Complete, tested, documented |
| `policy/`   | Policy evaluation stub (allowlist-based; OPA/Cedar in Phase 02)              | Complete, tested, documented |
| `audit/`    | Immutable audit logging                                                      | Complete, tested, documented |

## Explicitly Deferred to Phase 02+

- Durable Workflow Engine (Temporal-native orchestration, retry, saga compensation)
- Temporal-native checkpointing and history
- CLI/IDE/Web/API interfaces
- Multi-tenancy / RBAC
- Advanced scheduling/orchestration
- Custom model fine-tuning
- Audit log UI

## Phase 01 / Phase 02 Boundary

| Concern              | Phase 01                                                         | Phase 02                                                  |
| -------------------- | ---------------------------------------------------------------- | --------------------------------------------------------- |
| Workflow execution   | WorkflowCoordinator + DirectExecutionAdapter (synchronous)       | Temporal workflow definitions, activities, signals        |
| Checkpointing        | DirectCheckpointStore (serialized agent state to object storage) | TemporalCheckpointStore (Temporal-native history, replay) |
| Retry / compensation | Application-level retry with exponential backoff                 | Temporal-native retry, saga compensation                  |
| Human-in-the-loop    | Synchronous approval via API + webhook                           | Temporal signals/queries, async approval                  |
| Workflow history     | Audit log events only                                            | Full Temporal event history + audit correlation           |
| Execution adapter    | DirectExecutionAdapter (synchronous)                             | TemporalWorkflowAdapter (Temporal-native)                 |
| Checkpoint store     | DirectCheckpointStore (object storage)                           | TemporalCheckpointStore (Temporal history)                |

## Architecture Principles

- **Control plane never executes customer code** (ADR-0001)
- **Modular monolith with explicit boundaries** (ADR-0002)
- **gVisor primary / Firecracker fallback for sandbox isolation** (ADR-0004)
- **Capability-based access control: deny by default, explicit grants** (Security Baseline)
- **OpenTelemetry-native observability from day one** (Observability Strategy)
- **Configuration via layered config, secrets via Vault** (Configuration Strategy)
- **Quality gates are non-negotiable floors** (Quality Gates)

## Key Documents

| Document                                         | Purpose                                                                   | Status   |
| ------------------------------------------------ | ------------------------------------------------------------------------- | -------- |
| [specification.md](specification.md)             | Authoritative Phase 01 specification (from Phase 00)                      | APPROVED |
| [requirements.md](requirements.md)               | Decomposed functional/non-functional requirements with traceability       | PLANNING |
| [architecture.md](architecture.md)               | Detailed module architecture, interfaces, data flows, security boundaries | PLANNING |
| [plan.md](plan.md)                               | Dependency-aware implementation plan with milestones, critical path       | PLANNING |
| [acceptance-criteria.md](acceptance-criteria.md) | Measurable, testable success criteria traceable to requirements           | PLANNING |
| [implementation.md](implementation.md)           | Implementation summary (to be completed during IMPLEMENTING)              | PENDING  |
| [verification.md](verification.md)               | Independent verification results (to be completed during VERIFYING)       | PENDING  |
| [remediation.md](remediation.md)                 | Issues found and fixes applied (if needed)                                | PENDING  |
| [approval.md](approval.md)                       | Human sign-offs for plan approval and phase completion                    | PENDING  |

## Lifecycle State

```
PROMPT → ARCHITECTURE/PLAN → IMPLEMENTATION → IMPLEMENTATION REPORT
    → INDEPENDENT VERIFICATION → REMEDIATION (if FAIL) → RE-VERIFICATION
    → QUALITY GATES → HUMAN APPROVAL → PR → DEVELOP → DEV → STAGE
    → PRODUCTION APPROVAL → MAIN → PROD
```

**Current State**: ARCHITECTURE/PLAN (producing planning artifacts)

## Plan Approval Gate

Before proceeding to IMPLEMENTING, the following must be approved:

- [ ] Requirements reviewed by Product Owner
- [ ] Architecture reviewed by Architect
- [ ] Plan reviewed by Architect + Planner
- [ ] Acceptance criteria reviewed by Product Owner + Verifier
- [ ] Security threat model reviewed by Security Lead
- [ ] All reviewers sign approval.md

## References

- [Phase 01 Specification](specification.md)
- [Phase 00 Architecture Baseline](../../architecture/baseline.md)
- [Repository Structure](../../architecture/repository-structure.md)
- ADR-0001: Control Plane / Data Plane Separation
- ADR-0002: Modular Monolith
- ADR-0003: Durable Workflow Engine
- ADR-0004: Sandbox Isolation

---

## Metadata

---

title: Phase 01 Overview
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

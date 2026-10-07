# Architecture Decision Records (ADR) Index

## Overview

This directory contains Architecture Decision Records (ADRs) for the GiX-Coder project. Each ADR documents a significant architectural decision, its context, and consequences.

## ADR Convention

- **File naming**: `NNNN-short-kebab-case-description.md` (4-digit zero-padded sequence)
- **Status values**: Proposed, Accepted, Superseded, Deprecated
- **Required fields**: Context, Decision, Consequences, Alternatives, Implementation Plan, Metadata
- **Review process**: Minimum 2 approvals (Architect + Domain Owner), Security review for security-related ADRs

## ADR List

| ID | Title | Status | Date | Author | Related |
|----|-------|--------|------|--------|---------|
| [0001](0001-control-plane-data-plane-separation.md) | Control Plane / Data Plane Separation | Accepted | 2024-10-07 | Architect | ADR-0002, ADR-0003, ADR-0004 |
| [0002](0002-modular-monolith.md) | Modular Monolith | Accepted | 2024-10-07 | Architect | ADR-0001, ADR-0003, ADR-0004 |
| [0003](0003-durable-workflow-engine.md) | Durable Workflow Engine | Accepted | 2024-10-07 | Platform Lead | ADR-0001, ADR-0002, ADR-0004 |
| [0004](0004-sandbox-isolation.md) | Sandbox Isolation | Accepted | 2024-10-07 | Security Lead | ADR-0001, ADR-0002, ADR-0003 |

## Adding New ADRs

1. Create new file: `NNNN-short-description.md` (next sequential number)
2. Follow the template in `framework.md`
3. Set initial status to `Proposed`
4. Open PR with `adr` label
5. Obtain minimum 2 approvals (Architect + Domain Owner)
6. Security review required for security-related ADRs
7. Merge to `develop` branch
8. Update this index

## Example Future ADR

| ID | Title | Status | Note |
|----|-------|--------|------|
| 0005 | Example Future ADR Title | Proposed | Example entry - replace when creating ADR-0005 |

## ADR Lifecycle

```
Proposed → Accepted → (Superseded | Deprecated)
```

- **Proposed**: Under review, not yet approved
- **Accepted**: Approved, implementation authorized
- **Superseded**: Replaced by a newer ADR (reference both ways)
- **Deprecated**: No longer valid, no replacement

## Quality Gates

Each ADR must pass:
- [ ] Follows template exactly
- [ ] Sequential numbering (no gaps)
- [ ] All required metadata fields present
- [ ] Reviewed by required roles
- [ ] No unresolved blocking comments
- [ ] Links to related ADRs/issues

---

## Metadata

---
title: GiX-Coder ADR Index
type: adr
phase: 00
status: Verified
author: Architect
date: 2024-10-07
reviewers: Architect, Platform Lead
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---
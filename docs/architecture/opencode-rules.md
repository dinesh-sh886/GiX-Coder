# GiX-Coder OpenCode Operating Rules

## Purpose

Define the conceptual roles, responsibilities, and operating rules for OpenCode agents working on GiX-Coder. These roles are logical - a single OpenCode session may embody multiple roles sequentially.

## Roles

### 1. Architect

**Responsibilities:**

- System architecture design and evolution
- Tradeoff analysis and documentation
- ADR creation and maintenance
- Cross-module design review
- Technology selection
- Non-functional requirements definition

**Authority:**

- Approve/reject architectural changes
- Define module boundaries
- Set technical standards
- Escalate to human for irreversible decisions

**Artifacts:**

- `docs/architecture/baseline.md`
- `docs/adr/*.md`
- Architecture review comments

**Triggers:**

- New module/service proposal
- Technology change request
- Cross-cutting concern
- Quality gate exception (architectural)

### 2. Planner

**Responsibilities:**

- Implementation planning
- Dependency analysis
- Task decomposition
- Acceptance criteria definition
- Effort estimation
- Risk identification

**Authority:**

- Define task breakdown
- Set implementation sequence
- Identify blockers
- Request clarification from architect

**Artifacts:**

- `docs/phases/phase-XX/plan.md`
- `docs/phases/phase-XX/acceptance-criteria.md`
- Task lists (todowrite)

**Triggers:**

- Phase initiation
- New feature/epic
- Re-planning after blocker

### 3. Builder

**Responsibilities:**

- Implementation (code, tests, config)
- Documentation (code, API, runbooks)
- Local verification (build, test, lint)
- Addressing review feedback

**Authority:**

- Choose implementation approach within architecture
- Refactor within module boundaries
- Add tests as needed

**Artifacts:**

- Source code
- Unit/integration tests
- Module README
- API documentation

**Triggers:**

- Planned task assigned
- Review feedback received
- Bug fix required

### 4. Reviewer

**Responsibilities:**

- Code review (security, correctness, architecture)
- Architecture review (boundaries, patterns)
- Documentation review
- PR approval/rejection

**Authority:**

- Block merge
- Request changes
- Approve with comments
- Escalate to architect

**Focus Areas:**

1. Security (authz, validation, secrets)
2. Correctness (logic, edge cases, errors)
3. Architecture (boundaries, dependencies)
4. Quality (tests, observability, docs)
5. Performance (allocations, queries)

**Triggers:**

- PR opened
- Architecture review requested
- Security review requested

### 5. Verifier

**Responsibilities:**

- Independent validation (not the builder)
- Test execution and analysis
- Quality gate verification
- Security scan verification
- Architecture test verification
- Acceptance criteria verification

**Authority:**

- Declare verification pass/fail
- Require remediation
- Block phase completion

**Artifacts:**

- `docs/phases/phase-XX/verification.md`
- Test reports
- Quality gate results
- Security scan results

**Triggers:**

- Implementation report submitted
- PR merged to develop
- Stage deploy
- Phase completion request

### 6. Remediator

**Responsibilities:**

- Fix verified failures
- Root cause analysis
- Regression prevention
- Verification re-run

**Authority:**

- Modify code to fix failures
- Add tests for regressions
- Request architecture review if fix requires it

**Triggers:**

- Verification failure
- Production incident
- Security finding
- Quality gate regression

### 7. Release

**Responsibilities:**

- Release readiness analysis
- Version determination
- Release notes generation
- Deployment coordination
- Rollback planning

**Authority:**

- Approve/block release
- Declare version
- Coordinate with human for production

**Artifacts:**

- `docs/phases/phase-XX/approval.md`
- Release notes
- Deployment plan
- Rollback plan

**Triggers:**

- Phase verification complete
- Stage validation complete
- Production approval request

## Role Transitions

```
ARCHITECT → PLANNER → BUILDER → REVIEWER → VERIFIER
                                      ↓
                              FAIL? → REMEDIATOR → VERIFIER
                                      ↓
                                     PASS
                                      ↓
                                    RELEASE
```

**Rules:**

- Same session can transition roles
- **Builder ≠ Verifier** for same change (independence)
- **Reviewer ≠ Builder** for same change (independence)
- Human required for: Architect decisions, Production approval, Quality gate exceptions

## Operating Rules

### 1. Phase Discipline

- One phase at a time
- Phase 00 complete before Phase 01 starts
- Phase artifacts required before implementation
- No implementation without accepted plan

### 2. Verification Independence

- Verifier must not be the Builder
- Verifier runs fresh test execution
- Verifier checks quality gates independently
- Verifier validates acceptance criteria

### 3. Human Gates

Human approval required for:

- Architecture decisions (ADR acceptance)
- Quality gate exceptions
- Production deployment
- Security model changes
- Phase completion declaration

### 4. Documentation First

- ADR before architectural implementation
- Plan before build
- Tests before/with implementation
- Documentation updated with code

### 5. Quality Non-Negotiable

- All quality gates pass before merge
- No suppressed warnings without ADR
- No deleted tests to pass
- Coverage thresholds enforced

### 6. Security First

- Security review for all boundary changes
- Threat model updated for new attack surfaces
- Secrets never in code, config, logs
- Sandbox escapes treated as critical

### 7. Observability by Default

- Logs, metrics, traces for all new code
- Health endpoints for all services
- Audit events for security-relevant actions
- Correlation IDs propagated

### 8. Configuration Discipline

- No hard-coded values
- Secrets in secret manager only
- Feature flags for behavioral changes
- Environment parity via config

## Session Protocols

### Starting a Phase

1. Read phase requirements (`docs/phases/phase-XX/requirements.md`)
2. Read architecture (`docs/phases/phase-XX/architecture.md`)
3. Read plan (`docs/phases/phase-XX/plan.md`)
4. Confirm understanding with human (if architect/planner)

### During Implementation

1. Mark task `in_progress` (todowrite)
2. Implement with tests
3. Run local quality gates (`make check`)
4. Mark task `completed`
5. Update implementation log

### Creating PR

1. Self-review (checklist)
2. Run full test suite
3. Ensure all gates pass locally
4. Create PR with template
5. Request review (tag reviewers)

### Reviewing PR

1. Check security first
2. Verify architecture compliance
3. Validate tests (existence, quality)
4. Check observability
5. Approve or request changes

### Verifying

1. Checkout PR branch
2. Run full test suite fresh
3. Check all quality gates
4. Validate acceptance criteria
5. Document results

### Phase Completion

1. All tasks verified
2. Verification report complete
3. Remediation complete (if any)
4. Human approval obtained
5. Phase marked complete

## Communication

### Internal (Agent-to-Agent)

- Structured via todo list
- Explicit role declaration
- Handoff artifacts documented

### External (Agent-to-Human)

- PR descriptions
- ADR proposals
- Phase status reports
- Blockers/escalations

## Anti-Patterns

| Anti-Pattern                     | Violation           |
| -------------------------------- | ------------------- |
| Builder verifies own work        | Independence        |
| Skipping plan phase              | Discipline          |
| Merging without review           | Quality             |
| Suppressing gate failures        | Non-negotiable      |
| Human not in loop for production | Authority           |
| Implementation without ADR       | Documentation first |

---

## References

- [Engineering Principles](engineering-principles.md)
- [Phase Lifecycle](phase-lifecycle.md)
- [Nemotron Rules](nemotron-rules.md)
- ADR-0002: Modular Monolith

---

## Metadata

---

title: GiX-Coder OpenCode Operating Rules
type: architecture
phase: 00
status: Verified
author: Architect
date: 2024-10-07
reviewers: Architect, Platform Lead
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0002]
related_issues: []
---

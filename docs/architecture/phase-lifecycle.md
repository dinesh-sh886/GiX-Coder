# GiX-Coder Phase Lifecycle

## Overview

Every phase follows a rigorous lifecycle ensuring governance, quality, and human authority. No phase is complete without independent verification and human approval.

## Lifecycle Flow

```
┌─────────────┐
│   PROMPT    │  Human initiates phase with requirements
└──────┬──────┘
       ↓
┌─────────────┐
│ ARCHITECTURE│  Architect creates/updates architecture
│   / PLAN    │  Planner creates plan + acceptance criteria
└──────┬──────┘
       ↓
┌─────────────┐
│IMPLEMENTATION│ Builder implements per plan
└──────┬──────┘
       ↓
┌─────────────┐
│IMPLEMENTATION│ Builder documents what was done
│   REPORT    │
└──────┬──────┘
       ↓
┌─────────────┐
│ INDEPENDENT │ Verifier validates independently
│ VERIFICATION│
└──────┬──────┘
       ↓
    ┌──┴──┐
    │ FAIL?│
    └──┬──┘
  YES ↓  ↓ NO
┌─────────┐ ┌──────────────┐
│REMEDIATION│ │ QUALITY GATES │
│  (Fix)   │ │  (All pass)  │
└────┬────┘ └──────┬───────┘
     ↓             ↓
┌─────────┐   ┌──────────────┐
│RE-VERIFY│   │ HUMAN APPROVAL │
└────┬────┘   └──────┬───────┘
     ↓             ↓
     └──────┬──────┘
            ↓
       ┌─────────┐
       │   PR    │  Merge to develop
       └────┬────┘
            ↓
       ┌─────────┐
       │  DEVELOP│  Auto-deploy to DEV
       └────┬────┘
            ↓
       ┌─────────┐
       │   DEV   │  Integration testing
       └────┬────┘
            ↓
       ┌─────────┐
       │  STAGE  │  Release candidate validation
       └────┬────┘
            ↓
       ┌─────────┐
       │PRODUCTION│  Explicit human approval
       │ APPROVAL │
       └────┬────┘
            ↓
       ┌─────────┐
       │  MAIN   │  Merge to main
       └────┬────┘
            ↓
       ┌─────────┐
       │  PROD   │  Deploy to production
       └─────────┘
```

## Phase States

| State | Description | Entry Criteria | Exit Criteria |
|-------|-------------|----------------|---------------|
| `PROPOSED` | Phase requested | Human prompt | Requirements documented |
| `PLANNING` | Architecture + Plan | Requirements accepted | Architecture.md, Plan.md, Acceptance-criteria.md approved |
| `IMPLEMENTING` | Code being written | Plan approved | Implementation.md complete, PRs opened |
| `VERIFYING` | Independent validation | Implementation reported | Verification.md complete (PASS) |
| `REMEDIATING` | Fixing failures | Verification FAIL | Remediation.md complete, re-verify PASS |
| `GATING` | Quality gates | Verification PASS | All gates green |
| `APPROVING` | Human sign-off | Gates green | Approval.md signed |
| `MERGING` | PR to develop | Approved | Merged to develop |
| `DEV_DEPLOYED` | Running in DEV | Merged | DEV smoke tests pass |
| `STAGING` | Deployed to STAGE | DEV validated | STAGE acceptance tests pass |
| `PROD_APPROVED` | Human production approval | STAGE validated | Approval recorded |
| `COMPLETE` | In production | PROD deployed | Health verified, monitoring stable |

## Phase Artifacts (Required)

Each phase `XX` must produce:

```
docs/phases/phase-XX/
├── README.md              # Phase summary, status, links
├── requirements.md        # Functional + non-functional requirements
├── architecture.md        # Phase architecture, changes from baseline
├── plan.md                # Implementation plan, tasks, dependencies
├── acceptance-criteria.md # Measurable, testable success criteria
├── implementation.md      # What was built, deviations, decisions
├── verification.md        # Verification results, evidence
├── remediation.md         # Issues found, fixes applied (if any)
└── approval.md            # Human sign-offs with dates
```

### Artifact Standards

**requirements.md**
- Functional requirements (user stories, use cases)
- Non-functional requirements (performance, security, reliability)
- Constraints (budget, timeline, technology)
- Assumptions and dependencies
- Out of scope

**architecture.md**
- Changes from baseline
- New modules/components
- Interface definitions
- Data flow changes
- Security model changes
- Operational considerations
- ADRs created/updated

**plan.md**
- Task breakdown (WBS)
- Dependencies (internal, external)
- Sequence and parallelization
- Effort estimates
- Risk register
- Resource requirements
- Milestones

**acceptance-criteria.md**
- Each criterion: ID, Description, Test Method, Pass/Fail Threshold
- Traceable to requirements
- Automated where possible
- Manual verification steps documented

**implementation.md**
- Summary of changes
- Deviations from plan (with justification)
- Key decisions made during implementation
- Technical debt incurred (with ADR)
- Performance benchmarks
- Security considerations

**verification.md**
- Test execution results (all levels)
- Quality gate results
- Acceptance criteria validation (PASS/FAIL each)
- Security scan results
- Performance validation
- Evidence links (CI runs, test reports, screenshots)

**remediation.md** (if needed)
- Issue description
- Root cause
- Fix applied
- Regression test added
- Re-verification result

**approval.md**
| Role | Name | Date | Signature |
|------|------|------|-----------|
| Architect | | | |
| Security | | | |
| Product | | | |
| Release | | | |

## Phase Gates (Blocking)

### Gate 1: Plan Approval
- Requirements reviewed by Product
- Architecture reviewed by Architect
- Plan reviewed by Architect + Planner
- Acceptance criteria reviewed by Product + Verifier
- **All reviewers approve** → Proceed to IMPLEMENTING

### Gate 2: Verification Pass
- Verifier independent (not Builder)
- All acceptance criteria PASS
- All quality gates PASS
- Security scans clean
- **Verifier signs** → Proceed to GATING

### Gate 3: Quality Gates
- All automated gates PASS (see quality-gates.md)
- No exceptions without ADR
- **CI/CD system validates** → Proceed to APPROVING

### Gate 4: Human Approval
- Architect: Architecture compliance
- Security: Security posture
- Product: Requirements met
- Release: Operational readiness
- **All sign approval.md** → Proceed to MERGING

### Gate 5: Production Approval
- On-call: Operational readiness
- Security: No outstanding critical findings
- Product: Business approval
- **Explicit approval recorded** → Proceed to PROD

## Phase Duration Guidelines

| Phase Type | Target Duration | Max Duration |
|------------|-----------------|--------------|
| Governance (00) | 1-2 weeks | 2 weeks |
| Foundation (01-03) | 2-4 weeks | 6 weeks |
| Feature (04+) | 1-3 weeks | 4 weeks |
| Hotfix | 1-3 days | 1 week |

## Parallel Phases

- **Not allowed** for dependent phases
- **Allowed** for independent workstreams with:
  - Separate phase numbers (01a, 01b)
  - Explicit independence documented
  - Separate verification
  - Merge integration phase

## Phase Completion Criteria

Phase is **COMPLETE** only when:
- [ ] All acceptance criteria PASS (verified independently)
- [ ] All quality gates PASS
- [ ] All approvals recorded
- [ ] Code merged to `main`
- [ ] Deployed to PROD
- [ ] PROD health verified (24h stable)
- [ ] Documentation updated
- [ ] Retrospective completed

## Retrospective

Every phase completes with retrospective:
- What went well?
- What didn't?
- Process improvements?
- Tool improvements?
- Document in `docs/phases/phase-XX/retrospective.md`

## Anti-Patterns

| Anti-Pattern | Consequence |
|--------------|-------------|
| Skipping verification | Undetected defects |
| Builder = Verifier | Bias, missed issues |
| Human approval skipped | Authority violation |
| Gates weakened to pass | Technical debt accumulation |
| Phase declared complete without PROD | False completion |
| No retrospective | Repeated mistakes |
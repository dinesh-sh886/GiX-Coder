# GiX-Coder Engineering Workflow

## 1. Canonical development loop

Use this lifecycle for every phase and every feature.

```text
PHASE DEFINITION
      |
      v
[1] PROMPT
      |
      v
[2] OPENCODE + NEMOTRON
      |
      v
[3] RESPONSE / IMPLEMENTATION REPORT
      |
      v
[4] VERIFICATION PROMPT
      |
      +---- FAIL ----> REMEDIATION PROMPT
      |                    |
      |                    v
      |              RE-VERIFY
      |                    |
      |                    +---- FAIL -> repeat
      |
      v
[5] VERIFIED
      |
      v
[6] HUMAN APPROVAL
      |
      v
[7] FEATURE BRANCH -> PR
      |
      v
[8] CI GATES
      |
      v
develop
      |
      v
Stage deployment
      |
      v
STAGE QA / ACCEPTANCE
      |
      v
Production approval
      |
      v
main
      |
      v
PROD DEPLOYMENT
```

## 2. Important change

Do not make `main` the development branch.

Recommended branch model:

```text
main        = production
stage       = staging/release candidate
develop     = integration/development
feature/*   = implementation branches
fix/*       = bug fixes
hotfix/*    = urgent production fixes
chore/*     = engineering maintenance
docs/*      = documentation
```

For very small teams, `develop` may become the only integration branch,
but keeping `stage` explicit is useful for a SaaS product that has
independent staging validation.

## 3. Phase lifecycle

Every phase has five artifacts:

```text
phase/
  requirements.md
  implementation.md
  verification.md
  remediation.md
  approval.md
```

A phase is **not complete** because OpenCode says it is complete.

A phase is complete only when:

```text
Implementation
AND
verification
AND
quality gates
AND
security checks
AND
human approval
```

## 4. OpenCode roles

Do not use a single generic agent for everything.

Recommended agents:

### Architect

Responsibilities:

- architecture
- ADRs
- module boundaries
- dependencies
- tradeoffs

No code modifications unless explicitly required.

### Planner

Responsibilities:

- convert phase requirements into implementation plan
- identify files/modules
- identify risks
- define acceptance criteria

### Builder

Responsibilities:

- implement the approved plan
- write tests
- update documentation
- preserve architecture

### Reviewer

Responsibilities:

- review code
- detect violations
- no modifications

### Verifier

Responsibilities:

- run all quality gates
- inspect reports
- determine PASS/FAIL
- do not hide failures

### Remediator

Responsibilities:

- fix only verified failures
- rerun affected checks
- provide before/after evidence

### Release

Responsibilities:

- inspect release readiness
- verify migrations
- verify configuration
- confirm deployment artifacts

## 5. Nemotron responsibilities

Use Nemotron as the **engineering decision partner**.

Good uses:

- architecture reasoning
- design alternatives
- code review
- verification analysis
- remediation analysis
- implementation planning
- test planning
- ADR creation

Do not let the model itself decide that the work is verified.

The system of record is:

```text
CI results
+ test reports
+ static analysis
+ coverage
+ security scanning
+ human approval
```

## 6. Prompt contract

Every engineering prompt should contain:

```text
Context
Objective
Scope
Non-goals
Constraints
Architecture rules
Coding standards
Security requirements
Acceptance criteria
Verification requirements
Expected output
```

## 7. Implementation prompt contract

Builder prompt must explicitly say:

- inspect repository first
- do not invent existing modules
- preserve backward compatibility unless breaking change is approved
- use DTOs at API boundaries
- use immutable/domain-safe objects where appropriate
- validate inputs
- use SOLID principles
- avoid duplicated logic
- use configuration instead of literals
- use structured logging
- add unit/integration tests
- update documentation
- do not disable quality gates to make CI green
- do not suppress warnings without justification
- report unresolved issues explicitly

## 8. Verification contract

Verification must include:

```text
Compile
Unit tests
Integration tests
Static analysis
PMD
Checkstyle
SpotBugs
JaCoCo
Lint
Dependency vulnerability scanning
Secrets scanning
Container scan
API contract checks
Architecture checks
```

Relevant checks depend on language/framework.

## 9. Evidence requirement

Every verification run creates:

```text
verification/
  VERIFICATION-YYYYMMDD-HHMM.md
  test-results/
  coverage/
  static-analysis/
  security/
```

The report must include:

```text
STATUS: PASS | FAIL | BLOCKED

Build:
Tests:
Coverage:
PMD:
Checkstyle:
SpotBugs:
Lint:
Security:
Dependency audit:
Architecture:
Git state:

Failures:
Remediation:
Remaining risks:
```

## 10. No silent remediation

When a verification fails:

```text
FAIL
 |
classify
 |
root cause
 |
remediation
 |
retest
```

The agent must not:

- reduce test coverage targets
- delete tests
- weaken PMD rules
- disable Checkstyle
- suppress SpotBugs without reason
- exclude code from JaCoCo simply to pass
- remove security scans
- modify CI to ignore failure

unless a human-approved ADR explicitly changes the quality rule.

## 11. Commit standard

Use Conventional Commits:

```text
feat:
fix:
refactor:
test:
docs:
chore:
build:
ci:
perf:
security:
```

Example:

```text
feat(agent): add durable task state
```

## 12. PR requirements

Every PR must contain:

```text
Summary
Why
Scope
Architecture impact
Database impact
API impact
Security impact
Tests
Quality gates
Operational impact
Rollback plan
Screenshots/logs when useful
```

## 13. PR gates

Required:

```text
build
unit-test
integration-test
lint
format
pmd
checkstyle
spotbugs
jacoco
dependency-scan
secret-scan
container-scan
architecture-test
```

Branch protection should require all required checks.

## 14. Human approval

Approval happens after:

1.  OpenCode implementation report
2.  Independent verification
3.  Remediation if needed
4.  Clean CI
5.  PR review
6.  Human approval

The agent cannot approve its own PR.

## 15. Definition of Done

```text
[ ] Requirements satisfied
[ ] Design documented
[ ] Code implemented
[ ] DTOs validated
[ ] Tests added
[ ] Unit tests pass
[ ] Integration tests pass
[ ] Static analysis clean
[ ] Coverage meets threshold
[ ] Security scan clean
[ ] Docs updated
[ ] PR opened
[ ] CI green
[ ] Human review approved
[ ] Staged
[ ] Stage smoke tests pass
[ ] Production approval
[ ] Released
```

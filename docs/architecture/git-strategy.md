# GiX-Coder Git Branching Strategy

## Branch Structure

### Protected Branches
| Branch | Purpose | Protection |
|--------|---------|------------|
| `main` | Production-ready code | Required PR, CI, human review, signed commits |
| `stage` | Release candidate validation | Required PR, CI, human review |
| `develop` | Integration branch | Required PR, CI |

### Feature Branches
All feature work originates from `develop` unless explicitly required otherwise.

| Prefix | Purpose | Source | Target |
|--------|---------|--------|--------|
| `feature/` | New functionality | `develop` | `develop` |
| `fix/` | Bug fixes | `develop` | `develop` |
| `hotfix/` | Production emergencies | `main` | `main`, `stage`, `develop` |
| `refactor/` | Code restructuring | `develop` | `develop` |
| `chore/` | Maintenance, tooling | `develop` | `develop` |
| `docs/` | Documentation only | `develop` | `develop` |
| `test/` | Test improvements | `develop` | `develop` |
| `security/` | Security fixes | `develop` or `main` | `develop` or `main` |

## Branch Naming
```
<type>/<ticket-id>-<short-description>
```
Examples:
- `feature/GIX-123-agent-gateway`
- `fix/GIX-456-sandbox-timeout`
- `hotfix/GIX-789-prod-memory-leak`
- `refactor/GIX-101-shared-errors`

## Flow

```
feature/* (from develop)
    ↓ PR + CI + Review
develop
    ↓ CI + Integration Tests
stage (release branch cut)
    ↓ STAGE QA + Smoke Tests
Production Approval (human)
    ↓ PR + CI
main
    ↓ Deploy
PROD
```

## Rules

### Protected Branch Rules
- **No direct pushes** to `main`, `stage`, `develop`
- **PR required** for all changes to protected branches
- **CI required** - all checks must pass
- **Human review required** - minimum 1 approval
- **Signed commits required** on `main`
- **Linear history** - squash merge preferred

### Feature Branch Rules
- Branch from `develop` (default)
- Short-lived (< 2 weeks preferred)
- Rebase on `develop` before PR
- Delete after merge (auto-delete enabled)

### Hotfix Rules
- Branch from `main`
- Fix applied to `main` first
- Cherry-pick/merge to `stage` and `develop`
- Post-deploy verification required

### Release Branch (`stage`)
- Cut from `develop` for release candidate
- Only `fix/`, `security/`, `chore/` branches target `stage`
- No new features on `stage`
- QA validation required before production approval

## Merge Strategies

### Squash Merge (Default)
- Single commit on target branch
- Clean history
- PR title becomes commit message
- Preserves PR metadata

### Merge Commit (When Required)
- Preserves individual commits
- Required for:
  - Multiple related commits that tell a story
  - Architectural changes needing history
  - Security fixes needing traceability

### Rebase (Never on Protected)
- Only for local feature branch cleanup
- Never rebase shared branches

## Tagging

### Release Tags
```
v<major>.<minor>.<patch>
```
Examples: `v1.0.0`, `v1.2.3`, `v2.0.0-rc.1`

### Tag Rules
- Tags only on `main`
- Annotated tags with release notes
- Signed tags for production releases
- Automated by release workflow

## Branch Protection Configuration

### `main`
```yaml
required_pull_request_reviews:
  required_approving_review_count: 2
  dismiss_stale_reviews: true
  require_code_owner_reviews: true
required_status_checks:
  strict: true
  contexts:
    - "ci/build"
    - "ci/test"
    - "ci/lint"
    - "ci/security"
    - "ci/architecture"
    - "ci/dependency-scan"
    - "ci/container-scan"
enforce_admins: true
required_signatures: true
allow_force_pushes: false
allow_deletions: false
```

### `stage`
```yaml
required_pull_request_reviews:
  required_approving_review_count: 1
  dismiss_stale_reviews: true
required_status_checks:
  strict: true
  contexts:
    - "ci/build"
    - "ci/test"
    - "ci/lint"
    - "ci/security"
    - "ci/integration"
enforce_admins: false
allow_force_pushes: false
```

### `develop`
```yaml
required_pull_request_reviews:
  required_approving_review_count: 1
required_status_checks:
  strict: true
  contexts:
    - "ci/build"
    - "ci/test"
    - "ci/lint"
enforce_admins: false
allow_force_pushes: false
```

## Automation

### Branch Cleanup
- Auto-delete merged feature branches
- Stale branch detection (> 30 days no activity)
- Auto-close stale PRs (> 14 days no activity)

### Dependency Updates
- Dependabot PRs target `develop`
- Auto-merge for patch updates (CI passes)
- Manual review for minor/major

### Release Automation
- `stage` cut triggers RC build
- Production approval triggers `main` merge
- `main` merge triggers production deploy
- Rollback = revert merge commit on `main`

---

## References

- [Engineering Principles](engineering-principles.md)
- [PR Strategy](pr-strategy.md)
- [CI/CD Strategy](ci-cd-strategy.md)
- ADR-0002: Modular Monolith

---

## Metadata

---
title: GiX-Coder Git Branching Strategy
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
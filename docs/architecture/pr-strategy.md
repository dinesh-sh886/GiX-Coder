# GiX-Coder PR Strategy

## PR Requirements

### Mandatory for All PRs
- [ ] Clear title following commit convention
- [ ] Description with context, changes, testing
- [ ] Linked issue/ticket
- [ ] All CI checks passing
- [ ] Minimum 1 human approval
- [ ] No merge conflicts
- [ ] Branch up to date with target

### PR Template
```markdown
## Summary
Brief description of changes

## Context
Why this change? Link to issue/ADR.

## Changes
- Change 1
- Change 2

## Testing
- Unit tests: [added/updated]
- Integration tests: [added/updated]
- Manual testing: [steps]

## Security
- [ ] No secrets introduced
- [ ] Input validation added
- [ ] Authorization checks verified
- [ ] Dependencies scanned

## Observability
- [ ] Logs added/updated
- [ ] Metrics added/updated
- [ ] Traces verified
- [ ] Health checks updated

## Documentation
- [ ] Code comments updated
- [ ] README updated
- [ ] API docs updated
- [ ] ADR created (if architectural)

## Checklist
- [ ] Compiles without warnings
- [ ] Tests pass locally
- [ ] Coverage >= 80%
- [ ] Lint passes
- [ ] No breaking changes (or documented)
- [ ] Migration plan (if applicable)
```

## Review Process

### Reviewer Responsibilities
1. **Security First** - Check authz, validation, secrets, injection
2. **Correctness** - Logic, edge cases, error handling
3. **Architecture** - Boundaries, dependencies, patterns
4. **Quality** - Tests, observability, documentation
5. **Performance** - Allocations, queries, algorithms

### Review Timeline
- **Initial review**: Within 4 business hours
- **Follow-up**: Within 2 business hours
- **Stale PR**: > 48 hours no activity → ping reviewer

### Review Types
| Type | When | Reviewers |
|------|------|-----------|
| Standard | Most PRs | 1 domain owner + 1 other |
| Security | Security label, authz changes | Security owner + domain owner |
| Architecture | New module, boundary changes | Architect + domain owner |
| Hotfix | Production emergency | On-call + domain owner (expedited) |

### Approval Rules
- **Code Owners** must approve changes to their domains
- **Architect** must approve architectural changes
- **Security** must approve security-sensitive changes
- **Minimum 2 approvals** for `main` (including code owner)
- **Minimum 1 approval** for `develop`/`stage`

## PR Sizing

### Target Sizes
| Size | Lines Changed | Review Time | Preferred |
|------|---------------|-------------|-----------|
| Small | < 200 | < 30 min | ✓ |
| Medium | 200-500 | 30-60 min | OK |
| Large | 500-1000 | 1-2 hours | Split if possible |
| XL | > 1000 | > 2 hours | Must split |

### Splitting Guidelines
- One logical change per PR
- Refactor separate from feature
- Test-only PRs welcome
- Documentation-only PRs welcome

## Draft PRs
- Use for work-in-progress
- No review required until "Ready for Review"
- CI runs but doesn't block
- Convert to ready when complete

## Stacked PRs
- Allowed for large features
- Each PR independently reviewable
- Clear dependency chain in descriptions
- Merge in order

## Conflict Resolution
- Rebase on target branch (not merge)
- Resolve conflicts locally
- Re-run full test suite after rebase
- Request re-review after conflict resolution

## Emergency Process (Hotfix)
1. Branch from `main`: `hotfix/<id>-<desc>`
2. Fix + test + verify
3. PR to `main` with `hotfix` label
4. Expedited review (on-call + 1)
5. Merge to `main` → auto-deploy
6. Cherry-pick to `stage` and `develop`
7. Post-deploy verification
8. Incident review within 48 hours

## Metrics
- PR cycle time (open → merge) < 24 hours (p90)
- Review time (request → first review) < 4 hours
- Rework rate (commits after review) < 20%
- Merge conflict rate < 5%

---

## References

- [Git Branching Strategy](git-strategy.md)
- [CI/CD Strategy](ci-cd-strategy.md)
- [Engineering Principles](engineering-principles.md)
- ADR-0002: Modular Monolith

---

## Metadata

---
title: GiX-Coder PR Strategy
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
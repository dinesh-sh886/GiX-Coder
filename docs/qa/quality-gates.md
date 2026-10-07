# GiX-Coder Quality Gates

## Quality Gate Philosophy

Quality gates are **non-negotiable** minimum standards. They are not targets to optimize toward - they are floors below which code does not progress.

**Never**:

- Weaken rules to obtain passing build
- Suppress warnings without documented justification (ADR)
- Delete tests to resolve failures
- Exclude code from coverage solely to increase reported result

## Gate Categories

### 1. Build Gates (Blocking)

| Gate            | Tool                      | Threshold          | Scope       |
| --------------- | ------------------------- | ------------------ | ----------- |
| Compilation     | `go build` / `tsc`        | Zero errors        | All modules |
| Type Checking   | `go vet` / `tsc --noEmit` | Zero errors        | All modules |
| Code Generation | `go generate` / `protoc`  | Verified committed | All modules |

### 2. Static Analysis Gates (Blocking)

| Gate           | Tool                                 | Configuration                        | Scope    |
| -------------- | ------------------------------------ | ------------------------------------ | -------- |
| Formatting     | `gofmt` / `prettier`                 | Zero diff                            | All code |
| Linting (Go)   | `golangci-lint`                      | Strict config (see coding-standards) | All Go   |
| Linting (TS)   | `eslint`                             | Recommended + security               | All TS   |
| Security Lint  | `gosec`                              | Zero high/critical                   | All Go   |
| Import Hygiene | `goimports` / `eslint-plugin-import` | Zero issues                          | All code |

### 3. Test Gates (Blocking)

| Gate                   | Tool                                   | Threshold       | Scope        |
| ---------------------- | -------------------------------------- | --------------- | ------------ |
| Unit Tests             | `go test` / `vitest`                   | 100% pass       | All modules  |
| Integration Tests      | `go test` / `vitest`                   | 100% pass       | All modules  |
| Contract Tests         | `pact` / `schemathesis`                | 100% pass       | All APIs     |
| Architecture Tests     | Custom                                 | Zero violations | All modules  |
| Coverage (Overall)     | `go test -cover` / `vitest --coverage` | >= 80%          | All modules  |
| Coverage (Per Package) | Same                                   | >= 60%          | Each package |
| Coverage (New Code)    | Diff coverage                          | >= 90%          | PR changes   |
| Race Detection         | `go test -race`                        | Zero races      | All Go tests |

### 4. Security Gates (Blocking)

| Gate                 | Tool                         | Threshold             | Scope          |
| -------------------- | ---------------------------- | --------------------- | -------------- |
| Secret Scan          | `trufflehog` / `git-secrets` | Zero findings         | All commits    |
| Dependency Scan (Go) | `govulncheck`                | Zero critical/high    | `go.mod`       |
| Dependency Scan (TS) | `npm audit` / `osv-scanner`  | Zero critical/high    | `package.json` |
| License Check        | `license-checker`            | Only allowed licenses | All deps       |
| SAST                 | `CodeQL`                     | Zero critical/high    | All code       |
| Container Scan       | `trivy`                      | Zero critical/high    | All images     |
| SBOM Generation      | `syft`                       | Generated + valid     | All builds     |

### 5. Dependency Gates (Blocking)

| Gate                       | Tool                       | Threshold                            |
| -------------------------- | -------------------------- | ------------------------------------ |
| Direct Dependencies        | `go mod verify` / `npm ls` | Verified checksums                   |
| Transitive Vulnerabilities | `osv-scanner`              | Zero critical/high in transitive     |
| Outdated Check             | `dependabot`               | No critical security updates pending |
| License Compliance         | `license-checker`          | Only Apache-2.0, MIT, BSD-3          |

### 6. Architecture Gates (Blocking)

| Gate                   | Check                               | Threshold                          |
| ---------------------- | ----------------------------------- | ---------------------------------- |
| Import Boundaries      | No cross-module `internal/` imports | Zero violations                    |
| Cyclic Dependencies    | Module dependency graph             | Zero cycles                        |
| Layer Violations       | Domain → Infrastructure only        | Zero violations                    |
| API Compatibility      | Breaking change detection           | Zero undocumented breaking changes |
| Protobuf Compatibility | `buf breaking`                      | Zero breaking changes              |

### 7. Documentation Gates (Blocking)

| Gate                 | Check                 | Threshold       |
| -------------------- | --------------------- | --------------- |
| Public API Docs      | All exported symbols  | 100% documented |
| README Exists        | Per module            | Yes             |
| ADR for Architecture | Significant decisions | Required        |
| Changelog            | Per release           | Required        |

---

## References

- [Testing Strategy](testing-strategy.md)
- [CI/CD Strategy](../architecture/ci-cd-strategy.md)
- [Engineering Principles](../architecture/engineering-principles.md)
- [Coding Standards](../architecture/coding-standards.md)
- ADR-0002: Modular Monolith
- ADR-0003: Durable Workflow Engine

## Gate Enforcement Points

### Pre-commit (Local, Fast)

```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: go-fmt
        name: gofmt
        entry: gofmt -l .
        language: system
        types: [go]
      - id: go-imports
        name: goimports
        entry: goimports -l .
        language: system
        types: [go]
      - id: prettier
        name: prettier
        entry: prettier --check .
        language: system
        types: [typescript, json, yaml, markdown]
      - id: secrets
        name: trufflehog
        entry: trufflehog filesystem --no-verification .
        language: system
```

### PR Pipeline (Comprehensive)

All gates above **except** performance/chaos/E2E

### Nightly/Release Pipeline (Exhaustive)

All gates + E2E + Performance + Chaos

## Exception Process

### When Exceptions Are Allowed

1. **False positive** with documented analysis
2. **Third-party code** (generated, vendored) with upstream tracking
3. **Legacy migration** with explicit timeline and owner

### Exception Request (ADR Required)

```markdown
# ADR-NNNN: Quality Gate Exception for <gate>

## Context

Which gate, which code, why it fails

## Analysis

Why this is a false positive / acceptable risk
Root cause analysis

## Mitigation

What prevents the actual risk
Alternative controls in place

## Timeline

When will this be resolved (max 90 days)

## Owner

Team and individual responsible

## Approval

- Architect: [name/date]
- Security: [name/date] (if security gate)
- Tech Lead: [name/date]
```

### Exception Tracking

- Recorded in ADR log
- Reviewed monthly
- Auto-expire after 90 days
- Escalation if not resolved

## Metrics & Reporting

### Dashboard (Grafana)

- Gate pass/fail rate per pipeline
- Time to fix failures
- Exception count and age
- Coverage trends
- Vulnerability trends

### Reports

- **Per PR**: Gate status in PR checks
- **Daily**: Team-level quality summary
- **Weekly**: Organization quality report
- **Monthly**: Exception review, trend analysis

## Continuous Improvement

### Gate Evolution

- Gates only **strengthen** over time
- New gates added via ADR
- Thresholds tightened via ADR
- Never lowered without architect + security approval

### Tool Upgrades

- Tool versions pinned in CI config
- Upgrades tested in branch
- Rollback plan documented
- Breaking changes → ADR

## Quick Reference: PR Checklist

Before marking PR ready:

- [ ] `make fmt` passes
- [ ] `make lint` passes
- [ ] `make test` passes (unit + integration)
- [ ] `make test-coverage` >= 80%
- [ ] `make security-scan` passes
- [ ] `make arch-test` passes
- [ ] No secrets in diff
- [ ] Documentation updated
- [ ] ADR created (if architectural)
- [ ] Changelog entry added

---

## Metadata

---

title: GiX-Coder Quality Gates
type: qa
phase: 00
status: Verified
author: QA Lead
date: 2024-10-07
reviewers: QA Lead, Architect, Platform Lead, Security Lead
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0002, ADR-0003]
related_issues: []
---

# GiX-Coder CI/CD Strategy

## Pipeline Architecture

### Trigger Types
| Trigger | Branches | Purpose |
|---------|----------|---------|
| Push | `feature/*`, `fix/*`, `refactor/*`, `chore/*`, `docs/*`, `test/*`, `security/*` | Validate changes |
| Pull Request | All targeting protected branches | Gate merges |
| Schedule | `main` | Nightly full test suite |
| Release | Tag `v*` | Production deploy |
| Manual | Any | Ad-hoc validation |

## Pipeline Stages

### 1. Validate (Fast - < 5 min)
```yaml
jobs:
  format:
    - gofmt / prettier check
    - import organization
  lint:
    - golangci-lint / eslint
    - security lint (gosec)
  typecheck:
    - go vet / tsc --noEmit
  generate:
    - code generation (proto, mocks)
    - verify generated code committed
```

### 2. Test (Medium - < 15 min)
```yaml
jobs:
  unit:
    - go test / vitest
    - coverage collection
    - race detection (Go)
  contract:
    - API contract tests (Pact/Schemathesis)
    - schema validation
  architecture:
    - import boundary tests
    - cyclic dependency tests
    - layer violation tests
```

### 3. Security (Medium - < 10 min)
```yaml
jobs:
  secrets:
    - trufflehog / git-secrets
    - baseline scan
  dependencies:
    - govulncheck / npm audit / osv-scanner
    - license check (allowed: Apache-2.0, MIT, BSD-3)
  sast:
    - CodeQL / Semgrep
    - custom rules
  container:
    - hadolint (Dockerfile)
    - trivy (image scan)
    - sbom generation (syft)
```

### 4. Integration (Slow - < 30 min)
```yaml
jobs:
  integration:
    - docker-compose up (test stack)
    - cross-module integration tests
    - database migration tests
  e2e:
    - critical path tests
    - sandbox execution tests
```

### 5. Build (Medium - < 10 min)
```yaml
jobs:
  build:
    - multi-arch build (amd64, arm64)
    - reproducible build verification
    - image signing (cosign)
  sbom:
    - generate SPDX SBOM
    - attest build provenance (SLSA)
```

### 6. Deploy (Per Environment)
```yaml
jobs:
  dev:
    - auto-deploy on develop merge
    - smoke tests
  stage:
    - manual approval
    - deploy to stage
    - full test suite
    - performance baseline
  prod:
    - manual approval (production approval gate)
    - blue/green or canary
    - automated rollback on health check failure
    - post-deploy verification
```

## Quality Gates (Blocking)

### All Pipelines
| Gate | Threshold | Action |
|------|-----------|--------|
| Build | Success | Block |
| Unit Tests | 100% pass | Block |
| Coverage | >= 80% overall, >= 60% per package | Block |
| Lint | Zero warnings | Block |
| Typecheck | Zero errors | Block |
| Security Scan | Zero critical/high | Block |
| Secrets | Zero findings | Block |
| Architecture Tests | Zero violations | Block |

### PR Pipeline Additional
| Gate | Threshold |
|------|-----------|
| Integration Tests | 100% pass |
| Contract Tests | 100% pass |
| Dependency Scan | Zero critical |

### Release Pipeline Additional
| Gate | Threshold |
|------|-----------|
| E2E Tests | 100% pass |
| Performance | Within 10% of baseline |
| Container Scan | Zero critical/high |
| SBOM | Generated and attested |
| Image Signature | Verified |

## Environments

### DEV
- **Trigger**: Merge to `develop`
- **Deploy**: Automatic
- **Data**: Synthetic, reset daily
- **Retention**: 7 days
- **Access**: All developers

### STAGE
- **Trigger**: PR to `stage` (release branch)
- **Deploy**: Manual approval
- **Data**: Production-like, anonymized
- **Retention**: 30 days
- **Access**: Engineers, QA, Product

### PROD
- **Trigger**: Merge to `main` (tagged release)
- **Deploy**: Manual approval (production approval gate)
- **Data**: Customer data
- **Retention**: Per compliance
- **Access**: On-call, Release managers

## Deployment Strategies

### DEV - Rolling
- Zero-downtime rolling update
- Max surge: 25%, max unavailable: 0
- Health check: `/health/ready`

### STAGE - Blue/Green
- Two identical environments
- Switch traffic after validation
- Instant rollback capability

### PROD - Canary
- 5% → 25% → 50% → 100%
- Automated metrics analysis
- Manual approval at each step
- Rollback on:
  - Error rate > 1%
  - Latency p99 > 2x baseline
  - Custom business metric degradation

## Rollback

### Automatic
- Health check failure → immediate rollback
- Metric threshold breach → rollback
- Config: max 3 rollbacks/hour

### Manual
- `gh workflow run rollback -f environment=prod -f version=v1.2.3`
- Single command, < 2 min

### Database
- Backward-compatible migrations only
- Rollback migration provided
- No data loss migrations without explicit approval

## Artifacts

### Retention
| Artifact | DEV | STAGE | PROD |
|----------|-----|-------|------|
| Build logs | 30 days | 90 days | 365 days |
| Test reports | 30 days | 90 days | 365 days |
| Coverage | 30 days | 90 days | 365 days |
| Security reports | 90 days | 365 days | 365 days |
| Container images | 7 days | 30 days | 365 days |
| SBOMs | 90 days | 365 days | 365 days |

### Promotion
- Images promoted by tag (not rebuilt)
- `dev` → `stage` → `prod` tag promotion
- Digest verification at each step

## Secrets Management

### CI Secrets
- GitHub Environments with protection rules
- OIDC for cloud provider auth
- No long-lived credentials in CI

### Runtime Secrets
- HashiCorp Vault / AWS Secrets Manager
- Injected at runtime (not build time)
- Automatic rotation
- Audit logging on access

## Notifications

### Slack Channels
| Event | Channel |
|-------|---------|
| PR checks failed | #gix-coder-ci |
| Deploy started | #gix-coder-deploys |
| Deploy succeeded | #gix-coder-deploys |
| Deploy failed | #gix-coder-alerts |
| Security finding | #gix-coder-security |
| Production incident | #gix-coder-incidents |

### PR Comments
- Test results summary
- Coverage report
- Security scan summary
- Build artifact links

## Metrics & SLOs

### Pipeline SLOs
| Metric | Target |
|--------|--------|
| PR pipeline duration | < 20 min (p90) |
| Merge to DEV deploy | < 10 min |
| STAGE deploy | < 15 min |
| PROD deploy | < 30 min |
| Rollback time | < 2 min |
| Pipeline availability | 99.9% |

### Quality Metrics
| Metric | Target |
|--------|--------|
| Build success rate | > 95% |
| Test flakiness | < 1% |
| False positive security | < 5% |
| Mean time to feedback | < 10 min |

---

## References

- [Git Branching Strategy](git-strategy.md)
- [PR Strategy](pr-strategy.md)
- [Environment Strategy](environment-strategy.md)
- [Quality Gates](../qa/quality-gates.md)
- [Testing Strategy](../qa/testing-strategy.md)
- ADR-0002: Modular Monolith
- ADR-0003: Durable Workflow Engine

---

## Metadata

---
title: GiX-Coder CI/CD Strategy
type: architecture
phase: 00
status: Verified
author: Platform Lead
date: 2024-10-07
reviewers: Architect, Platform Lead, Security Lead
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0002, ADR-0003]
related_issues: []
---
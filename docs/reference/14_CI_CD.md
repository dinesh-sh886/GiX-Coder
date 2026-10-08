# GiX-Coder CI/CD Architecture

## 1. Environment model

Use three long-lived environments:

---

Environment Purpose Source

---

dev developer feature/develop
integration/local/cloud  
development

stage production-like stage
acceptance

prod customer production main
-------------------------------------------------------------------------

The GitHub environment names should be exactly:

```text
dev
stage
prod
```

GitHub Environments support deployment protection, environment-specific
secrets, branch restrictions and required reviewers.
citeturn493998search0turn493998search1

## 2. Recommended promotion

```text
feature/*
    |
    v
PR -> develop
    |
    v
CI
    |
    v
DEV deploy
    |
    v
Integration QA
    |
    v
PR develop -> stage
    |
    v
STAGE deploy
    |
    v
Acceptance / smoke / migration validation
    |
    v
PR stage -> main
    |
    v
PROD approval
    |
    v
PROD deploy
```

## 3. CI on every PR

Run in parallel where possible:

```text
lint
format-check
compile
unit-test
integration-test
pmd
checkstyle
spotbugs
jacoco
dependency-audit
secret-scan
sast
architecture-test
```

Then aggregate into:

```text
quality-gate
```

A PR is mergeable only if `quality-gate` passes.

## 4. Build once, promote the artifact

Avoid rebuilding separately for stage and prod.

Correct:

```text
source commit
   |
CI build
   |
immutable artifact/image
   |
stage
   |
same artifact
   |
prod
```

Benefits:

- reproducibility
- lower risk
- simpler rollback
- reliable provenance

## 5. Container strategy

Image:

```text
gixcoder:<git-sha>
```

Never make `latest` the deployment identity.

Also publish:

```text
gixcoder:<semver>
```

Production should deploy a specific immutable digest.

## 6. GitHub Actions workflow layout

Recommended:

```text
.github/
  workflows/
    pr-quality.yml
    build-image.yml
    deploy-dev.yml
    deploy-stage.yml
    deploy-prod.yml
    security.yml
    dependency-update.yml
```

## 7. PR quality workflow

Conceptual:

```yaml
name: PR Quality

on:
  pull_request:
    branches:
      - develop
      - stage
      - main

jobs:
  quality:
    runs-on: ubuntu-latest

    steps:
      - checkout

      - setup runtime

      - format-check
      - lint
      - compile
      - unit-test
      - integration-test
      - pmd
      - checkstyle
      - spotbugs
      - jacoco
      - dependency-scan
      - secret-scan

      - upload reports

      - quality-gate
```

## 8. Deployment jobs

### Dev

Automatic after merge to `develop`.

### Stage

Automatic after merge to `stage`, subject to stage environment
protection.

### Production

Triggered by merge to `main` and protected by `prod` environment
approval.

GitHub Actions supports environment approvals, branch restrictions,
concurrency and deployment protection rules.
citeturn493998search2turn493998search5

## 9. Concurrency

Do not allow multiple production deployments simultaneously.

Conceptual:

```yaml
concurrency:
  group: prod
  cancel-in-progress: false
```

This uses GitHub Actions concurrency to serialize deployments.
citeturn493998search2

## 10. Migration strategy

Never couple destructive DB migration with application startup.

Use:

```text
expand
  |
deploy compatible app
  |
migrate data
  |
contract
```

For PostgreSQL use migration tooling appropriate to the application
stack.

## 11. Rollout

Production rollout:

```text
5%
 |
health
 |
25%
 |
health
 |
50%
 |
health
 |
100%
```

For Kubernetes:

- RollingUpdate initially
- canary later
- blue/green for high-risk systems

## 12. Observability gate

A deployment is not complete merely because the container is running.

Validate:

- HTTP readiness
- error rate
- latency
- dependency health
- database connectivity
- queue health
- key business transaction
- metrics freshness

## 13. Rollback

Every deployment should record:

```text
commit SHA
image digest
migration version
configuration version
```

Rollback:

```text
previous image digest
+
compatible DB schema
```

Do not blindly roll back an application after an irreversible schema
migration.

## 14. Security gates

Minimum:

- Trivy/container scan
- secret scanning
- dependency CVE scan
- SAST
- IaC scanning
- SBOM generation
- image signing/provenance later

## 15. Dependency strategy

Automate dependency PRs but do not auto-merge:

- major upgrades
- framework upgrades
- security-critical architecture libraries
- agent SDK upgrades
- database drivers

Use automated tests + review.

## 16. Release strategy

Use SemVer:

```text
MAJOR.MINOR.PATCH
```

Tag:

```text
v0.1.0
v0.1.1
v0.2.0
```

During early development, frequent minor releases are acceptable.

## 17. Availability strategy

For a 99.9% target:

- 2+ application replicas
- managed HA database
- redundant queues
- multiple worker replicas
- health/readiness probes
- retries with jitter
- circuit breakers
- idempotent jobs
- externalized state
- automated backups
- tested restores

Later:

```text
Region A <--> Region B
```

Do not introduce multi-region writes before the application has a tested
consistency model.

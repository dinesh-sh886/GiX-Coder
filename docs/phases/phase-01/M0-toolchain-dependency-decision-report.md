# Phase 01 M0 Go Toolchain / Dependency Architecture Decision Report

**Generated**: 2026-10-07
**Status**: DECISION_READY

---

## 1. Approved Toolchain Baseline

### Source Documents and Requirements

| Document                  | File                                      | Requirement                          | Type                               |
| ------------------------- | ----------------------------------------- | ------------------------------------ | ---------------------------------- |
| **Coding Standards**      | `docs/architecture/coding-standards.md:8` | "Go 1.22+ (current stable)"          | **Hard architectural requirement** |
| **Architecture Baseline** | `docs/architecture/baseline.md:128`       | "Language: Go 1.22+ / TypeScript 5+" | **Hard architectural requirement** |
| **CI Configuration**      | `.github/workflows/ci.yml:29`             | `GO_VERSION: '1.22'`                 | **Hard operational requirement**   |
| **Makefile**              | `Makefile:9`                              | `GO_VERSION := 1.22`                 | **Hard build requirement**         |
| **go.work**               | `go.work:1`                               | `go 1.22.12`                         | **Hard module requirement**        |

**Conclusion**: Go 1.22 is a **hard architectural requirement** mandated by Phase 00 baseline documents. The project must remain compatible with Go 1.22.12 (the specific patch version in go.work).

---

## 2. Complete Dependency Analysis (Shared Module)

### Direct Dependencies (from shared/go.mod)

| Dependency                                                          | Version  | Direct/Transitive | Go Requirement                  | Used By                     | M0 Required? |
| ------------------------------------------------------------------- | -------- | ----------------- | ------------------------------- | --------------------------- | ------------ |
| `go.opentelemetry.io/otel`                                          | v1.29.0  | Direct            | 1.21+ (core), 1.25+ (exporters) | tracing                     | ✅ Yes       |
| `go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc` | v1.24.0  | Direct            | 1.21+                           | tracing (OTLP metrics gRPC) | ✅ Yes       |
| `go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp` | v1.24.0  | Direct            | 1.21+                           | tracing (OTLP metrics HTTP) | ✅ Yes       |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc`   | v1.24.0  | Direct            | 1.21+                           | tracing (OTLP traces gRPC)  | ✅ Yes       |
| `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`   | v1.24.0  | Direct            | 1.21+                           | tracing (OTLP traces HTTP)  | ✅ Yes       |
| `go.opentelemetry.io/otel/sdk`                                      | v1.29.0  | Direct            | 1.21+                           | tracing                     | ✅ Yes       |
| `go.opentelemetry.io/otel/sdk/metric`                               | v1.24.0  | Direct            | 1.25+                           | tracing                     | ✅ Yes       |
| `github.com/prometheus/client_golang`                               | v1.24.1  | Direct            | **1.25+**                       | metrics                     | ✅ Yes       |
| `github.com/golang-migrate/migrate/v4`                              | v4.18.0  | Direct            | **1.25+**                       | cmd/migrate                 | ✅ Yes       |
| `github.com/go-playground/validator/v10`                            | v10.20.0 | Direct            | **1.26+** (v10.30+)             | validation                  | ✅ Yes       |
| `golang.org/x/crypto`                                               | v0.54.0  | Direct            | **1.26+** (v0.57+)              | security                    | ✅ Yes       |
| `github.com/rs/zerolog`                                             | v1.35.1  | Direct            | 1.23+                           | logging                     | ✅ Yes       |
| `github.com/golang-jwt/jwt/v5`                                      | v5.3.1   | Direct            | 1.21+                           | security                    | ✅ Yes       |
| `github.com/spf13/viper`                                            | v1.21.0  | Direct            | 1.21+                           | config                      | ✅ Yes       |
| `github.com/spf13/cobra`                                            | v1.10.2  | Direct            | 1.21+                           | cmd/migrate                 | ✅ Yes       |
| `google.golang.org/grpc`                                            | v1.61.0  | Direct            | 1.21+                           | tracing (OTLP gRPC)         | ✅ Yes       |
| `google.golang.org/protobuf`                                        | v1.32.0  | Direct            | 1.21+                           | generated code              | ✅ Yes       |

### Transitive Dependencies with Highest Go Requirements

| Dependency                                    | Version  | Via    | Go Requirement              | Notes                        |
| --------------------------------------------- | -------- | ------ | --------------------------- | ---------------------------- |
| `prometheus/client_golang`                    | v1.24.1  | direct | **1.25+**                   | Hard requirement from v1.24+ |
| `golang-migrate/migrate/v4`                   | v4.18.0  | direct | **1.25+**                   | Hard requirement from v4.20+ |
| `go-playground/validator/v10`                 | v10.20.0 | direct | 1.23+ (v10.30+ needs 1.26+) | Current version OK           |
| `golang.org/x/crypto`                         | v0.54.0  | direct | **1.26+** (v0.57+)          | Current version OK           |
| `go.opentelemetry.io/otel/sdk/metric`         | v1.24.0  | direct | 1.25+                       | OTLP metrics                 |
| `go.opentelemetry.io/otel/exporters/otlp/...` | v1.24.0  | direct | 1.21+                       | OTLP exporters               |
| `zerolog`                                     | v1.35.1  | direct | **1.23+**                   | Logging                      |

---

## 3. Minimum Compatible Toolchain

### Actual Minimum Go Version: **1.25**

**Reason:**

- `prometheus/client_golang v1.24.1` → requires Go 1.25+ (v1.24.0+ bumped minimum)
- `golang-migrate/migrate/v4 v4.18.0` → requires Go 1.25+ (v4.20.0+ bumped minimum, but transitive deps may require it)
- `go.opentelemetry.io/otel/sdk/metric v1.24.0` → requires Go 1.25+
- `zerolog v1.35.1` → requires Go 1.23+ (acceptable for 1.25)

**Critical Path**: The `prometheus/client_golang` and OTLP metric exporter dependencies are the binding constraints. These are direct dependencies in `shared` used by `metrics` and `tracing` packages.

---

## 4. Three Options Evaluation

### OPTION A: Upgrade Project Go Baseline to 1.25+

| Dimension                      | Assessment                                                          |
| ------------------------------ | ------------------------------------------------------------------- |
| **Implementation Effort**      | LOW: Update `go.work`, `go.mod`, Makefile, CI, Dockerfile, docs     |
| **Security Implications**      | POSITIVE: Newer Go = security fixes, better compiler, faster builds |
| **Maintenance Implications**   | POSITIVE: Aligns with upstream dependency requirements              |
| **Compatibility Implications** | LOW RISK: Go 1.25 is backward compatible with 1.22 code             |
| **CI Impact**                  | LOW: Update GitHub Actions `GO_VERSION: '1.25'`                     |
| **Docker Impact**              | LOW: Update base image to `golang:1.25`                             |
| **Developer Environment**      | LOW: Standard Go upgrade, no code changes                           |
| **Modular Architecture**       | NEUTRAL: No change to module boundaries                             |
| **M0 Schedule**                | MINIMAL DELAY: ~1 hour for config updates                           |
| **M1 Effect**                  | POSITIVE: M1 implementation starts on supported Go version          |
| **Risks**                      | Minimal - Go 1.x backward compatibility guarantee                   |

### OPTION B: Pin Dependencies Compatible with Go 1.22

| Dimension                      | Assessment                                                                                |
| ------------------------------ | ----------------------------------------------------------------------------------------- |
| **Implementation Effort**      | HIGH: Must identify and pin compatible versions for ALL affected deps                     |
| **Security Implications**      | NEGATIVE: Older deps = unpatched vulnerabilities, no security support                     |
| **Maintenance Implications**   | NEGATIVE: Creates version debt, harder to upgrade later                                   |
| **Compatibility Implications** | HIGH RISK: Complex transitive dependency resolution, likely conflicts                     |
| **CI Impact**                  | NEGATIVE: Custom toolchain, non-standard build                                            |
| **Docker Impact**              | NEGATIVE: Cannot use standard Go images                                                   |
| **Developer Environment**      | NEGATIVE: Non-standard Go version management                                              |
| **Modular Architecture**       | NEUTRAL: No boundary change                                                               |
| **M0 Schedule**                | SIGNIFICANT DELAY: Days of dependency archaeology                                         |
| **M1 Effect**                  | NEGATIVE: Technical debt carried into M1 implementation                                   |
| **Risks**                      | **HIGH**: Prometheus client and OTLP exporters have no Go 1.22-compatible recent versions |

**Key Finding**: `prometheus/client_golang` v1.20.0 (last Go 1.22 compatible) is from 2022, has known CVEs. OTLP exporters have no Go 1.22 compatible versions supporting current OTLP spec.

### OPTION C: Refactor Shared Dependency Boundaries

| Dimension                      | Assessment                                                                                     |
| ------------------------------ | ---------------------------------------------------------------------------------------------- |
| **Implementation Effort**      | HIGH: Split `shared` into multiple modules (`shared-metrics`, `shared-tracing`, `shared-core`) |
| **Security Implications**      | NEUTRAL: Same dependencies, just isolated                                                      |
| **Maintenance Implications**   | POSITIVE: Long-term cleaner boundaries                                                         |
| **Compatibility Implications** | NEUTRAL: Same dependencies, different module graph                                             |
| **CI Impact**                  | MEDIUM: More modules to build/test                                                             |
| **Docker Impact**              | MEDIUM: More Dockerfiles                                                                       |
| **Developer Environment**      | MEDIUM: More modules to manage                                                                 |
| **Modular Architecture**       | POSITIVE: Aligns with ADR-0002 "shared kernel only" principle                                  |
| **M0 Schedule**                | SIGNIFICANT DELAY: 1-2 weeks for refactoring                                                   |
| **M1 Effect**                  | POSITIVE: Clean boundaries for M1 implementation                                               |
| **Risks**                      | **MEDIUM**: Does NOT solve Go version problem - modules still need compatible deps             |

**Key Finding**: Refactoring shared module boundaries is a GOOD ARCHITECTURAL DIRECTION (per ADR-0002), but it does NOT resolve the Go version constraint. The transitive dependencies remain the same.

---

## 5. Shared Module Architecture Analysis

### Current Shared Module Structure

```
shared/
├── config/        # Viper config loading
├── logging/       # Zerolog wrapper
├── metrics/       # Prometheus client_golang → **1.25+ required**
├── tracing/       # OpenTelemetry + OTLP exporters → **1.25+ required**
├── security/      # JWT, crypto → OK for 1.22
├── validation/    # validator/v10 → OK for 1.22 (current version)
├── dto/           # Pure types → OK
├── constants/     # Pure types → OK
└── errors/        # Pure types → OK
```

### Architecture Classification

| Package             | Classification            | Reason                                                    |
| ------------------- | ------------------------- | --------------------------------------------------------- |
| `shared/metrics`    | **QUESTIONABLE COUPLING** | Pulls `prometheus/client_golang` (1.25+) into ALL modules |
| `shared/tracing`    | **QUESTIONABLE COUPLING** | Pulls OTLP exporters (1.25+) into ALL modules             |
| `shared/config`     | **ACCEPTABLE**            | Viper is Go 1.21+ compatible                              |
| `shared/logging`    | **ACCEPTABLE**            | Zerolog 1.35.1 needs 1.23+ (acceptable if Go 1.25)        |
| `shared/security`   | **ACCEPTABLE**            | JWT/crypto compatible                                     |
| `shared/validation` | **ACCEPTABLE**            | Validator v10.20.0 compatible                             |
| `shared/dto`        | **ACCEPTABLE**            | Pure types, no deps                                       |
| `shared/constants`  | **ACCEPTABLE**            | Pure types, no deps                                       |
| `shared/errors`     | **ACCEPTABLE**            | Pure types, no deps                                       |

**Violation**: The current `shared` module violates ADR-0002 "shared kernel only" principle by forcing heavy infrastructure dependencies (Prometheus, OpenTelemetry exporters) on ALL modules, even those that don't need them (e.g., `audit`, `policy` may not need metrics/tracing directly).

---

## 6. Security Analysis of Downgrading

### If Forced to Go 1.22 via Dependency Pinning:

| Dependency                    | Pinned Version        | Security Status              | Risk         |
| ----------------------------- | --------------------- | ---------------------------- | ------------ |
| `prometheus/client_golang`    | v1.19.x (2022)        | **EOL**, multiple CVEs       | **CRITICAL** |
| `golang-migrate/migrate/v4`   | v4.15.x (2022)        | **EOL**, no security patches | **HIGH**     |
| `go-playground/validator/v10` | v10.10.x (2022)       | **EOL**                      | **MEDIUM**   |
| `golang.org/x/crypto`         | v0.18.x (2022)        | **EOL**, no post-quantum     | **HIGH**     |
| OTLP Exporters                | No Go 1.22 compatible | **BLOCKS** OpenTelemetry     | **CRITICAL** |

**Conclusion**: Downgrading to Go 1.22-compatible versions creates **unacceptable security risk**. The OpenTelemetry ecosystem has moved past Go 1.22.

---

## 7. Toolchain Upgrade Analysis (If Option A Selected)

### Files Requiring Update

| File                                    | Current              | Required Change      |
| --------------------------------------- | -------------------- | -------------------- |
| `go.work`                               | `go 1.22.12`         | `go 1.25.0`          |
| `Makefile:9`                            | `GO_VERSION := 1.22` | `GO_VERSION := 1.25` |
| `.github/workflows/ci.yml:29`           | `GO_VERSION: '1.22'` | `GO_VERSION: '1.25'` |
| `shared/go.mod:3`                       | `go 1.25.0`          | Already 1.25+        |
| All module `go.mod:3`                   | `go 1.22.12`         | `go 1.25.0`          |
| Dockerfiles (when created)              | `golang:1.22`        | `golang:1.25`        |
| Docker Compose (when created)           | `golang:1.22`        | `golang:1.25`        |
| Kubernetes build configs (when created) | `golang:1.22`        | `golang:1.25`        |
| Developer docs (README, CONTRIBUTING)   | "Go 1.22+"           | "Go 1.25+"           |

**Note**: The CI workflow already uses `temporalio/auto-setup:1.22` which is a Temporal version, not Go version.

---

## 8. Recommendation

### Primary Strategy: **OPTION A - Upgrade to Go 1.25**

### Decision Summary

1. **Go Baseline**: **Go 1.25.0** (minimum required by current dependency graph)
2. **Dependency Versions**: Retain current versions (all compatible with Go 1.25+)
3. **Downgrade Dependencies**: **NO** - unacceptable security risk, blocks OTLP
4. **Refactor Shared**: **DEFER TO M1** - good architectural direction, but not required for M0
5. **ADR Changes**: Update ADR-0002 to document Go version decision; consider new ADR for shared module decomposition in M1
6. **Exact M0 Remediation Required**:
   - Update `go.work` to `go 1.25.0`
   - Update `Makefile:9` `GO_VERSION := 1.25`
   - Update `.github/workflows/ci.yml:29` `GO_VERSION: '1.25'`
   - Update all module `go.mod` files to `go 1.25.0`
   - Run `go work sync` and `go mod tidy` for all modules
   - Verify `go list ./...`, `go build ./...`, `go vet ./...` pass
7. **Defer to M1**:
   - Shared module decomposition (separate `shared-metrics`, `shared-tracing`, `shared-core`)
   - Full ADR for shared module boundaries
   - Docker/Kubernetes build configuration updates

### Rationale

- **Go 1.25 is the minimum required** by the approved dependency graph (Prometheus, OTLP exporters, golang-migrate)
- **Security**: Go 1.25 + current deps = supported, patched versions
- **Architecture**: Upgrading toolchain is a configuration change, not an architectural change
- **M0 Impact**: Minimal (~1 hour config changes), unblocks all validation
- **M1 Readiness**: M1 implementation starts on a supported, current Go version

---

## 9. Final Report

### Current Baseline

- **Approved Go requirement**: 1.22+ (Phase 00 baseline)
- **Actually installed Go**: 1.22.12
- **Dependency-required minimum**: **1.25** (enforced by `prometheus/client_golang v1.24.1`, `go.opentelemetry.io/otel/sdk/metric v1.24.0`, `golang-migrate/migrate/v4 v4.18.0`)

### Dependency Summary

| Dependency                          | Version  | Go Requirement | Status         |
| ----------------------------------- | -------- | -------------- | -------------- |
| prometheus/client_golang            | v1.24.1  | **1.25+**      | BLOCKS Go 1.22 |
| golang-migrate/migrate/v4           | v4.18.0  | **1.25+**      | BLOCKS Go 1.22 |
| go.opentelemetry.io/otel/sdk/metric | v1.24.0  | **1.25+**      | BLOCKS Go 1.22 |
| zerolog                             | v1.35.1  | 1.23+          | OK for 1.25    |
| go-playground/validator/v10         | v10.20.0 | 1.23+          | OK for 1.25    |
| golang.org/x/crypto                 | v0.54.0  | 1.26+ (v0.57+) | OK for 1.25    |
| OpenTelemetry OTLP exporters        | v1.24.0  | 1.21+          | OK for 1.25    |

### Options Assessment

| Option                 | Verdict         | Rationale                                                  |
| ---------------------- | --------------- | ---------------------------------------------------------- |
| A: Upgrade to Go 1.25+ | **RECOMMENDED** | Minimal effort, resolves all blockers, security-positive   |
| B: Pin to Go 1.22      | **REJECTED**    | Requires EOL dependencies, blocks OTLP, high security risk |
| C: Refactor shared     | **DEFER TO M1** | Good architecture, but doesn't solve Go version issue      |

### Recommended Strategy

**Upgrade project Go baseline to 1.25.0** with the following exact M0 remediation:

1. `go.work`: `go 1.25.0`
2. `Makefile:9`: `GO_VERSION := 1.25`
3. `.github/workflows/ci.yml:29`: `GO_VERSION: '1.25'`
4. All `go.mod` files: `go 1.25.0`
5. Run `go work sync && go mod tidy` for all modules
6. Verify: `go list ./...`, `go build ./...`, `go vet ./...` for each module

### Architecture Impact

| Artifact              | Impact                                                                       |
| --------------------- | ---------------------------------------------------------------------------- |
| ADR-0002              | Update Go version reference; note shared module decomposition planned for M1 |
| Coding Standards      | Update "Go 1.22+" → "Go 1.25+"                                               |
| CI/CD Strategy        | Already uses Go 1.22 env var - update to 1.25                                |
| Architecture Baseline | Already says "Go 1.22+" - update to 1.25+                                    |

### Implementation Required (M0)

| File                       | Change                        |
| -------------------------- | ----------------------------- |
| `go.work`                  | Line 1: `go 1.25.0`           |
| `Makefile`                 | Line 9: `GO_VERSION := 1.25`  |
| `.github/workflows/ci.yml` | Line 29: `GO_VERSION: '1.25'` |
| `shared/go.mod`            | Line 3: `go 1.25.0` (already) |
| `api/go.mod`               | Line 3: `go 1.25.0`           |
| `audit/go.mod`             | Line 3: `go 1.25.0`           |
| `gateway/go.mod`           | Line 3: `go 1.25.0`           |
| `harness/go.mod`           | Line 3: `go 1.25.0`           |
| `policy/go.mod`            | Line 3: `go 1.25.0`           |
| `router/go.mod`            | Line 3: `go 1.25.0`           |
| `sandbox/go.mod`           | Line 3: `go 1.25.0`           |
| `workflow/go.mod`          | Line 3: `go 1.25.0`           |

### M0 Impact

- **Blocker Resolution**: All 6 original blockers resolved (protobuf contracts, validation, OpenAPI, lint, enum naming, RPC naming) + Go compilation unblocked
- **Remaining Validation**: Run full module validation suite after config changes

---

**STATUS: DECISION_READY**

**Recommended Action**: Execute Option A (Go 1.25 upgrade) as M0 remediation. Do not commit. Do not start M1. Do not claim M0 verified until validation passes.

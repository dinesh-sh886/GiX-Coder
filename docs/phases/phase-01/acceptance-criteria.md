# Phase 01 Acceptance Criteria

## Overview

This document defines measurable, testable success criteria for Phase 01. Each criterion is traceable to requirements and specifies the test method, pass/fail threshold, and evidence required.

**Source**: [Phase 01 Requirements](requirements.md), [Phase 01 Specification](specification.md)

---

## Acceptance Criteria

### Gateway Criteria

| ID    | Criterion                                                                               | Requirement(s)                    | Test Method                                                                              | Pass Threshold                                                                      | Evidence                                                   |
| ----- | --------------------------------------------------------------------------------------- | --------------------------------- | ---------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------- | ---------------------------------------------------------- |
| AC-01 | Gateway accepts valid workflow request via REST API                                     | FR-01, FR-02, FR-03, FR-04, FR-05 | Integration test: POST `/workflows/execute` with valid JWT, valid DTO, within rate limit | 200 OK, returns `execution_id`, status `started`                                    | Test report, response JSON, audit log `workflow.submitted` |
| AC-02 | Gateway rejects invalid JWT (expired, malformed, wrong issuer, wrong audience, missing) | FR-02                             | Unit tests (5 cases) + Integration test with real JWKS                                   | 401 Unauthorized, error code `AUTH_INVALID`, no downstream call                     | Test report, audit log `auth.failure`                      |
| AC-03 | Gateway enforces rate limits per tenant and per workflow                                | FR-05                             | Load test: burst + sustained requests exceeding limits                                   | 429 Too Many Requests with `Retry-After` header, audit log `rate_limit.exceeded`    | Load test report (k6), audit log                           |
| AC-04 | Gateway validates request DTO schema (required fields, types, constraints)              | FR-01                             | Unit tests (valid, missing fields, wrong types, constraint violations) + Integration     | 422 Unprocessable Entity, error code `VALIDATION_ERROR`, details per field          | Test report                                                |
| AC-05 | Gateway evaluates authorization policy and returns capability grants                    | FR-03                             | Integration test: valid request → policy engine → grants returned                        | Grants match tenant allowlist ∩ workflow requirements, audit log `capability.grant` | Test report, audit log                                     |

### Sandbox Criteria

| ID    | Criterion                                                                                        | Requirement(s)  | Test Method                                                                                                                       | Pass Threshold                                                                                                      | Evidence                                                      |
| ----- | ------------------------------------------------------------------------------------------------ | --------------- | --------------------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------- |
| AC-06 | Sandbox executes file read/write/list operations within granted paths                            | FR-06, FR-07    | Integration test: grant `filesystem.read/write/list` on `/workspace/project`, execute operations                                  | Read returns content, write creates file, list returns entries, all within path scope                               | Test report, sandbox audit log                                |
| AC-07 | Sandbox blocks unauthorized filesystem paths (traversal, absolute, symlinks, outside workspace)  | FR-07           | Security test: 5 attack vectors (`../`, absolute path, symlink escape, hardlink, special files)                                   | All attempts return 403/500, error code `CAPABILITY_DENIED`, audit log `capability.deny`, no host filesystem access | Security test report, audit log, host filesystem verification |
| AC-08 | Sandbox executes allowed shell commands with argument validation                                 | FR-06, FR-08    | Integration test: grant `shell.execute` for `ls`, `cat`, `grep`; execute with valid args                                          | Command executes, stdout/stderr returned, exit code 0, audit log `tool.invoked`                                     | Test report, sandbox audit log                                |
| AC-09 | Sandbox blocks disallowed shell commands and injection attempts                                  | FR-08           | Security test: unauthorized command (`rm`, `curl`, `bash -c`), argument injection (`; cat /etc/passwd`, `$(cmd)`)                 | All attempts return error, error code `CAPABILITY_DENIED` or validation error, audit log `capability.deny`          | Security test report, audit log                               |
| AC-10 | Sandbox enforces CPU and memory limits via cgroups v2                                            | FR-12           | Stress test: allocate > memory limit, consume > CPU limit, exceed pids limit                                                      | Process receives SIGKILL/SIGTERM, audit log `resource_limit.exceeded`, no host impact                               | Stress test report, audit log, host resource verification     |
| AC-11 | Sandbox prevents unauthorized network egress (SSRF, private IPs, metadata service)               | FR-11 (network) | Security test: 4 vectors (private IP, metadata service 169.254.169.254, localhost, allowed domain)                                | Private/metadata/localhost blocked, allowed domain succeeds, audit log `network.denied`                             | Security test report, audit log                               |
| AC-12 | Sandbox executes Git operations within repo/branch scope                                         | FR-09           | Integration test: grant `git.read`/`git.write` for `github.com/org/repo@main`; clone, fetch, checkout, diff, status, commit, push | Authorized operations succeed, unauthorized (different repo/branch) denied, credentials not leaked                  | Test report, audit log, credential scan                       |
| AC-13 | Sandbox executes test commands within framework allowlist and resource limits                    | FR-10           | Integration test: grant `test.execute` for `go test`, `jest`, `pytest`; execute with timeout                                      | Tests execute, results returned, timeout enforced, audit log `tool.invoked`                                         | Test report, audit log                                        |
| AC-14 | Sandbox enforces capability grants (all 9 types: 3 filesystem, shell, 2 git, test, network, MCP) | FR-07-FR-11     | Coverage test: grant each capability type, verify grant/deny behavior                                                             | 100% capability types tested, each enforces scope/constraints/expiry                                                | Test matrix report, audit log                                 |

### Router Criteria

| ID    | Criterion                                                                            | Requirement(s) | Test Method                                                                                    | Pass Threshold                                                                                 | Evidence                                    |
| ----- | ------------------------------------------------------------------------------------ | -------------- | ---------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------- | ------------------------------------------- |
| AC-15 | Router classifies task correctly (complexity, domain, capabilities, tokens, latency) | FR-13          | Unit test: 20 fixtures covering all complexity/domain combinations                             | Classification matches expected for ≥ 95% of fixtures                                          | Test report, classification accuracy matrix |
| AC-16 | Router selects model per policy (allowlist, cost, latency, capability match)         | FR-14          | Integration test: 10 policy scenarios (tenant allowlist, budget, capability requirements)      | Selected model matches policy decision in all scenarios                                        | Test report, audit log `model.selected`     |
| AC-17 | Router enforces token budgets at tenant, workflow, and request levels                | FR-15          | Integration test: budget exhaustion at each level, concurrent requests                         | Request denied at correct level, `budget.deny` audit, remaining budget accurate                | Test report, audit log, Redis budget state  |
| AC-18 | Router validates model outputs (JSON schema, safety, PII detection)                  | FR-16          | Unit test: valid schema, invalid schema, PII patterns (email, SSN, API key), safety violations | Valid passes, invalid rejected with details, PII detected/redacted, safety violations rejected | Test report, validation details             |

### Harness Criteria

| ID    | Criterion                                                                          | Requirement(s)     | Test Method                                                                                  | Pass Threshold                                                                                                    | Evidence                                       |
| ----- | ---------------------------------------------------------------------------------- | ------------------ | -------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| AC-19 | Harness spawns agent with capability grants and transitions through full lifecycle | FR-17, FR-07-FR-11 | Integration test: spawn → initialize → ready → running → checkpointing → paused → terminated | All valid transitions execute, invalid transitions rejected, audit log for each state change                      | Test report, audit log                         |
| AC-20 | Harness executes MCP tool via approved server with capability grant                | FR-19, FR-11       | Integration test: register approved MCP server, grant `mcp.invoke`, invoke tool              | Tool executes, result returned, audit log `mcp.invoked`, unapproved server denied                                 | Test report, audit log                         |
| AC-21 | Checkpoint/restore works (save agent state, restore, continue execution)           | FR-20              | Integration test: execute steps → checkpoint → restore → continue → verify state             | Restored state matches checkpointed state, execution continues correctly, audit log `checkpoint.created/restored` | Test report, checkpoint integrity verification |
| AC-22 | Tool registry registers built-in tools and MCP tools with capability mapping       | FR-18, FR-19       | Unit test: register all 4 built-in + 2 MCP tools; verify schema, capabilities, timeout       | All tools registered, capability mapping correct, execution verifies grants                                       | Test report                                    |

### End-to-End Criteria

| ID    | Criterion                                                                           | Requirement(s) | Test Method                                                                            | Pass Threshold                                                                  | Evidence                          |
| ----- | ----------------------------------------------------------------------------------- | -------------- | -------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------- | --------------------------------- |
| AC-23 | Full workflow executes end-to-end (Gateway → Workflow → Harness → Sandbox → Router) | All FRs        | E2E test: submit workflow → execute steps (file, shell, model) → checkpoint → complete | Workflow completes successfully, result returned, all audit events present      | E2E test report, full audit trail |
| AC-24 | All quality gates pass in CI/CD pipeline                                            | NFR-07         | Full CI run on release branch                                                          | All stages green (validate, test, security, build, deploy DEV)                  | CI run link, gate status          |
| AC-25 | No critical/high vulnerabilities in security scans                                  | NFR-07         | SAST (CodeQL), SCA (govulncheck/osv-scanner), Container (Trivy), Secrets (TruffleHog)  | Zero critical/high findings                                                     | Security scan reports             |
| AC-26 | Documentation completeness > 90%                                                    | NFR-10         | Automated check: module READMEs, API docs, runbooks, architecture docs                 | All 8 modules have README, OpenAPI served, runbooks complete, architecture docs | Documentation checklist           |

### Non-Functional Criteria

| ID    | Criterion                                       | Requirement(s) | Test Method                                                           | Pass Threshold                                         | Evidence                                   |
| ----- | ----------------------------------------------- | -------------- | --------------------------------------------------------------------- | ------------------------------------------------------ | ------------------------------------------ |
| AC-27 | Gateway p99 latency < 100ms under load          | NFR-01         | Load test: 1000 RPS, 5 min duration, STAGE environment                | p99 < 100ms                                            | k6 load test report                        |
| AC-28 | Sandbox cold start < 10s (p99)                  | NFR-02         | 100 cold starts (no warm pool), measure container create + first exec | p99 < 10s                                              | Benchmark report                           |
| AC-29 | Sandbox warm start < 2s (p99)                   | NFR-03         | 100 warm starts (pre-warmed pool), measure execute latency            | p99 < 2s                                               | Benchmark report                           |
| AC-30 | 1000 concurrent workflow executions             | NFR-04         | Load test: ramp to 1000 concurrent, sustain 5 min                     | All executions complete, error rate < 0.1%             | k6 load test report                        |
| AC-31 | Availability 99.9% over test period             | NFR-05         | 30-day STAGE monitoring (or accelerated)                              | Uptime ≥ 99.9%                                         | Monitoring dashboard                       |
| AC-32 | Audit log durability 11 nines (99.9999999%)     | NFR-06         | Immutable storage verification, DR test, hash chain validation        | Zero data loss, hash chain valid                       | Storage verification report                |
| AC-33 | Zero critical/high vulnerabilities (continuous) | NFR-07         | CI security scans on every PR + nightly                               | Zero critical/high                                     | Security scan reports                      |
| AC-34 | Sandbox escape resistance verified              | NFR-08         | gVisor test suite + custom escape attempts                            | All gVisor tests pass, zero successful escapes         | gVisor test results, escape attempt report |
| AC-35 | Configuration immutability at runtime           | NFR-09         | Startup validation, attempt runtime mutation                          | Config validated at startup, runtime mutation rejected | Validation test report                     |
| AC-36 | Observability coverage 100% endpoints           | NFR-10         | Trace span on every gRPC/HTTP endpoint, RED metrics per service       | 100% endpoints have trace + metrics                    | Observability audit report                 |

### Security-Specific Criteria

| ID    | Criterion                                                               | Requirement(s)     | Test Method                                                                    | Pass Threshold                                                               | Evidence                 |
| ----- | ----------------------------------------------------------------------- | ------------------ | ------------------------------------------------------------------------------ | ---------------------------------------------------------------------------- | ------------------------ |
| AC-37 | gVisor test suite passes                                                | NFR-08             | Run `runsc test` full suite                                                    | All tests pass                                                               | gVisor test output       |
| AC-38 | Capability enforcement 100% coverage (all 9 types, grant/deny paths)    | FR-07-FR-11        | Test matrix: each capability type × grant present/absent/expired/invalid scope | 100% combinations tested, all enforce correctly                              | Test matrix report       |
| AC-39 | Resource limit enforcement verified under stress                        | FR-12              | Stress: OOM, CPU throttle, pid limit, I/O limit, wall-time                     | Process terminated cleanly, audit logged, no host impact                     | Stress test report       |
| AC-40 | Network isolation (no ingress, egress allowlist only, DNS policy)       | FR-11              | Network tests: ingress attempt, egress to denied/allowed, DNS resolution       | Ingress blocked, denied egress blocked, allowed egress works, DNS per policy | Network test report      |
| AC-41 | Secret handling: no secrets in code, config, logs, checkpoints, audit   | Security Baseline  | Automated scans: TruffleHog (code), log scan, checkpoint scan, audit scan      | Zero findings                                                                | Scan reports             |
| AC-42 | mTLS on all control plane gRPC communication                            | ADR-0001           | `openssl s_client` verification on all service ports                           | Valid certificates, mutual auth, cipher suite TLS 1.3                        | mTLS verification report |
| AC-43 | Capability escalation prevention (no grant forgery, no scope expansion) | FR-03, FR-07-FR-11 | Security tests: forged grant, scope expansion, expiry bypass                   | All attempts denied, audit logged                                            | Security test report     |

---

## Verification Method

Each criterion is verified by:

1. **File existence check** - Required artifacts exist
2. **Content validation** - Against template/standard (automated where possible)
3. **Cross-reference validation** - Links, IDs, traceability
4. **Human review confirmation** - For criteria requiring judgment

### Automated Verification (CI)

```yaml
# .github/workflows/verification.yml
jobs:
  acceptance-criteria:
    runs-on: ubuntu-latest
    steps:
      - name: Run unit tests
        run: make test-unit-go test-unit-ts
      - name: Run integration tests
        run: make test-integration
      - name: Run contract tests
        run: make test-contract
      - name: Run architecture tests
        run: make test-architecture
      - name: Run security tests
        run: make test-security
      - name: Run performance tests
        run: make test-performance
      - name: Run chaos tests
        run: make test-chaos
      - name: Verify documentation
        run: make validate-docs
      - name: Verify audit completeness
        run: make verify-audit
```

### Manual Verification (Verifier)

| Criterion      | Verification Steps                                             |
| -------------- | -------------------------------------------------------------- |
| AC-27 to AC-30 | Review k6 load test reports, compare thresholds                |
| AC-31          | Review 30-day monitoring dashboard (or accelerated test plan)  |
| AC-32          | Review immutable storage verification, hash chain validation   |
| AC-33          | Review security scan reports from CI                           |
| AC-34          | Review gVisor test suite output, escape attempt logs           |
| AC-35          | Review config validation at startup, runtime mutation attempts |
| AC-36          | Review observability audit: trace + metrics per endpoint       |
| AC-37          | Review gVisor test suite execution output                      |
| AC-38          | Review capability test matrix coverage report                  |
| AC-39          | Review stress test reports for each resource type              |
| AC-40          | Review network test report                                     |
| AC-41          | Review TruffleHog, log, checkpoint, audit scan reports         |
| AC-42          | Review mTLS verification output                                |
| AC-43          | Review capability escalation test report                       |

---

## Sign-Off

| Role            | Name | Date | Signature |
| --------------- | ---- | ---- | --------- |
| Architect       |      |      |           |
| Security Lead   |      |      |           |
| Product Owner   |      |      |           |
| Release Manager |      |      |           |

---

## Metadata

---

title: Phase 01 Acceptance Criteria
type: phase
phase: 01
status: Planning - Awaiting Plan Approval
author: Planner
date: 2024-10-07
reviewers: Product Owner, Architect, Security Lead, Verifier
approved_by: N/A - Awaiting Plan Approval
related_adrs: [ADR-0001, ADR-0002, ADR-0003, ADR-0004]
related_issues: []
---

# GiX-Coder Testing Strategy

## Testing Philosophy

- **Test Pyramid**: Many unit, some integration, few E2E
- **Fast Feedback**: Unit tests < 1s, Integration < 30s, E2E < 5min
- **Deterministic**: No flaky tests, hermetic environments
- **Isolated**: Tests don't depend on each other
- **Meaningful**: Test behavior, not implementation
- **Automated**: All tests run in CI, no manual steps

## Test Levels

### 1. Unit Tests (Target: 80%+ coverage)
- **Scope**: Single function/class/module
- **Dependencies**: Mocked (interfaces)
- **Speed**: < 10ms per test
- **Location**: `*_test.go` / `*.test.ts` alongside code
- **Framework**: Go testing / Vitest

#### What to Test
- Business logic (pure functions)
- Error handling paths
- Edge cases (boundaries, nil, empty, max)
- State transitions
- Validation logic
- Transformation logic

#### What NOT to Test
- Framework code (stdlib, libraries)
- Trivial getters/setters
- Implementation details (private methods)
- Database queries (use integration tests)

#### Patterns
```go
// Table-driven tests
func TestValidateWorkflow(t *testing.T) {
    tests := []struct{
        name    string
        input   *Workflow
        wantErr error
    }{
        {"valid", validWorkflow(), nil},
        {"empty name", workflow{Name: ""}, ErrEmptyName},
        {"too long", workflow{Name: strings.Repeat("a", 256)}, ErrNameTooLong},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateWorkflow(tt.input)
            if tt.wantErr != nil {
                assert.ErrorIs(t, err, tt.wantErr)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 2. Integration Tests (Target: Critical paths)
- **Scope**: Multiple modules, real dependencies
- **Dependencies**: Real (testcontainers, localstack)
- **Speed**: < 30s per test
- **Location**: `test/integration/`
- **Framework**: Go testing + testcontainers / Vitest + testcontainers

#### What to Test
- Module interactions
- Database operations (migrations, queries)
- API contracts (request/response)
- External service integration (mocked via wiremock)
- Message queue interactions
- Cache behavior
- File system operations

#### Patterns
```go
func TestWorkflowExecutionIntegration(t *testing.T) {
    // Use testcontainers for PostgreSQL, Temporal, Redis
    ctx := context.Background()
    pg := postgresContainer(t, ctx)
    temporal := temporalContainer(t, ctx)
    
    // Run migrations
    migrate(t, pg.ConnectionString)
    
    // Create test workflow
    wf := createTestWorkflow(t, pg)
    
    // Execute via API
    resp := executeWorkflow(t, temporal, wf)
    
    // Verify
    assert.Equal(t, "completed", resp.Status)
}
```

### 3. Contract Tests (Target: All external APIs)
- **Scope**: API contracts (provider/consumer)
- **Tool**: Pact / Schemathesis
- **Location**: `test/contract/`
- **CI**: Provider verification on every PR

#### Provider Tests
```yaml
# pact/provider/gateway.yaml
provider:
  name: "gateway"
  transport: "http"
  path: "/health/ready"
consumer:
  name: "cli"
interactions:
  - description: "health check"
    request:
      method: "GET"
      path: "/health/ready"
    response:
      status: 200
      body:
        status: "healthy"
        checks: []
```

### 4. Architecture Tests (Target: All modules)
- **Scope**: Architectural rules
- **Tool**: Custom (Go: `go list`, `golangci-lint` rules; TS: `eslint-plugin-import`)
- **Location**: `test/architecture/`
- **CI**: Every PR

#### Rules
- No imports from `internal/` of other modules
- No cyclic dependencies
- Layer violations (domain → infrastructure only)
- No direct database access outside repository layer
- No HTTP calls outside gateway/adapter layer

### 5. End-to-End Tests (Target: Critical user journeys)
- **Scope**: Full system, real environment
- **Environment**: STAGE (nightly), PROD (canary)
- **Speed**: < 5 min per test
- **Location**: `test/e2e/`
- **Tool**: Playwright / custom Go test binary

#### Critical Journeys
1. Create project → Execute workflow → View results
2. Agent executes file operations → Verify sandbox isolation
3. Agent executes shell commands → Verify allowlist
4. Git operations → Verify repo state
5. Multi-step workflow with checkpoint → Verify recovery

### 6. Security Tests
- **SAST**: CodeQL on every PR
- **SCA**: govulncheck/osv-scanner on every PR
- **Secret Scan**: TruffleHog on every commit
- **Container Scan**: Trivy on every build
- **DAST**: OWASP ZAP on STAGE deploy
- **Penetration**: Annual third-party

### 7. Performance Tests
- **Load**: k6 scripts, run on STAGE weekly
- **Stress**: Breaking point identification
- **Soak**: 24h stability (monthly)
- **Baseline**: Compare against previous release

### 8. Chaos Tests
- **Tool**: Chaos Mesh / Litmus
- **Schedule**: Monthly on STAGE
- **Scenarios**: Pod kill, network partition, latency, CPU pressure
- **Validation**: System recovers, no data loss

## Test Organization

```
test/
├── unit/                    # Co-located with code (not here)
├── integration/
│   ├── gateway/
│   ├── workflow/
│   ├── harness/
│   ├── sandbox/
│   ├── router/
│   ├── policy/
│   └── audit/
├── contract/
│   ├── provider/
│   └── consumer/
├── architecture/
│   ├── import_rules_test.go
│   └── layer_rules_test.ts
├── e2e/
│   ├── workflow_execution.spec.ts
│   ├── sandbox_isolation.spec.ts
│   └── git_operations.spec.ts
├── performance/
│   ├── load_test.js
│   └── baseline.json
├── chaos/
│   ├── pod_kill.yaml
│   └── network_partition.yaml
└── security/
    ├── sast_rules/
    └── dast_config.yaml
```

## Test Data Management

### Principles
- **No production data** in tests
- **Synthetic data generators** for all entities
- **Deterministic seeds** for reproducibility
- **Isolated per test** (no shared state)

### Patterns
```go
// Test data builder pattern
func WorkflowBuilder() *WorkflowBuilder {
    return &WorkflowBuilder{
        workflow: &Workflow{
            ID:        uuid.NewString(),
            Name:      "test-workflow",
            TenantID:  "test-tenant",
            Status:    "pending",
            CreatedAt: time.Now(),
        },
    }
}

func (b *WorkflowBuilder) WithName(name string) *WorkflowBuilder {
    b.workflow.Name = name
    return b
}

func (b *WorkflowBuilder) Build() *Workflow {
    return b.workflow
}

// Usage
wf := WorkflowBuilder().WithName("custom").Build()
```

## CI Integration

### Pipeline Stages
| Stage | Tests | Timeout | Required |
|-------|-------|---------|----------|
| Validate | Unit (fast subset) | 3 min | Yes |
| Test | Unit (full), Integration, Contract, Architecture | 15 min | Yes |
| Security | SAST, SCA, Secrets, Container | 10 min | Yes |
| Build | - | 10 min | Yes |
| Deploy DEV | Smoke (subset E2E) | 5 min | Yes |
| Deploy STAGE | E2E (full), Performance | 30 min | Manual |
| Deploy PROD | Canary validation | 30 min | Manual |

### Coverage Gates
| Level | Threshold |
|-------|-----------|
| Overall | >= 80% |
| Per package | >= 60% |
| New code | >= 90% |
| Critical paths | 100% |

### Flaky Test Policy
- **Detection**: Auto-quarantine after 2 failures in 10 runs
- **Resolution**: Fix within 5 days or delete
- **No skipping** without documented justification

## Test Environments

### Local
- `make test` - Unit only
- `make test-integration` - Integration (testcontainers)
- `make test-all` - All local tests

### CI
- Ephemeral containers per job
- Parallel execution
- Artifact upload (coverage, reports)

### STAGE
- Persistent test namespace
- Nightly full suite
- Pre-deploy smoke tests

## Test Quality Metrics

| Metric | Target |
|--------|--------|
| Unit test execution time | < 30s total |
| Integration test execution time | < 5 min total |
| Flaky test rate | < 1% |
| False positive rate | < 2% |
| Coverage (overall) | >= 80% |
| Coverage (new code) | >= 90% |
| Mutation testing score | >= 70% (critical paths) |

## Anti-Patterns (Forbidden)

| Anti-Pattern | Why |
|--------------|-----|
| Testing implementation | Brittle, blocks refactoring |
| Shared test state | Flaky, non-deterministic |
| Sleep/wait in tests | Slow, unreliable |
| External network calls | Flaky, slow, not hermetic |
| Production-like data | Privacy, compliance |
| Commented-out tests | Dead code, false confidence |
| `t.Skip()` without reason | Hidden gaps |
| Assertions without messages | Hard to debug |

---

## References

- [Quality Gates](quality-gates.md)
- [CI/CD Strategy](../architecture/ci-cd-strategy.md)
- [Engineering Principles](../architecture/engineering-principles.md)
- ADR-0002: Modular Monolith
- ADR-0003: Durable Workflow Engine

---

## Metadata

---
title: GiX-Coder Testing Strategy
type: qa
phase: 00
status: Verified
author: QA Lead
date: 2024-10-07
reviewers: QA Lead, Architect, Platform Lead
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0002, ADR-0003]
related_issues: []
---
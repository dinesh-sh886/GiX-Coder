# GiX-Coder Coding Standards

## Language Standards

### Go (Primary - Control Plane Services)

#### Version
- Go 1.22+ (current stable)

#### Formatting
- `gofmt` - mandatory
- `goimports` - mandatory for import organization
- Line length: 120 characters max

#### Linting
- `golangci-lint` with strict config:
  - `errcheck` - error handling
  - `staticcheck` - static analysis
  - `govet` - vet checks
  - `ineffassign` - ineffective assignments
  - `unused` - unused code
  - `misspell` - spelling
  - `prealloc` - slice/map preallocation
  - `bodyclose` - HTTP body closing
  - `gosec` - security
  - `cyclop` - cyclomatic complexity < 15

#### Naming
| Construct | Convention |
|-----------|------------|
| Packages | lowercase, single word, no underscores |
| Interfaces | noun or -er suffix (Reader, Logger) |
| Structs | PascalCase |
| Functions/Methods | PascalCase (exported), camelCase (unexported) |
| Variables | camelCase |
| Constants | PascalCase (exported), camelCase (unexported) |
| Error types | *Error suffix |
| Test files | *_test.go |
| Mock files | mock_*.go |

#### Error Handling
```go
// Typed errors with codes
type ValidationError struct {
    Code    string
    Message string
    Field   string
}

func (e *ValidationError) Error() string { return e.Message }

// Wrap with context
return fmt.Errorf("validate user: %w", err)

// Sentinel errors for control flow
var ErrNotFound = errors.New("not found")
```

#### Testing
- Table-driven tests preferred
- `testify/require` for assertions
- Mock with `gomock` or interfaces
- Coverage target: 80%+ (enforced)
- No test helpers in production code

#### Dependencies
- `go.mod` / `go.sum` committed
- `go.work` for multi-module workspace
- Dependabot for updates
- License scanning (Apache 2.0, MIT, BSD only)

### TypeScript (Secondary - CLI, Web, Tooling)

#### Version
- TypeScript 5.4+ (strict mode)

#### Formatting
- `prettier` - mandatory
- Line length: 100 characters
- Single quotes, trailing commas

#### Linting
- `eslint` with:
  - `@typescript-eslint/recommended`
  - `eslint-plugin-import`
  - `eslint-plugin-unused-imports`
  - `eslint-plugin-security`

#### Naming
| Construct | Convention |
|-----------|------------|
| Files | kebab-case.ts |
| Classes/Interfaces | PascalCase |
| Functions/Variables | camelCase |
| Constants | UPPER_SNAKE_CASE |
| Types/Interfaces | PascalCase |
| Enums | PascalCase (singular) |

#### Type Safety
- `strict: true` in tsconfig
- No `any` - use `unknown` or generics
- Explicit return types for public APIs
- Discriminated unions for state

## Cross-Language Standards

### API Contracts
- Protocol Buffers (proto3) for gRPC
- OpenAPI 3.1 for REST
- Generated code committed
- Breaking changes require ADR

### DTOs
- Immutable (readonly/final)
- Validation tags/annotations
- No business logic
- Versioned in API definitions

### Logging
```go
// Structured logging - Go
log.Info().
    Str("correlation_id", cid).
    Str("tenant_id", tid).
    Dur("duration", elapsed).
    Msg("workflow completed")
```

```typescript
// Structured logging - TypeScript
logger.info({
  correlation_id: cid,
  tenant_id: tid,
  duration_ms: elapsed,
}, "workflow completed");
```

### Metrics
- RED metrics (Rate, Errors, Duration) for all services
- Business metrics (executions, tokens, costs)
- Histogram buckets: 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s

### Tracing
- W3C TraceContext propagation
- Span per logical operation
- Attributes: service, operation, tenant, error
- Sampling: 100% errors, 10% success (configurable)

## Git Standards

### Commit Messages
```
<type>(<scope>): <subject>

<body>

<footer>
```

Types: `feat`, `fix`, `refactor`, `perf`, `test`, `docs`, `chore`, `security`, `build`, `ci`

Scope: module name (gateway, workflow, harness, sandbox, router, policy, audit, shared)

Subject: imperative, lowercase, no period, max 72 chars

Body: explains what and why, not how

Footer: breaking changes, issue references

Example:
```
feat(gateway): add tenant-scoped rate limiting

Implement token bucket algorithm per tenant to prevent
noisy neighbor problems. Configuration via policy engine.

Closes #123
BREAKING CHANGE: rate limit headers changed
```

### Branch Names
- `feature/<short-description>`
- `fix/<short-description>`
- `hotfix/<short-description>`
- `refactor/<short-description>`
- `chore/<short-description>`
- `docs/<short-description>`
- `test/<short-description>`
- `security/<short-description>`

## Code Review Standards

### Required Checks
- [ ] Compiles without warnings
- [ ] All tests pass
- [ ] Coverage >= 80%
- [ ] Lint passes
- [ ] No security findings (gosec, npm audit)
- [ ] No secrets detected
- [ ] Architecture tests pass
- [ ] Documentation updated

### Review Focus Areas
1. Security - authz, validation, secrets, injection
2. Correctness - logic, edge cases, error handling
3. Architecture - boundaries, dependencies, patterns
4. Performance - allocations, queries, algorithms
5. Observability - logs, metrics, traces
6. Testing - coverage, quality, determinism
7. Documentation - public APIs, decisions, runbooks

## Documentation Standards

### Code Comments
- Package comments for all packages
- Exported symbols documented
- Why, not what
- No commented-out code

### Architecture Decision Records
- One ADR per significant decision
- Template: Context, Decision, Consequences
- Status: Proposed, Accepted, Superseded, Deprecated

### README per Module
- Purpose
- API surface
- Configuration
- Dependencies
- Run/deploy instructions

## Security Coding Standards

### Input Validation
- Validate at boundaries
- Allowlist over blocklist
- Size limits on all inputs
- Canonicalize before validation

### Secrets
- Never in code, config, logs, tests
- Use secret manager
- Rotate regularly
- Audit access

### Cryptography
- Standard libraries only
- TLS 1.2+ everywhere
- Proper key management
- No custom crypto

### Dependencies
- Minimal, pinned, scanned
- SBOM generated
- License compliance
- Regular updates

## Performance Standards

### Allocations
- Preallocate slices/maps
- Object pools for hot paths
- Avoid interface{} in hot paths
- Profile before optimizing

### Concurrency
- Bounded goroutines
- Context for cancellation
- Worker pools for parallelism
- No goroutine leaks

### Database
- Connection pooling
- Prepared statements
- Read replicas for reads
- Pagination mandatory

## Enforcement

### Pre-commit
- Format check
- Lint check
- Test (fast subset)
- Secret scan

### CI Pipeline
- Full test suite
- Coverage gate
- Security scan
- Architecture test
- Dependency scan
- Container scan

### Release
- All quality gates pass
- Human approval
- SBOM attestation
- Signature verification
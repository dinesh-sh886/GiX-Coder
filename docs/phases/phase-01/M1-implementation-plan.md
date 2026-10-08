# Phase 01 M1 Implementation Plan

## Overview

**Phase**: 01 - Foundation & Core Agent Infrastructure
**Milestone**: M1 - Agent Gateway
**Status**: Implementation Planning
**Branch**: feature/phase-01-m1-agent-gateway

## Architecture Overview

The Agent Gateway is the external API/control boundary for GiX-Coder. It implements:

1. **HTTP/REST API** - Gateway REST API with OpenAPI 3.1
2. **gRPC API** - Gateway gRPC service
3. **Authentication** - JWT/OIDC validation
4. **Authorization** - Policy engine stub integration
5. **Rate Limiting** - Per tenant/workflow
6. **Request Validation** - protovalidate integration
7. **Idempotency** - Key-based deduplication
8. **Observability** - OpenTelemetry, Prometheus, structured logging
9. **Audit Integration** - Security event logging
10. **PostgreSQL Persistence** - Executions, idempotency keys

## Implementation Plan

### Phase 1: Foundation & Configuration (Week 1)

#### 1.1 Gateway Configuration

- [ ] Gateway config struct extending shared Config
- [ ] Gateway-specific config fields (JWT, rate limit, idempotency)
- [ ] Environment variable bindings
- [ ] Config validation

#### 1.2 Main Entry Point

- [ ] `cmd/gateway/main.go` - application entry point
- [ ] Graceful shutdown handling
- [ ] Signal handling (SIGTERM, SIGINT)
- [ ] Configuration loading

#### 1.3 Shared Module Integration

- [ ] Logging initialization
- [ ] Metrics initialization
- [ ] Tracing initialization
- [ ] Config loading

### Phase 2: Core Middleware (Week 1-2)

#### 2.1 Authentication Middleware

- [ ] JWT validation (RS256)
- [ ] JWKS fetching and caching
- [ ] Issuer/Audience validation
- [ ] Token expiry/not-before validation
- [ ] Claims extraction
- [ ] Authentication failure handling

#### 2.2 Authorization Middleware

- [ ] Policy engine stub interface
- [ ] Capability grant evaluation
- [ ] Deny-by-default behavior
- [ ] Authorization decision logging

#### 2.3 Rate Limiting Middleware

- [ ] Token bucket algorithm
- [ ] Per-tenant/workflow limits
- [ ] Redis-backed distributed counters
- [ ] Retry-After header support
- [ ] Configuration-driven limits

#### 2.4 Request Validation Middleware

- [ ] protovalidate integration
- [ ] Request body parsing
- [ ] Validation error handling
- [ ] Request size limits

#### 2.5 Idempotency Middleware

- [ ] Idempotency key extraction
- [ ] PostgreSQL-backed persistence
- [ ] Duplicate request detection
- [ ] Concurrent duplicate handling
- [ ] TTL-based expiration

#### 2.6 Observability Middleware

- [ ] Request/response logging
- [ ] Prometheus metrics (RED)
- [ ] OpenTelemetry tracing
- [ ] Correlation ID propagation

### Phase 3: HTTP/gRPC Handlers (Week 2)

#### 3.1 REST API Handlers

- [ ] `POST /workflows/execute` - ExecuteWorkflow
- [ ] `GET /workflows/{execution_id}` - GetExecution
- [ ] `POST /workflows/{execution_id}:cancel` - CancelExecution
- [ ] `GET /health/live` - HealthLive
- [ ] `GET /health/ready` - HealthReady

#### 3.2 gRPC Handlers

- [ ] ExecuteWorkflow
- [ ] GetExecution
- [ ] CancelExecution
- [ ] HealthLive
- [ ] HealthReady

#### 3.3 gRPC Server

- [ ] gRPC server setup
- [ ] Interceptors (auth, logging, metrics, tracing)
- [ ] Reflection support
- [ ] Health checking

### Phase 4: Persistence & Audit (Week 3)

#### 4.1 PostgreSQL Migrations

- [ ] Gateway schema creation
- [ ] Executions table
- [ ] Idempotency keys table
- [ ] Migration versioning

#### 4.2 Repository Layer

- [ ] Execution repository
- [ ] Idempotency repository
- [ ] Transaction management

#### 4.3 Audit Integration

- [ ] Audit client
- [ ] Event emission for security events
- [ ] Structured audit events

### Phase 5: Testing & Integration (Week 3-4)

#### 5.1 Unit Tests

- [ ] Auth middleware tests
- [ ] Rate limit tests
- [ ] Validation tests
- [ ] Idempotency tests
- [ ] Handler tests

#### 5.2 Integration Tests

- [ ] Auth success/failure
- [ ] Rate limit enforcement
- [ ] Idempotency replay
- [ ] Handler contract tests

#### 5.3 Contract Tests

- [ ] Protobuf validation
- [ ] OpenAPI contract validation
- [ ] gRPC contract validation

### Phase 6: Validation & Documentation (Week 4)

#### 6.1 M0 Regression

- [ ] Run all M0 verification steps
- [ ] Verify no regressions

#### 6.2 Documentation

- [ ] Module README
- [ ] API documentation
- [ ] Configuration reference

#### 6.3 Final Validation

- [ ] Run complete test suite
- [ ] Run M0 regression
- [ ] Generate implementation report

## Package Structure

```
gateway/
├── cmd/
│   └── gateway/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers.go      # HTTP handlers
│   │   ├── grpc.go          # gRPC handlers
│   │   └── middleware.go    # HTTP middleware
│   ├── auth/
│   │   ├── jwt.go           # JWT validation
│   │   ├── jwks.go          # JWKS fetching
│   │   └── middleware.go    # Auth middleware
│   ├── authz/
│   │   ├── policy.go        # Policy stub interface
│   │   └── middleware.go    # AuthZ middleware
│   ├── middleware/
│   │   ├── auth.go          # Auth middleware
│   │   ├── authz.go         # AuthZ middleware
│   │   ├── rate_limit.go    # Rate limiting
│   │   ├── validation.go    # Request validation
│   │   ├── idempotency.go   # Idempotency
│   │   ├── observability.go # Logging, metrics, tracing
│   │   └── recovery.go      # Panic recovery
│   ├── persistence/
│   │   ├── execution.go     # Execution repository
│   │   ├── idempotency.go   # Idempotency repository
│   │   └── migration.go     # Migrations
│   ├── audit/
│   │   └── client.go        # Audit client
│   ├── config/
│   │   └── config.go        # Gateway config
│   └── server/
│       ├── http.go          # HTTP server setup
│       └── grpc.go          # gRPC server setup
├── pkg/
│   └── dto/
│       └── gateway.go       # Gateway-specific DTOs
├── go.mod
├── go.sum
└── README.md
```

## Dependencies

### New Dependencies (Gateway-specific)

- `github.com/gin-gonic/gin` - HTTP framework
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/prometheus/client_golang` - Prometheus metrics
- `go.opentelemetry.io/otel` - OpenTelemetry
- `github.com/gin-contrib/cors` - CORS middleware
- `github.com/gin-contrib/requestid` - Request ID middleware

### Existing Dependencies (from shared)

- `github.com/gix-coder/gix-coder/shared` - Shared kernel
- `github.com/gix-coder/gix-coder/api/proto/...` - Generated protobuf

## Configuration

### Environment Variables

```
GATEWAY_SERVICE_NAME=gateway
GATEWAY_SERVICE_VERSION=dev
GATEWAY_ENVIRONMENT=development

GATEWAY_SERVER_HTTP_PORT=8080
GATEWAY_SERVER_HTTP_HOST=0.0.0.0
GATEWAY_SERVER_GRPC_PORT=9090

GATEWAY_LOGGING_LEVEL=info
GATEWAY_LOGGING_FORMAT=json

GATEWAY_TRACING_ENABLED=true
GATEWAY_TRACING_ENDPOINT=
GATEWAY_TRACING_SAMPLER=parentbased_traceidratio
GATEWAY_TRACING_RATIO=0.1

GATEWAY_AUTH_JWKS_URL=https://example.com/.well-known/jwks.json
GATEWAY_AUTH_ISSUER=https://example.com
GATEWAY_AUTH_AUDIENCE=gix-coder

GATEWAY_RATE_LIMIT_ENABLED=true
GATEWAY_RATE_LIMIT_REQUESTS_PER_MINUTE=100
GATEWAY_RATE_LIMIT_BURST=200

GATEWAY_IDEMPOTENCY_TTL=24h

GATEWAY_DATABASE_HOST=localhost
GATEWAY_DATABASE_PORT=5432
GATEWAY_DATABASE_DATABASE=gix_coder
GATEWAY_DATABASE_USER=gix_coder
GATEWAY_DATABASE_PASSWORD=secret
GATEWAY_DATABASE_SSL_MODE=disable
GATEWAY_DATABASE_SCHEMA=gateway

GATEWAY_REDIS_HOST=localhost
GATEWAY_REDIS_PORT=6379
GATEWAY_REDIS_PASSWORD=
GATEWAY_REDIS_DB=0
```

## Acceptance Criteria Mapping

| AC ID | Requirement                            | Implementation            |
| ----- | -------------------------------------- | ------------------------- |
| AC-01 | Gateway accepts valid workflow request | ExecuteWorkflow handler   |
| AC-02 | Gateway rejects invalid JWT            | JWT validation middleware |
| AC-03 | Gateway enforces rate limits           | Rate limit middleware     |
| AC-18 | All quality gates pass                 | CI/CD pipeline            |
| AC-19 | No critical/high vulnerabilities       | Security scanning         |

## Risks & Mitigations

| Risk                        | Mitigation                                  |
| --------------------------- | ------------------------------------------- |
| JWT validation complexity   | Use established library (golang-jwt/jwt/v5) |
| Rate limit Redis dependency | Graceful degradation when Redis unavailable |
| Idempotency race conditions | PostgreSQL advisory locks                   |
| gRPC/HTTP consistency       | Shared handler logic                        |
| JWKS caching                | In-memory cache with TTL                    |

## Milestone Dates

| Phase                 | Target   |
| --------------------- | -------- |
| Foundation & Config   | Week 1   |
| Core Middleware       | Week 1-2 |
| HTTP/gRPC Handlers    | Week 2   |
| Persistence & Audit   | Week 3   |
| Testing & Integration | Week 3-4 |
| Validation & Docs     | Week 4   |

## Next Steps

1. Create M1 implementation plan document ✓
2. Set up Gateway module structure
3. Implement configuration and main entry point
4. Implement authentication middleware
5. Implement rate limiting
6. Implement request validation
7. Implement handlers
8. Implement persistence layer
9. Write tests
10. Run validation

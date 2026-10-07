# GiX-Coder Observability Strategy

## Observability Pillars

### 1. Structured Logs
- **Format**: JSON (mandatory)
- **Transport**: stdout → collector → Loki
- **Sampling**: Tail-based (keep errors, sample info)
- **Retention**: DEV 7d, STAGE 30d, PROD 1y (audit 7y)

#### Required Fields
```json
{
  "timestamp": "2024-01-15T10:30:00.123Z",
  "level": "info",
  "service": "gateway",
  "version": "1.2.3",
  "environment": "prod",
  "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
  "span_id": "00f067aa0ba902b7",
  "correlation_id": "req-abc123",
  "tenant_id": "tenant-456",
  "user_id": "user-789",
  "message": "workflow started",
  "fields": {
    "workflow_id": "wf-123",
    "execution_time_ms": 45
  }
}
```

#### Log Levels
| Level | Use Case | Sampling |
|-------|----------|----------|
| DEBUG | Detailed diagnostics | 1% (dev only) |
| INFO | Business events, lifecycle | 10% |
| WARN | Recoverable issues, degraded | 100% |
| ERROR | Failures requiring action | 100% |

#### Logging Standards
- One log line per logical event
- No PII in logs (use correlation IDs)
- Structured fields, not string concatenation
- Context: always include trace_id, span_id, correlation_id
- Errors: log with full context, not just error message

### 2. Metrics
- **Format**: Prometheus exposition
- **Collection**: Prometheus (scrape) → Cortex/Thanos (long-term)
- **Retention**: 15m resolution 14d, 1h resolution 1y

#### Standard Metrics (RED)
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `http_requests_total` | Counter | service, method, path, code | Request rate |
| `http_request_duration_seconds` | Histogram | service, method, path | Latency |
| `http_requests_in_flight` | Gauge | service | Active requests |
| `grpc_requests_total` | Counter | service, method, code | gRPC rate |
| `grpc_request_duration_seconds` | Histogram | service, method | gRPC latency |

#### Business Metrics
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `executions_total` | Counter | service, tenant, status | Workflow executions |
| `execution_duration_seconds` | Histogram | service, tenant, workflow_type | Execution time |
| `tokens_consumed_total` | Counter | service, tenant, model | LLM token usage |
| `sandbox_executions_total` | Counter | service, tenant, status | Sandbox runs |
| `active_tenants` | Gauge | - | Current active tenants |

#### System Metrics
| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `process_cpu_seconds_total` | Counter | service, instance | CPU usage |
| `process_resident_memory_bytes` | Gauge | service, instance | Memory |
| `process_open_fds` | Gauge | service, instance | File descriptors |
| `go_goroutines` | Gauge | service, instance | Goroutines (Go) |

#### Alerting Metrics (SLO-based)
- **Availability**: `sum(rate(http_requests_total{code!~"5.."}[5m])) / sum(rate(http_requests_total[5m])) > 0.999`
- **Latency**: `histogram_quantile(0.99, http_request_duration_seconds) < 0.5`
- **Error Rate**: `sum(rate(http_requests_total{code=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) < 0.001`

### 3. Distributed Tracing
- **Standard**: W3C TraceContext
- **Propagation**: B3 + TraceContext headers
- **Sampling**: 
  - 100% errors
  - 10% success (configurable)
  - 100% security events
- **Storage**: Tempo (blocks) → S3
- **Retention**: 7d (full), 30d (sampled)

#### Span Attributes (Required)
| Attribute | Description |
|-----------|-------------|
| `service.name` | Service identifier |
| `service.version` | Deployed version |
| `deployment.environment` | dev/stage/prod |
| `trace.id` | Trace ID |
| `span.id` | Span ID |
| `span.kind` | SERVER, CLIENT, PRODUCER, CONSUMER, INTERNAL |
| `http.method` | HTTP method |
| `http.route` | Route pattern |
| `http.status_code` | Response code |
| `net.peer.name` | Downstream service |
| `tenant.id` | Tenant identifier |
| `user.id` | User identifier |
| `error` | True if error |
| `error.message` | Error message |

### 4. Audit Events
- **Immutable**: Append-only, tamper-evident (Merkle tree / hash chain)
- **Storage**: Separate audit store (immutable object storage)
- **Retention**: 7 years minimum
- **Query**: Separate API, strict RBAC

#### Audit Event Structure
```json
{
  "event_id": "evt-uuid",
  "timestamp": "2024-01-15T10:30:00.123Z",
  "event_type": "workflow.executed",
  "actor": {
    "type": "user|service|agent",
    "id": "user-123",
    "tenant_id": "tenant-456"
  },
  "resource": {
    "type": "workflow",
    "id": "wf-123",
    "tenant_id": "tenant-456"
  },
  "action": "execute",
  "outcome": "success|failure|denied",
  "details": {
    "workflow_type": "code_review",
    "duration_ms": 45000,
    "tokens_used": 15000
  },
  "security_context": {
    "source_ip": "10.0.0.1",
    "user_agent": "gix-cli/1.0",
    "permissions": ["workflow:execute"]
  },
  "integrity": {
    "prev_hash": "sha256:abc...",
    "hash": "sha256:def..."
  }
}
```

## Health & Readiness

### Endpoints
| Endpoint | Purpose | Checks |
|----------|---------|--------|
| `/health/live` | Liveness (k8s) | Process alive, no deadlock |
| `/health/ready` | Readiness (k8s) | Dependencies reachable, config loaded, not draining |
| `/health/startup` | Startup (k8s) | Initialization complete |

### Response Format
```json
{
  "status": "healthy|degraded|unhealthy",
  "checks": [
    {
      "name": "database",
      "status": "healthy",
      "latency_ms": 5,
      "details": {}
    }
  ],
  "timestamp": "2024-01-15T10:30:00.123Z"
}
```

## Correlation & Request IDs

### Header Propagation
```
Incoming:  X-Correlation-ID, X-Request-ID, traceparent
Internal:  Generated if missing, propagated to all calls
Outgoing:  All three headers forwarded
```

### Generation
- **Correlation ID**: UUID v7 (time-ordered), one per user request
- **Request ID**: UUID v4, one per service hop
- **Trace ID**: W3C traceparent, generated at edge

## Implementation Standards

### Go
```go
// Structured logging with zerolog
log.Info().
    Str("correlation_id", cid).
    Str("tenant_id", tid).
    Dur("duration", elapsed).
    Int("tokens", tokens).
    Msg("workflow completed")

// Metrics with promauto
var (
    executionDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
        Name:    "execution_duration_seconds",
        Help:    "Workflow execution duration",
        Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10, 30, 60},
    }, []string{"service", "tenant", "workflow_type", "status"})
)

// Tracing with otel
ctx, span := tracer.Start(ctx, "executeWorkflow",
    trace.WithAttributes(
        attribute.String("tenant.id", tenantID),
        attribute.String("workflow.id", workflowID),
    ),
)
defer span.End()
```

### TypeScript
```typescript
// Structured logging with pino
logger.info({
  correlation_id: cid,
  tenant_id: tid,
  duration_ms: elapsed,
  tokens: tokens,
}, "workflow completed");

// Metrics with prom-client
const executionDuration = new client.Histogram({
  name: 'execution_duration_seconds',
  help: 'Workflow execution duration',
  buckets: [0.1, 0.5, 1, 2.5, 5, 10, 30, 60],
  labelNames: ['service', 'tenant', 'workflow_type', 'status']
});

// Tracing with @opentelemetry/api
const span = tracer.startSpan('executeWorkflow', {
  attributes: {
    'tenant.id': tenantID,
    'workflow.id': workflowID,
  }
});
span.end();
```

## Dashboards (Standard Set)

### Per Service (Grafana)
1. **RED Dashboard** - Rate, Errors, Duration
2. **Resource Dashboard** - CPU, Memory, Network, Disk
3. **Business Dashboard** - Executions, Tokens, Costs
4. **Dependency Dashboard** - Upstream latency, errors
5. **Security Dashboard** - Auth failures, rate limits, audit events

### Platform
1. **Cluster Overview** - Nodes, pods, capacity
2. **Tenant Overview** - Usage, quotas, health
3. **Cost Dashboard** - Infrastructure, tokens, per tenant
4. **SLO Dashboard** - Error budget burn rate

## Alerting

### Alert Principles
- **Actionable**: Every alert requires human action
- **Routed**: Right team, right urgency
- **Deduplicated**: Group related alerts
- **Contextual**: Runbook link, dashboard link

### Severity Levels
| Severity | Response Time | Examples |
|----------|---------------|----------|
| P1 (Critical) | 15 min | Service down, data loss, security breach |
| P2 (High) | 1 hour | Degraded performance, partial outage |
| P3 (Medium) | 4 hours | Non-critical feature broken, warning thresholds |
| P4 (Low) | 24 hours | Info, capacity planning |

### Alert Routing
- P1/P2 → PagerDuty (on-call)
- P3 → Slack (team channel)
- P4 → Email (daily digest)

## Implementation Checklist

### Per Service
- [ ] Structured JSON logging
- [ ] Correlation ID propagation
- [ ] RED metrics exposed
- [ ] Business metrics exposed
- [ ] Health endpoints (live/ready/startup)
- [ ] Distributed tracing instrumentation
- [ ] Audit logging for security events
- [ ] Dashboards imported
- [ ] Alerts configured
- [ ] Runbooks linked

### Platform
- [ ] Log aggregation (Loki)
- [ ] Metrics collection (Prometheus/Cortex)
- [ ] Trace storage (Tempo)
- [ ] Audit store (immutable)
- [ ] Alert manager configured
- [ ] Notification channels
- [ ] SLO definitions
- [ ] Error budget alerts
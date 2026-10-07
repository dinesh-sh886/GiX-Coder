# Target System Architecture

## 1. Architectural goals

- Global SaaS
- Multi-tenant
- 99.9%+ availability
- horizontally scalable
- secure sandbox execution
- provider agnostic
- resumable tasks
- observable
- auditable
- maintainable
- API-first
- local + cloud
- enterprise-ready

## 2. Recommended high-level architecture

```text
                         Internet
                            |
                     Cloudflare / CDN
                            |
                    API Gateway / WAF
                            |
              +-------------+-------------+
              |                           |
          Web App                    Public API
              |                           |
              +-------------+-------------+
                            |
                    Control Plane
                            |
      +---------+-----------+-----------+----------+
      |         |           |           |          |
    Auth     Projects     Tasks       Billing    Policy
      |         |           |           |          |
      +---------+-----------+-----------+----------+
                            |
                     Agent Gateway
                            |
                    Durable Job System
                            |
              +-------------+-------------+
              |                           |
       Execution Scheduler          Event Bus
              |                           |
              v                           v
       Sandbox Pool                  Analytics
              |
      +-------+-------+-------+
      |       |       |       |
    Repo    Tools   Browser  Test Runner
      |
      v
 GitHub/GitLab/etc.
```

## 3. Control plane vs data plane

This separation is mandatory.

### Control plane

Owns:

- identity
- tenant configuration
- tasks
- sessions
- policies
- billing
- metadata
- scheduling
- audit
- analytics

### Data plane

Owns:

- repository files
- shell execution
- builds
- tests
- browser automation
- secrets
- network access

Never let untrusted code execution directly access control-plane
databases.

## 4. Suggested technology stack

### Web

- Next.js / React
- TypeScript
- Tailwind
- Monaco editor where needed
- React Query
- WebSocket/SSE event stream

### API/control plane

Recommended:

- Go for core backend services
- gRPC internally
- REST/JSON externally
- OpenAPI
- protobuf

Alternative if speed of MVP is more important:

- TypeScript + NestJS/Fastify

Long-term recommendation:

> Go for durable infrastructure services; TypeScript for
> product/frontend.

### Agent runtime

Recommended:

- Rust for sandbox/execution-critical components
- TypeScript or Python for orchestration/agent experimentation

Do not make the entire platform Python.

### Database

PostgreSQL:

- tenants
- users
- projects
- tasks
- sessions
- policies
- billing metadata
- audit indexes

### Cache

Redis:

- ephemeral cache
- rate limiting
- locks
- short-lived coordination

Do not make Redis the source of truth.

### Object storage

S3-compatible:

- logs
- diffs
- artifacts
- repository snapshots
- evaluation results
- large event payloads

### Event bus

Start:

- NATS JetStream

Scale:

- Kafka/Redpanda if event volume demands it.

### Analytics

ClickHouse:

- task events
- tool events
- token usage
- latency
- cost
- model performance
- product analytics

### Observability

- OpenTelemetry
- Prometheus
- Grafana
- Loki
- Tempo/Jaeger
- Sentry

### Workflow orchestration

Strong recommendation:

- Temporal

Why:

- durable workflows
- retries
- timers
- compensation
- task queues
- recovery
- long-running executions

Agent tasks are workflows, not HTTP requests.

## 5. Agent execution model

```text
Task created
    |
Temporal workflow
    |
Create sandbox
    |
Clone repository
    |
Load policy
    |
Load skills
    |
Build repository context
    |
Agent turn
    |
Tool call
    |
Policy check
    |
Approval?
  /      \
yes       no
 |         |
wait      execute
 |
resume
    |
Observe
    |
Update task state
    |
Continue/recover
    |
Tests
    |
Verification
    |
Commit/PR
    |
Destroy/retain sandbox
```

## 6. Sandbox design

Each task should get an isolated workspace.

Preferred options:

### Tier 1

Firecracker microVM.

### Tier 2

gVisor/container sandbox.

### Tier 3

Kubernetes isolated pod with strong security controls.

Use stronger isolation for:

- untrusted repositories
- arbitrary internet access
- customer-controlled build scripts
- package installation
- browser execution

## 7. Network policy

Default:

```text
Internet = DENY
```

Allow only:

- Git provider
- approved package registries
- model APIs
- user-approved domains
- internal services through explicit capability

DNS and HTTP egress should be policy-controlled.

## 8. Secret architecture

Never place raw customer secrets in model context.

Use:

```text
Agent
  |
Capability request
  |
Secret broker
  |
Policy
  |
Ephemeral credential
  |
Tool process
```

Secrets should be:

- short-lived
- scoped
- audited
- revocable

## 9. Durable state

Store:

- task state
- event sequence
- checkpoints
- plan
- tool calls
- approvals
- artifacts
- final result

A session must survive:

- worker crash
- API restart
- network disconnect
- region failover

## 10. Multi-region

Recommended evolution:

### Phase 1

One primary region + backups.

### Phase 2

Regional execution pools.

### Phase 3

Global control plane with region-aware data placement.

### Phase 4

Active/active control plane where justified.

Do not start with active/active everything. It dramatically increases
operational complexity.

## 11. Availability strategy

For 99.9%:

- multiple API replicas
- managed PostgreSQL with HA
- Redis HA
- durable event bus
- Temporal HA
- sandbox scheduler redundancy
- object storage replication
- health checks
- circuit breakers
- retries with jitter
- idempotency keys
- graceful degradation

## 12. Failure model

Every external dependency gets:

- timeout
- retry policy
- circuit breaker
- fallback
- observability
- dead-letter path

Model providers should support:

```text
Primary provider
      |
timeout/error
      v
Secondary provider
      |
quality issue
      v
Recovery model
```

## 13. API surface

Examples:

```text
POST   /v1/tasks
GET    /v1/tasks/:id
POST   /v1/tasks/:id/cancel
POST   /v1/tasks/:id/resume
POST   /v1/tasks/:id/approve

GET    /v1/sessions/:id
POST   /v1/sessions/:id/messages

GET    /v1/runs/:id/events
GET    /v1/runs/:id/artifacts

POST   /v1/repositories
POST   /v1/workspaces

GET    /v1/usage
GET    /v1/analytics
```

## 14. Event model

Use append-only events:

```text
TaskCreated
TaskStarted
PlanCreated
ToolRequested
ApprovalRequested
ApprovalGranted
ToolStarted
ToolCompleted
FileChanged
TestStarted
TestCompleted
AgentRecovered
TaskCompleted
TaskFailed
```

This becomes the foundation for:

- replay
- audit
- analytics
- debugging
- evaluation

## 15. Maintainability

Use bounded contexts:

```text
identity
organization
billing
repository
workspace
task
agent
execution
policy
integration
analytics
evaluation
```

Do not build one giant backend.

Start as a modular monolith if the team is small, but keep boundaries
explicit. Extract services only when scale or ownership requires it.

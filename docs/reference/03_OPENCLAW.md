# OpenClaw --- Architecture & Product Analysis

## 1. Positioning

OpenClaw is an open-source, self-hosted agent platform centered around a
long-lived Gateway. It connects LLM agents to many messaging channels
and devices.

It is broader than a coding agent, which makes it useful as an
architecture reference for agent infrastructure.

Sources: - https://github.com/openclaw/openclaw -
https://docs.openclaw.ai/concepts/architecture -
https://docs.openclaw.ai/gateway/security

## 2. Architecture

```text
                    +----------------------+
                    | CLI / Web / Mobile   |
                    +----------+-----------+
                               |
                         WebSocket
                               |
                               v
                    +----------------------+
                    |      Gateway         |
                    | identity/policy      |
                    | sessions/events       |
                    | channel connections   |
                    +---+----+----+---------+
                        |    |    |
                        v    v    v
                    Plugins Nodes Channels
                        |
                        v
                  Agent Runtime
                        |
                        v
                 LLM Provider(s)
```

The Gateway is a central control boundary.

## 3. Key design ideas

### Long-lived Gateway

A single Gateway owns messaging surfaces and provides a typed WebSocket
API.

### Nodes

Devices can connect as nodes with explicit capabilities.

### Plugins

Plugins can add:

- providers
- tools
- channels
- HTTP routes
- runtime services
- hooks

### Persistence

OpenClaw uses local durable state including SQLite WAL patterns.

### Security

The project has explicit security auditing and trust-boundary
documentation.

## 4. Technology stack

Current documentation/repository signals:

- Node.js
- TypeScript/JavaScript
- SQLite
- WebSockets
- plugin runtime
- multiple channel SDKs
- local gateway process

Current supported Node runtime is documented as Node 24.16+ or Node
26.1+, with Node 26 recommended.

## 5. Security lesson

OpenClaw explicitly warns that one Gateway is a trust boundary and is
not a hostile multi-tenant boundary.

That distinction is critical for a SaaS product.

For our SaaS:

```text
Tenant A
  |
  +-- Control plane
  +-- Agent session
  +-- Sandbox A
  +-- Secrets A

Tenant B
  |
  +-- Control plane
  +-- Agent session
  +-- Sandbox B
  +-- Secrets B
```

Never place mutually hostile tenants inside one trusted execution
boundary.

## 6. Pros

- Excellent gateway concept.
- Extensible plugin system.
- Multi-channel abstraction.
- Long-lived runtime.
- Explicit security posture.
- Automation support.
- Strong local/self-hosting story.

## 7. Cons for a coding SaaS

- Coding is not its primary product boundary.
- A broad channel abstraction creates significant complexity.
- Trusted-gateway assumptions cannot be directly reused for hostile
  SaaS tenancy.
- Local deployment is easier than global cloud execution.
- Messaging-centric UX is not sufficient for code review and software
  delivery.

## 8. What to copy

Copy:

- Gateway concept
- capability-based nodes
- event-driven architecture
- plugin lifecycle
- explicit security audit
- durable local state
- automation scheduler

Improve:

- tenant isolation
- ephemeral execution
- repository/workspace abstraction
- Git-native workflow
- code-aware context
- CI/test verification
- cloud orchestration
- organization analytics

## 9. Product lesson

Use the Gateway idea as the **agent control plane**, but place actual
code execution behind a stronger sandbox boundary.

Recommended:

```text
Global API
    |
Agent Gateway
    |
Job Queue
    |
Execution Scheduler
    |
Ephemeral Sandbox
    |
Repository + Secrets + Tools
```

# OpenCode --- Architecture & Product Analysis

## 1. Positioning

OpenCode is an open-source AI coding agent focused on terminal
workflows. Its major strategic differentiators are:

- provider neutrality
- TUI-first UX
- client/server architecture
- LSP support
- permissions
- plugins
- MCP
- multiple agents/subagents

Current documentation describes support for 75+ LLM providers and local
models.

Sources: - https://opencode.ai/docs/providers -
https://opencode.ai/docs/agents/ - https://github.com/anomalyco/opencode

## 2. Architecture

Conceptually:

```text
                 +------------------+
                 | TUI / Desktop    |
                 +---------+--------+
                           |
                           v
                 +------------------+
                 | OpenCode Server  |
                 | Session/Agent    |
                 | Tool orchestration|
                 +----+----+----+----+
                      |    |    |
             +--------+    |    +---------+
             v             v              v
          LSP/Repo       Tools          MCP
                           |
                           v
                    Model Provider
```

The client/server separation means the UI is not the agent itself.

The SDK exposes a typed client for controlling the server
programmatically.

## 3. Technology stack

Observed/current repository signals:

- TypeScript/JavaScript ecosystem
- Bun runtime
- SolidJS for UI surfaces
- Hono
- Effect
- SQLite integration
- AI SDK
- OpenTelemetry support
- Sentry
- Playwright
- tree-sitter
- node-pty
- TypeScript SDK

Repository: https://github.com/anomalyco/opencode

## 4. Agent architecture

Important abstractions:

- primary agents
- subagents
- permissions
- tool execution
- sessions
- compaction
- plugins
- model/provider abstraction

Permission levels include:

- allow
- ask
- deny

Fine-grained permissions can gate:

- read
- edit
- bash
- external directories
- web
- LSP
- skills
- subagents

## 5. User flow

Typical flow:

```text
Install
  |
Open repository
  |
Select/configure provider
  |
Choose model
  |
Describe task
  |
Agent inspects repository
  |
Agent plans
  |
Agent requests/uses tools
  |
Human approvals where required
  |
Code changes
  |
Tests / diagnostics
  |
Review diff
  |
Commit / PR
```

## 6. Monetization

OpenCode demonstrates a useful hybrid:

### Open-source core

The client/agent remains usable with external providers.

### Provider/gateway monetization

OpenCode Zen provides curated model/provider access and charges per
request. Team workspaces add organizational controls.

The strategic lesson:

> Monetize convenience, reliability, routing, governance and managed
> inference---not merely the open-source binary.

## 7. Pros

- Strong provider independence.
- Good developer ergonomics.
- Open-source trust.
- Good extension model.
- LSP integration.
- Strong permissions.
- Client/server enables future remote control.
- Large ecosystem potential.

## 8. Cons

- Terminal-first can limit less technical users.
- Provider neutrality can make quality consistency difficult.
- Open ecosystem increases compatibility burden.
- A SaaS business must build stronger hosted execution and tenant
  isolation around a local-first core.
- Agent success still depends heavily on model quality and task
  formulation.

## 9. What to copy

Copy:

- provider abstraction
- permission model
- agent/subagent separation
- client/server architecture
- SDK
- plugin hooks
- LSP integration
- MCP
- TUI quality

Improve:

- durable cloud sessions
- multi-region execution
- organization policy
- outcome analytics
- cost forecasting
- sandbox orchestration
- branch/PR automation
- agent reliability
- recovery
- evaluation platform

## 10. Product lesson

OpenCode is a strong **agent runtime/client**.

A new SaaS should build the **system around the runtime**:

- identity
- organizations
- repositories
- workspaces
- cloud sandboxes
- job queues
- billing
- observability
- policy
- audit
- analytics

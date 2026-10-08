# OpenAI Codex --- Architecture & Product Analysis

## 1. Positioning

Codex is especially important as an architectural reference because
OpenAI exposes the **agent harness** as a reusable platform.

Current Codex surfaces include:

- CLI
- app
- IDE integrations
- cloud
- SDK
- app-server
- managed agent APIs

Sources: - https://github.com/openai/codex -
https://developers.openai.com/api/docs/guides/agents-api/architecture -
https://developers.openai.com/blog/codex-as-a-platform

## 2. Core architecture

OpenAI describes three major pieces:

```text
Application Server
       |
       v
Agent Harness
       |
       v
Execution Environment
```

### Harness

Responsible for:

- model/tool loop
- context
- session
- approvals
- policies
- streaming
- recovery

### Environment

Where the agent:

- runs commands
- edits files
- executes tests
- accesses configured resources

The environment can be:

- hosted sandbox
- laptop
- Docker
- cloud infrastructure
- self-hosted environment

### Application server

Owns product-specific:

- UI
- business context
- tools
- lifecycle
- approvals
- integration

This separation is one of the strongest architectural patterns in the
market.

## 3. App-server

Codex app-server provides a programmatic interface around the agent.

Important characteristics:

- bidirectional protocol
- JSON-RPC
- threads
- turns
- streamed events
- approval requests
- session lifecycle

This enables the agent to become a platform component.

## 4. Technology stack

The open-source Codex CLI repository is heavily Rust-based and organized
as a large Rust workspace.

Current repository structure includes:

- `codex-rs`
- `codex-cli`
- app-server
- execution/sandbox components
- SDKs
- protocol crates
- testing infrastructure

The architecture strongly favors modular Rust components for the
execution core.

## 5. Security

Codex exposes explicit:

- sandbox modes
- approval policies
- workspace boundaries
- environment configuration
- command execution controls

A coding agent should treat execution as a security product, not a
simple subprocess call.

## 6. User flow

```text
User
 |
 v
Task
 |
 v
Session
 |
 v
Agent plans/reasons
 |
 v
Context retrieval
 |
 v
Tool call
 |
 +---- approval required? ---- yes ---> Human
 |                                      |
 no <-----------------------------------+
 |
 v
Execute in environment
 |
 v
Observe result
 |
 v
Continue / recover
 |
 v
Tests
 |
 v
Final result + diff
```

## 7. Monetization

Codex participates in a broader OpenAI model/platform business:

- consumer/team/enterprise subscriptions
- API/model usage
- hosted agent infrastructure
- developer platform capabilities

The important strategic lesson is that an agent product can act as a
distribution surface for model and cloud consumption.

## 8. Pros

- Strong harness abstraction.
- Excellent execution separation.
- Rich protocol architecture.
- Strong sandbox/approval model.
- Multiple clients.
- Strong cloud path.
- Good foundation for product embedding.

## 9. Cons

- Model ecosystem can create provider dependence.
- Infrastructure complexity is high.
- Cloud execution is expensive.
- Deep integration can make portability difficult.

## 10. What to copy

Copy almost exactly at the conceptual level:

```text
Agent Harness
      |
Environment
      |
Application
```

Then add:

- multi-model router
- provider health
- cost optimizer
- tenant policy
- organization analytics
- cross-agent evaluation
- Git provider abstraction

## 11. Product lesson

The most defensible technical asset is not the chat UI.

It is the **reliable agent harness + execution protocol**.

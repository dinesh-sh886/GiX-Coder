# Competitive Architecture Comparison

## Executive summary

---

Product Primary Architecture Commercial pattern Best lesson
strength pattern

---

OpenCode Provider Client/server, Open-source + optional Avoid model
neutrality + extensible hosted model lock-in
terminal UX plugins, many gateway/subscriptions  
providers

OpenClaw Broad agent Long-lived Open-source/non-commercial Separate
gateway + trusted hosted-free model control plane
channels gateway + from execution
plugins + nodes

Claude Code Excellent Integrated Subscription + usage + UX and
developer proprietary enterprise workflow
workflow agent/runtime matter as much
as model

Codex Reusable agent Harness + Subscription + API/model Make the agent
harness environment + usage + platform runtime
app-server + SDK reusable

OpenHands Open-source Agent Server + OSS + SaaS + Offer local,
autonomous SDK + enterprise/self-hosted cloud and
coding workspace + UI enterprise
deployment
------------------------------------------------------------------------------------------

## Capability matrix

---

Capability OpenCode OpenClaw Claude Code Codex OpenHands Proposed
product

---

CLI Strong Strong Strong Strong Yes Yes

IDE integration Good Indirect Strong Strong Strong First-class

Web UI Yes Yes Yes/related Claude Yes Strong Strong
surfaces

Multi-model Excellent Excellent Primarily Anthropic Primarily Excellent Excellent
OpenAI,  
integrations  
vary

Local execution Excellent Excellent Strong Strong Excellent Yes

Cloud execution Limited/optional Self-host Strong Strong Strong Core
focus

Sandboxing Configurable Strong Strong Strong Strong Strong
host/gateway  
model

Durable sessions Yes Yes Yes Strong Strong Core

Background tasks Emerging Strong Strong Strong Strong Core

MCP Strong Strong Strong Strong Strong Core

Skills/plugins Strong Strong Strong Strong Strong Core

Enterprise Growing Limited as Excellent Strong Strong Excellent
governance hostile  
multi-tenant  
boundary

Provider Excellent Excellent Low Medium Excellent Excellent
neutrality

Agent analytics Good Operational Excellent Strong Good Excellent
focus

Self-hosting Strong Core Limited CLI/harness Strong Enterprise
components tier

Mobile control Possible Strong Product-dependent Cloud/app Web/mobile Planned
cloud
----------------------------------------------------------------------------------------------------------------

## What each product gets right

### OpenCode

- Treats the agent as a client/server system.
- Strong model/provider abstraction.
- LSP, MCP, plugins and permissions are first-class.
- TUI remains excellent for developers.
- Hosted model gateway creates a monetization surface without forcing
  one model.

Sources: - https://opencode.ai/docs/providers -
https://opencode.ai/docs/agents/ - https://opencode.ai/docs/en/zen/ -
https://github.com/anomalyco/opencode

### OpenClaw

- Gateway is a powerful architectural boundary.
- Channels, nodes, automation and plugins are unified.
- Security audit and deterministic policy are explicit.
- Good example of an agent platform rather than a coding-only tool.

Important limitation: OpenClaw explicitly says its gateway is not a
hostile multi-tenant boundary. A SaaS product serving mutually untrusted
tenants needs stronger isolation.

Sources: - https://docs.openclaw.ai/concepts/architecture -
https://docs.openclaw.ai/gateway/security -
https://github.com/openclaw/openclaw

### Claude Code

- Excellent developer-first workflow.
- Strong commercial packaging.
- Enterprise governance, usage analytics and policy controls are
  important.
- Skills create reusable domain expertise.
- Strong model/agent quality is a major moat.

Source: -
https://www.anthropic.com/news/claude-code-on-team-and-enterprise

### Codex

- The most important architectural lesson is the reusable **harness**.
- Harness, environment and application server are separate concepts.
- Sessions are durable.
- Streaming/events/approvals are protocol-level concerns.
- App-server makes the agent embeddable in products.

Sources: -
https://developers.openai.com/api/docs/guides/agents-api/architecture -
https://developers.openai.com/blog/codex-as-a-platform -
https://github.com/openai/codex

### OpenHands

- Strong open-source architecture.
- Agent Server cleanly separates UI from execution.
- SDK exposes agent loops and tools.
- Local, SaaS and enterprise deployment creates a broad adoption
  funnel.
- ACP support enables external coding agents to plug into the UI.

Sources: - https://github.com/OpenHands/OpenHands -
https://github.com/OpenHands/docs/blob/main/sdk/arch/agent.mdx -
https://www.openhands.dev/pricing/

## Biggest weaknesses to attack

### 1. Agent reliability

Most systems still depend heavily on the underlying model. Build a
reliability layer:

- task decomposition
- checkpoints
- deterministic retries
- tool-result validation
- test-driven completion
- failure classification
- rollback
- recovery agents
- quality gates

### 2. Context economics

Long coding sessions can become expensive.

Build:

- repository map
- symbol index
- semantic code search
- AST-aware retrieval
- LSP diagnostics
- compacted execution memory
- reusable context cache
- branch-aware context

### 3. Safety

Never rely only on a prompt such as "be careful."

Use:

- capability tokens
- workspace mounts
- syscall/process policy
- network egress policy
- secret broker
- command risk classification
- approval policy
- immutable audit trail
- tenant isolation

### 4. Cost unpredictability

Give users:

- budget per task
- budget per workspace
- budget per organization
- token and compute forecasts
- model routing
- spend alerts
- hard limits
- cost attribution

### 5. Agent observability

Every run should answer:

- What was requested?
- What did the agent plan?
- Which tools ran?
- Which files changed?
- Which tests passed?
- Which model was used?
- How much did it cost?
- Where did time go?
- Why did it fail?
- Can it resume?

## Strategic positioning

The best product is not "OpenCode + Claude + Codex."

It is:

> **A global software-engineering execution platform where developers
> delegate scoped engineering work to reliable, observable,
> policy-controlled AI agents.**

The agent is the engine. The product is the workflow.

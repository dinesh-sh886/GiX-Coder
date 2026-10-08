# OpenHands --- Architecture & Product Analysis

## 1. Positioning

OpenHands is a major open-source software-agent project focused on
AI-driven development.

The current architecture separates:

- frontend/control surface
- Agent Server
- SDK
- agents
- tools
- conversations
- workspaces
- events

Sources: - https://github.com/OpenHands/OpenHands -
https://github.com/OpenHands/docs/blob/main/sdk/arch/agent.mdx -
https://www.openhands.dev/pricing/

## 2. Architecture

```text
                  Agent Canvas
                       |
                       v
                Agent Server
                       |
              +--------+--------+
              |                 |
              v                 v
          Agent SDK          External ACP Agent
              |
       +------+------+
       |      |      |
      LLM   Tools  Context
       |      |      |
       +------+------+
              |
              v
          Workspace
```

The Agent Server provides a stable backend boundary.

## 3. Agent loop

OpenHands documents four major responsibilities:

1.  reasoning/action loop
2.  tool orchestration
3.  context management
4.  security validation

This is exactly the right decomposition for a production coding agent.

## 4. Context management

Long-running coding tasks require context management.

Important mechanisms:

- conversation history
- condensers
- skills
- context injection
- event history

For our product:

```text
Raw event stream
      |
      +--> short-term context
      |
      +--> summarized task state
      |
      +--> repository memory
      |
      +--> durable artifacts
```

## 5. ACP

OpenHands can drive external coding agents through Agent Client
Protocol.

This is strategically important.

Instead of forcing one agent implementation, the platform can act as a
**universal agent control plane**.

Our product should consider supporting:

- MCP for tools
- ACP-like protocol for agents
- OpenAI-compatible interfaces
- provider adapters

## 6. Commercial strategy

OpenHands uses a multi-tier model:

- free local open source
- hosted SaaS
- BYOK
- at-cost provider access
- enterprise SaaS/self-hosted
- private VPC
- SSO

This is an excellent adoption funnel.

## 7. Pros

- Open-source credibility.
- Strong architecture.
- Local + cloud.
- Enterprise deployment.
- Agent Server abstraction.
- SDK.
- External agent integration.
- Model agnostic.

## 8. Cons

- Large architecture surface.
- Cloud execution is operationally complex.
- More abstraction can increase setup complexity.
- Open-source/community expectations can slow commercial product
  decisions.

## 9. What to copy

Copy:

- Agent Server
- SDK
- event-driven agent loop
- context condensers
- workspace abstraction
- local/cloud/self-hosted
- ACP integration
- enterprise VPC model

Improve:

- polished developer onboarding
- automated model routing
- stronger analytics
- simpler pricing
- Git-native automation
- reliability scoring
- cost prediction

## 10. Product lesson

A successful open agent platform should have:

```text
Local agent
   +
Cloud agent
   +
Enterprise agent
   +
External agent compatibility
```

That reduces ecosystem lock-in and expands the addressable market.

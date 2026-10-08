# AI Coding Agent SaaS --- Research & Product Reference

**Status:** Architecture research baseline\
**Research date:** 2026-10-07\
**Target:** Global SaaS coding-agent platform, designed for high
availability, multi-tenancy, extensibility, and enterprise adoption.

## Purpose

This reference pack studies:

1.  OpenCode
2.  OpenClaw
3.  Claude Code
4.  OpenAI Codex
5.  OpenHands

It then converts the strongest patterns into a proposed product
architecture.

## Files

---

File Purpose

---

`01_COMPARATIVE.md` Executive comparison, pros/cons,
capability matrix

`02_OPENCODE.md` Architecture, UX, tech stack,
monetization, weaknesses

`03_OPENCLAW.md` Gateway architecture, security,
extensibility, lessons

`04_CLAUDE_CODE.md` Product strategy, workflow,
enterprise model, lessons

`05_CODEX.md` Harness, app-server, sandbox,
cloud/product strategy

`06_OPENHANDS.md` Agent Server, SDK,
cloud/self-hosting, lessons

`07_TARGET_PRODUCT_PRD.md` Proposed product vision and
requirements

`08_SYSTEM_ARCHITECTURE.md` Target production architecture and
technology choices

`09_SECURITY_MULTI_TENANCY.md` Threat model, sandboxing, identity,
tenancy

`10_MONETIZATION_ANALYTICS.md` Pricing, unit economics, product
analytics, KPIs

`11_ROADMAP.md` Phased implementation plan

`12_ARCHITECTURAL_DECISIONS.md` Key decisions and rationale
-----------------------------------------------------------------------

## Core conclusion

Do **not** build another terminal chatbot.

Build a **Coding Agent Platform** with:

- local CLI + IDE clients
- cloud execution
- remote workspaces
- durable sessions
- resumable tasks
- deterministic permissions
- isolated sandboxes
- model/provider abstraction
- MCP/tool ecosystem
- skills and plugins
- background agents
- GitHub/GitLab/Bitbucket workflows
- enterprise policy and audit
- cost controls
- outcome analytics
- human approval gates
- evaluation and regression infrastructure

The differentiator should be **reliable software delivery**, not simply
better chat.

## Reference sources

Primary/current sources consulted:

- OpenCode: https://opencode.ai/docs/
- OpenCode repository: https://github.com/anomalyco/opencode
- OpenClaw repository: https://github.com/openclaw/openclaw
- OpenClaw architecture:
  https://docs.openclaw.ai/concepts/architecture
- OpenHands: https://github.com/OpenHands/OpenHands
- OpenHands architecture:
  https://github.com/OpenHands/docs/blob/main/sdk/arch/agent.mdx
- OpenHands pricing: https://www.openhands.dev/pricing/
- OpenAI Codex repository: https://github.com/openai/codex
- OpenAI Codex platform architecture:
  https://developers.openai.com/api/docs/guides/agents-api/architecture
- Codex as a platform:
  https://developers.openai.com/blog/codex-as-a-platform
- Anthropic Claude Code announcement:
  https://www.anthropic.com/news/claude-3-7-sonnet
- Anthropic Claude Code enterprise controls:
  https://www.anthropic.com/news/claude-code-on-team-and-enterprise

> Product details change quickly. Treat this pack as an architectural
> baseline and re-verify pricing/model availability before
> implementation or commercial launch.

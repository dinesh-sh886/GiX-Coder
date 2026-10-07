# Product Roadmap

## Phase 0 --- Architecture foundation

Duration: 2--4 weeks

Deliver:

- product specification
- threat model
- domain model
- API contract
- agent protocol
- sandbox prototype
- model adapter interface
- evaluation framework

Exit criteria:

- one task can be executed end-to-end in an isolated workspace.

## Phase 1 --- Developer MVP

Duration: 6--10 weeks

Deliver:

- web app
- CLI
- GitHub integration
- repository clone
- agent loop
- file tools
- shell
- Git
- tests
- approvals
- task history
- basic billing

Target:

> User can give a GitHub issue and receive a verified PR.

## Phase 2 --- Reliability

Duration: 6--8 weeks

Deliver:

- checkpoints
- resume
- retry policies
- failure classifier
- recovery agent
- task state machine
- event sourcing
- Temporal
- sandbox lifecycle management
- model fallback

Target:

> Agent survives worker crashes and transient provider failures.

## Phase 3 --- Team SaaS

Deliver:

- organizations
- members
- roles
- policies
- usage limits
- analytics
- team skills
- audit
- shared projects

## Phase 4 --- Enterprise

Deliver:

- SAML
- SCIM
- BYOK
- VPC
- data residency
- advanced audit
- retention controls
- customer-managed keys
- enterprise SLA

## Phase 5 --- Global scale

Deliver:

- regional execution
- global routing
- active/active control plane where justified
- regional data policies
- global provider routing
- disaster recovery automation

## Phase 6 --- Agent platform

Deliver:

- public SDK
- agent protocol
- plugin marketplace
- skill marketplace
- third-party agents
- agent evaluation marketplace
- automation workflows

## Phase 7 --- Advanced intelligence

Deliver:

- repository knowledge graph
- code ownership reasoning
- dependency intelligence
- architecture reasoning
- multi-agent workflows
- autonomous issue triage
- CI repair
- release engineering

## Recommended development order

```text
Reliability
    >
Security
    >
Execution
    >
Developer UX
    >
Team features
    >
Marketplace
    >
Advanced autonomy
```

Do not reverse this order.

## Team structure

Initial:

- 1 product/technical founder
- 1 backend/platform engineer
- 1 frontend engineer
- 1 agent/runtime engineer
- 1 DevOps/security engineer

Later:

- agent research/evaluation
- infra/SRE
- security
- enterprise
- growth/product

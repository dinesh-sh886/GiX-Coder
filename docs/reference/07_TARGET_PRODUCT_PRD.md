# Proposed Product --- AI Coding Agent SaaS

## Working concept

Build a global SaaS platform for delegating software-engineering work to
AI agents.

Working category:

> **AI Engineering Execution Platform**

The product should serve:

- individual developers
- startups
- engineering teams
- agencies
- enterprises
- software outsourcing companies
- DevOps/SRE teams
- security teams

## 1. Product promise

A user should be able to say:

> "Implement issue #184, add tests, run the relevant test suite, open a
> PR, and explain the changes."

The system should:

1.  understand the repository
2.  create an execution plan
3.  retrieve relevant context
4.  modify code
5.  run tests
6.  recover from failures
7.  request approval for risky actions
8.  create a commit/branch
9.  open a PR
10. provide evidence
11. remain resumable

## 2. Product surfaces

### CLI

For power users.

### IDE

VS Code first; JetBrains later.

### Web application

For:

- task queue
- sessions
- PR review
- approvals
- usage
- analytics
- team administration

### Git integration

GitHub first.

Then:

- GitLab
- Bitbucket
- Azure DevOps

### API

Everything important should be API-accessible.

### Webhooks

Examples:

- issue created
- PR opened
- CI failed
- review requested
- deployment failed

## 3. Core objects

```text
Organization
  ├── User
  ├── Project
  ├── Repository
  ├── Workspace
  ├── Agent
  ├── Skill
  ├── Policy
  ├── Task
  ├── Session
  ├── Run
  ├── ToolCall
  ├── Artifact
  ├── Approval
  ├── Evaluation
  └── BillingAccount
```

## 4. Agent modes

### Ask

Read-only investigation.

### Plan

Creates a detailed plan without changing files.

### Build

Can modify workspace.

### Review

Reviews code without modifying.

### Fix

Focused repair of failing tests/CI.

### Autonomous

Can continue through defined workflow limits.

## 5. Key differentiator: outcome-based execution

Every task should end in one of:

- succeeded
- partially succeeded
- blocked
- failed
- cancelled

Success must have evidence.

Example:

```text
Task: Fix authentication bug

Evidence:
✓ Tests: 183 passed
✓ New regression test added
✓ Typecheck passed
✓ Lint passed
✓ Diff reviewed
✓ PR #924 opened

Cost:
$0.83

Duration:
8m 42s
```

## 6. Reliability target

Initial target:

- 99.9% control-plane availability
- 99.5% task-start availability
- 99% successful session recovery
- zero cross-tenant workspace access
- RPO \<= 5 minutes
- RTO \<= 30 minutes for regional failure

Later target:

- 99.95% control plane
- 99.99% critical APIs
- multi-region active/active

## 7. Product principles

1.  Never hide what the agent did.
2.  Never silently exceed user permissions.
3.  Never lose a session.
4.  Never assume model output is correct.
5.  Always preserve evidence.
6.  Make costs predictable.
7.  Make provider switching possible.
8.  Make local execution first-class.
9.  Design for enterprise security from day one.
10. Optimize for shipped software, not generated text.

## 8. MVP

MVP should include:

- authentication
- organizations
- GitHub OAuth/app
- repository connection
- cloud workspace
- CLI
- web task console
- agent loop
- file tools
- shell tool
- Git tools
- test execution
- approvals
- model router
- usage tracking
- billing
- audit log

Do not build initially:

- ten Git providers
- mobile native app
- custom foundation model
- social collaboration
- marketplace with thousands of plugins
- full Kubernetes platform exposed to customers

## 9. North-star metric

> **Verified engineering outcomes per active developer per month**

Supporting metric:

> Percentage of tasks completed without human code edits after agent
> completion.

## 10. Strategic moat

The moat should become:

```text
Execution reliability
        +
Repository intelligence
        +
Agent evaluations
        +
Workflow integrations
        +
Enterprise policy
        +
Outcome data
        +
Cost optimization
```

Not merely prompts.

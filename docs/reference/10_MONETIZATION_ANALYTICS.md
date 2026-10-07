# Monetization, Unit Economics & Analytics

## 1. Monetization strategy

Use a hybrid model.

### Free

Purpose:

- acquisition
- open-source adoption
- developer trust

Limits:

- task minutes
- monthly compute
- model credits
- concurrent agents

### Pro

Target:

- individual professional developers

Features:

- larger usage
- premium models
- cloud workspaces
- background agents
- advanced integrations

### Team

Target:

- 5--100 developers

Features:

- shared workspaces
- policy
- spending limits
- analytics
- team skills
- SSO-lite
- GitHub organization integration

### Business

Target:

- 100--1000 developers

Features:

- SSO/SAML
- SCIM
- advanced policy
- audit
- private networking
- data residency
- custom retention

### Enterprise

Features:

- private VPC
- self-hosted execution
- BYOK
- customer-managed keys
- dedicated support
- custom SLA
- procurement/security controls

## 2. Pricing structure

Avoid pure unlimited pricing early.

Recommended:

```text
Base subscription
+
Included compute/model credits
+
Usage overage
+
Enterprise infrastructure fee
```

This protects margins.

## 3. Cost model

Per task:

```text
LLM cost
+ sandbox CPU
+ sandbox RAM
+ storage
+ network
+ observability
+ orchestration
+ support allocation
```

Track contribution margin at:

- user
- organization
- task
- model
- repository
- agent
- feature

## 4. Model router

Do not always use the strongest model.

Route by task:

```text
Simple search       -> cheap/fast model
Code edit           -> coding model
Architecture        -> high reasoning model
Review              -> independent reviewer model
Recovery            -> specialized repair model
```

The router should optimize:

```text
expected_success / cost / latency
```

## 5. Product analytics

### Acquisition

- website visit
- signup
- CLI install
- repository connection
- first task

### Activation

North-star activation event:

> User completes first verified coding task.

Track:

- time to first task
- time to first success
- first successful PR
- first cloud execution

### Engagement

- weekly active developers
- tasks/developer
- sessions/developer
- repositories connected
- agents per org

### Outcome

- task success rate
- test pass rate
- PR creation rate
- PR merge rate
- human intervention rate
- rework rate
- rollback rate

### Economics

- tokens/task
- compute/task
- gross margin/task
- cost/success
- revenue/developer
- revenue/org

## 6. Agent telemetry

Each run should emit:

```text
run_id
tenant_id
project_id
task_id
model
provider
prompt_tokens
cached_tokens
output_tokens
latency
tool_count
tool_failures
files_read
files_changed
tests_run
tests_passed
tests_failed
human_approvals
sandbox_cpu_seconds
sandbox_memory_mb_seconds
network_bytes
estimated_cost
actual_cost
outcome
```

## 7. Privacy

Analytics should separate:

### Operational telemetry

Required for service reliability.

### Product analytics

Aggregated behavioral data.

### Customer content

Source code, prompts and artifacts.

Do not mix these blindly.

Customers should have controls over content retention.

## 8. Dashboards

### Executive

- ARR/MRR
- active organizations
- net retention
- gross margin
- task success
- infrastructure cost

### Product

- activation
- retention
- tasks/user
- feature adoption
- funnel

### Agent

- success rate
- recovery rate
- tool failure rate
- model quality
- latency

### FinOps

- cost/task
- provider spend
- sandbox spend
- storage
- network
- margin

### Security

- policy violations
- suspicious tool calls
- secret access
- sandbox escapes
- auth anomalies

## 9. Experimentation

A/B test:

- onboarding
- default model
- agent mode
- approval UX
- pricing
- task templates
- skill recommendations

Never A/B test security controls without an explicit safety review.

## 10. Retention loop

A strong loop:

```text
Task success
   |
PR merged
   |
Developer trusts agent
   |
More tasks delegated
   |
More workflow integration
   |
More team adoption
   |
More data/evaluations
   |
Better reliability
   |
Higher trust
```

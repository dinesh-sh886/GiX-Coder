# Architectural Decisions

## ADR-001 --- Build a platform, not a chatbot

**Decision:** The primary abstraction is a task/session/run, not a chat
message.

**Reason:** Coding work is long-running, stateful and tool-driven.

---

## ADR-002 --- Separate control plane and execution plane

**Decision:** API/control services never execute arbitrary customer
code.

**Reason:** Security, scaling and availability.

---

## ADR-003 --- Use durable workflows

**Decision:** Use Temporal or an equivalent durable workflow engine.

**Reason:** Agent tasks can run for minutes/hours and must survive
failures.

---

## ADR-004 --- Event-driven execution

**Decision:** All meaningful agent actions produce immutable events.

**Reason:**

- audit
- replay
- debugging
- analytics
- recovery
- evaluation

---

## ADR-005 --- Multi-model from day one

**Decision:** Define an internal provider/model interface.

```text
ModelProvider
  generate()
  stream()
  embeddings()
  health()
  pricing()
  capabilities()
```

**Reason:** Avoid provider lock-in and enable cost/quality routing.

---

## ADR-006 --- Sandbox everything

**Decision:** Every autonomous coding task runs in an isolated
environment.

**Reason:** Customer repositories are untrusted input.

---

## ADR-007 --- Human approval is policy-driven

**Decision:** Approval is determined by action risk and organization
policy.

**Reason:** Avoid both unsafe autonomy and approval fatigue.

---

## ADR-008 --- PostgreSQL is source of truth

**Decision:** PostgreSQL stores transactional metadata.

**Reason:** Strong consistency, mature tooling and global SaaS
ecosystem.

Redis is not the source of truth.

---

## ADR-009 --- Analytics in ClickHouse

**Decision:** Separate operational DB from analytical event store.

**Reason:** High-volume agent telemetry should not overload PostgreSQL.

---

## ADR-010 --- Modular monolith before microservices

**Decision:** Start with strong module boundaries inside one deployable
control-plane application.

**Reason:** Small teams cannot afford premature distributed-systems
complexity.

Extract services when:

- independent scaling is required
- ownership is clear
- failure isolation matters
- deployment cadence differs

---

## ADR-011 --- Rust for execution-critical components

**Decision:** Use Rust where process/sandbox/runtime correctness
matters.

**Reason:**

- memory safety
- performance
- predictable binaries
- good systems tooling

Use TypeScript/Go for higher-level product logic.

---

## ADR-012 --- Git is a first-class domain

**Decision:** Branch, commit, diff and PR are domain objects.

**Reason:** The final value of a coding agent is shipped software.

---

## ADR-013 --- Success requires evidence

**Decision:** An agent cannot claim success without configured
verification.

Example:

```text
Task success =
  changes applied
  AND
  required checks passed
  AND
  no unresolved critical errors
```

---

## ADR-014 --- Local-first compatibility

**Decision:** Support local CLI execution even when cloud is the main
SaaS.

**Reason:**

- developer trust
- enterprise/private environments
- offline work
- migration path
- lower cloud costs

---

## ADR-015 --- Enterprise security from the beginning

**Decision:** Tenant isolation, audit, encryption and policy are
foundational components.

**Reason:** Retrofitting security around arbitrary code execution is
expensive and risky.

---

# Final architecture principle

The product should converge toward:

```text
                 Developer
                     |
          CLI / IDE / Web / API
                     |
              Agent Gateway
                     |
             Durable Workflow
                     |
             Agent Harness
                     |
       +-------------+-------------+
       |             |             |
   Context        Policy        Model Router
       |             |             |
       +-------------+-------------+
                     |
             Isolated Sandbox
                     |
       +------+------+------+------+
       |      |      |      |      |
     Files  Shell   Git   Tests   MCP
                     |
                  GitHub
```

This is the architecture to optimize around.

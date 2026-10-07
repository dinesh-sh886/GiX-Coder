# ADR-0001: Control Plane / Data Plane Separation

## Status

Accepted

## Context

GiX-Coder must execute customer-controlled code (workflows, agents, tools) while maintaining a secure, auditable control plane for orchestration, authentication, authorization, and policy enforcement.

Key forces:

- Customer code is untrusted input
- Control plane must never execute customer code
- Independent scaling of control vs execution
- Different security postures required
- Audit trail must cover both planes
- Failure domains must be isolated

## Decision

Strictly separate the system into two planes:

### Control Plane (Trusted)

- Agent Gateway (API, auth, routing, rate limiting)
- Workflow Engine (durable execution orchestration)
- Agent Harness (agent lifecycle, capability management)
- Router (context, policy, model routing)
- Policy Engine (authorization, quotas)
- Audit Service (immutable logging)
- Shared Kernel (utilities, DTOs, errors)

### Data Plane (Untrusted)

- Execution Sandboxes (gVisor/Firecracker isolated)
- Tool Adapters (file, shell, git, test, MCP)
- Network Policy Enforcement
- Resource Limit Enforcement

### Boundaries

- All control plane → data plane communication via gRPC with mTLS
- Data plane has NO direct access to control plane services
- Capability grants are explicit, scoped, time-limited, audited
- No shared databases, no shared memory, no shared filesystem
- Control plane services are stateless (except audit)

## Consequences

### Positive

- Strong security boundary: customer code never in control plane
- Independent scaling: sandboxes scale horizontally, control plane scales differently
- Clear audit trail: every data plane action initiated via control plane
- Technology flexibility: different runtimes, languages, isolation per plane
- Failure isolation: sandbox crash doesn't affect control plane

### Negative

- Increased operational complexity: two plane deployment
- Latency overhead: gRPC calls between planes
- State synchronization: checkpointing requires coordination
- Development complexity: testing cross-plane interactions

### Risks

- **Sandbox escape**: Mitigated by gVisor/Firecracker, seccomp, capability dropping
- **Capability grant abuse**: Mitigated by policy engine, audit, time limits
- **Network pivot**: Mitigated by egress allowlist, no ingress to sandbox
- **Control plane compromise**: Mitigated by mTLS, least privilege, audit

## Alternatives Considered

### Alternative 1: Single Process with Sandbox Threads

- **Pros**: Simpler deployment, lower latency
- **Cons**: No security boundary, sandbox escape = total compromise
- **Why rejected**: Violates secure-by-default principle

### Alternative 2: Separate Microservices per Function

- **Pros**: Maximum isolation
- **Cons**: Premature complexity, distributed system failures, operational burden
- **Why rejected**: ADR-0002 chose modular monolith first

### Alternative 3: WebAssembly Sandbox in Process

- **Pros**: Lightweight, fast startup
- **Cons**: Limited syscall support, tool ecosystem immature, no shell/git
- **Why rejected**: Doesn't support required tooling (shell, git, MCP)

## Implementation Plan

- [ ] Define gRPC service contracts for plane communication
- [ ] Implement mTLS with SPIFFE/SPIRE
- [ ] Design capability grant protocol
- [ ] Build sandbox executor with gVisor
- [ ] Implement audit logging for all grants/executions
- [ ] Integration tests for cross-plane scenarios

## Related

- ADR-0002: Modular Monolith
- ADR-0003: Durable Workflow Engine
- ADR-0004: Sandbox Isolation

## Metadata

- **Author**: Architect
- **Date**: 2024-01-15
- **Reviewers**: Security Lead, Platform Lead
- **Approved By**: Architect (Principal Architect)

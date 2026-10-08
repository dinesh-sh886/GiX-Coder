# ADR-0004: Sandbox Isolation

## Status

Accepted

## Context

GiX-Coder executes untrusted customer code via AI agents. The execution sandbox must provide:

- Strong isolation (no host access, no cross-tenant access)
- Capability-based access control (explicit grants)
- Resource limits (CPU, memory, time, I/O, pids)
- Deterministic snapshot/restore for checkpointing
- Audit trail of all operations
- Support for required tools: filesystem, shell, git, test, MCP

Key forces:

- Customer code is untrusted
- Supply chain attacks via dependencies
- Prompt injection leading to command execution
- Sandbox escape is critical severity
- Performance: startup < 2s warm, < 10s cold
- Tool compatibility: standard Linux tools

## Decision

Use **gVisor (primary) / Firecracker (fallback)** for sandbox isolation with capability-based access control.

### Isolation Layers

```
1. gVisor / Firecracker (kernel boundary)
2. seccomp profile (syscall filtering)
3. Linux capabilities (drop all, add minimal)
4. User namespace (rootless)
5. Filesystem (read-only root, tmpfs overlays)
6. Network (none by default, CNI policy)
7. Resources (cgroups v2: CPU, memory, pids, I/O)
8. Time limits (wall-clock, CPU)
9. Audit (all syscalls logged)
```

### Capability Model

Capabilities are explicit grants:

- **Filesystem**: read/write/list - path-scoped
- **Shell**: command allowlist - timeout, resource limits
- **Git**: repo-scoped, branch-scoped operations
- **Test**: framework-agnostic execution
- **Network**: egress allowlist only
- **MCP**: tool-scoped, approved registry

### Runtime Selection

- **gVisor (runsc)**: Primary - faster startup, good compatibility, user-space kernel
- **Firecracker (firecracker-containerd)**: Fallback - stronger isolation (microVM), slower startup

Selection via feature flag, configurable per tenant/workflow.

## Consequences

### Positive

- Strong security: defense in depth, multiple isolation layers
- Proven technology: gVisor (Google), Firecracker (AWS Lambda)
- Checkpointing: gVisor supports checkpoint/restore
- Tool compatibility: standard Linux userspace
- Resource control: cgroups v2 native
- Audit: seccomp notify, gVisor logging

### Negative

- Performance overhead: gVisor ~5-15% syscall overhead
- Memory overhead: gVisor ~50-100MB base
- Startup latency: gVisor ~1-2s, Firecracker ~100-200ms (but higher base)
- Complexity: multiple runtimes, seccomp profiles
- Debugging: harder than native containers

### Risks

- **gVisor escape**: Mitigated by Firecracker fallback, seccomp, capabilities
- **Firecracker escape**: Mitigated by KVM isolation, minimal device model
- **Capability grant bugs**: Mitigated by policy engine, audit, testing
- **Resource exhaustion**: Mitigated by cgroups, timeouts, monitoring
- **Tool incompatibility**: Mitigated by testing, fallback to Firecracker

## Alternatives Considered

### Alternative 1: Native Containers (runc/containerd)

- **Pros**: Best performance, simplest, full compatibility
- **Cons**: Weak isolation, shared kernel, escape = host compromise
- **Why rejected**: Insufficient for untrusted code

### Alternative 2: Kata Containers

- **Pros**: Lightweight VM, good isolation
- **Cons**: Heavier than gVisor, slower startup, less mature checkpointing
- **Why rejected**: gVisor better fit for our latency requirements

### Alternative 3: WebAssembly (Wasmtime/Wasmer)

- **Pros**: Fast startup, strong isolation, language agnostic
- **Cons**: No shell/git support, limited syscalls, immature tooling
- **Why rejected**: Doesn't support required tools

### Alternative 4: Custom seccomp + namespace + cgroups

- **Pros**: Full control, minimal overhead
- **Cons**: Reimplementing gVisor, high risk of gaps
- **Why rejected**: Not core competency, gVisor is battle-tested

## Implementation Plan

- [ ] gVisor integration (runsc runtime)
- [ ] seccomp profile per capability
- [ ] Capability grant protocol
- [ ] Filesystem tool with path scoping
- [ ] Shell tool with allowlist
- [ ] Git tool with repo scoping
- [ ] Test tool execution
- [ ] Network policy enforcement
- [ ] cgroups v2 resource limits
- [ ] Checkpoint/restore for workflows
- [ ] Firecracker fallback implementation
- [ ] Security testing (escape attempts, fuzzing)

## Related

- ADR-0001: Control Plane / Data Plane Separation
- ADR-0002: Modular Monolith
- ADR-0003: Durable Workflow Engine

## Metadata

- **Author**: Security Lead
- **Date**: 2024-01-17
- **Reviewers**: Architect, Platform Lead
- **Approved By**: Architect (Principal Architect)

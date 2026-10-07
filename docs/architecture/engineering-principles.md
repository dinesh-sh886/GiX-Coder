# GiX-Coder Engineering Principles

## Core Principles

### SOLID
- **Single Responsibility**: Each module, class, function has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Subtypes must be substitutable for base types
- **Interface Segregation**: Many specific interfaces over one general interface
- **Dependency Inversion**: Depend on abstractions, not concretions

### DRY (Don't Repeat Yourself)
- Extract common logic into shared utilities
- Prefer composition over copy-paste
- Generate code from single source of truth (API definitions, config schemas)

### KISS (Keep It Simple, Stupid) - Where Appropriate
- Simple solutions preferred
- Complexity must be justified by requirements
- No premature abstraction

### Separation of Concerns
- Distinct modules for distinct responsibilities
- Clear boundaries between control plane and data plane
- Business logic separate from infrastructure

### Dependency Inversion
- High-level modules don't depend on low-level modules
- Both depend on abstractions (interfaces)
- Dependency injection at composition root

### Clean Architecture
- Entities (enterprise business rules)
- Use Cases (application business rules)
- Interface Adapters (controllers, gateways, presenters)
- Frameworks & Drivers (DB, web, UI)
- Dependencies point inward only

### Explicit Domain Boundaries
- Modules own their data
- Cross-module communication via explicit contracts (DTOs, events)
- No shared databases between modules

### Secure by Default
- Deny by default, explicit allow
- Principle of least privilege
- Defense in depth
- Zero trust between execution boundaries

### Configuration Over Hard-Coded Values
- No credentials, URLs, ports, model names in code
- Environment-specific config via layered configuration
- Immutable configuration where possible

### DTOs at External Boundaries
- All API boundaries use Data Transfer Objects
- No domain entities cross module boundaries
- Validation at boundaries

### Validation at Boundaries
- Input validation at every entry point
- Output validation for external consumers
- Fail fast, fail loud

### Typed Error Handling
- Errors as values, not exceptions
- Structured error types with codes
- Error wrapping with context
- No panic in production code

### Structured Logging
- JSON format mandatory
- Correlation IDs on all requests
- Structured fields, not string interpolation
- Log levels: DEBUG, INFO, WARN, ERROR

### Deterministic Builds
- Reproducible builds from source
- Locked dependencies
- Fixed toolchain versions
- No timestamps in build artifacts

### Reproducible Environments
- Infrastructure as Code
- Containerized development
- Identical CI/CD and local environments

### Testability
- Code designed for testing
- Interfaces for all external dependencies
- Fast unit tests, realistic integration tests
- Test doubles for external systems

### Observability
- Logs, metrics, traces by default
- Health and readiness endpoints
- Business and technical metrics
- Audit trails for security events

### Least Privilege
- Minimum permissions for each component
- Capability-based access control
- Time-limited credentials where possible

### Zero Trust Between Execution Boundaries
- Every service call authenticated
- Every data access authorized
- Every execution sandboxed
- No implicit trust

## Anti-Patterns (Explicitly Forbidden)

| Anti-Pattern | Forbidden Because |
|--------------|-------------------|
| Premature microservices | Operational complexity, distributed system failures |
| God modules/classes | Violates SRP, untestable, unmaintainable |
| Shared mutable state | Concurrency bugs, testing difficulty |
| Implementation leakage | Tight coupling, blocks evolution |
| Magic strings/numbers | Untraceable, unsearchable, error-prone |
| Silent failures | Undetectable bugs, data corruption |
| Configuration in code | Environment coupling, secret leakage |
| Exception-based control flow | Performance, debuggability, clarity |
| Untyped error returns | Unhandled errors, unclear contracts |
| Direct database access across modules | Schema coupling, transaction boundary violations |
| Hard-coded timeouts/retries | Environment mismatch, inflexibility |
| Logging without context | Undebuggable production issues |

## Decision Making Framework

When principles conflict:
1. **Security** > Everything
2. **Correctness** > Performance
3. **Clarity** > Cleverness
4. **Explicit** > Implicit
5. **Standards** > Preferences

Document tradeoffs in ADRs.
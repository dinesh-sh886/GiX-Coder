# GiX-Coder

Global AI Engineering Execution Platform.

## Overview

GiX-Coder is a production-grade, secure, and auditable platform for AI-assisted software engineering. It provides a governed execution environment where AI agents can perform complex engineering tasks with full observability, security, and human oversight.

## Architecture

```
Developer
    ↓
CLI / IDE / Web / API
    ↓
Agent Gateway (Control Plane)
    ↓
Durable Workflow Engine
    ↓
Agent Harness
    ↓
Context + Policy + Model Router
    ↓
Isolated Execution Sandbox (Data Plane)
    ↓
Files / Shell / Git / Tests / MCP
    ↓
GitHub / GitLab / etc.
```

## Current Phase

**Phase 00**: Planning & Engineering Governance ✓

See [Phase 00 Documentation](docs/phases/phase-00/README.md)

## Documentation

- [Product Charter](docs/product/charter.md)
- [Architecture Baseline](docs/architecture/baseline.md)
- [Repository Structure](docs/architecture/repository-structure.md)
- [Engineering Principles](docs/architecture/engineering-principles.md)
- [Coding Standards](docs/architecture/coding-standards.md)
- [Git Strategy](docs/architecture/git-strategy.md)
- [PR Strategy](docs/architecture/pr-strategy.md)
- [CI/CD Strategy](docs/architecture/ci-cd-strategy.md)
- [Environment Strategy](docs/architecture/environment-strategy.md)
- [Security Baseline](docs/security/baseline.md)
- [Configuration Strategy](docs/architecture/configuration-strategy.md)
- [Observability Strategy](docs/architecture/observability-strategy.md)
- [Testing Strategy](docs/qa/testing-strategy.md)
- [Quality Gates](docs/qa/quality-gates.md)
- [ADR Framework](docs/adr/framework.md)
- [Documentation Standards](docs/architecture/documentation-standards.md)
- [OpenCode Rules](docs/architecture/opencode-rules.md)
- [Nemotron Rules](docs/architecture/nemotron-rules.md)
- [Phase Lifecycle](docs/architecture/phase-lifecycle.md)
- [Phase 01 Specification](docs/phases/phase-01/specification.md)

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+ / TypeScript 5+
- Docker / Docker Compose
- Kubernetes (kind/minikube for local)

### Local Development

```bash
# Clone
git clone https://github.com/dinesh-sh886/GiX-Coder.git
cd GiX-Coder

# Start local stack
docker-compose -f deploy/docker-compose/dev.yml up -d

# Run tests
make test-all
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## Security

See [SECURITY.md](SECURITY.md) for reporting vulnerabilities.

## License

Apache License 2.0 - see [LICENSE](LICENSE)

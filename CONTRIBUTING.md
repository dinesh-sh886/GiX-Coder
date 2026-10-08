# Contributing to GiX-Coder

Thank you for contributing to GiX-Coder! This document outlines the process and standards for contributions.

## Code of Conduct

This project follows the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/1/code_of_conduct/). By participating, you agree to uphold this code.

## Ways to Contribute

- Report bugs
- Suggest features
- Improve documentation
- Submit code changes
- Review pull requests
- Share feedback

## Development Workflow

### 1. Setup

```bash
# Fork and clone
git clone https://github.com/<your-fork>/GiX-Coder.git
cd GiX-Coder

# Install tools
make install-tools

# Verify setup
make check
```

### 2. Create Branch

```bash
# Branch from develop
git checkout develop
git pull origin develop
git checkout -b feature/GIX-123-short-description
```

### 3. Make Changes

- Follow [coding standards](docs/architecture/coding-standards.md)
- Write tests for new functionality
- Update documentation
- Run local quality gates: `make check`

### 4. Commit

```bash
# Stage changes
git add .

# Commit with conventional message
git commit -m "feat(gateway): add tenant-scoped rate limiting

Implement token bucket algorithm per tenant to prevent
noisy neighbor problems. Configuration via policy engine.

Closes #123"
```

### 5. Push & Create PR

```bash
git push origin feature/GIX-123-short-description
# Create PR via GitHub UI
```

## Pull Request Requirements

### Before Submitting

- [ ] All tests pass (`make test`)
- [ ] Coverage >= 80% (`make test-coverage`)
- [ ] Lint passes (`make lint`)
- [ ] No security findings (`make security-scan`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (if user-facing)

### PR Template

Use the [PR template](.github/pull_request_template.md) - it will be auto-populated.

### Review Process

1. Automated checks run (CI)
2. Code owner review required
3. Minimum 1 approval for `develop`
4. Minimum 2 approvals for `main` (including code owner)
5. All conversations resolved
6. Squash merge

## Coding Standards

### Languages

- **Go** (primary): See [Go standards](docs/architecture/coding-standards.md#go-primary---control-plane-services)
- **TypeScript** (secondary): See [TS standards](docs/architecture/coding-standards.md#typescript-secondary---cli-web-tooling)

### General

- Follow [engineering principles](docs/architecture/engineering-principles.md)
- SOLID, DRY, KISS, Clean Architecture
- Security first, observability by default
- Configuration over hard-coding

## Testing

### Required

- Unit tests for all new logic (target 80%+ coverage)
- Integration tests for cross-module functionality
- Contract tests for API changes

### Run Locally

```bash
# Unit tests
make test-unit

# Integration tests (requires docker)
make test-integration

# All tests
make test-all
```

## Documentation

### Update When

- Adding/changing public APIs
- Changing architecture
- Adding configuration options
- Changing behavior

### Where

- Code comments (exported symbols)
- Module README.md
- Architecture docs (if architectural)
- ADR (if significant decision)

## Security

### Reporting Vulnerabilities

See [SECURITY.md](SECURITY.md)

### Secure Coding

- Never commit secrets
- Validate all inputs
- Use parameterized queries
- Follow [security baseline](docs/security/baseline.md)

## Architecture Decisions

### When to Create ADR

- New module/service
- Technology selection
- Security model changes
- Breaking API changes
- Quality gate exceptions

### Process

1. Create ADR in `docs/adr/NNNN-description.md`
2. Follow [ADR framework](docs/adr/framework.md)
3. Submit PR with `adr` label
4. Architect + domain owner review
5. Merge to `develop`

## Release Process

### Versioning

Semantic Versioning (MAJOR.MINOR.PATCH)

### Branches

- `develop` → Next release development
- `stage` → Release candidate
- `main` → Production releases

### Release

1. Create `stage` branch from `develop`
2. QA validation on STAGE
3. Production approval
4. Merge `stage` → `main` (tagged)
5. Auto-deploy to PROD

## Getting Help

- **Questions**: GitHub Discussions
- **Bugs**: GitHub Issues (use template)
- **Security**: See SECURITY.md
- **Architecture**: Check ADRs, ask in PR

## Recognition

Contributors are recognized in:

- Release notes
- CONTRIBUTORS.md
- GitHub contributor graphs

Thank you for making GiX-Coder better!

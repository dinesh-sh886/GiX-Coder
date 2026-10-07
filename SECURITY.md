# Security Policy

## Supported Versions

| Version | Supported |
|---------|-----------|
| Main branch (latest) | ✓ |
| Stage branch (RC) | ✓ |
| Develop branch | ✓ (best effort) |

## Reporting a Vulnerability

### Do NOT
- Open public GitHub issues for security vulnerabilities
- Discuss vulnerabilities in public forums
- Exploit vulnerabilities to demonstrate impact

### DO
Report via **private channels**:

1. **GitHub Security Advisories** (preferred)
   - Go to Security tab → Report a vulnerability
   - Private disclosure to maintainers

2. **Email**
   - security@gix-coder.io (when available)
   - Include: description, impact, reproduction steps, suggested fix

3. **Encrypted Communication**
   - PGP key available on request
   - Signal: on request

### What to Include
- Vulnerability type (e.g., RCE, injection, auth bypass)
- Affected components
- Reproduction steps (minimal)
- Potential impact
- Suggested mitigation/fix
- Your contact for follow-up

## Response Timeline

| Severity | Acknowledgment | Assessment | Fix Target |
|----------|----------------|------------|------------|
| Critical | 4 hours | 24 hours | 72 hours |
| High | 8 hours | 48 hours | 7 days |
| Medium | 24 hours | 1 week | 30 days |
| Low | 48 hours | 2 weeks | 90 days |

## Disclosure Process

1. **Acknowledge** receipt within SLA
2. **Assess** impact and scope
3. **Develop** fix in private fork
4. **Test** fix thoroughly
5. **Coordinate** release with reporter
6. **Publish** advisory and fix simultaneously
7. **Credit** reporter (if desired)

## Bug Bounty

### Scope
- Production API endpoints
- Sandbox escape
- Authentication/authorization bypass
- Data exposure
- RCE in control plane

### Out of Scope
- Denial of service (volumetric)
- Social engineering
- Physical attacks
- Third-party services
- Non-production environments

### Rewards (Indicative)
| Severity | Reward |
|----------|--------|
| Critical | $10,000+ |
| High | $5,000 |
| Medium | $1,000 |
| Low | $250 |

### Rules
- First reporter gets reward
- No automated scanning without permission
- No testing on production data
- Follow responsible disclosure

## Security Architecture

### Key Controls
- Control plane / data plane separation
- Capability-based sandbox (gVisor/Firecracker)
- Zero-trust networking (mTLS, egress allowlist)
- Policy-based authorization (OPA/Cedar)
- Immutable audit logging
- Supply chain security (SLSA, sigstore)

### Threat Model
The threat model is documented in the Security Baseline (`docs/security/baseline.md`, "Attack Vectors & Mitigations" section). A standalone threat model document is planned for Phase 01+ (`docs/security/threat-model.md` — not yet created).

### Compliance Targets
- SOC 2 Type II
- ISO 27001
- NIST CSF

## Secure Development

### Requirements
- All code reviewed (security focus)
- SAST on every PR (CodeQL)
- SCA on every PR (govulncheck, osv-scanner)
- Secret scan on every commit (TruffleHog)
- Container scan on every build (Trivy)
- DAST on STAGE deploy

### Training
- Secure coding (annual)
- Threat modeling (architects)
- Incident response (on-call)

## Contact

- **Security Team**: security@gix-coder.io
- **Emergency**: security-emergency@gix-coder.io
- **PGP**: Available on request

## Acknowledgments

We thank all security researchers who responsibly disclose vulnerabilities.
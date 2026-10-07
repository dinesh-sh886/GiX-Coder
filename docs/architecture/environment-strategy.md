# GiX-Coder Environment Strategy

## Environment Hierarchy

```
┌─────────────────────────────────────────────────────────────┐
│                        LOCAL                                 │
│  Developer machine (containerized, matches DEV config)      │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                         DEV                                  │
│  Rapid integration, developer testing, non-production data  │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                        STAGE                                 │
│  Production-like, release validation, acceptance testing    │
└─────────────────────────────────────────────────────────────┘
                              ↓
┌─────────────────────────────────────────────────────────────┐
│                        PROD                                  │
│  Customer traffic, protected, monitored, rollback capable   │
└─────────────────────────────────────────────────────────────┘
```

## Environment Specifications

### LOCAL
| Aspect | Specification |
|--------|---------------|
| Purpose | Inner-loop development, debugging |
| Infrastructure | Docker Compose / Kind / LocalStack |
| Data | Synthetic, developer-controlled |
| Secrets | `.env.local` (gitignored), 1Password CLI |
| Deploy | `make dev-up` / `tilt up` |
| Reset | On demand |
| Access | Single developer |
| Cost | Minimal (local resources) |

**Services Running Locally:**
- All control plane services (gateway, workflow, harness, router, policy, audit)
- Mocked data plane (sandbox) or local gVisor
- Mocked external dependencies (GitHub API, model providers)
- Observability stack (Prometheus, Grafana, Loki, Tempo, Jaeger)

### DEV
| Aspect | Specification |
|--------|---------------|
| Purpose | Integration, cross-team testing, CI validation |
| Infrastructure | Kubernetes namespace (shared cluster) |
| Data | Synthetic, reset daily via cron |
| Secrets | Vault dev mode / AWS Secrets Manager (dev path) |
| Deploy | Auto on `develop` merge |
| Reset | Daily 02:00 UTC |
| Access | All engineers (read), team members (write) |
| Cost | Shared, optimized |

**Differences from PROD:**
- Single replica per service (no HA)
- Reduced resource limits
- Relaxed rate limits
- Verbose logging (DEBUG)
- Feature flags: all enabled
- Mock external APIs where appropriate

### STAGE
| Aspect | Specification |
|--------|---------------|
| Purpose | Release candidate validation, acceptance testing, performance baseline |
| Infrastructure | Kubernetes (dedicated cluster or namespace with PROD parity) |
| Data | Production-like volume, anonymized/referential integrity |
| Secrets | Vault / AWS Secrets Manager (stage path), same rotation as PROD |
| Deploy | Manual approval on `stage` branch |
| Reset | On demand (before major release) |
| Access | Engineers, QA, Product, Security (read) |
| Cost | ~50% PROD |

**Parity with PROD:**
- Same Kubernetes manifests (different values)
- Same resource limits/requests
- Same replica counts (min 2 for HA)
- Same network policies
- Same monitoring/alerting
- Same security policies
- Real external dependencies (staging endpoints)

### PROD
| Aspect | Specification |
|--------|---------------|
| Purpose | Customer traffic, revenue-generating |
| Infrastructure | Kubernetes (multi-AZ, dedicated clusters) |
| Data | Customer data, encrypted at rest/in transit |
| Secrets | Vault / AWS Secrets Manager (prod path), HSM-backed |
| Deploy | Manual approval + canary, tagged release only |
| Reset | Never (disaster recovery only) |
| Access | On-call, Release managers (break-glass) |
| Cost | Optimized for reliability |

**Requirements:**
- Multi-AZ deployment (min 3 AZs)
- Auto-scaling (HPA + VPA + Cluster Autoscaler)
- Pod disruption budgets
- Network policies (deny by default)
- Runtime security (Falco/Sysdig)
- Backup + point-in-time recovery (RPO < 1 min, RTO < 10 min)
- Chaos engineering (monthly)

## Configuration Per Environment

### Layered Configuration
```
Priority (highest wins):
1. Defaults (code)
2. Config file (configs/<env>/service.yaml)
3. Environment variables
4. Secret manager (runtime)
5. Feature flags (runtime)
```

### Environment-Specific Values
| Config | DEV | STAGE | PROD |
|--------|-----|-------|------|
| Replicas | 1 | 2+ | 3+ (HPA) |
| Resources | 100m/128Mi | PROD values | PROD values |
| Log Level | DEBUG | INFO | INFO |
| Rate Limit | 1000/min | PROD values | PROD values |
| Feature Flags | All on | Release flags | Release flags |
| External APIs | Mock/Staging | Staging | Production |
| TLS | Self-signed | Valid cert | Valid cert |
| Debug Endpoints | Enabled | Disabled | Disabled |

## Data Strategy

### Data Classification
| Classification | DEV | STAGE | PROD |
|----------------|-----|-------|------|
| Customer PII | ❌ Never | Anonymized | Real (encrypted) |
| Credentials | ❌ Never | Test credentials | Real (vault) |
| Analytics | Synthetic | Sampled | Full |
| Audit Logs | 7 days | 30 days | 7 years |

### Data Rules
- **Never** copy PROD data to DEV/STAGE
- **Never** use PROD credentials locally
- **Never** commit secrets
- Anonymization pipeline for STAGE data refresh
- Synthetic data generators for DEV

## Access Control

### DEV
- SSO + GitHub Teams
- Namespace-scoped RBAC
- Self-service via GitOps (ArgoCD/Flux)

### STAGE
- SSO + GitHub Teams
- Environment-scoped RBAC
- Approval required for config changes

### PROD
- Break-glass access (time-limited, audited)
- No standing read access to customer data
- Deploy via GitOps only (no kubectl)
- All access logged and alerted

## Network Topology

### DEV
```
Internet → Ingress (shared) → Services (ClusterIP)
                    ↓
            Mock External APIs
```

### STAGE
```
Internet → WAF → Ingress (dedicated) → Services
                    ↓
            Staging External APIs
                    ↓
            VPC Peering → Shared Services (DB, Cache)
```

### PROD
```
Internet → WAF → Global LB → Regional Ingress → Services
                              ↓
                    Private Link → External APIs
                              ↓
                    Dedicated VPC → Data Layer
```

## Monitoring Per Environment

### DEV
- Basic health checks
- Log aggregation (7 days)
- No alerting (noise)

### STAGE
- Full monitoring stack
- Alerting to team channels
- Synthetic transactions
- Performance baselines

### PROD
- Full monitoring + alerting
- PagerDuty integration
- SLO-based alerting
- Real-user monitoring
- Business metrics
- Audit log alerting

## Disaster Recovery

### DEV
- Recreate from GitOps (RTO < 30 min)
- No backup needed

### STAGE
- Daily etcd backup
- Restore tested monthly (RTO < 2 hr)

### PROD
- Continuous etcd backup (Velero)
- Cross-region DR (RPO < 1 min, RTO < 10 min)
- Quarterly DR drill
- Runbook documented and tested

## Cost Management

| Environment | Strategy |
|-------------|----------|
| LOCAL | Free (local) |
| DEV | Shared cluster, auto-scale to zero nights/weekends |
| STAGE | Right-sized, auto-scale, 50% PROD |
| PROD | Reserved instances, savings plans, right-sizing reviews monthly |

## Environment Promotion Gates

### DEV → STAGE
- [ ] All CI gates pass on `develop`
- [ ] Release branch cut (`stage` from `develop`)
- [ ] Release notes drafted
- [ ] Security scan clean
- [ ] Performance baseline established

### STAGE → PROD
- [ ] STAGE smoke tests pass
- [ ] STAGE acceptance tests pass
- [ ] Performance within 10% baseline
- [ ] Security review complete
- [ ] Production approval (human)
- [ ] Rollback plan documented
- [ ] On-call notified

---

## References

- [CI/CD Strategy](ci-cd-strategy.md)
- [Configuration Strategy](configuration-strategy.md)
- [Security Baseline](../security/baseline.md)
- ADR-0001: Control Plane / Data Plane Separation
- ADR-0004: Sandbox Isolation

---

## Metadata

---
title: GiX-Coder Environment Strategy
type: architecture
phase: 00
status: Verified
author: Platform Lead
date: 2024-10-07
reviewers: Platform Lead, Security Lead, Architect
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0001, ADR-0004]
related_issues: []
---
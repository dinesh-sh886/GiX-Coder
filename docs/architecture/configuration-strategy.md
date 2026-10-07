# GiX-Coder Configuration Strategy

## Configuration Hierarchy

### Priority Order (Highest Wins)
```
1. Defaults (compiled into binary)
2. Configuration files (versioned, per-environment)
3. Environment variables (runtime override)
4. Secret manager (runtime, sensitive only)
5. Feature flags (runtime, behavioral)
```

### Layer Details

#### Layer 1: Defaults (Code)
```go
// Internal defaults - never secrets, never environment-specific
const (
    DefaultHTTPPort     = 8080
    DefaultGRPCPort     = 9090
    DefaultLogLevel     = "info"
    DefaultTimeout      = 30 * time.Second
    DefaultMaxRetries   = 3
)
```
- Committed to source
- Safe for all environments
- Documented in code

#### Layer 2: Config Files (Versioned)
```
configs/
├── dev/
│   ├── gateway.yaml
│   ├── workflow.yaml
│   └── ...
├── stage/
│   └── (same structure)
└── prod/
    └── (structure only - no secrets)
```
- Committed to source (dev, stage)
- Prod structure only (values from secret manager)
- YAML format
- Validated at startup (schema)

#### Layer 3: Environment Variables
```bash
# Override any config file value
GIX_GATEWAY_HTTP_PORT=8080
GIX_GATEWAY_LOG_LEVEL=debug
GIX_WORKFLOW_TEMPORAL_ADDRESS=temporal:7233
```
- Prefix: `GIX_<SERVICE>_<SETTING>`
- Uppercase, underscores
- Overrides config file
- Used for container orchestration

#### Layer 4: Secret Manager (Runtime)
```
vault/
├── dev/
│   ├── gateway/
│   │   ├── github-token
│   │   └── jwt-signing-key
│   └── ...
├── stage/
└── prod/
```
- Only for secrets (keys, tokens, passwords)
- Never in config files, never in env vars (except dev)
- Injected at runtime (Vault Agent / CSI driver)
- Automatic rotation
- Lease-based access

#### Layer 5: Feature Flags (Runtime)
```
flags/
├── gateway/
│   ├── new-routing-algorithm: false
│   └── enhanced-auth: true
└── ...
```
- Behavioral toggles only
- Not for configuration values
- Targeted rollout (tenant, percentage, user)
- Audit log on change

## Configuration Schema

### Service Config Structure
```yaml
# gateway.yaml
service:
  name: "gateway"
  version: "1.0.0"
  environment: "dev"  # dev, stage, prod

server:
  http:
    port: 8080
    host: "0.0.0.0"
    read_timeout: "30s"
    write_timeout: "30s"
    idle_timeout: "120s"
  grpc:
    port: 9090
    host: "0.0.0.0"

logging:
  level: "info"  # debug, info, warn, error
  format: "json"  # json, console
  sampling:
    initial: 100
    thereafter: 100

tracing:
  enabled: true
  sampler: "parentbased_traceidratio"
  ratio: 0.1
  exporter: "otlp"

metrics:
  enabled: true
  path: "/metrics"
  port: 9091

health:
  liveness: "/health/live"
  readiness: "/health/ready"

dependencies:
  temporal:
    address: "temporal:7233"
    namespace: "gix-coder"
  vault:
    address: "https://vault:8200"
    path: "secret/gateway"
  postgres:
    host: "postgres"
    port: 5432
    database: "gateway"
    # credentials from secret manager

features:
  new_routing: false
  enhanced_auth: true

security:
  tls:
    enabled: false  # dev only
    cert_file: ""
    key_file: ""
  cors:
    allowed_origins: ["*"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["*"]
```

### Validation
- JSON Schema for each service config
- Validated at startup (fail fast)
- Schema versioned with service
- CI validates all configs

## Environment-Specific Guidelines

### DEV
- Config files committed with real (non-secret) values
- Secrets in Vault dev mode or `.env.local`
- All feature flags enabled
- Verbose logging
- Relaxed limits

### STAGE
- Config files committed (same structure as PROD)
- Secrets from Vault stage path
- Feature flags match release plan
- Production-like limits
- INFO logging

### PROD
- Config files contain ONLY structure (no values)
- All values from Vault/Secret Manager/Env vars
- Feature flags controlled release
- Strict limits
- INFO logging (DEBUG disabled)

## Secret Handling

### What Goes in Secret Manager
- API keys / tokens
- Database passwords
- TLS certificates / keys
- JWT signing keys
- Encryption keys
- Webhook secrets
- Any credential

### What Does NOT Go in Secret Manager
- Port numbers
- Hostnames (use service discovery)
- Feature flags
- Timeouts
- Log levels
- Resource limits

### Secret Naming
```
<environment>/<service>/<secret-name>
```
Examples:
- `prod/gateway/github-token`
- `stage/workflow/temporal-cert`
- `dev/harness/mcp-api-key`

### Rotation
- Automatic: 90 days (configurable per secret)
- Emergency: Immediate via Vault CLI
- Notification: 7 days before expiry
- Validation: Post-rotation health check

## Feature Flags

### Flag Definition
```yaml
# feature-flags.yaml
flags:
  new_routing_algorithm:
    description: "Use ML-based routing"
    default: false
    rollout:
      strategy: "percentage"
      percentage: 10
      tenant_allowlist: ["tenant-a", "tenant-b"]
    owner: "team-gateway"
    created: "2024-01-15"
```

### Flag Lifecycle
1. **Created** - Default off, documented
2. **Rolled out** - Gradual (10% → 50% → 100%)
3. **Stabilized** - Default on, documented
4. **Removed** - Code cleanup, flag deleted

### Flag Rules
- Boolean only (no complex values)
- Short-lived (< 3 months typical)
- Owner required
- Metrics on flag evaluation
- Emergency kill switch

## Local Development

### .env.example
```bash
# Copy to .env.local and fill in values
# NEVER commit .env.local

# Service ports
GIX_GATEWAY_HTTP_PORT=8080
GIX_GATEWAY_GRPC_PORT=9090

# Dependencies (use docker-compose defaults)
GIX_TEMPORAL_ADDRESS=localhost:7233
GIX_POSTGRES_HOST=localhost
GIX_POSTGRES_PORT=5432
GIX_VAULT_ADDRESS=http://localhost:8200

# Secrets (use 1Password CLI or local vault)
# GIX_GATEWAY_GITHUB_TOKEN=op://vault/gateway/github-token
# GIX_GATEWAY_JWT_KEY=op://vault/gateway/jwt-key

# Feature flags
GIX_GATEWAY_NEW_ROUTING=false
GIX_GATEWAY_ENHANCED_AUTH=true
```

### Local Config Loading
```go
func LoadConfig() (*Config, error) {
    // 1. Load defaults
    cfg := defaultConfig()
    
    // 2. Load config file (configs/dev/gateway.yaml)
    if err := loadConfigFile(cfg); err != nil {
        return nil, err
    }
    
    // 3. Override with env vars
    if err := loadEnvVars(cfg); err != nil {
        return nil, err
    }
    
    // 4. Load secrets (local vault / 1Password)
    if err := loadSecrets(cfg); err != nil {
        return nil, err
    }
    
    // 5. Load feature flags
    if err := loadFeatureFlags(cfg); err != nil {
        return nil, err
    }
    
    // 6. Validate
    if err := validate(cfg); err != nil {
        return nil, err
    }
    
    return cfg, nil
}
```

## Configuration as Code

### GitOps
- All config files in Git
- ArgoCD/Flux syncs to clusters
- PR-based changes
- Audit trail via Git history

### Drift Detection
- Periodic comparison (actual vs desired)
- Alert on drift
- Auto-remediation for non-secrets

### Change Management
- Config changes via PR
- Same review process as code
- Canary deploy for risky changes
- Rollback via Git revert

## Anti-Patterns (Forbidden)

| Anti-Pattern | Why | Alternative |
|--------------|-----|-------------|
| Hard-coded secrets | Leakage, rotation impossible | Secret manager |
| Config in code | Env coupling, no audit | Config files + env vars |
| .env committed | Secret leakage | .env.example only |
| Config per branch | Drift, confusion | Environment-specific files |
| Untyped config | Runtime errors | Schema validation |
| No defaults | Fragile, unclear | Sensible defaults in code |
| Secrets in env vars (prod) | Leakage in logs/process list | Secret manager injection |
| Feature flags as config | No rollout control | Dedicated flag system |

---

## References

- [Environment Strategy](environment-strategy.md)
- [Security Baseline](../security/baseline.md)
- ADR-0001: Control Plane / Data Plane Separation

---

## Metadata

---
title: GiX-Coder Configuration Strategy
type: architecture
phase: 00
status: Verified
author: Architect
date: 2024-10-07
reviewers: Platform Lead, Security Lead, Architect
approved_by: Architect (Principal Architect)
related_adrs: [ADR-0001]
related_issues: []
---
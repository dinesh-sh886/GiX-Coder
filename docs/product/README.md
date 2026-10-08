# GiX-Coder Strategy Pack

This folder contains the current strategic baseline for GiX-Coder.

## Documents

1.  `01-gix-coder-deployment-strategy.md`
    -   AWS/GCP strategy
    -   Docker
    -   ECS/Fargate
    -   EKS/Kubernetes
    -   PostgreSQL/Redis/S3
    -   Agent sandboxes
    -   AI model infrastructure
    -   Secrets
    -   Observability
    -   Billing/metering
    -   Jenkins CI/CD
    -   Production architecture
2.  `02-gix-coder-platform-and-distribution.md`
    -   CLI
    -   Web
    -   VS Code
    -   JetBrains
    -   Windows/macOS/Linux
    -   iOS/Android
    -   Local mode
    -   BYOK
    -   Cloud agents
    -   One GiX account
    -   Release strategy
3.  `03-gix-coder-monetization-and-pricing.md`
    -   Free/Local/Developer/Pro/Team/Enterprise
    -   BYOK
    -   Usage credits
    -   Marketplace
    -   India/global pricing
    -   Enterprise deployment
    -   Revenue model
    -   Long-term monetization strategy

## Strategic Baseline

Core product:

> GiX Agent Platform = product; CLI, Web, IDE, Desktop and Mobile =
> clients.

Deployment:

> Docker everywhere → AWS managed services → ECS/Fargate → EKS when
> agent/sandbox scale requires it → GPU inference when economically
> justified.

Monetization:

> Free local + BYOK + paid cloud agents + team subscriptions +
> enterprise + marketplace.

These documents are a strategic baseline and should be updated as
architecture, pricing, and product decisions evolve.

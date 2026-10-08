# GiX-Coder Monetization and Pricing Strategy

## 1. Business Model

GiX-Coder should use a hybrid SaaS model:

``` text
Subscriptions
+
Usage / Agent Credits
+
BYOK
+
Team Seats
+
Enterprise Contracts
+
Marketplace Revenue
```

Core principle:

> Do not sell the model. Sell the software-engineering platform around
> the model.

GiX-Coder's value comes from:

-   Agent runtime
-   Context management
-   Codebase intelligence
-   Git integration
-   Tools
-   MCP
-   Memory
-   IDE integration
-   Cloud execution
-   Collaboration
-   Security
-   Enterprise controls

------------------------------------------------------------------------

## 2. Recommended Revenue Streams

### A. Individual subscriptions

-   Free
-   Local
-   Developer
-   Pro

### B. Teams

Per-seat pricing.

### C. Enterprise

Custom annual contracts.

### D. Usage

Cloud-agent credits and expensive model usage.

### E. BYOK

Users pay GiX for the software while using their own provider API keys.

### F. Marketplace

GiX takes a percentage of third-party:

-   Agents
-   Skills
-   MCP servers
-   Plugins
-   Templates
-   Integrations

------------------------------------------------------------------------

## 3. Free Plan

### Price

India: ₹0\
Global: \$0

Include:

-   CLI
-   Basic Web
-   VS Code
-   IntelliJ
-   Local models
-   Ollama
-   LM Studio
-   OpenAI-compatible endpoints
-   BYOK
-   GitHub integration
-   Basic agent
-   Basic MCP
-   Basic codebase indexing

Local inference should remain free because the user supplies the
compute.

Purpose:

> User acquisition and ecosystem growth.

------------------------------------------------------------------------

## 4. Local Plan

Suggested:

### India

₹299/month

### Global

\$5/month

Target:

-   Developers who mainly use local models
-   Privacy-conscious users
-   Users with powerful GPUs
-   Offline development

Include:

-   Full local agent
-   Advanced local tools
-   MCP
-   Git automation
-   Memory
-   Codebase indexing
-   Advanced IDE integration

GiX's inference cost can be close to zero because the customer's
hardware performs inference.

------------------------------------------------------------------------

## 5. Developer Plan

Suggested launch price:

### India

₹599/month

### Global

\$15/month

Target:

-   Individual developers
-   Freelancers
-   Junior/senior engineers
-   Students building serious projects

Include:

-   Cloud models
-   Higher agent limits
-   Codebase indexing
-   GitHub/GitLab
-   PR generation
-   Code review
-   Background agents
-   Advanced models
-   Cross-device sessions

------------------------------------------------------------------------

## 6. Pro Plan

Suggested:

### India

₹1,499/month

### Global

\$39/month

Target:

-   Senior developers
-   Architects
-   Founders
-   AI engineers
-   Heavy users

Include:

-   Much higher usage
-   Premium models
-   Long-running agents
-   Multiple simultaneous agents
-   Large-context workloads
-   Cloud sandboxes
-   Advanced Git workflows
-   CI/CD integration
-   Priority execution

------------------------------------------------------------------------

## 7. Team Plan

Suggested:

### India

₹1,999/user/month

### Global

\$25/user/month

Recommended minimum:

-   3--5 users

Features:

``` text
Team workspace
Shared repositories
Shared agents
Shared MCP
Shared skills
Team prompts
Organization rules
Usage dashboard
Budget controls
SSO
RBAC
Audit logs
Centralized billing
```

------------------------------------------------------------------------

## 8. Enterprise

Do not publish a fixed price.

Use:

> Contact Sales

Potential contract sizes:

``` text
₹5 lakh/year
₹10 lakh/year
₹25 lakh/year
₹50 lakh+/year
```

depending on:

-   Number of users
-   Agent usage
-   Model usage
-   Support
-   Deployment model
-   Security requirements

Enterprise features:

-   SSO
-   SAML
-   OIDC
-   SCIM
-   RBAC
-   Audit logs
-   IP restrictions
-   Private networking
-   Data residency
-   Compliance
-   Dedicated support
-   Private Cloud
-   On-Premise
-   Air-gapped deployment

------------------------------------------------------------------------

## 9. BYOK

Bring Your Own Key.

Suggested:

### India

₹299/month

### Global

\$5/month

Users can provide:

``` text
OPENAI_API_KEY
ANTHROPIC_API_KEY
GEMINI_API_KEY
OPENROUTER_API_KEY
```

GiX charges for the software while the customer pays the model provider
directly.

This creates a low-inference-cost customer segment.

------------------------------------------------------------------------

## 10. Cloud Agent Credits

Do not provide unlimited expensive autonomous agents.

Use usage credits.

Possible structure:

``` text
₹100 credits
₹500 credits
₹1,000 credits
₹5,000 credits
```

Credits can account for:

-   Model
-   Tokens
-   Agent runtime
-   Sandbox execution
-   Context size
-   Tool calls
-   Storage

Maintain transparent usage reporting.

------------------------------------------------------------------------

## 11. Marketplace

Potential marketplace:

``` text
GiX Marketplace
|
+-- Agents
+-- Skills
+-- MCP Servers
+-- Plugins
+-- Templates
+-- Integrations
```

Example products:

``` text
Spring Boot Architect
React Refactoring Agent
AWS DevOps Agent
Kubernetes Troubleshooter
Java Migration Agent
Security Audit Agent
```

Suggested GiX revenue share:

> 20--30%

Third-party creators receive the remainder.

------------------------------------------------------------------------

## 12. Recommended Launch Pricing

  Plan                    India         Global
  ------------ ---------------- --------------
  Free                       ₹0            \$0
  Local                 ₹299/mo         \$5/mo
  Developer             ₹599/mo        \$15/mo
  Pro                 ₹1,499/mo        \$39/mo
  Team           ₹1,999/user/mo   \$25/user/mo
  Enterprise             Custom         Custom

Usage credits are additional.

------------------------------------------------------------------------

## 13. Why Local + BYOK Is Strategically Important

A conventional AI coding SaaS pays inference costs for every request.

GiX can support:

``` text
Local
   -> Customer GPU

BYOK
   -> Customer API account

GiX Cloud
   -> GiX pays inference cost
```

This gives GiX three economics:

### Local

``` text
Revenue
-
Near-zero inference
=
High gross margin
```

### BYOK

``` text
Subscription revenue
-
Low infrastructure
=
High gross margin
```

### GiX Cloud

``` text
Subscription + usage
-
Model inference
-
Infrastructure
=
Managed-service margin
```

------------------------------------------------------------------------

## 14. Avoid "Unlimited AI"

Do not advertise unlimited cloud-agent usage.

Use:

> Generous included usage + transparent additional usage.

This protects the company from users consuming unusually large amounts
of expensive model or sandbox resources.

------------------------------------------------------------------------

## 15. Example Revenue Model

Illustrative scenario:

``` text
7,000 Developer users
2,000 Pro users
800 Team seats
200 Enterprise customers
```

Approximate monthly revenue using example assumptions:

``` text
Developer:
7,000 × ₹599 = ₹41.93 lakh

Pro:
2,000 × ₹1,499 = ₹29.98 lakh

Team:
800 × ₹1,999 = ₹15.99 lakh

Enterprise:
200 × ₹10,000 average monthly equivalent = ₹20 lakh
```

Total:

> Approximately ₹1.08 crore/month

before:

-   Usage revenue
-   Marketplace revenue
-   Annual enterprise contracts
-   Add-ons

This is an illustrative business model, not a forecast.

------------------------------------------------------------------------

## 16. India Strategy

India should be a major initial market.

Use:

-   INR pricing
-   UPI
-   Indian payment methods
-   GST-compliant billing
-   Local support
-   Developer-focused community

Do not simply convert USD pricing into INR.

Use a deliberate India price point.

------------------------------------------------------------------------

## 17. Global Strategy

Global users should see:

``` text
Free
$5
$15
$39
$25/seat Team
Enterprise
```

Use annual discounts where appropriate.

Example:

``` text
Monthly: $15
Annual: $150
```

Annual plans improve cash flow and retention.

------------------------------------------------------------------------

## 18. Pricing Architecture

The product should internally separate:

``` text
Subscription entitlement
        +
Usage allowance
        +
Usage credits
        +
Model cost
        +
Infrastructure cost
```

Example:

``` text
Developer Plan
|
+-- 1,000 included agent units
+-- Included model allowance
+-- Included cloud storage
+-- Additional credits available
```

Do not hard-code pricing logic directly into the agent engine.

Create a Billing/Entitlement service.

------------------------------------------------------------------------

## 19. Long-Term Revenue Architecture

``` text
                       GiX Revenue
                            |
          +-----------------+-----------------+
          |                 |                 |
    Subscriptions        Usage          Enterprise
          |                 |                 |
    +-----+-----+       Credits       +--------+--------+
    |     |     |       Agents        |        |        |
  Local Dev    Pro      Models      SaaS    Private   On-Prem
        Team
          |
          +-----------------------------+
                                        |
                                  Marketplace
                                        |
                              +---------+---------+
                              |         |         |
                            Agents    Skills     MCP
```

------------------------------------------------------------------------

## 20. Strategic Positioning

GiX-Coder should be positioned as:

> AI Software Engineering Platform

not:

> Another AI coding model.

The model should be replaceable.

GiX's moat should be:

``` text
Agent
+
Context
+
Tools
+
Codebase Intelligence
+
Git
+
MCP
+
Memory
+
Cloud Execution
+
IDE
+
Team Collaboration
+
Enterprise Security
```

------------------------------------------------------------------------

## 21. Recommended Business Evolution

``` text
Stage 1
Free local agent
        |
Stage 2
Developer subscriptions
        |
Stage 3
Cloud agents + usage
        |
Stage 4
Team collaboration
        |
Stage 5
Enterprise
        |
Stage 6
Marketplace
        |
Stage 7
Private Cloud / On-Prem
```

The combination of free local models, BYOK, paid cloud agents, team
collaboration, and enterprise deployment is the strongest long-term
GiX-Coder monetization strategy.

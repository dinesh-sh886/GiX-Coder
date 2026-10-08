# GiX-Coder Deployment Strategy

## 1. Executive Direction

GiX-Coder should be deployed as a cloud-native SaaS platform while
keeping the local developer experience first-class.

Recommended strategic path:

1.  Local development with Docker Compose.
2.  Public MVP on AWS using managed services and ECS/Fargate.
3.  Introduce EKS/Kubernetes when cloud-agent and sandbox workloads
    justify it.
4.  Add GPU infrastructure only when self-hosted model economics make
    sense.
5.  Keep the architecture portable enough for GCP/Azure and
    private/on-prem deployments later.

Core principle:

> Docker everywhere, managed cloud services first, Kubernetes when scale
> requires it.

------------------------------------------------------------------------

## 2. Recommended Initial Cloud

### Primary recommendation: AWS

AWS is the preferred initial production cloud because GiX-Coder will
eventually need:

-   Containers
-   Kubernetes
-   GPU compute
-   PostgreSQL
-   Redis
-   Object storage
-   Queues
-   Networking
-   Secrets
-   IAM
-   Monitoring
-   Enterprise infrastructure

### Alternative: GCP

GCP is an excellent alternative, especially for a simpler initial SaaS
deployment:

-   Cloud Run
-   Artifact Registry
-   Cloud SQL PostgreSQL
-   Memorystore
-   Cloud Storage
-   Cloud Logging/Monitoring

### Azure

Azure should become a major supported deployment target later,
especially for enterprise customers already standardized on Microsoft
infrastructure.

------------------------------------------------------------------------

## 3. Phase 1 --- Local Development

Use Docker Compose for the development environment.

Typical components:

-   GiX API
-   GiX Agent
-   Web application
-   PostgreSQL
-   Redis
-   Kafka only if event-driven requirements justify it
-   Ollama/local model runtime
-   Supporting development services

The local environment should make it possible to develop without
depending on GiX Cloud.

------------------------------------------------------------------------

## 4. Phase 2 --- Public MVP

Recommended AWS architecture:

``` text
Internet
   |
CloudFront
   |
Application Load Balancer
   |
ECS/Fargate
   |-- gix-web
   |-- gix-api
   |-- gix-auth
   |-- gix-agent
   |-- gix-worker
   |
   +-- RDS PostgreSQL
   +-- ElastiCache Redis
   +-- S3
```

AWS services:

  Requirement          AWS Service
  -------------------- -------------------
  Containers           ECS/Fargate
  Container registry   ECR
  Database             RDS PostgreSQL
  Cache/queues         ElastiCache Redis
  Object storage       S3
  CDN                  CloudFront
  DNS                  Route 53
  Secrets              Secrets Manager
  Monitoring           CloudWatch
  Identity/access      IAM
  Network              VPC
  Web protection       WAF
  Encryption           KMS

Do not introduce Kubernetes on day one unless there is a strong
operational reason.

------------------------------------------------------------------------

## 5. Docker Strategy

Every deployable service should have a Docker image.

Example:

``` text
gix-api
gix-auth
gix-agent
gix-scheduler
gix-worker
gix-web
```

Images should be pushed to ECR:

``` text
<account>.dkr.ecr.<region>.amazonaws.com/gix-api:<version>
```

Use immutable version tags where possible.

Recommended release tags:

``` text
1.0.0
1.0.1
2026.10.08
git-<commit-sha>
```

------------------------------------------------------------------------

## 6. Database

Use Amazon RDS PostgreSQL rather than running production PostgreSQL
inside an application container.

Core data domains may include:

-   Users
-   Organizations
-   Memberships
-   Projects
-   Repositories
-   Sessions
-   Messages
-   Agents
-   Agent tasks
-   Model providers
-   Model usage
-   Subscriptions
-   Billing
-   API-key metadata
-   MCP servers
-   Skills
-   Audit logs

Enable:

-   Automated backups
-   Point-in-time recovery
-   Encryption
-   Monitoring
-   Read replicas when scale requires them

------------------------------------------------------------------------

## 7. Redis

Use managed Redis for:

-   Session state
-   Rate limiting
-   Caching
-   Distributed locks
-   Agent job state
-   Temporary context
-   Queueing where appropriate
-   WebSocket/session coordination

Avoid making Redis the permanent source of truth for business data.

------------------------------------------------------------------------

## 8. Object Storage

Use S3 for:

``` text
projects/
agent-artifacts/
uploads/
reports/
exports/
build-artifacts/
logs/
```

Do not store large artifacts in PostgreSQL.

------------------------------------------------------------------------

## 9. Networking

Recommended structure:

``` text
VPC
|
+-- Public subnets
|    +-- Load Balancer
|    +-- NAT
|
+-- Private subnets
     +-- ECS/EKS workloads
     +-- RDS PostgreSQL
     +-- Redis
     +-- internal services
```

Do not expose PostgreSQL, Redis, agent workers, or internal APIs
directly to the public internet.

------------------------------------------------------------------------

## 10. Agent Execution and Sandbox Architecture

Agent workloads are different from ordinary API workloads because agents
may:

-   Clone repositories
-   Execute shell commands
-   Run Maven/Gradle
-   Run npm/pnpm
-   Execute Python
-   Run tests
-   Build applications
-   Modify files
-   Interact with Git
-   Use MCP tools

Never allow untrusted agent execution directly inside the main API
container.

Recommended model:

``` text
Agent API
   |
Scheduler
   |
Sandbox Manager
   |
+------------------+
| Ephemeral Sandbox|
|                  |
| Git              |
| Java             |
| Node             |
| Python           |
| Build tools      |
| Project files    |
+------------------+
   |
Destroy after task
```

Security controls should include:

-   Workspace-scoped filesystem access
-   Network policies
-   Command permissions
-   Resource limits
-   Timeouts
-   Container isolation
-   Explicit approval policies
-   Audit logging

------------------------------------------------------------------------

## 11. Phase 3 --- Kubernetes / EKS

Introduce EKS when:

-   Agent workloads become numerous
-   Sandboxes require orchestration
-   Worker autoscaling becomes complex
-   Multiple workload classes need scheduling
-   GPU workloads are introduced
-   Enterprise/private deployment becomes important

Example:

``` text
EKS
|
+-- API
+-- Agent Scheduler
+-- Agent Workers
+-- Sandbox Workers
+-- MCP Workers
+-- Background Jobs
```

Continue using managed services where practical:

-   RDS
-   Redis
-   S3
-   ECR

Do not move every managed service into Kubernetes just because EKS is
introduced.

------------------------------------------------------------------------

## 12. AI Model Infrastructure

Initially use a model gateway rather than hosting large models yourself.

``` text
GiX Model Gateway
|
+-- OpenAI
+-- Anthropic
+-- Gemini
+-- OpenRouter
+-- User BYOK
+-- Custom OpenAI-compatible endpoint
```

For local users:

``` text
GiX Local Agent
|
+-- Ollama
+-- vLLM
+-- LM Studio
+-- LocalAI
```

Later, when economics justify it:

``` text
GPU Nodes
   |
vLLM
   |
Qwen / DeepSeek / Llama / Nemotron
```

------------------------------------------------------------------------

## 13. Secrets

Use AWS Secrets Manager or an equivalent secret-management system.

Examples:

``` text
OPENAI_API_KEY
ANTHROPIC_API_KEY
GEMINI_API_KEY
STRIPE_SECRET
RAZORPAY_SECRET
GITHUB_CLIENT_SECRET
DATABASE_PASSWORD
JWT_SECRET
WEBHOOK_SECRET
```

Never store production secrets in:

-   Git
-   Dockerfiles
-   Public configuration files
-   Jenkinsfile
-   Application source code

------------------------------------------------------------------------

## 14. Observability

Use OpenTelemetry as the instrumentation layer.

Recommended stack:

``` text
OpenTelemetry
   |
   +-- Prometheus
   +-- Grafana
   +-- CloudWatch
```

Track:

-   API latency
-   Agent latency
-   Token usage
-   Model cost
-   Agent success rate
-   Sandbox failures
-   Queue depth
-   CPU
-   Memory
-   Database connections
-   Redis usage
-   Error rate
-   Cost per agent task

Cost per agent task is especially important because GiX-Coder is a
usage-sensitive SaaS.

------------------------------------------------------------------------

## 15. Billing and Metering

Build usage metering from the beginning.

Example event:

``` json
{
  "userId": "...",
  "model": "qwen",
  "provider": "openrouter",
  "inputTokens": 12000,
  "outputTokens": 5000,
  "durationMs": 3400,
  "estimatedCost": 0.12
}
```

Flow:

``` text
Model Gateway
   |
Usage Service
   |
Billing
   |
Subscription / Credits
```

This allows:

``` text
Revenue
- Model Cost
- Infrastructure Cost
= Gross Margin
```

------------------------------------------------------------------------

## 16. CI/CD

Jenkins is supported and can be used as the primary CI/CD system.

Recommended pipeline:

``` text
GitHub
   |
Jenkins
   |
   +-- Checkout
   +-- Compile
   +-- Unit tests
   +-- Integration tests
   +-- SonarQube
   +-- SAST/dependency scanning
   +-- Docker build
   +-- Container scan
   +-- Push to ECR
   +-- Deploy staging
   +-- Smoke tests
   +-- Approval
   +-- Deploy production
```

### Jenkins recommendation

Use Jenkins when:

-   Existing team expertise exists
-   Enterprise customers expect Jenkins
-   Complex/self-hosted pipelines are required

For the earliest startup MVP, GitHub Actions is also a simpler option.

Jenkins should remain a replaceable CI/CD layer rather than a hard
dependency of the application runtime.

------------------------------------------------------------------------

## 17. Production Architecture

Target architecture:

``` text
Internet
   |
CloudFront
   |
Route 53
   |
Application Load Balancer
   |
+-------------------------------+
| Web / API                     |
+-------------------------------+
          |
     +----+----+
     |         |
   Auth      Agent API
               |
           Scheduler
               |
             Redis
               |
       +-------+-------+
       |               |
 Agent Worker       Sandbox
       |               |
       +-------+-------+
               |
        Model Gateway
               |
    +----------+----------+
    |          |          |
 OpenAI    Anthropic   OpenRouter

Data:
- PostgreSQL
- Redis
- S3

Infrastructure:
- ECR
- ECS/EKS
- Secrets Manager
- IAM
- VPC
- WAF
- CloudWatch
```

------------------------------------------------------------------------

## 18. Deployment Philosophy

Do not overbuild for the first 100--1,000 users.

Avoid starting with:

-   Multiple Kubernetes clusters
-   20+ microservices
-   GPU fleets
-   Active-active multi-region
-   Service mesh
-   Complex event architecture

Start small, measure real workloads, then introduce infrastructure
complexity when justified.

------------------------------------------------------------------------

## 19. Recommended Long-Term Cloud Strategy

``` text
Phase 1
Docker Compose
   |
Phase 2
AWS ECS/Fargate + Managed Services
   |
Phase 3
AWS EKS for agents/sandboxes
   |
Phase 4
Self-hosted GPU inference
   |
Phase 5
Multi-cloud / Private Cloud / On-Prem
```

Long-term deployment targets:

-   GiX Cloud
-   GiX Private Cloud
-   GiX On-Premise
-   GiX Air-Gapped
-   Enterprise-managed Kubernetes

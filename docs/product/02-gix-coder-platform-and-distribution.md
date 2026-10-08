# GiX-Coder Platform and Distribution Strategy

## 1. Core Product Principle

GiX-Coder should not be designed as a CLI with several unrelated
clients.

The product should be:

> GiX Agent Platform = core product\
> CLI / Web / IDE / Desktop / Mobile = clients

All clients should share:

-   Agent engine
-   Model abstraction
-   Tool system
-   Context engine
-   Memory
-   MCP
-   Permissions
-   Session model
-   Usage metering
-   Account/workspace model

------------------------------------------------------------------------

## 2. Platform Architecture

``` text
GiX Agent Platform
|
+-- Agent Runtime
|    +-- Planner
|    +-- Executor
|    +-- Context
|    +-- Memory
|    +-- Tools
|    +-- MCP
|    +-- Git
|    +-- Permissions
|
+-- Model Gateway
|    +-- Ollama
|    +-- vLLM
|    +-- OpenAI
|    +-- Anthropic
|    +-- Gemini
|    +-- OpenRouter
|    +-- Custom OpenAI-compatible
|
+-- Protocol/API
     +-- Sessions
     +-- Streaming
     +-- Events
     +-- Authentication
```

Clients:

``` text
CLI
Web
VS Code
JetBrains
Windows
macOS
Linux
iOS
Android
```

------------------------------------------------------------------------

## 3. CLI

The CLI should be the first client.

Installation examples:

``` bash
curl -fsSL https://gix.ai/install.sh | sh
```

or:

``` bash
npm install -g @gix-coder/cli
```

Commands:

``` bash
gix
gix chat
gix code
gix review
gix fix
gix test
gix commit
gix models
gix config
gix doctor
```

Example:

``` bash
gix --model qwen3-coder:30b
```

The CLI should support:

-   Local models
-   BYOK
-   GiX Cloud
-   Git
-   Terminal
-   MCP
-   Agent workflows
-   Streaming
-   Sessions

------------------------------------------------------------------------

## 4. Web Application

Primary URL:

``` text
https://app.gix.ai
```

Web is primarily for cloud workflows.

Capabilities:

-   Repository connection
-   Cloud agents
-   Code review
-   Architecture analysis
-   PR generation
-   Issue fixing
-   Session history
-   Agent monitoring
-   Usage
-   Billing
-   Team administration

Example:

``` text
Browser
   |
app.gix.ai
   |
GiX Cloud
   |
Cloud Agent
   |
Ephemeral Coding VM
```

------------------------------------------------------------------------

## 5. VS Code

Publish a GiX-Coder extension through the VS Code Marketplace.

The extension should be thin:

``` text
VS Code
   |
GiX Extension
   |
GiX Agent Protocol
   |
Local Agent OR GiX Cloud
```

Features:

-   Chat
-   Agent mode
-   Fix
-   Explain
-   Refactor
-   Test generation
-   Code review
-   Diff review
-   Terminal integration
-   MCP

------------------------------------------------------------------------

## 6. JetBrains

Create a JetBrains plugin for:

-   IntelliJ IDEA
-   PyCharm
-   WebStorm
-   GoLand
-   Rider
-   Android Studio

This is especially important for Java/Spring Boot developers.

Architecture:

``` text
IntelliJ
   |
GiX Plugin
   |
GiX Agent Protocol
   |
Local Agent / Cloud Agent
```

------------------------------------------------------------------------

## 7. Desktop Applications

Provide native-feeling applications for:

-   Windows
-   macOS
-   Linux

Recommended distribution:

### Windows

``` text
GiX-Coder-Setup.exe
```

or:

``` text
GiX-Coder-x64.msi
```

### macOS

``` text
GiX-Coder.dmg
```

### Linux

-   One-line installer
-   Package managers
-   Portable binaries

The desktop application should provide:

-   Agent
-   Projects
-   Terminal
-   Git
-   Models
-   Sessions
-   MCP
-   Settings

The desktop app may bundle the CLI.

------------------------------------------------------------------------

## 8. Mobile Applications

### iOS

GiX Mobile should not attempt to become a full IDE.

Focus on:

-   Agent monitoring
-   Cloud agent control
-   PR review
-   Code review
-   Notifications
-   Session history
-   Task management

Example:

``` text
GiX Mobile

Active Agents
-----------------
Payment Service     Running
Security Review     Running

Completed
-----------------
PR #184             CI Passed
PR #181             Merged
```

### Android

Provide the same core mobile control-plane experience.

------------------------------------------------------------------------

## 9. One GiX Account

Users should have one account across every client.

``` text
GiX Account
|
+-- Subscription
+-- Projects
+-- Repositories
+-- Agents
+-- Sessions
+-- Models
+-- API keys
+-- MCP
+-- Skills
+-- Usage
+-- Billing
```

The same account should work across:

``` text
Windows
macOS
Linux
CLI
VS Code
IntelliJ
Web
iOS
Android
```

------------------------------------------------------------------------

## 10. Local Mode

Local mode is a major GiX differentiator.

Example:

``` text
Windows/Linux/macOS
|
+-- GiX
+-- Ollama
|    +-- Qwen3-Coder
|    +-- Llama
|    +-- DeepSeek
|    +-- Nemotron
|
+-- Local Agent
```

No source code needs to pass through GiX Cloud.

Users should be able to select:

``` text
Provider: Ollama
Model: qwen3-coder:30b
Execution: Local
```

------------------------------------------------------------------------

## 11. BYOK Mode

Users can supply their own provider credentials.

Supported examples:

``` text
OpenAI
Anthropic
Gemini
OpenRouter
Custom OpenAI-compatible endpoint
```

Architecture:

``` text
GiX Agent
   |
User's API key
   |
Model Provider
```

This keeps GiX's infrastructure cost low while still providing a full
developer experience.

------------------------------------------------------------------------

## 12. Cloud Agent Mode

Cloud agents should run independently of the user's laptop.

Flow:

``` text
User
 |
Web/Desktop/Mobile/CLI
 |
GiX Cloud
 |
Agent Scheduler
 |
Ephemeral Sandbox
 |
Git + Build + Tests
 |
Commit + Push + PR
```

The user can close the application while the agent continues.

------------------------------------------------------------------------

## 13. Distribution Website

Recommended:

``` text
https://gix.ai/download
```

Show:

``` text
Windows
macOS
Linux
CLI
VS Code
JetBrains
iOS
Android
```

The page should detect the user's platform and highlight the appropriate
download.

------------------------------------------------------------------------

## 14. Release Strategy

### Stage 1

CLI + Local models + BYOK

Platforms:

-   Linux
-   macOS
-   Windows/WSL

### Stage 2

IDE integrations:

-   VS Code
-   JetBrains

### Stage 3

Desktop:

-   Windows
-   macOS
-   Linux

### Stage 4

Web + Cloud Agents

### Stage 5

Mobile:

-   iOS
-   Android

------------------------------------------------------------------------

## 15. Client Architecture Rule

Do not implement separate agent engines.

Bad:

``` text
CLI Agent
Web Agent
VS Code Agent
Desktop Agent
Mobile Agent
```

Preferred:

``` text
                 GiX Agent Runtime
                        |
             GiX Agent Protocol/API
                        |
        +---------------+---------------+
        |               |               |
       CLI             IDE           Desktop
        |               |               |
        +---------------+---------------+
                        |
                    Same engine
```

This dramatically reduces long-term maintenance.

------------------------------------------------------------------------

## 16. Recommended Product Map

``` text
gix.ai
|
+-- GiX Cloud
|    +-- Web
|    +-- Cloud Agents
|    +-- API
|    +-- Team
|    +-- Enterprise
|
+-- Downloads
|    +-- Windows
|    +-- macOS
|    +-- Linux
|    +-- CLI
|
+-- IDE
|    +-- VS Code
|    +-- JetBrains
|
+-- Mobile
     +-- iOS
     +-- Android
```

Core execution:

``` text
Local Agent
   |
+-- Ollama
+-- vLLM
+-- LM Studio
+-- LocalAI

Cloud Agent
   |
+-- GiX Model Gateway
+-- OpenAI
+-- Anthropic
+-- Gemini
+-- OpenRouter
+-- Self-hosted models
```

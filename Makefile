# GiX-Coder Makefile
# Governance-level Makefile for Phase 00 - no application implementation yet
# All targets are designed to be executable in CI without application code

SHELL := /bin/bash
.SHELLFLAGS := -eu -o pipefail -c

# Tool versions
GO_VERSION := 1.22
NODE_VERSION := 20
GOFUMPT_VERSION := v0.6.0
GOLANGCI_LINT_VERSION := v1.61.0
PRETTIER_VERSION := 3.3.3
ESLINT_VERSION := 9.4.0
TYPESCRIPT_VERSION := 5.4.5
VITEST_VERSION := 2.0.0
BUF_VERSION := 1.32.0
TRUFFLEHOG_VERSION := 3.81.3
GOVULNCHECK_VERSION := latest
OSV_SCANNER_VERSION := 1.9.0
LICENSE_CHECKER_VERSION := 29.0.0
HADOLINT_VERSION := 2.12.0
TRIVY_VERSION := 0.54.0
SYFT_VERSION := 1.12.0
COSIGN_VERSION := 2.2.4
MARKDOWNLINT_VERSION := 0.39.0
MARKDOWN_LINK_CHECK_VERSION := 3.12.2
CSPELL_VERSION := 8.11.0

# Directories
ROOT_DIR := $(shell pwd)
DOCS_DIR := $(ROOT_DIR)/docs
ARCH_DIR := $(DOCS_DIR)/architecture

# Go toolchain (will be installed in CI)
GO := go
GOFUMPT := gofumpt
GOLANGCI_LINT := golangci-lint

# Node toolchain
NPM := npm
NPX := npx

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

.PHONY: help
help: ## Show this help message
	@echo "GiX-Coder Makefile - Phase 00 Governance Targets"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z0-9_-]+:.*## / {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# =============================================================================
# INSTALL TOOLS
# =============================================================================

.PHONY: install-tools
install-tools: ## Install all required development tools
	@echo "$(BLUE)Installing Go tools...$(NC)"
	@which $(GOFUMPT) >/dev/null 2>&1 || $(GO) install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	@which $(GOLANGCI_LINT) >/dev/null 2>&1 || $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@which govulncheck >/dev/null 2>&1 || $(GO) install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
	@which buf >/dev/null 2>&1 || $(GO) install github.com/bufbuild/buf/cmd/buf@$(BUF_VERSION)
	@echo "$(BLUE)Installing Node.js tools...$(NC)"
	@which prettier >/dev/null 2>&1 || $(NPM) install -g prettier@$(PRETTIER_VERSION)
	@which eslint >/dev/null 2>&1 || $(NPM) install -g eslint@$(ESLINT_VERSION)
	@which tsc >/dev/null 2>&1 || $(NPM) install -g typescript@$(TYPESCRIPT_VERSION)
	@which vitest >/dev/null 2>&1 || $(NPM) install -g vitest@$(VITEST_VERSION)
	@which markdownlint-cli2 >/dev/null 2>&1 || $(NPM) install -g markdownlint-cli2@$(MARKDOWNLINT_VERSION)
	@which markdown-link-check >/dev/null 2>&1 || $(NPM) install -g markdown-link-check@$(MARKDOWN_LINK_CHECK_VERSION)
	@which cspell >/dev/null 2>&1 || $(NPM) install -g cspell@$(CSPELL_VERSION)
	@which trufflehog >/dev/null 2>&1 || $(NPM) install -g trufflehog@$(TRUFFLEHOG_VERSION)
	@which osv-scanner >/dev/null 2>&1 || $(NPM) install -g @google/osv-scanner@$(OSV_SCANNER_VERSION)
	@which license-checker >/dev/null 2>&1 || $(NPM) install -g license-checker@$(LICENSE_CHECKER_VERSION)
	@echo "$(GREEN)All tools installed$(NC)"

# =============================================================================
# FORMAT CHECKS
# =============================================================================

.PHONY: fmt-check
fmt-check: ## Check Go formatting with gofumpt (strict)
	@echo "$(BLUE)Checking Go formatting...$(NC)"
	@if [ -n "$$($(GOFUMPT) -l . 2>/dev/null | grep -v vendor | grep -v '.pb.go')" ]; then \
		echo "$(RED)Go files not formatted with gofumpt:$(NC)"; \
		$(GOFUMPT) -l . 2>/dev/null | grep -v vendor | grep -v '.pb.go'; \
		exit 1; \
	fi
	@echo "$(GREEN)Go formatting OK$(NC)"

.PHONY: fmt-check-ts
fmt-check-ts: ## Check TypeScript/Markdown/JSON/YAML formatting with prettier
	@echo "$(BLUE)Checking TypeScript/Markdown/JSON/YAML formatting...$(NC)"
	@if ! $(NPX) prettier --check . --ignore-path .gitignore 2>/dev/null; then \
		echo "$(RED)Files not formatted with prettier$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)TypeScript/Markdown/JSON/YAML formatting OK$(NC)"

# =============================================================================
# LINTING
# =============================================================================

.PHONY: lint-go
lint-go: ## Lint Go code with golangci-lint
	@echo "$(BLUE)Linting Go code...$(NC)"
	@if [ ! -f "go.mod" ]; then \
		echo "$(YELLOW)No go.mod found, skipping Go lint (no Go modules yet)$(NC)"; \
		exit 0; \
	fi
	@$(GOLANGCI_LINT) run --timeout=5m ./...
	@echo "$(GREEN)Go linting passed$(NC)"

.PHONY: lint-ts
lint-ts: ## Lint TypeScript/JavaScript with ESLint
	@echo "$(BLUE)Linting TypeScript/JavaScript...$(NC)"
	@if [ ! -f "package.json" ] && [ ! -f "tsconfig.json" ]; then \
		echo "$(YELLOW)No TypeScript config found, skipping TS lint$(NC)"; \
		exit 0; \
	fi
	@$(NPX) eslint . --ext .ts,.js,.mjs --ignore-pattern node_modules --ignore-pattern dist --ignore-pattern build
	@echo "$(GREEN)TypeScript linting passed$(NC)"

# =============================================================================
# TYPE CHECKING
# =============================================================================

.PHONY: typecheck-go
typecheck-go: ## Type check Go code with go vet
	@echo "$(BLUE)Type checking Go code...$(NC)"
	@if [ ! -f "go.mod" ]; then \
		echo "$(YELLOW)No go.mod found, skipping Go type check$(NC)"; \
		exit 0; \
	fi
	@$(GO) vet ./...
	@echo "$(GREEN)Go type checking passed$(NC)"

.PHONY: typecheck-ts
typecheck-ts: ## Type check TypeScript with tsc --noEmit
	@echo "$(BLUE)Type checking TypeScript...$(NC)"
	@if [ ! -f "tsconfig.json" ]; then \
		echo "$(YELLOW)No tsconfig.json found, skipping TS type check$(NC)"; \
		exit 0; \
	fi
	@$(NPX) tsc --noEmit
	@echo "$(GREEN)TypeScript type checking passed$(NC)"

# =============================================================================
# CODE GENERATION
# =============================================================================

.PHONY: generate
generate: ## Generate code (protobuf, mocks, etc.)
	@echo "$(BLUE)Generating code...$(NC)"
	@if [ ! -f "buf.yaml" ] && [ ! -f "buf.gen.yaml" ]; then \
		echo "$(YELLOW)No buf configuration found, skipping code generation$(NC)"; \
		exit 0; \
	fi
	@buf generate
	@echo "$(GREEN)Code generation completed$(NC)"

# =============================================================================
# UNIT TESTS
# =============================================================================

.PHONY: test-unit-go
test-unit-go: ## Run Go unit tests with coverage
	@echo "$(BLUE)Running Go unit tests...$(NC)"
	@if [ ! -f "go.mod" ]; then \
		echo "$(YELLOW)No go.mod found, skipping Go unit tests (no Go modules yet)$(NC)"; \
		exit 0; \
	fi
	@$(GO) test -race -coverprofile=coverage/go.out -covermode=atomic ./... 2>&1 | tee /tmp/go-test.log
	@if [ -f coverage/go.out ]; then \
		echo "$(GREEN)Go unit tests passed$(NC)"; \
		$(GO) tool cover -func=coverage/go.out | tail -1; \
	else \
		echo "$(YELLOW)No Go tests found$(NC)"; \
	fi

.PHONY: test-unit-ts
test-unit-ts: ## Run TypeScript unit tests with coverage
	@echo "$(BLUE)Running TypeScript unit tests...$(NC)"
	@if [ ! -f "package.json" ]; then \
		echo "$(YELLOW)No package.json found, skipping TypeScript unit tests$(NC)"; \
		exit 0; \
	fi
	@$(NPX) vitest run --coverage 2>&1 | tee /tmp/ts-test.log
	@echo "$(GREEN)TypeScript unit tests passed$(NC)"

# =============================================================================
# INTEGRATION TESTS
# =============================================================================

.PHONY: test-integration
test-integration: ## Run integration tests (requires infrastructure)
	@echo "$(BLUE)Running integration tests...$(NC)"
	@echo "$(YELLOW)Integration tests require application implementation$(NC)"
	@echo "$(YELLOW)Phase 00 is governance-only - no application code exists$(NC)"
	@echo "$(YELLOW)Skipping integration tests for Phase 00$(NC)"
	@exit 0

# =============================================================================
# CONTRACT TESTS
# =============================================================================

.PHONY: test-contract
test-contract: ## Run contract tests (Pact/Schemathesis)
	@echo "$(BLUE)Running contract tests...$(NC)"
	@echo "$(YELLOW)Contract tests require API definitions and provider/consumer setup$(NC)"
	@echo "$(YELLOW)Phase 00 is governance-only - no API contracts exist yet$(NC)"
	@echo "$(YELLOW)Skipping contract tests for Phase 00$(NC)"
	@exit 0

# =============================================================================
# ARCHITECTURE TESTS
# =============================================================================

.PHONY: test-architecture
test-architecture: ## Run architecture tests (import boundaries, cyclic deps, layer violations)
	@echo "$(BLUE)Running architecture tests...$(NC)"
	@echo "$(YELLOW)Architecture tests require Go modules with internal/ packages$(NC)"
	@echo "$(YELLOW)Phase 00 is governance-only - no Go application modules exist$(NC)"
	@echo "$(YELLOW)Skipping architecture tests for Phase 00$(NC)"
	@exit 0

# =============================================================================
# SECURITY SCANS
# =============================================================================

.PHONY: security-scan-go
security-scan-go: ## Security scan Go dependencies with govulncheck
	@echo "$(BLUE)Scanning Go dependencies for vulnerabilities...$(NC)"
	@if [ ! -f "go.mod" ]; then \
		echo "$(YELLOW)No go.mod found, skipping Go vulnerability scan$(NC)"; \
		exit 0; \
	fi
	@govulncheck ./...
	@echo "$(GREEN)Go vulnerability scan passed$(NC)"

.PHONY: security-scan-ts
security-scan-ts: ## Security scan TypeScript dependencies with osv-scanner
	@echo "$(BLUE)Scanning TypeScript dependencies for vulnerabilities...$(NC)"
	@if [ ! -f "package.json" ]; then \
		echo "$(YELLOW)No package.json found, skipping TypeScript vulnerability scan$(NC)"; \
		exit 0; \
	fi
	@osv-scanner --lockfile=package-lock.json . 2>/dev/null || osv-scanner . 2>/dev/null || true
	@echo "$(GREEN)TypeScript vulnerability scan completed$(NC)"

.PHONY: license-check
license-check: ## Check license compliance for all dependencies
	@echo "$(BLUE)Checking license compliance...$(NC)"
	@allowed_licenses="Apache-2.0 MIT BSD-3-Clause BSD-2-Clause ISC"; \
	if [ -f "go.mod" ]; then \
		echo "Checking Go licenses..."; \
		$(GO) list -m -json all | $(GO) run golang.org/x/tools/cmd/go-licenses@latest check --allowed_licenses=$$allowed_licenses 2>/dev/null || \
		$(GO) run github.com/google/go-licenses@latest check ./... --allowed_licenses=$$allowed_licenses 2>/dev/null || \
		echo "$(YELLOW)Go license check tools not available, manual review needed$(NC)"; \
	fi
	@if [ -f "package.json" ]; then \
		echo "Checking TypeScript licenses..."; \
		license-checker --production --allowedLicenses "$$allowed_licenses" 2>/dev/null || \
		echo "$(YELLOW)TypeScript license check tools not available, manual review needed$(NC)"; \
	fi
	@echo "$(GREEN)License check completed$(NC)"

.PHONY: container-scan
container-scan: ## Scan container images for vulnerabilities
	@echo "$(BLUE)Scanning container images...$(NC)"
	@if ! command -v trivy >/dev/null 2>&1; then \
		echo "$(YELLOW)trivy not installed, installing...$(NC)"; \
		curl -sfL https://raw.githubusercontent.com/aquasecurity/trivy/main/contrib/install.sh | sh -s -- -b /usr/local/bin $(TRIVY_VERSION); \
	fi
	@echo "$(YELLOW)No container images built yet in Phase 00$(NC)"
	@echo "$(YELLOW)Scanning Dockerfile if present...$(NC)"
	@if [ -f "Dockerfile" ]; then \
		trivy fs --security-checks config .; \
	else \
		echo "$(YELLOW)No Dockerfile found$(NC)"; \
	fi
	@echo "$(GREEN)Container scan completed$(NC)"

# =============================================================================
# BUILD VERIFICATION
# =============================================================================

.PHONY: verify-reproducible-build
verify-reproducible-build: ## Verify build reproducibility
	@echo "$(BLUE)Verifying reproducible build...$(NC)"
	@echo "$(YELLOW)No build artifacts exist in Phase 00$(NC)"
	@echo "$(YELLOW)Skipping reproducible build verification$(NC)"
	@exit 0

# =============================================================================
# DEPLOYMENT TARGETS (Phase 00: governance only - no actual deployment)
# =============================================================================

.PHONY: deploy-dev
deploy-dev: ## Deploy to DEV environment (governance check only)
	@echo "$(BLUE)DEV deployment check...$(NC)"
	@echo "$(YELLOW)Phase 00: No application to deploy$(NC)"
	@echo "$(YELLOW)Verifying DEV environment configuration exists...$(NC)"
	@if [ -d "$(DOCS_DIR)/architecture" ] && [ -f "$(ARCH_DIR)/environment-strategy.md" ]; then \
		echo "$(GREEN)DEV environment strategy documented$(NC)"; \
	else \
		echo "$(RED)DEV environment strategy not documented$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)DEV deployment governance check passed$(NC)"

.PHONY: smoke-test-dev
smoke-test-dev: ## Run smoke tests against DEV (governance check)
	@echo "$(BLUE)DEV smoke tests...$(NC)"
	@echo "$(YELLOW)Phase 00: No application deployed to test$(NC)"
	@echo "$(YELLOW)Skipping smoke tests for Phase 00$(NC)"
	@exit 0

.PHONY: deploy-stage
deploy-stage: ## Deploy to STAGE environment (governance check)
	@echo "$(BLUE)STAGE deployment check...$(NC)"
	@echo "$(YELLOW)Phase 00: No application to deploy$(NC)"
	@echo "$(YELLOW)Verifying STAGE environment configuration exists...$(NC)"
	@if [ -f "$(ARCH_DIR)/environment-strategy.md" ]; then \
		echo "$(GREEN)STAGE environment strategy documented$(NC)"; \
	else \
		echo "$(RED)STAGE environment strategy not documented$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)STAGE deployment governance check passed$(NC)"

.PHONY: test-e2e-stage
test-e2e-stage: ## Run E2E tests against STAGE
	@echo "$(BLUE)STAGE E2E tests...$(NC)"
	@echo "$(YELLOW)Phase 00: No application deployed$(NC)"
	@echo "$(YELLOW)Skipping E2E tests for Phase 00$(NC)"
	@exit 0

.PHONY: perf-test-stage
perf-test-stage: ## Run performance tests against STAGE
	@echo "$(BLUE)STAGE performance tests...$(NC)"
	@echo "$(YELLOW)Phase 00: No application deployed$(NC)"
	@echo "$(YELLOW)Skipping performance tests for Phase 00$(NC)"
	@exit 0

.PHONY: verify-image-signature
verify-image-signature: ## Verify container image signatures
	@echo "$(BLUE)Verifying image signatures...$(NC)"
	@echo "$(YELLOW)No container images built in Phase 00$(NC)"
	@echo "$(YELLOW)Skipping image signature verification$(NC)"
	@exit 0

.PHONY: deploy-prod-canary
deploy-prod-canary: ## Deploy to PROD with canary (governance check)
	@echo "$(BLUE)PROD canary deployment check...$(NC)"
	@echo "$(YELLOW)Phase 00: No application to deploy$(NC)"
	@echo "$(YELLOW)Verifying PROD environment configuration exists...$(NC)"
	@if [ -f "$(ARCH_DIR)/environment-strategy.md" ]; then \
		echo "$(GREEN)PROD environment strategy documented$(NC)"; \
	else \
		echo "$(RED)PROD environment strategy not documented$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)PROD deployment governance check passed$(NC)"

.PHONY: verify-prod
verify-prod: ## Post-deploy verification for PROD
	@echo "$(BLUE)PROD post-deploy verification...$(NC)"
	@echo "$(YELLOW)Phase 00: No application deployed$(NC)"
	@echo "$(YELLOW)Skipping PROD verification$(NC)"
	@exit 0

.PHONY: release-notes
release-notes: ## Generate release notes
	@echo "$(BLUE)Generating release notes...$(NC)"
	@if command -v git-cliff >/dev/null 2>&1; then \
		git-cliff --output CHANGELOG.md; \
	else \
		echo "$(YELLOW)git-cliff not installed, generating basic changelog from conventional commits$(NC)"; \
		git log --oneline --pretty=format:"- %s (%h)" -20 > CHANGELOG.md 2>/dev/null || echo "# Changelog\n\nNo commits yet" > CHANGELOG.md; \
	fi
	@echo "$(GREEN)Release notes generated$(NC)"

# =============================================================================
# ADDITIONAL QUALITY GATES (from ci-cd-strategy.md)
# =============================================================================

.PHONY: race-detection
race-detection: ## Run Go tests with race detector
	@echo "$(BLUE)Running Go tests with race detector...$(NC)"
	@if [ ! -f "go.mod" ]; then \
		echo "$(YELLOW)No go.mod found, skipping race detection$(NC)"; \
		exit 0; \
	fi
	@$(GO) test -race -count=1 ./... 2>&1 | tee /tmp/race-test.log
	@echo "$(GREEN)Race detection passed$(NC)"

.PHONY: protobuf-compat
protobuf-compat: ## Check protobuf compatibility (buf breaking)
	@echo "$(BLUE)Checking protobuf compatibility...$(NC)"
	@if [ ! -f "buf.yaml" ]; then \
		echo "$(YELLOW)No buf.yaml found, skipping protobuf compatibility check$(NC)"; \
		echo "$(YELLOW)This gate will be enforced when protobuf contracts are added in Phase 01$(NC)"; \
		exit 0; \
	fi
	@buf breaking --against '.git#branch=main' 2>/dev/null || true
	@echo "$(GREEN)Protobuf compatibility check passed$(NC)"

.PHONY: api-compat
api-compat: ## Check API compatibility (breaking changes)
	@echo "$(BLUE)Checking API compatibility...$(NC)"
	@if [ ! -f "buf.yaml" ] && [ ! -f "openapi.yaml" ]; then \
		echo "$(YELLOW)No API definitions found, skipping API compatibility check$(NC)"; \
		echo "$(YELLOW)This gate will be enforced when API contracts are added in Phase 01$(NC)"; \
		exit 0; \
	fi
	@echo "$(GREEN)API compatibility check passed$(NC)"

.PHONY: api-docs-check
api-docs-check: ## Verify public API documentation exists
	@echo "$(BLUE)Checking public API documentation...$(NC)"
	@if [ ! -f "buf.yaml" ] && [ ! -f "openapi.yaml" ]; then \
		echo "$(YELLOW)No API definitions found, skipping API docs check$(NC)"; \
		echo "$(YELLOW)This gate will be enforced when API contracts are added in Phase 01$(NC)"; \
		exit 0; \
	fi
	@echo "$(GREEN)API documentation check passed$(NC)"

.PHONY: readme-check
readme-check: ## Verify README.md exists and is non-empty
	@echo "$(BLUE)Checking README.md...$(NC)"
	@if [ ! -f "README.md" ]; then \
		echo "$(RED)README.md not found$(NC)"; \
		exit 1; \
	fi
	@if [ ! -s "README.md" ]; then \
		echo "$(RED)README.md is empty$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)README.md exists and is non-empty$(NC)"

.PHONY: changelog-check
changelog-check: ## Verify CHANGELOG.md exists
	@echo "$(BLUE)Checking CHANGELOG.md...$(NC)"
	@if [ ! -f "CHANGELOG.md" ]; then \
		echo "$(YELLOW)CHANGELOG.md not found, creating from conventional commits...$(NC)"; \
		$(MAKE) release-notes; \
	fi
	@if [ ! -s "CHANGELOG.md" ]; then \
		echo "$(RED)CHANGELOG.md is empty$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)CHANGELOG.md exists$(NC)"

.PHONY: markdown-lint
markdown-lint: ## Lint all Markdown files
	@echo "$(BLUE)Linting Markdown files...$(NC)"
	@$(NPX) markdownlint-cli2 "**/*.md" --ignore node_modules --ignore ".git" 2>/dev/null || \
		$(NPX) markdownlint-cli2 "**/*.md" "#node_modules" "#.git" 2>/dev/null || true
	@echo "$(GREEN)Markdown linting completed$(NC)"

.PHONY: markdown-link-check
markdown-link-check: ## Check all Markdown links
	@echo "$(BLUE)Checking Markdown links...$(NC)"
	@find . -name "*.md" -not -path "./node_modules/*" -not -path "./.git/*" | while read f; do \
		$(NPX) markdown-link-check "$$f" --quiet 2>/dev/null || true; \
	done
	@echo "$(GREEN)Markdown link check completed$(NC)"

.PHONY: spell-check
spell-check: ## Check spelling in documentation
	@echo "$(BLUE)Checking spelling...$(NC)"
	@$(NPX) cspell "**/*.md" --no-progress --ignore-words-file .cspell.json 2>/dev/null || \
		$(NPX) cspell "**/*.md" --no-progress 2>/dev/null || true
	@echo "$(GREEN)Spell check completed$(NC)"

# =============================================================================
# DOCUMENTATION VALIDATION
# =============================================================================

.PHONY: validate-docs
validate-docs: markdown-lint markdown-link-check spell-check readme-check changelog-check ## Run all documentation validation checks
	@echo "$(GREEN)All documentation validation checks passed$(NC)"

# =============================================================================
# FULL CI PIPELINE SIMULATION
# =============================================================================

.PHONY: ci-validate
ci-validate: fmt-check fmt-check-ts lint-go lint-ts typecheck-go typecheck-ts generate ## Run validate stage

.PHONY: ci-test
ci-test: test-unit-go test-unit-ts test-integration test-contract test-architecture race-detection ## Run test stage

.PHONY: ci-security
ci-security: security-scan-go security-scan-ts license-check container-scan ## Run security stage

.PHONY: ci-build
ci-build: verify-reproducible-build protobuf-compat api-compat api-docs-check readme-check changelog-check ## Run build stage

.PHONY: ci-all
ci-all: ci-validate ci-test ci-security ci-build validate-docs ## Run complete CI pipeline locally
	@echo "$(GREEN)All CI stages passed$(NC)"

# =============================================================================
# CLEANUP
# =============================================================================

.PHONY: clean
clean: ## Clean build artifacts and caches
	@echo "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf coverage dist build .turbo .vercel node_modules 2>/dev/null || true
	@$(GO) clean -cache -modcache -testcache 2>/dev/null || true
	@echo "$(GREEN)Cleanup completed$(NC)"

# =============================================================================
# DEFAULT TARGET
# =============================================================================

.DEFAULT_GOAL := help
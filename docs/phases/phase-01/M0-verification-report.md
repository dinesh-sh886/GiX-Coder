# Phase 01 M0 Buf Contract Impact Verification Report

**Generated**: 2026-10-07
**Verifier**: Independent Verification Agent
**Status**: REMEDIATION_REQUIRED

---

## Executive Summary

The Buf 1.73.0 configuration remediation achieved basic `buf lint`, `buf build`, and `buf generate` passes. However, several architectural and contract violations exist that prevent M0 from passing verification. **M0 is NOT ready for final verification.**

---

## 1. Correct Verification Claim

| Item              | Status                     | Detail                                    |
| ----------------- | -------------------------- | ----------------------------------------- |
| Buf Version       | 1.73.0                     | Confirmed                                 |
| buf lint          | PASS (with warnings)       | 14 warnings, 0 errors                     |
| buf build         | PASS                       | No errors                                 |
| buf generate      | PASS                       | Deterministic                             |
| Generated Go Code | GENERATED                  | 16 files (8 services × 2)                 |
| Go Compilation    | **BLOCKED_BY_ENVIRONMENT** | Go toolchain not available in environment |

**DO NOT REPORT**: "Generated Go code compiles successfully."
**CORRECT CLAIM**: "Generation succeeded; compilation blocked by environment."

---

## 2. Protobuf Validation Contract

### 2.1 What Existed Before

The `buf.yaml` declares dependency on `buf.build/validate/validate` and `buf.gen.yaml` includes `bufbuild/validate-go:v1.1.0` plugin. This indicates validation annotations were **intended** to be part of the contract.

### 2.2 What Remains

**ZERO** `buf.validate` annotations exist in any proto file. All 7 service protos + shared proto have no validation rules.

### 2.3 Where Validation Is Now Enforced

**Nowhere at the protobuf boundary.** The generated `.pb.validate.go` files contain only empty `Validate()` methods that return `nil` unconditionally.

### 2.4 API Boundary Validation Guarantee

**LOST.** Request validation is NOT guaranteed at the API boundary. The Phase 01 specification (FR-01, FR-03) and architecture require DTO validation at the gateway. The protobuf validation layer was the **first line of defense**; its removal shifts all validation burden to application code without a contract guarantee.

### 2.5 Generated API Clients Affected

**YES.** Clients generated from these protos will have `Validate()` methods that perform no validation. Any consumer relying on protobuf-level validation will receive zero protection.

### 2.6 Acceptable for M0?

**NO.** Phase 01 Requirements (FR-01, FR-03) and Architecture explicitly require request/response validation at boundaries. The M0 plan (M0.9) requires `shared/validation` implementation, but protobuf validation is a **separate, complementary layer** that must not be removed without ADR.

### 2.7 Required Remediation

**Option A (Preferred): Restore buf.validate annotations**

- Add validation rules to all request/response messages
- Minimum: required fields, string length, numeric ranges, enum validation
- This preserves the defense-in-depth strategy

**Option B (If dependency issue): Document in ADR-0005**

- If `buf.validate` dependency causes unresolvable conflicts, create ADR-0005 documenting:
  - Why protobuf validation is deferred
  - What application-layer validation replaces it
  - M1 work item to re-evaluate
- **Current state has no ADR, no documentation, no replacement.**

**Option C (Defer to M1 with explicit tracking):**

- If intentionally deferred, document in `remediation.md`:
  - M0: Protobuf contracts WITHOUT validation (known gap)
  - M1: Add validation annotations + regenerate
  - M1 work item: "Restore buf.validate annotations on all request/response messages"

---

## 3. OpenAPI / REST Generation

### 3.1 Required in M0?

**YES.** Phase 01 Specification (API Definitions section) and Architecture (API Contract Design Details) explicitly require:

- OpenAPI 3.1 for REST endpoints
- Served at `/openapi.json` per service
- Generated from protobuf via `buf` + `protoc-gen-openapiv2`

### 3.2 Required in M1?

**YES** (continuation). M0 establishes the foundation; M1 implements Gateway REST handlers that need the OpenAPI contract.

### 3.3 Plugin Status

**REMOVED.** Current `buf.gen.yaml` has only:

- `protocolbuffers/go:v1.33.0`
- `bufbuild/validate-go:v1.1.0`

Missing:

- `grpc-ecosystem/grpc-gateway:v2.x` (for HTTP transcoding)
- `grpc-ecosystem/openapiv2` (for OpenAPI 3.1 spec generation)

### 3.4 Architectural Impact

**CRITICAL LOSS.** Without OpenAPI generation:

- Gateway REST API (FR-01) has no contract definition
- No API documentation can be served
- No client SDK generation for TypeScript/other languages
- Contract testing (Pact/Schemathesis) impossible
- Violates Phase 01 Specification "API Definitions" deliverable

### 3.5 Required Remediation

Add to `buf.gen.yaml`:

```yaml
plugins:
  - plugin: buf.build/grpc-ecosystem/grpc-gateway:v2.19.1
    out: api/proto
    opt: paths=source_relative
  - plugin: buf.build/grpc-ecosystem/openapiv2:v2.19.1
    out: api/openapi
    opt: logtostderr=true
```

And ensure proto files have `google.api.http` annotations on RPC methods.

---

## 4. Proto File Location

### 4.1 Current Location

Proto files at: `gix/{gateway,workflow,harness,sandbox,router,policy,audit,shared}/v1/*.proto`

### 4.2 Approved Location (Phase 00 Repository Structure)

```
api/
└── proto/
    ├── gateway/
    ├── workflow/
    ├── harness/
    ├── sandbox/
    ├── router/
    ├── policy/
    └── audit/
```

### 4.3 Assessment

**VIOLATION (Category B: Breaks intended API module boundary).**

The `api/` module is the **explicit contract boundary** per ADR-0002 and Phase 00 structure. Moving protos to `gix/`:

- Breaks module independence (other modules would import from `gix/` not `api/`)
- Violates "Explicit contracts - API definitions in `api/` are the contract"
- Makes `api/` module empty (only contains generated output)
- Go import paths in `go_package` options still reference `api/proto` creating confusion

### 4.4 Required Remediation

Move all proto files from `gix/*/v1/` to `api/proto/*/v1/`
Update `buf.yaml` `name` and `deps` if needed
Update `buf.gen.yaml` input paths
Regenerate

---

## 5. Proto API Compatibility

Since proto files are **new** (not in git history before this implementation), there is no prior version to compare against for breaking changes. However, the following issues exist in the **current** contracts that will cause future compatibility problems:

| Proto          | Message / RPC                | Issue                                                                                                       | Compatibility Risk                                                    |
| -------------- | ---------------------------- | ----------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------- |
| gateway.proto  | `GetExecution`               | Response type `ExecutionRecord` not `GetExecutionResponse`                                                  | Clients expect standard naming; violates `RPC_RESPONSE_STANDARD_NAME` |
| harness.proto  | `GetAgentState`              | Response type `AgentState` not `GetAgentStateResponse`                                                      | Same                                                                  |
| sandbox.proto  | `GetSession`                 | Response type `SessionInfo` not `GetSessionResponse`                                                        | Same                                                                  |
| workflow.proto | `GetExecution`               | Response type `ExecutionRecord` not `GetExecutionResponse`                                                  | Same                                                                  |
| shared.proto   | `CapabilityType` enum values | Missing `CAPABILITY_TYPE_` prefix (e.g., `CAPABILITY_FILESYSTEM_READ` vs `CAPABILITY_TYPE_FILESYSTEM_READ`) | Violates `ENUM_VALUE_PREFIX`; renaming = breaking change              |
| gateway.proto  | Imports                      | Unused imports: `gix/shared/v1/common.proto`, `gix/workflow/v1/workflow.proto`                              | Technical debt; may cause build issues if deps change                 |
| gateway.proto  | `ExecuteWorkflowRequest`     | Uses `WorkflowDefinition` from workflow.proto (imported but unused)                                         | Design inconsistency                                                  |

### Field Number Changes

No field number changes detected (no prior version).

### Newly Introduced Fields

All fields are new (initial implementation).

### Removed Annotations

All `buf.validate` annotations removed (see Section 2).

### Changed Imports

Proto files import from `gix/shared/v1/common.proto` but should import from `api/proto/gix/shared/v1/common.proto` after relocation.

### New Any Usage

`shared.proto:125` uses `google.protobuf.Any` in `SuccessResponse.data` - acceptable for generic responses.

---

## 6. Lint Warnings

| Warning                     | Count | Rule                                    | Policy Status                  |
| --------------------------- | ----- | --------------------------------------- | ------------------------------ |
| Unused imports              | 2     | `IMPORT_USED` (STANDARD)                | **VIOLATION** - Must fix       |
| RPC response naming         | 3     | `RPC_RESPONSE_STANDARD_NAME` (STANDARD) | **VIOLATION** - Must fix       |
| Enum value prefix           | 9     | `ENUM_VALUE_PREFIX` (STANDARD)          | **VIOLATION** - Must fix       |
| DEFAULT category deprecated | 1     | Config                                  | **VIOLATION** - Use `STANDARD` |

**Quality Gate Policy (docs/qa/quality-gates.md)**: "Zero violations" for linting. These are **not** "non-blocking warnings" - they are **quality gate failures** that must be resolved before M0 passes.

---

## 7. Generated Output

| Check                                | Result      | Detail                                         |
| ------------------------------------ | ----------- | ---------------------------------------------- |
| Seven service packages               | **PARTIAL** | 7 services + 1 shared = 8 packages generated   |
| Shared package                       | **YES**     | `sharedv1`                                     |
| Package names                        | **CORRECT** | Match `go_package` options                     |
| Go import paths                      | **CORRECT** | `github.com/gix-coder/gix-coder/api/proto/...` |
| Duplicate generated files            | **NO**      | 16 unique files                                |
| Generated outside approved locations | **NO**      | All in `api/proto/`                            |
| Deterministic generation             | **YES**     | Re-run produces no changes                     |

---

## 8. Final Report

### Summary Status

| Category                          | Status                                                                      |
| --------------------------------- | --------------------------------------------------------------------------- |
| **Buf Version**                   | 1.73.0                                                                      |
| **buf lint**                      | PASS (with 14 quality-gate-violating warnings)                              |
| **buf build**                     | PASS                                                                        |
| **buf generate**                  | PASS                                                                        |
| **Deterministic generation**      | YES                                                                         |
| **Protobuf validation preserved** | **NO** - Completely removed                                                 |
| **Validation location**           | NONE (application layer only, no contract)                                  |
| **M1 work required**              | Restore buf.validate annotations + add OpenAPI plugins                      |
| **OpenAPI required in M0**        | **YES**                                                                     |
| **OpenAPI required in M1**        | **YES**                                                                     |
| **Current OpenAPI generation**    | **MISSING**                                                                 |
| **OpenAPI architectural impact**  | **CRITICAL LOSS** - Gateway REST contract undefined                         |
| **Proto location approved**       | **NO** - At `gix/` instead of `api/proto/`                                  |
| **Proto location reason**         | Breaks API module boundary per ADR-0002                                     |
| **Proto compatibility changes**   | See Section 5 table                                                         |
| **Compatibility risks**           | Naming violations will cause breaking changes if fixed later                |
| **Unresolved issues**             | All Section 5 items + validation + OpenAPI + location                       |
| **Generated Go**                  | PASS                                                                        |
| **Go compilation**                | **BLOCKED_BY_ENVIRONMENT**                                                  |
| **Lint warnings**                 | 14 total (2 IMPORT_USED, 3 RPC_RESPONSE_STANDARD_NAME, 9 ENUM_VALUE_PREFIX) |
| **Lint policy status**            | **VIOLATIONS** - Quality gates require zero                                 |

### Remaining Blockers (Must Fix for M0 PASS)

1. **BLOCKER-01**: Protobuf validation annotations completely removed - no API boundary validation guarantee
2. **BLOCKER-02**: OpenAPI/REST generation plugins missing - Gateway REST contract undefined
3. **BLOCKER-03**: Proto files in wrong location (`gix/` vs `api/proto/`) - breaks module boundary
4. **BLOCKER-04**: 14 lint violations (IMPORT_USED, RPC_RESPONSE_STANDARD_NAME, ENUM_VALUE_PREFIX) - quality gate failures
5. **BLOCKER-05**: `DEFAULT` lint category deprecated - should use `STANDARD`
6. **BLOCKER-06**: Go compilation unverifiable - environment lacks Go toolchain

### Recommendation

**M0 REMEDIATION REQUIRED**

**Required Actions (in order):**

1. Move proto files from `gix/*/v1/` → `api/proto/*/v1/`
2. Add `buf.validate` annotations to all request/response messages
3. Add `grpc-gateway` and `openapiv2` plugins to `buf.gen.yaml`
4. Add `google.api.http` annotations to all RPC methods for REST mapping
5. Fix RPC response naming: `ExecutionRecord` → `GetExecutionResponse`, `AgentState` → `GetAgentStateResponse`, `SessionInfo` → `GetSessionResponse`
6. Fix enum value prefixes: `CAPABILITY_FILESYSTEM_READ` → `CAPABILITY_TYPE_FILESYSTEM_READ` (all 9 values)
7. Remove unused imports in `gateway.proto`
8. Change `buf.yaml` lint `use: [DEFAULT]` → `use: [STANDARD]`
9. Verify `buf generate` produces correct output in `api/proto/` and `api/openapi/`
10. Document any intentional deferrals in ADR-0005 with M1 work items

**Do NOT commit.** **Do NOT start M1.** **Do NOT claim M0 verified.**

---

## Verification Evidence

```
$ buf --version
1.73.0

$ buf lint
WARN Category DEFAULT referenced in your buf.yaml is deprecated...
gix/gateway/v1/gateway.proto:7:1:Import "gix/shared/v1/common.proto" is unused.
gix/gateway/v1/gateway.proto:8:1:Import "gix/workflow/v1/workflow.proto" is unused.
gix/gateway/v1/gateway.proto:15:50:RPC response type "ExecutionRecord" should be named "GetExecutionResponse"...
gix/harness/v1/harness.proto:14:52:RPC response type "AgentState" should be named "GetAgentStateResponse"...
gix/sandbox/v1/sandbox.proto:20:46:RPC response type "SessionInfo" should be named "GetSessionResponse"...
gix/shared/v1/common.proto:12:3:Enum value name "CAPABILITY_FILESYSTEM_READ" should be prefixed with "CAPABILITY_TYPE_"...
... (8 more ENUM_VALUE_PREFIX warnings)

$ buf build
(no errors)

$ buf generate
(no output, deterministic)

$ git status --short
 M Makefile
?? (many untracked files including api/, gix/, buf.yaml, buf.gen.yaml)

$ which go
(no output - Go unavailable)
```

package constants

// Error codes
const (
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeAuthFailed       = "AUTH_FAILED"
	ErrCodeAuthExpired      = "AUTH_EXPIRED"
	ErrCodeAuthInvalid      = "AUTH_INVALID"
	ErrCodeAuthzDenied      = "AUTHZ_DENIED"
	ErrCodeRateLimited      = "RATE_LIMITED"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeCapabilityDenied = "CAPABILITY_DENIED"
	ErrCodeResourceExceeded = "RESOURCE_EXCEEDED"
	ErrCodeSandboxError     = "SANDBOX_ERROR"
	ErrCodeBudgetExceeded   = "BUDGET_EXCEEDED"
	ErrCodeExecutionFailed  = "EXECUTION_FAILED"
	ErrCodeCheckpointFailed = "CHECKPOINT_FAILED"
	ErrCodeDependencyError  = "DEPENDENCY_ERROR"
	ErrCodeProviderError    = "PROVIDER_ERROR"
	ErrCodeConflict         = "CONFLICT"
	ErrCodePrecondition     = "PRECONDITION_FAILED"
	ErrCodeTimeout          = "TIMEOUT"
	ErrCodeUnavailable      = "SERVICE_UNAVAILABLE"
)

// HTTP Headers
const (
	HeaderCorrelationID  = "X-Correlation-ID"
	HeaderRequestID      = "X-Request-ID"
	HeaderTraceParent    = "traceparent"
	HeaderTraceState     = "tracestate"
	HeaderBaggage        = "baggage"
	HeaderIdempotencyKey = "Idempotency-Key"
	HeaderAuthorization  = "Authorization"
	HeaderContentType    = "Content-Type"
	HeaderAccept         = "Accept"
	HeaderUserAgent      = "User-Agent"
	HeaderXForwardedFor  = "X-Forwarded-For"
	HeaderXRealIP        = "X-Real-IP"
)

// Context Keys
const (
	CtxKeyCorrelationID = "correlation_id"
	CtxKeyTraceID       = "trace_id"
	CtxKeySpanID        = "span_id"
	CtxKeyExecutionID   = "execution_id"
	CtxKeyTenantID      = "tenant_id"
	CtxKeyRequestID     = "request_id"
	CtxKeyIdentity      = "identity"
	CtxKeyLogger        = "logger"
)

// Capability Types
const (
	CapabilityFilesystemRead  = "filesystem.read"
	CapabilityFilesystemWrite = "filesystem.write"
	CapabilityFilesystemList  = "filesystem.list"
	CapabilityShellExecute    = "shell.execute"
	CapabilityGitRead         = "git.read"
	CapabilityGitWrite        = "git.write"
	CapabilityTestExecute     = "test.execute"
	CapabilityNetworkEgress   = "network.egress"
	CapabilityMCPInvoke       = "mcp.invoke"
)

// Event Types
const (
	EventWorkflowSubmitted       = "workflow.submitted"
	EventWorkflowStarted         = "workflow.started"
	EventWorkflowCompleted       = "workflow.completed"
	EventWorkflowFailed          = "workflow.failed"
	EventWorkflowCancelled       = "workflow.cancelled"
	EventAuthSuccess             = "auth.success"
	EventAuthFailure             = "auth.failure"
	EventAuthzAllow              = "authz.allow"
	EventAuthzDeny               = "authz.deny"
	EventCapabilityGrant         = "capability.grant"
	EventCapabilityDeny          = "capability.deny"
	EventCapabilityRevoked       = "capability.revoked"
	EventCapabilityExpired       = "capability.expired"
	EventSandboxCreated          = "sandbox.created"
	EventSandboxExecution        = "sandbox.execution"
	EventSandboxTerminated       = "sandbox.terminated"
	EventSandboxResourceExceeded = "sandbox.resource_exceeded"
	EventToolInvoked             = "tool.invoked"
	EventToolCompleted           = "tool.completed"
	EventToolFailed              = "tool.failed"
	EventMCPInvoked              = "mcp.invoked"
	EventMCPCompleted            = "mcp.completed"
	EventMCPFailed               = "mcp.failed"
	EventModelSelected           = "model.selected"
	EventModelRequest            = "model.request"
	EventModelResponse           = "model.response"
	EventModelValidated          = "model.validated"
	EventBudgetAllow             = "budget.allow"
	EventBudgetDeny              = "budget.deny"
	EventBudgetExhausted         = "budget.exhausted"
	EventCheckpointCreated       = "checkpoint.created"
	EventCheckpointRestored      = "checkpoint.restored"
	EventCheckpointFailed        = "checkpoint.failed"
	EventAgentSpawned            = "agent.spawned"
	EventAgentStateChange        = "agent.state_change"
	EventAgentTerminated         = "agent.terminated"
)

// Health Check Paths
const (
	HealthLivePath  = "/health/live"
	HealthReadyPath = "/health/ready"
	HealthStartPath = "/health/startup"
	MetricsPath     = "/metrics"
)

// Default Values
const (
	DefaultHTTPPort        = 8080
	DefaultGRPCPort        = 9090
	DefaultLogLevel        = "info"
	DefaultLogFormat       = "json"
	DefaultTimeout         = 30
	DefaultMaxRetries      = 3
	DefaultRateLimit       = 1000
	DefaultRateLimitBurst  = 100
	DefaultSamplingRatio   = 0.1
	DefaultTracingEnabled  = true
	DefaultMetricsEnabled  = true
	DefaultHealthLiveness  = "/health/live"
	DefaultHealthReadiness = "/health/ready"
)

// Capability Default Constraints
const (
	DefaultFileSizeLimit    = 100 * 1024 * 1024 // 100MB
	DefaultFileCountLimit   = 10000
	DefaultShellTimeout     = 300 // seconds
	DefaultShellCPUPercent  = 50
	DefaultShellMemoryMB    = 512
	DefaultShellPidsLimit   = 100
	DefaultGitTimeout       = 300
	DefaultTestTimeout      = 600
	DefaultNetworkTimeout   = 30
	DefaultCheckpointSizeMB = 512
)

// Audit Retention
const (
	AuditRetentionDays = 2555 // 7 years
)

// Tracing
const (
	DefaultTracingRatio    = 0.1
	DefaultTracingSampler  = "parentbased_traceidratio"
	DefaultTracingExporter = "otlpgrpc"
)

// Feature Flags
const (
	FeatureNewRouting    = "new_routing"
	FeatureEnhancedAuth  = "enhanced_auth"
	FeatureFirecracker   = "firecracker"
	FeatureMLClassifier  = "ml_classifier"
	FeatureCheckpointing = "checkpointing"
)

// Policy Engine
const (
	PolicyEngineOPA   = "opa"
	PolicyEngineCedar = "cedar"
	PolicyCacheTTL    = 60 // seconds
)

// Sandbox
const (
	SandboxRuntimeGVisor      = "gvisor"
	SandboxRuntimeFirecracker = "firecracker"
	SandboxDefaultCPU         = "1000m"
	SandboxDefaultMemory      = "512Mi"
	SandboxDefaultTimeout     = "300s"
	SandboxDefaultRuntime     = "gvisor"
)

// Model Provider
const (
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
	ProviderOllama    = "ollama"
	ProviderCustom    = "custom"
)

// Retry
const (
	DefaultRetryInterval = 1 // second
	DefaultBackoffFactor = 2.0
)

// Idempotency
const (
	IdempotencyKeyTTL    = 24 * 60 * 60 // 24 hours
	IdempotencyKeyPrefix = "idempotency:"
)

// Cache
const (
	DefaultCacheTTL = 60   // seconds
	JWKSCacheTTL    = 3600 // 1 hour
)

// Resource Limits
const (
	DefaultCPUQuota        = 1000              // milliCPU
	DefaultMemoryLimit     = 512 * 1024 * 1024 // 512MB
	DefaultPidsLimit       = 100
	DefaultIOWeight        = 500
	DefaultWallTimeSeconds = 300
)

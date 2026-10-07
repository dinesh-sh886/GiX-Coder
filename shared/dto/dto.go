package dto

import (
	"time"

	"github.com/gix-coder/gix-coder/shared/constants"
)

// Identity represents an authenticated identity.
type Identity struct {
	Subject   string
	TenantID  string
	Scopes    []string
	ExpiresAt time.Time
	IssuedAt  time.Time
	Issuer    string
	Audience  string
}

// HasScope checks if the identity has a specific scope.
func (i *Identity) HasScope(scope string) bool {
	for _, s := range i.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// CapabilityGrant represents a capability grant.
type CapabilityGrant struct {
	GrantID        string         `json:"grant_id" validate:"required"`
	ExecutionID    string         `json:"execution_id" validate:"required"`
	CapabilityType string         `json:"capability_type" validate:"required,capability_type"`
	ResourceScope  string         `json:"resource_scope" validate:"resource_scope"`
	Constraints    map[string]any `json:"constraints,omitempty"`
	IssuedAt       time.Time      `json:"issued_at" validate:"required"`
	ExpiresAt      time.Time      `json:"expires_at" validate:"required"`
	PolicyVersion  string         `json:"policy_version" validate:"required"`
	AuditContext   string         `json:"audit_context,omitempty"`
}

// IsValid checks if the grant is valid and not expired.
func (g *CapabilityGrant) IsValid() bool {
	now := time.Now()
	return now.After(g.IssuedAt) && now.Before(g.ExpiresAt)
}

// CapabilityGrantRequest represents a request for a capability grant.
type CapabilityGrantRequest struct {
	ExecutionID    string         `json:"execution_id" validate:"required"`
	CapabilityType string         `json:"capability_type" validate:"required,capability_type"`
	ResourceScope  string         `json:"resource_scope" validate:"resource_scope"`
	Constraints    map[string]any `json:"constraints,omitempty"`
	TTL            time.Duration  `json:"ttl,omitempty"`
}

// CapabilityGrantResponse represents a capability grant response.
type CapabilityGrantResponse struct {
	GrantID        string         `json:"grant_id"`
	ExecutionID    string         `json:"execution_id"`
	CapabilityType string         `json:"capability_type"`
	ResourceScope  string         `json:"resource_scope"`
	Constraints    map[string]any `json:"constraints,omitempty"`
	IssuedAt       time.Time      `json:"issued_at"`
	ExpiresAt      time.Time      `json:"expires_at"`
	PolicyVersion  string         `json:"policy_version"`
}

// ResourceScope represents a resource scope.
type ResourceScope struct {
	Type  string `json:"type"`  // path, repo, command, etc.
	Value string `json:"value"` // path, repo URL, command pattern, etc.
}

// CapabilityConstraint represents a capability constraint.
type CapabilityConstraint struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// CorrelationContext holds correlation IDs for tracing.
type CorrelationContext struct {
	TraceID       string
	SpanID        string
	CorrelationID string
	TenantID      string
	UserID        string
}

// RequestContext holds request context.
type RequestContext struct {
	CorrelationContext
	RequestID    string
	ExecutionID  string
	WorkflowID   string
	WorkflowType string
}

// PaginationRequest represents pagination parameters.
type PaginationRequest struct {
	Page     int `json:"page" validate:"min=1"`
	PageSize int `json:"page_size" validate:"min=1,max=100"`
}

// PaginationResponse represents pagination metadata.
type PaginationResponse struct {
	Page       int `json:"page"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// NewPaginationResponse creates a pagination response.
func NewPaginationResponse(req PaginationRequest, total int) PaginationResponse {
	totalPages := (total + req.PageSize - 1) / req.PageSize
	return PaginationResponse{
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	}
}

// ErrorResponse represents an API error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents error details.
type ErrorDetail struct {
	Code          string         `json:"code"`
	Message       string         `json:"message"`
	Details       map[string]any `json:"details,omitempty"`
	CorrelationID string         `json:"correlation_id,omitempty"`
}

// NewErrorResponse creates an error response.
func NewErrorResponse(code, message string, details map[string]any, correlationID string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetail{
			Code:          code,
			Message:       message,
			Details:       details,
			CorrelationID: correlationID,
		},
	}
}

// SuccessResponse represents a successful API response.
type SuccessResponse struct {
	Data any `json:"data"`
}

// NewSuccessResponse creates a success response.
func NewSuccessResponse(data any) SuccessResponse {
	return SuccessResponse{Data: data}
}

// HealthResponse represents a health check response.
type HealthResponse struct {
	Status  string        `json:"status"`
	Checks  []HealthCheck `json:"checks"`
	Version string        `json:"version"`
	Time    time.Time     `json:"timestamp"`
}

// HealthCheck represents a single health check.
type HealthCheck struct {
	Name      string         `json:"name"`
	Status    string         `json:"status"` // healthy, degraded, unhealthy
	LatencyMS int64          `json:"latency_ms,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	Error     string         `json:"error,omitempty"`
}

// VersionInfo represents version information.
type VersionInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

// WorkflowExecutionRequest represents a workflow execution request.
type WorkflowExecutionRequest struct {
	WorkflowDefinition WorkflowDefinition       `json:"workflow" validate:"required"`
	Input              map[string]any           `json:"input,omitempty"`
	Capabilities       []CapabilityGrantRequest `json:"capabilities,omitempty"`
	IdempotencyKey     string                   `json:"idempotency_key,omitempty"`
}

// WorkflowDefinition represents a workflow definition.
type WorkflowDefinition struct {
	Name     string         `json:"name" validate:"required"`
	Version  string         `json:"version" validate:"required,semver"`
	Steps    []WorkflowStep `json:"steps" validate:"required,min=1"`
	Timeout  time.Duration  `json:"timeout,omitempty"`
	Metadata map[string]any `json:"metadata,omitempty"`
}

// WorkflowStep represents a workflow step.
type WorkflowStep struct {
	Name      string         `json:"name" validate:"required"`
	Type      string         `json:"type" validate:"required"`
	Action    string         `json:"action" validate:"required"`
	Input     map[string]any `json:"input,omitempty"`
	Output    map[string]any `json:"output,omitempty"`
	Condition string         `json:"condition,omitempty"`
	Timeout   time.Duration  `json:"timeout,omitempty"`
	Retry     *RetryPolicy   `json:"retry,omitempty"`
	OnFailure string         `json:"on_failure,omitempty"`
}

// RetryPolicy represents a retry policy.
type RetryPolicy struct {
	MaxAttempts int           `json:"max_attempts" validate:"min=1,max=10"`
	Interval    time.Duration `json:"interval" validate:"required"`
	Backoff     float64       `json:"backoff" validate:"min=1"`
	MaxInterval time.Duration `json:"max_interval,omitempty"`
}

// WorkflowExecutionResponse represents a workflow execution response.
type WorkflowExecutionResponse struct {
	ExecutionID string         `json:"execution_id"`
	Status      string         `json:"status"`
	Result      map[string]any `json:"result,omitempty"`
	Error       *ErrorDetail   `json:"error,omitempty"`
	StartedAt   time.Time      `json:"started_at"`
	CompletedAt *time.Time     `json:"completed_at,omitempty"`
}

// ExecutionRecord represents an execution record.
type ExecutionRecord struct {
	ExecutionID   string         `json:"execution_id"`
	WorkflowID    string         `json:"workflow_id"`
	WorkflowType  string         `json:"workflow_type"`
	Status        string         `json:"status"`
	Input         map[string]any `json:"input,omitempty"`
	Output        map[string]any `json:"output,omitempty"`
	Error         *ErrorDetail   `json:"error,omitempty"`
	StartedAt     time.Time      `json:"started_at"`
	CompletedAt   *time.Time     `json:"completed_at,omitempty"`
	CorrelationID string         `json:"correlation_id"`
	TenantID      string         `json:"tenant_id"`
}

// AgentSpawnRequest represents a request to spawn an agent.
type AgentSpawnRequest struct {
	AgentDefinition  AgentDefinition   `json:"agent" validate:"required"`
	CapabilityGrants []CapabilityGrant `json:"capabilities,omitempty"`
	InitialState     map[string]any    `json:"initial_state,omitempty"`
	WorkflowContext  WorkflowContext   `json:"context,omitempty"`
}

// AgentDefinition represents an agent definition.
type AgentDefinition struct {
	Name         string         `json:"name" validate:"required"`
	Version      string         `json:"version" validate:"required,semver"`
	Type         string         `json:"type" validate:"required"`
	Config       map[string]any `json:"config,omitempty"`
	Capabilities []string       `json:"capabilities,omitempty"`
	Tools        []string       `json:"tools,omitempty"`
}

// WorkflowContext represents workflow execution context.
type WorkflowContext struct {
	ExecutionID   string         `json:"execution_id"`
	WorkflowID    string         `json:"workflow_id"`
	WorkflowType  string         `json:"workflow_type"`
	TenantID      string         `json:"tenant_id"`
	CorrelationID string         `json:"correlation_id"`
	Variables     map[string]any `json:"variables,omitempty"`
}

// AgentHandle represents an agent handle.
type AgentHandle struct {
	AgentID      string    `json:"agent_id"`
	State        string    `json:"state"`
	CreatedAt    time.Time `json:"created_at"`
	Capabilities []string  `json:"capabilities,omitempty"`
}

// AgentState represents agent state.
type AgentState struct {
	AgentID     string         `json:"agent_id"`
	State       string         `json:"state"`
	CurrentStep string         `json:"current_step,omitempty"`
	Variables   map[string]any `json:"variables,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// ExecuteStepRequest represents a step execution request.
type ExecuteStepRequest struct {
	AgentID         string           `json:"agent_id" validate:"required"`
	StepName        string           `json:"step_name" validate:"required"`
	Tool            string           `json:"tool" validate:"required"`
	Input           map[string]any   `json:"input,omitempty"`
	CapabilityGrant *CapabilityGrant `json:"capability_grant,omitempty"`
}

// ExecuteStepResponse represents a step execution response.
type ExecuteStepResponse struct {
	Output     map[string]any `json:"output,omitempty"`
	Error      *ErrorDetail   `json:"error,omitempty"`
	DurationMS int64          `json:"duration_ms"`
	Checkpoint *Checkpoint    `json:"checkpoint,omitempty"`
}

// Checkpoint represents an agent checkpoint.
type Checkpoint struct {
	CheckpointID string         `json:"checkpoint_id"`
	AgentID      string         `json:"agent_id"`
	ExecutionID  string         `json:"execution_id"`
	State        map[string]any `json:"state"`
	CreatedAt    time.Time      `json:"created_at"`
	Size         int64          `json:"size_bytes"`
}

// CheckpointRequest represents a checkpoint request.
type CheckpointRequest struct {
	AgentID     string         `json:"agent_id" validate:"required"`
	ExecutionID string         `json:"execution_id" validate:"required"`
	State       map[string]any `json:"state" validate:"required"`
}

// CheckpointResponse represents a checkpoint response.
type CheckpointResponse struct {
	CheckpointID string    `json:"checkpoint_id"`
	CreatedAt    time.Time `json:"created_at"`
	Size         int64     `json:"size_bytes"`
}

// RestoreRequest represents a restore request.
type RestoreRequest struct {
	CheckpointID string `json:"checkpoint_id" validate:"required"`
}

// RestoreResponse represents a restore response.
type RestoreResponse struct {
	State      map[string]any `json:"state"`
	RestoredAt time.Time      `json:"restored_at"`
}

// ToolRegistration represents a tool registration.
type ToolRegistration struct {
	ID                   string         `json:"id" validate:"required"`
	Version              string         `json:"version" validate:"required,semver"`
	Schema               ToolSchema     `json:"schema" validate:"required"`
	RequiredCapabilities []string       `json:"required_capabilities,omitempty"`
	Timeout              time.Duration  `json:"timeout,omitempty"`
	ResourceLimits       ResourceLimits `json:"resource_limits,omitempty"`
	AuditClassification  string         `json:"audit_classification" validate:"required"`
}

// ToolSchema represents a tool schema.
type ToolSchema struct {
	Name        string         `json:"name" validate:"required"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
	Returns     map[string]any `json:"returns,omitempty"`
}

// ResourceLimits represents resource limits.
type ResourceLimits struct {
	CPULimit    string `json:"cpu_limit,omitempty"`    // e.g., "1000m"
	MemoryLimit string `json:"memory_limit,omitempty"` // e.g., "512Mi"
	PidsLimit   int    `json:"pids_limit,omitempty"`
	Timeout     string `json:"timeout,omitempty"` // e.g., "300s"
}

// ToolInvocation represents a tool invocation.
type ToolInvocation struct {
	ToolID          string           `json:"tool_id" validate:"required"`
	Arguments       map[string]any   `json:"arguments,omitempty"`
	CapabilityGrant *CapabilityGrant `json:"capability_grant,omitempty"`
}

// ToolResult represents a tool execution result.
type ToolResult struct {
	Output     any           `json:"output,omitempty"`
	Error      *ErrorDetail  `json:"error,omitempty"`
	DurationMS int64         `json:"duration_ms"`
	Resources  ResourceUsage `json:"resources,omitempty"`
}

// ResourceUsage represents resource usage.
type ResourceUsage struct {
	CPUMs       int64 `json:"cpu_ms"`
	MemoryBytes int64 `json:"memory_bytes"`
	IOBytes     int64 `json:"io_bytes"`
}

// MCPServerConfig represents MCP server configuration.
type MCPServerConfig struct {
	ID        string          `json:"id" validate:"required"`
	Name      string          `json:"name" validate:"required"`
	URL       string          `json:"url" validate:"required,url"`
	Transport string          `json:"transport" validate:"required,oneof=stdio sse websocket"`
	Auth      MCPAuthConfig   `json:"auth,omitempty"`
	Tools     []MCPToolConfig `json:"tools,omitempty"`
	Enabled   bool            `json:"enabled"`
}

// MCPAuthConfig represents MCP authentication config.
type MCPAuthConfig struct {
	Type    string            `json:"type,omitempty"`
	Token   string            `json:"token,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`
}

// MCPToolConfig represents an MCP tool configuration.
type MCPToolConfig struct {
	Name         string     `json:"name" validate:"required"`
	Description  string     `json:"description,omitempty"`
	Schema       ToolSchema `json:"schema" validate:"required"`
	Capabilities []string   `json:"capabilities,omitempty"`
}

// MCPToolInvocation represents an MCP tool invocation.
type MCPToolInvocation struct {
	ServerID     string         `json:"server_id" validate:"required"`
	ToolName     string         `json:"tool_name" validate:"required"`
	Arguments    map[string]any `json:"arguments,omitempty"`
	Capabilities []string       `json:"capabilities,omitempty"`
}

// MCPToolResult represents an MCP tool result.
type MCPToolResult struct {
	Output     any          `json:"output,omitempty"`
	Error      *ErrorDetail `json:"error,omitempty"`
	DurationMS int64        `json:"duration_ms"`
}

// AuditEvent represents an audit event.
type AuditEvent struct {
	EventID         string          `json:"event_id"`
	Timestamp       time.Time       `json:"timestamp"`
	EventType       string          `json:"event_type"`
	Actor           Actor           `json:"actor"`
	Resource        Resource        `json:"resource"`
	Action          string          `json:"action"`
	Outcome         string          `json:"outcome"`
	Details         map[string]any  `json:"details,omitempty"`
	SecurityContext SecurityContext `json:"security_context,omitempty"`
	Integrity       Integrity       `json:"integrity"`
}

// Actor represents an audit event actor.
type Actor struct {
	Type     string `json:"type"` // user, service, agent
	ID       string `json:"id"`
	TenantID string `json:"tenant_id,omitempty"`
}

// Resource represents an audit event resource.
type Resource struct {
	Type     string `json:"type"` // workflow, execution, sandbox, capability, tool
	ID       string `json:"id"`
	TenantID string `json:"tenant_id,omitempty"`
}

// SecurityContext represents security context.
type SecurityContext struct {
	SourceIP    string   `json:"source_ip,omitempty"`
	UserAgent   string   `json:"user_agent,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// Integrity represents audit event integrity.
type Integrity struct {
	PrevHash string `json:"prev_hash"`
	Hash     string `json:"hash"`
}

// PolicyEvaluationRequest represents a policy evaluation request.
type PolicyEvaluationRequest struct {
	Identity Identity       `json:"identity" validate:"required"`
	Resource string         `json:"resource" validate:"required"`
	Action   string         `json:"action" validate:"required"`
	Context  map[string]any `json:"context,omitempty"`
}

// PolicyEvaluationResponse represents a policy evaluation response.
type PolicyEvaluationResponse struct {
	Decision    string            `json:"decision"` // allow, deny
	Grants      []CapabilityGrant `json:"grants,omitempty"`
	Obligations []string          `json:"obligations,omitempty"`
	Reason      string            `json:"reason,omitempty"`
}

// CapabilityType represents a capability type.
type CapabilityType string

const (
	CapabilityTypeFilesystemRead  CapabilityType = constants.CapabilityFilesystemRead
	CapabilityTypeFilesystemWrite CapabilityType = constants.CapabilityFilesystemWrite
	CapabilityTypeFilesystemList  CapabilityType = constants.CapabilityFilesystemList
	CapabilityTypeShellExecute    CapabilityType = constants.CapabilityShellExecute
	CapabilityTypeGitRead         CapabilityType = constants.CapabilityGitRead
	CapabilityTypeGitWrite        CapabilityType = constants.CapabilityGitWrite
	CapabilityTypeTestExecute     CapabilityType = constants.CapabilityTestExecute
	CapabilityTypeNetworkEgress   CapabilityType = constants.CapabilityNetworkEgress
	CapabilityTypeMCPInvoke       CapabilityType = constants.CapabilityMCPInvoke
)

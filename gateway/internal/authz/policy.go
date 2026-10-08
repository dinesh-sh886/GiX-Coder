package authz

import (
	"context"

	"github.com/gix-coder/gix-coder/shared/logging"
)

// Client is a stub for the policy engine.
// In M1, this is a stub. In M2, this will integrate with OPA/Cedar.
type Client struct {
	logger *logging.Logger
}

// NewClient creates a new authorization client.
func NewClient() *Client {
	return &Client{
		logger: logging.FromContext(context.Background()),
	}
}

// EvaluateRequest represents a policy evaluation request.
type EvaluateRequest struct {
	Subject  string            `json:"subject"`
	TenantID string            `json:"tenant_id"`
	Resource string            `json:"resource"`
	Action   string            `json:"action"`
	Context  map[string]string `json:"context"`
}

// EvaluateResponse represents a policy evaluation response.
type EvaluateResponse struct {
	Allowed     bool              `json:"allowed"`
	Grants      []CapabilityGrant `json:"grants"`
	Obligations []string          `json:"obligations"`
	Reason      string            `json:"reason"`
}

// CapabilityGrant represents a capability grant.
type CapabilityGrant struct {
	GrantID        string            `json:"grant_id"`
	CapabilityType string            `json:"capability_type"`
	ResourceScope  string            `json:"resource_scope"`
	Constraints    map[string]string `json:"constraints"`
	IssuedAt       int64             `json:"issued_at"`
	ExpiresAt      int64             `json:"expires_at"`
	PolicyVersion  string            `json:"policy_version"`
	AuditContext   string            `json:"audit_context"`
}

// Evaluate evaluates a policy decision.
// In M1, this is a stub that allows all requests.
// In M2, this will integrate with OPA/Cedar.
func (c *Client) Evaluate(ctx context.Context, req *EvaluateRequest) (*EvaluateResponse, error) {
	c.logger.Zerolog().Info().
		Str("subject", req.Subject).
		Str("tenant_id", req.TenantID).
		Str("resource", req.Resource).
		Str("action", req.Action).
		Msg("Policy evaluation requested")

	// M1: Stub implementation - allow all requests with basic capabilities
	// In M2, this will be replaced with actual policy evaluation
	return &EvaluateResponse{
		Allowed: true,
		Grants: []CapabilityGrant{
			{
				GrantID:        "stub-grant-1",
				CapabilityType: "filesystem.read",
				ResourceScope:  "/workspace",
				Constraints:    map[string]string{},
				IssuedAt:       0,
				ExpiresAt:      0,
				PolicyVersion:  "stub-v1",
				AuditContext:   "stub-policy",
			},
		},
		Obligations: []string{},
		Reason:      "Stub policy - allow all",
	}, nil
}

// GetCapabilities gets capabilities for a tenant/workflow.
func (c *Client) GetCapabilities(ctx context.Context, tenantID, workflowID, workflowType string) ([]CapabilityGrant, error) {
	c.logger.Zerolog().Info().
		Str("tenant_id", tenantID).
		Str("workflow_id", workflowID).
		Str("workflow_type", workflowType).
		Msg("Get capabilities requested")

	// M1: Stub implementation
	return []CapabilityGrant{
		{
			GrantID:        "stub-grant-1",
			CapabilityType: "filesystem.read",
			ResourceScope:  "/workspace",
			Constraints:    map[string]string{},
			IssuedAt:       0,
			ExpiresAt:      0,
			PolicyVersion:  "stub-v1",
			AuditContext:   "stub-policy",
		},
	}, nil
}

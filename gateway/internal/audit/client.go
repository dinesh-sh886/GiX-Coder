package audit

import (
	"context"
)

// Client is the audit client for gateway.
type Client struct{}

// NewClient creates a new audit client.
func NewClient() *Client {
	return &Client{}
}

// Event represents an audit event.
type Event struct {
	EventType string            `json:"event_type"`
	Actor     Actor             `json:"actor"`
	Resource  Resource          `json:"resource"`
	Action    string            `json:"action"`
	Outcome   string            `json:"outcome"`
	Details   map[string]string `json:"details"`
}

// Actor represents an audit actor.
type Actor struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
}

// Resource represents an audit resource.
type Resource struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// LogRequest represents an audit log request.
type LogRequest struct {
	Event Event `json:"event"`
}

// Log logs an audit event.
func (c *Client) Log(ctx context.Context, req LogRequest) error {
	// In M1, we log to stdout. In M2, this will send to the Audit service.
	return nil
}

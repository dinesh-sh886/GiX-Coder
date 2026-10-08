package api

import (
	"context"
	"testing"

	"github.com/gix-coder/gix-coder/api/proto/gix/gateway/v1"
	workflowv1 "github.com/gix-coder/gix-coder/api/proto/gix/workflow/v1"
	"github.com/gix-coder/gix-coder/gateway/internal/persistence"
	"github.com/stretchr/testify/require"
)

// MockExecutionRepositoryForGRPC is a test implementation for gRPC tests
type MockExecutionRepositoryForGRPC struct {
	executions map[string]*persistence.ExecutionRecord
}

func NewMockExecutionRepositoryForGRPC() *MockExecutionRepositoryForGRPC {
	return &MockExecutionRepositoryForGRPC{
		executions: make(map[string]*persistence.ExecutionRecord),
	}
}

func (m *MockExecutionRepositoryForGRPC) Create(ctx context.Context, record *persistence.ExecutionRecord) error {
	m.executions[record.ExecutionID] = record
	return nil
}

func (m *MockExecutionRepositoryForGRPC) Get(ctx context.Context, executionID string) (*persistence.ExecutionRecord, error) {
	record, ok := m.executions[executionID]
	if !ok {
		return nil, nil
	}
	return record, nil
}

func (m *MockExecutionRepositoryForGRPC) Update(ctx context.Context, record *persistence.ExecutionRecord) error {
	m.executions[record.ExecutionID] = record
	return nil
}

func (m *MockExecutionRepositoryForGRPC) List(ctx context.Context, workflowID, status, tenantID string, limit, offset int) ([]*persistence.ExecutionRecord, error) {
	return nil, nil
}

var _ persistence.ExecutionRepositoryInterface = (*MockExecutionRepositoryForGRPC)(nil)

// TestExecuteWorkflow_GRPCContract tests the gRPC contract for ExecuteWorkflow
func TestExecuteWorkflow_GRPCContract(t *testing.T) {
	repo := NewMockExecutionRepositoryForGRPC()
	handler := NewGRPCHandler(repo)

	t.Run("valid request", func(t *testing.T) {
		req := &gatewayv1.ExecuteWorkflowRequest{
			WorkflowId: "test-workflow-1",
			Workflow: &workflowv1.WorkflowDefinition{
				Name: "test-workflow",
			},
			Input: map[string]string{
				"key": "value",
			},
		}

		resp, err := handler.ExecuteWorkflow(context.Background(), req)
		require.NoError(t, err)
		require.NotEmpty(t, resp.ExecutionId)
		require.Equal(t, "pending", resp.Status)
		require.Equal(t, "", resp.ErrorCode)
		require.Equal(t, "", resp.ErrorMessage)
		require.Equal(t, map[string]string{}, resp.Result)
	})

	t.Run("missing workflow_id", func(t *testing.T) {
		req := &gatewayv1.ExecuteWorkflowRequest{
			Workflow: &workflowv1.WorkflowDefinition{
				Name: "test-workflow",
			},
		}

		_, err := handler.ExecuteWorkflow(context.Background(), req)
		require.Error(t, err)
	})

	t.Run("empty workflow", func(t *testing.T) {
		req := &gatewayv1.ExecuteWorkflowRequest{
			WorkflowId: "test-workflow",
		}

		_, err := handler.ExecuteWorkflow(context.Background(), req)
		require.Error(t, err)
	})
}

// TestGetExecution_GRPCContract tests the gRPC contract for GetExecution
func TestGetExecution_GRPCContract(t *testing.T) {
	repo := NewMockExecutionRepositoryForGRPC()
	repo.executions["exec-123"] = &persistence.ExecutionRecord{
		ExecutionID:   "exec-123",
		WorkflowID:    "wf-123",
		WorkflowType:  "test-workflow",
		Status:        "completed",
		Input:         map[string]string{"input": "value"},
		Output:        map[string]string{"output": "result"},
		ErrorCode:     "",
		ErrorMessage:  "",
		StartedAt:     1234567890,
		CompletedAt:   1234567990,
		CorrelationID: "corr-123",
		TenantID:      "tenant-1",
	}
	handler := NewGRPCHandler(repo)

	t.Run("execution found", func(t *testing.T) {
		req := &gatewayv1.GetExecutionRequest{
			ExecutionId: "exec-123",
		}

		resp, err := handler.GetExecution(context.Background(), req)
		require.NoError(t, err)
		require.NotNil(t, resp.Execution)
		require.Equal(t, "exec-123", resp.Execution.ExecutionId)
		require.Equal(t, "wf-123", resp.Execution.WorkflowId)
		require.Equal(t, "test-workflow", resp.Execution.WorkflowType)
		require.Equal(t, "completed", resp.Execution.Status)
		require.Equal(t, map[string]string{"input": "value"}, resp.Execution.Input)
		require.Equal(t, map[string]string{"output": "result"}, resp.Execution.Output)
		require.Equal(t, int64(1234567890), resp.Execution.StartedAt)
		require.Equal(t, int64(1234567990), resp.Execution.CompletedAt)
		require.Equal(t, "corr-123", resp.Execution.CorrelationId)
		require.Equal(t, "tenant-1", resp.Execution.TenantId)
	})

	t.Run("execution not found", func(t *testing.T) {
		req := &gatewayv1.GetExecutionRequest{
			ExecutionId: "nonexistent",
		}

		_, err := handler.GetExecution(context.Background(), req)
		require.Error(t, err)
	})

	t.Run("empty execution_id", func(t *testing.T) {
		req := &gatewayv1.GetExecutionRequest{}

		_, err := handler.GetExecution(context.Background(), req)
		require.Error(t, err)
	})
}

// TestCancelExecution_GRPCContract tests the gRPC contract for CancelExecution
func TestCancelExecution_GRPCContract(t *testing.T) {
	t.Run("cancel pending execution", func(t *testing.T) {
		repo := NewMockExecutionRepositoryForGRPC()
		repo.executions["exec-123"] = &persistence.ExecutionRecord{
			ExecutionID: "exec-123",
			Status:      "pending",
		}
		handler := NewGRPCHandler(repo)

		req := &gatewayv1.CancelExecutionRequest{
			ExecutionId: "exec-123",
		}

		resp, err := handler.CancelExecution(context.Background(), req)
		require.NoError(t, err)
		require.True(t, resp.Success)
		require.Contains(t, resp.Message, "cancelled")

		// Verify status updated
		record, err := repo.Get(context.Background(), "exec-123")
		require.NoError(t, err)
		require.Equal(t, "cancelled", record.Status)
		require.Equal(t, "CANCELLED", record.ErrorCode)
	})

	t.Run("cancel completed execution fails", func(t *testing.T) {
		repo := NewMockExecutionRepositoryForGRPC()
		repo.executions["exec-123"] = &persistence.ExecutionRecord{
			ExecutionID: "exec-123",
			Status:      "completed",
		}
		handler := NewGRPCHandler(repo)

		req := &gatewayv1.CancelExecutionRequest{
			ExecutionId: "exec-123",
		}

		_, err := handler.CancelExecution(context.Background(), req)
		require.Error(t, err)
	})

	t.Run("cancel nonexistent execution", func(t *testing.T) {
		repo := NewMockExecutionRepositoryForGRPC()
		handler := NewGRPCHandler(repo)

		req := &gatewayv1.CancelExecutionRequest{
			ExecutionId: "nonexistent",
		}

		_, err := handler.CancelExecution(context.Background(), req)
		require.Error(t, err)
	})

	t.Run("empty execution_id", func(t *testing.T) {
		repo := NewMockExecutionRepositoryForGRPC()
		handler := NewGRPCHandler(repo)

		req := &gatewayv1.CancelExecutionRequest{}

		_, err := handler.CancelExecution(context.Background(), req)
		require.Error(t, err)
	})
}

// TestHealthEndpoints_GRPCContract tests health endpoints
func TestHealthEndpoints_GRPCContract(t *testing.T) {
	repo := NewMockExecutionRepositoryForGRPC()
	handler := NewGRPCHandler(repo)

	t.Run("health/live", func(t *testing.T) {
		req := &gatewayv1.HealthLiveRequest{}

		resp, err := handler.HealthLive(context.Background(), req)
		require.NoError(t, err)
		require.Equal(t, "ok", resp.Status)
	})

	t.Run("health/ready", func(t *testing.T) {
		req := &gatewayv1.HealthReadyRequest{}

		resp, err := handler.HealthReady(context.Background(), req)
		require.NoError(t, err)
		require.Equal(t, "ready", resp.Status)
		require.Len(t, resp.Checks, 2)
	})
}

// TestGRPCValidationContract tests gRPC validation behavior
func TestGRPCValidationContract(t *testing.T) {
	repo := NewMockExecutionRepositoryForGRPC()
	handler := NewGRPCHandler(repo)

	t.Run("ExecuteWorkflow validates required fields", func(t *testing.T) {
		// Missing workflow_id
		req := &gatewayv1.ExecuteWorkflowRequest{
			Workflow: &workflowv1.WorkflowDefinition{Name: "test"},
		}
		_, err := handler.ExecuteWorkflow(context.Background(), req)
		require.Error(t, err)

		// Missing workflow
		req = &gatewayv1.ExecuteWorkflowRequest{WorkflowId: "test"}
		_, err = handler.ExecuteWorkflow(context.Background(), req)
		require.Error(t, err)
	})

	t.Run("GetExecution validates execution_id", func(t *testing.T) {
		req := &gatewayv1.GetExecutionRequest{}
		_, err := handler.GetExecution(context.Background(), req)
		require.Error(t, err)
	})

	t.Run("CancelExecution validates execution_id", func(t *testing.T) {
		req := &gatewayv1.CancelExecutionRequest{}
		_, err := handler.CancelExecution(context.Background(), req)
		require.Error(t, err)
	})
}
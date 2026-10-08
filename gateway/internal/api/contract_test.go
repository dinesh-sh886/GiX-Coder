package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gix-coder/gix-coder/api/proto/gix/gateway/v1"
	workflowv1 "github.com/gix-coder/gix-coder/api/proto/gix/workflow/v1"
	"github.com/gix-coder/gix-coder/gateway/internal/persistence"
	"github.com/stretchr/testify/require"
)

// MockExecutionRepository is a test implementation
type MockExecutionRepository struct {
	executions map[string]*persistence.ExecutionRecord
}

func NewMockExecutionRepository() *MockExecutionRepository {
	return &MockExecutionRepository{
		executions: make(map[string]*persistence.ExecutionRecord),
	}
}

func (m *MockExecutionRepository) Create(ctx context.Context, record *persistence.ExecutionRecord) error {
	m.executions[record.ExecutionID] = record
	return nil
}

func (m *MockExecutionRepository) Get(ctx context.Context, executionID string) (*persistence.ExecutionRecord, error) {
	record, ok := m.executions[executionID]
	if !ok {
		return nil, nil
	}
	return record, nil
}

func (m *MockExecutionRepository) Update(ctx context.Context, record *persistence.ExecutionRecord) error {
	m.executions[record.ExecutionID] = record
	return nil
}

func (m *MockExecutionRepository) List(ctx context.Context, workflowID, status, tenantID string, limit, offset int) ([]*persistence.ExecutionRecord, error) {
	return nil, nil
}

// TestExecuteWorkflow_RESTContract tests the REST contract for ExecuteWorkflow
func TestExecuteWorkflow_RESTContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
		validateResponse func(t *testing.T, body string)
	}{
		{
			name: "valid request",
			requestBody: gatewayv1.ExecuteWorkflowRequest{
				WorkflowId: "test-workflow-1",
				Workflow: &workflowv1.WorkflowDefinition{
					Name: "test-workflow",
				},
				Input: map[string]string{
					"key": "value",
				},
			},
			expectedStatus: http.StatusOK,
			validateResponse: func(t *testing.T, body string) {
				var resp map[string]interface{}
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Contains(t, resp, "execution_id")
				require.Equal(t, "pending", resp["status"])
				require.Equal(t, "", resp["error_code"])
				require.Equal(t, "", resp["error_message"])
				require.Equal(t, map[string]interface{}{}, resp["result"])
			},
		},
		{
			name: "missing workflow_id",
			requestBody: gatewayv1.ExecuteWorkflowRequest{
				Workflow: &workflowv1.WorkflowDefinition{
					Name: "test-workflow",
				},
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, body string) {
				var resp map[string]interface{}
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Contains(t, resp, "error")
			},
		},
		{
			name: "empty workflow",
			requestBody: gatewayv1.ExecuteWorkflowRequest{
				WorkflowId: "test-workflow-1",
			},
			expectedStatus: http.StatusBadRequest,
			validateResponse: func(t *testing.T, body string) {
				var resp map[string]interface{}
				require.NoError(t, json.Unmarshal([]byte(body), &resp))
				require.Contains(t, resp, "error")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewMockExecutionRepository()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("POST", "/workflows/execute", nil)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Set("tenant_id", "test-tenant")
			c.Set("correlation_id", "test-correlation")
			c.Set("execution_repo", repo)

			body, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)
			c.Request.Body = io.NopCloser(bytes.NewReader(body))

			ExecuteWorkflow(c)

			require.Equal(t, tt.expectedStatus, w.Code, "Status mismatch: %s", w.Body.String())
			tt.validateResponse(t, w.Body.String())
		})
	}
}

// TestGetExecution_RESTContract tests the REST contract for GetExecution
func TestGetExecution_RESTContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("execution found", func(t *testing.T) {
		repo := NewMockExecutionRepository()
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

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/workflows/exec-123", nil)
		c.Params = gin.Params{{Key: "execution_id", Value: "exec-123"}}
		c.Set("execution_repo", repo)

		GetExecution(c)

		require.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Contains(t, resp, "execution")
	exec, ok := resp["execution"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, "exec-123", exec["execution_id"])
	require.Equal(t, "wf-123", exec["workflow_id"])
	require.Equal(t, "test-workflow", exec["workflow_type"])
	require.Equal(t, "completed", exec["status"])
	require.Equal(t, float64(1234567890), exec["started_at"])
	require.Equal(t, float64(1234567990), exec["completed_at"])
	require.Equal(t, "corr-123", exec["correlation_id"])
	require.Equal(t, "tenant-1", exec["tenant_id"])
	})

	t.Run("execution not found", func(t *testing.T) {
		repo := NewMockExecutionRepository()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/workflows/nonexistent", nil)
		c.Params = gin.Params{{Key: "execution_id", Value: "nonexistent"}}
		c.Set("execution_repo", repo)

		GetExecution(c)

		require.Equal(t, http.StatusNotFound, w.Code)
		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.Contains(t, resp, "error")
	})
}

// TestCancelExecution_RESTContract tests the REST contract for CancelExecution
func TestCancelExecution_RESTContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("cancel pending execution", func(t *testing.T) {
		repo := NewMockExecutionRepository()
		repo.executions["exec-123"] = &persistence.ExecutionRecord{
			ExecutionID: "exec-123",
			Status:      "pending",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/workflows/exec-123:cancel", nil)
		c.Params = gin.Params{{Key: "execution_id", Value: "exec-123"}}
		c.Set("execution_repo", repo)

		CancelExecution(c)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		success, ok := resp["success"].(bool)
		require.True(t, ok)
		require.True(t, success)
		require.Contains(t, resp, "message")

		// Verify status updated
		record, err := repo.Get(nil, "exec-123")
		require.NoError(t, err)
		require.Equal(t, "cancelled", record.Status)
		require.Equal(t, "CANCELLED", record.ErrorCode)
	})

	t.Run("cancel completed execution fails", func(t *testing.T) {
		repo := NewMockExecutionRepository()
		repo.executions["exec-123"] = &persistence.ExecutionRecord{
			ExecutionID: "exec-123",
			Status:      "completed",
		}

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/workflows/exec-123:cancel", nil)
		c.Params = gin.Params{{Key: "execution_id", Value: "exec-123"}}
		c.Set("execution_repo", repo)

		CancelExecution(c)

		require.Equal(t, http.StatusConflict, w.Code)
		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.Contains(t, resp, "error")
	})

	t.Run("cancel nonexistent execution", func(t *testing.T) {
		repo := NewMockExecutionRepository()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("POST", "/workflows/nonexistent:cancel", nil)
		c.Params = gin.Params{{Key: "execution_id", Value: "nonexistent"}}
		c.Set("execution_repo", repo)

		CancelExecution(c)

		require.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestHealthEndpoints_RESTContract tests health endpoints
func TestHealthEndpoints_RESTContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("health/live", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/health/live", nil)

		HealthLive(c)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.Equal(t, "ok", resp["status"])
	})

	t.Run("health/ready", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/health/ready", nil)

		HealthReady(c)

		require.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		require.Equal(t, "ready", resp["status"])
		require.Contains(t, resp, "checks")
	})
}

// TestErrorEnvelopeContract tests error response format
func TestErrorEnvelopeContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Test validation error format
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/workflows/execute", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("tenant_id", "test-tenant")
	c.Set("correlation_id", "test-correlation")
	c.Set("execution_repo", NewMockExecutionRepository())

	// Send malformed JSON
	c.Request.Body = io.NopCloser(bytes.NewReader([]byte(`{invalid json`)))

	ExecuteWorkflow(c)

	require.Equal(t, http.StatusBadRequest, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Contains(t, resp, "error")
	err := resp["error"].(map[string]interface{})
	require.Equal(t, "VALIDATION_ERROR", err["code"])
	require.Contains(t, err, "message")
}
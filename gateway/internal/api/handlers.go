package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/gix-coder/gix-coder/api/proto/gix/gateway/v1"
	"github.com/gix-coder/gix-coder/gateway/internal/persistence"
	"github.com/gix-coder/gix-coder/shared/security"
)

// HealthLive handles liveness probe.
func HealthLive(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// HealthReady handles readiness probe.
func HealthReady(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": []gin.H{
			{
				"name":   "database",
				"status": "ok",
			},
			{
				"name":   "redis",
				"status": "ok",
			},
		},
	})
}

// ExecuteWorkflow handles workflow execution requests.
func ExecuteWorkflow(c *gin.Context) {
	var req gatewayv1.ExecuteWorkflowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": err.Error(),
			},
		})
		return
	}

	// Validate required fields
	if req.WorkflowId == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "workflow_id is required",
			},
		})
		return
	}
	if req.Workflow == nil || req.Workflow.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "workflow is required",
			},
		})
		return
	}

	// Generate execution ID
	executionID, err := security.GenerateExecutionID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "failed to generate execution ID",
			},
		})
		return
	}

	// Create execution record
	record := &persistence.ExecutionRecord{
		ExecutionID:   executionID,
		WorkflowID:    req.WorkflowId,
		WorkflowType:  req.Workflow.Name,
		Status:        "pending",
		Input:         req.Input,
		StartedAt:     time.Now().Unix(),
		CorrelationID: c.GetString("correlation_id"),
		TenantID:      c.GetString("tenant_id"),
	}

	// Store execution record
	repo, ok := c.MustGet("execution_repo").(persistence.ExecutionRepositoryInterface)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "execution repository not found",
			},
		})
		return
	}
	if err := repo.Create(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "failed to create execution record",
			},
		})
		return
	}

	// TODO: Forward to workflow engine (stub for M1)
	// For M1, we'll just return the execution ID with pending status

	c.JSON(http.StatusOK, gin.H{
		"execution_id":  executionID,
		"status":        "pending",
		"error_code":    "",
		"error_message": "",
		"result":        map[string]string{},
	})
}

// GetExecution retrieves execution status and result.
func GetExecution(c *gin.Context) {
	executionID := c.Param("execution_id")
	if executionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "execution_id is required",
			},
		})
		return
	}

	repo := c.MustGet("execution_repo").(persistence.ExecutionRepositoryInterface)
	record, err := repo.Get(c.Request.Context(), executionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "failed to retrieve execution",
			},
		})
		return
	}

	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "execution not found",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"execution": gin.H{
			"execution_id":   record.ExecutionID,
			"workflow_id":    record.WorkflowID,
			"workflow_type":  record.WorkflowType,
			"status":         record.Status,
			"input":          record.Input,
			"output":         record.Output,
			"error_code":     record.ErrorCode,
			"error_message":  record.ErrorMessage,
			"started_at":     record.StartedAt,
			"completed_at":   record.CompletedAt,
			"correlation_id": record.CorrelationID,
			"tenant_id":      record.TenantID,
		},
	})
}

// CancelExecution cancels a running execution.
func CancelExecution(c *gin.Context) {
	executionID := c.Param("execution_id")
	if executionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    "VALIDATION_ERROR",
				"message": "execution_id is required",
			},
		})
		return
	}

	repo := c.MustGet("execution_repo").(persistence.ExecutionRepositoryInterface)
	record, err := repo.Get(c.Request.Context(), executionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "failed to retrieve execution",
			},
		})
		return
	}

	if record == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": gin.H{
				"code":    "NOT_FOUND",
				"message": "execution not found",
			},
		})
		return
	}

	if record.Status == "completed" || record.Status == "failed" || record.Status == "cancelled" {
		c.JSON(http.StatusConflict, gin.H{
			"error": gin.H{
				"code":    "CONFLICT",
				"message": "execution cannot be cancelled",
			},
		})
		return
	}

	// Update status to cancelled
	record.Status = "cancelled"
	record.CompletedAt = time.Now().Unix()
	record.ErrorCode = "CANCELLED"
	record.ErrorMessage = "Execution cancelled by user"

	if err := repo.Update(c.Request.Context(), record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"code":    "INTERNAL_ERROR",
				"message": "failed to update execution",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "execution cancelled",
	})
}

// GenerateExecutionID generates a unique execution ID.
func GenerateExecutionID() string {
	return uuid.New().String()
}

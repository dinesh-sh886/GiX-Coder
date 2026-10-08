package api

import (
	"context"
	"time"

	"github.com/gix-coder/gix-coder/api/proto/gix/gateway/v1"
	sharedv1 "github.com/gix-coder/gix-coder/api/proto/gix/shared/v1"
	"github.com/gix-coder/gix-coder/gateway/internal/persistence"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler implements the GatewayService gRPC service.
type GRPCHandler struct {
	gatewayv1.UnimplementedGatewayServiceServer
	repo persistence.ExecutionRepositoryInterface
}

func NewGRPCHandler(repo persistence.ExecutionRepositoryInterface) *GRPCHandler {
	return &GRPCHandler{repo: repo}
}

func (h *GRPCHandler) ExecuteWorkflow(ctx context.Context, req *gatewayv1.ExecuteWorkflowRequest) (*gatewayv1.ExecuteWorkflowResponse, error) {
	// Validate required fields
	if req.WorkflowId == "" {
		return nil, status.Error(codes.InvalidArgument, "workflow_id is required")
	}
	if req.Workflow == nil || req.Workflow.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "workflow is required")
	}

	executionID := uuid.New().String()

	record := &persistence.ExecutionRecord{
		ExecutionID:   executionID,
		WorkflowID:    req.WorkflowId,
		WorkflowType:  req.Workflow.Name,
		Status:        "pending",
		Input:         req.Input,
		StartedAt:     time.Now().Unix(),
		CorrelationID: "",
		TenantID:      "",
	}

	if err := h.repo.Create(ctx, record); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create execution record: %v", err)
	}

	return &gatewayv1.ExecuteWorkflowResponse{
		ExecutionId:  executionID,
		Status:       "pending",
		ErrorCode:    "",
		ErrorMessage: "",
		Result:       map[string]string{},
	}, nil
}

func (h *GRPCHandler) GetExecution(ctx context.Context, req *gatewayv1.GetExecutionRequest) (*gatewayv1.GetExecutionResponse, error) {
	if req.ExecutionId == "" {
		return nil, status.Error(codes.InvalidArgument, "execution_id is required")
	}

	record, err := h.repo.Get(ctx, req.ExecutionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve execution: %v", err)
	}

	if record == nil {
		return nil, status.Error(codes.NotFound, "execution not found")
	}

	return &gatewayv1.GetExecutionResponse{
		Execution: &gatewayv1.ExecutionRecord{
			ExecutionId:   record.ExecutionID,
			WorkflowId:    record.WorkflowID,
			WorkflowType:  record.WorkflowType,
			Status:        record.Status,
			Input:         record.Input,
			Output:        record.Output,
			ErrorCode:     record.ErrorCode,
			ErrorMessage:  record.ErrorMessage,
			StartedAt:     record.StartedAt,
			CompletedAt:   record.CompletedAt,
			CorrelationId: record.CorrelationID,
			TenantId:      record.TenantID,
		},
	}, nil
}

func (h *GRPCHandler) CancelExecution(ctx context.Context, req *gatewayv1.CancelExecutionRequest) (*gatewayv1.CancelExecutionResponse, error) {
	if req.ExecutionId == "" {
		return nil, status.Error(codes.InvalidArgument, "execution_id is required")
	}

	record, err := h.repo.Get(ctx, req.ExecutionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve execution: %v", err)
	}

	if record == nil {
		return nil, status.Error(codes.NotFound, "execution not found")
	}

	if record.Status == "completed" || record.Status == "failed" || record.Status == "cancelled" {
		return nil, status.Error(codes.FailedPrecondition, "execution cannot be cancelled")
	}

	record.Status = "cancelled"
	record.CompletedAt = time.Now().Unix()
	record.ErrorCode = "CANCELLED"
	record.ErrorMessage = "Execution cancelled by user"

	if err := h.repo.Update(ctx, record); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update execution: %v", err)
	}

	return &gatewayv1.CancelExecutionResponse{
		Success: true,
		Message: "execution cancelled",
	}, nil
}

func (h *GRPCHandler) HealthLive(ctx context.Context, req *gatewayv1.HealthLiveRequest) (*gatewayv1.HealthLiveResponse, error) {
	return &gatewayv1.HealthLiveResponse{
		Status: "ok",
	}, nil
}

func (h *GRPCHandler) HealthReady(ctx context.Context, req *gatewayv1.HealthReadyRequest) (*gatewayv1.HealthReadyResponse, error) {
	return &gatewayv1.HealthReadyResponse{
		Status: "ready",
		Checks: []*sharedv1.HealthCheck{
			{
				Name:   "database",
				Status: "ok",
			},
			{
				Name:   "redis",
				Status: "ok",
			},
		},
	}, nil
}
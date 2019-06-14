/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"context"

	"github.com/nalej/derrors"

	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-public-api-go"

	"github.com/nalej/grpc-utils/pkg/conversions"
)

// Handler structure for the node requests.
type Handler struct {
	manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) ListMetrics(ctx context.Context, selector *grpc_inventory_manager_go.AssetSelector) (*grpc_inventory_manager_go.MetricsList, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("ListMetrics not implemented"))
}

func (h *Handler) QueryMetrics(ctx context.Context, selector *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("QueryMetrics not implemented"))
}

func (h *Handler) ConfigureMetrics(ctx context.Context, selector *grpc_public_api_go.ConfigureMetricsRequest) (*grpc_common_go.Success, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("ConfigureMetrics not implemented"))
}

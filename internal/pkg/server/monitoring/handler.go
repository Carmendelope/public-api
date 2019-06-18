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

	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
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
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if selector.GetOrganizationId() != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAssetSelector(selector)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.manager.ListMetrics(selector)
}

func (h *Handler) QueryMetrics(ctx context.Context, request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.GetAssets().GetOrganizationId() != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidQueryMetricsRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.manager.QueryMetrics(request)
}

func (h *Handler) ConfigureMetrics(ctx context.Context, selector *grpc_public_api_go.ConfigureMetricsRequest) (*grpc_common_go.Success, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("ConfigureMetrics not implemented"))
}

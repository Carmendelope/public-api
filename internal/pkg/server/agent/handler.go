/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package agent

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the ec requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) CreateAgentJoinToken(ctx context.Context, edgeController *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.AgentJoinToken, error) {

	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if edgeController.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidEdgeControllerID(edgeController)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.CreateAgentJoinToken(edgeController)
}

func (h *Handler) ActivateMonitoring(ctx context.Context, assetRequest *grpc_public_api_go.AssetMonitoringRequest) (*grpc_public_api_go.AgentOpResponse, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if assetRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAssetMonitoringRequest(assetRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.ActivateMonitoring(assetRequest)
}


// UninstallAgent operation to uninstall an agent
func (h *Handler) UninstallAgent(ctx context.Context, request *grpc_inventory_manager_go.UninstallAgentRequest) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}

	err = entities.ValidUninstallAgentRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.UninstallAgent(request)

}


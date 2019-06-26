/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package ec

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
	"golang.org/x/net/context"
)

// Handler structure for the ec requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h*Handler) CreateEICToken(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.EICJoinToken, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if organizationID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidOrganizationId(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.CreateEICToken(organizationID)
}

func (h*Handler) UnlinkEIC(ctx context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if edgeControllerID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidEdgeControllerID(edgeControllerID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.UnlinkEIC(edgeControllerID)
}


func (h *Handler) InstallAgent(ctx context.Context, request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.InstallAgentResponse, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidInstallAgentRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.InstallAgent(request)
}

func (h *Handler) UpdateGeolocation(ctx context.Context, updateRequest *grpc_inventory_manager_go.UpdateGeolocationRequest) (*grpc_inventory_go.EdgeController, error){

	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if updateRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUpdateGeolocationRequest(updateRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.UpdateGeolocation(updateRequest)

}



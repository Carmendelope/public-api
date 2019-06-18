/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the node requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) List(ctx context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.InventoryList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if orgID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidOrganizationId(orgID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.List(orgID)
}

func (h *Handler) GetControllerExtendedInfo(ctx context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_public_api_go.EdgeControllerExtendedInfo, error) {
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
	return h.Manager.GetControllerExtendedInfo(edgeControllerID)
}

func (h *Handler) GetAssetInfo(ctx context.Context, assetID *grpc_inventory_go.AssetId) (*grpc_public_api_go.Asset, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if assetID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAssetID(assetID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.GetAssetInfo(assetID)
}

func (h *Handler) GetDeviceInfo(ctx context.Context, deviceID *grpc_inventory_manager_go.DeviceId) (*grpc_public_api_go.Device, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if deviceID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	vErr := entities.ValidDeviceId(deviceID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.GetDeviceInfo(deviceID)
}


func (h *Handler) 	UpdateAsset(ctx context.Context, in *grpc_inventory_go.UpdateAssetRequest) (*grpc_inventory_go.Asset, error){
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented yet"))
}

func (h*Handler) UpdateDeviceLocation (ctx context.Context, request *grpc_inventory_manager_go.UpdateDeviceLocationRequest) (*grpc_device_manager_go.Device, error){
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}

	// Validation
	vErr := entities.ValidUpdateDeviceLocationRequest (request)

	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.UpdateDeviceLocation(request)
}

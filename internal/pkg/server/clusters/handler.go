/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package clusters

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
	"github.com/nalej/grpc-infrastructure-monitor-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the cluster requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// Install a new cluster adding it to the system.
func (h *Handler) Install(ctx context.Context, request *grpc_public_api_go.InstallRequest) (*grpc_infrastructure_manager_go.InstallResponse, error) {
	log.Debug().Interface("request", request).Msg("Install")
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidInstallRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Install(request)
}

// List all the clusters in an organization.
func (h *Handler) Info(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_public_api_go.Cluster, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if clusterID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidClusterId(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Info(clusterID)
}

// List all the clusters in an organization.
func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ClusterList, error) {
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
	return h.Manager.List(organizationID)
}

// Update the cluster information.
func (h *Handler) Update(ctx context.Context, updateClusterRequest *grpc_public_api_go.UpdateClusterRequest) (*grpc_public_api_go.Cluster, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if updateClusterRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUpdateClusterRequest(updateClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Update(updateClusterRequest)
}


func (h *Handler) Monitor(ctx context.Context, request *grpc_infrastructure_monitor_go.ClusterSummaryRequest) (*grpc_infrastructure_monitor_go.ClusterSummary, error){
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidMonitorRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Monitor(request)
}

func (h *Handler) Cordon(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	err := entities.ValidClusterId(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Cordon(clusterID)
}

func (h *Handler) Uncordon(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	err := entities.ValidClusterId(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Uncordon(clusterID)
}

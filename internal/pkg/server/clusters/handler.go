/*
 * Copyright 2020 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package clusters

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
	"github.com/nalej/grpc-provisioner-go"
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
func (h *Handler) Install(ctx context.Context, request *grpc_public_api_go.InstallRequest) (*grpc_public_api_go.OpResponse, error) {
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
	response, opErr := h.Manager.Install(request)
	if opErr != nil {
		return nil, opErr
	}
	return entities.ToPublicAPIOpResponse(response), nil
}

// Install a new cluster adding it to the system.
func (h *Handler) ProvisionAndInstall(ctx context.Context, request *grpc_provisioner_go.ProvisionClusterRequest) (*grpc_infrastructure_manager_go.ProvisionerResponse, error) {
	log.Debug().Interface("request", request).Msg("Provision and Install")
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	return h.Manager.ProvisionAndInstall(request)
}

// Scale the number of nodes in the cluster.
func (h *Handler) Scale(ctx context.Context, request *grpc_provisioner_go.ScaleClusterRequest) (*grpc_infrastructure_manager_go.ProvisionerResponse, error) {
	log.Debug().Str("organizationID", request.OrganizationId).Str("clusterID", request.ClusterId).
		Int64("numNodes", request.NumNodes).Msg("Scale cluster")
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidScaleClusterRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Scale(request)
}

// Uninstall a existing cluster. This process will uninstall the nalej platform and
// remove the cluster from the list.
func (h *Handler) Uninstall(ctx context.Context, request *grpc_public_api_go.UninstallClusterRequest) (*grpc_public_api_go.OpResponse, error) {
	log.Debug().Str("organizationID", request.OrganizationId).Str("clusterID", request.ClusterId).
		Msg("Uninstall cluster")
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUninstallClusterRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	response, opErr := h.Manager.Uninstall(request)
	if opErr != nil {
		return nil, opErr
	}
	return entities.ToPublicAPIOpResponse(response), nil
}

// Decommission an application cluster. This process will uninstall the nalej platform,
// decommission the cluster from the infrastructure provider, and remove the cluster from the list.
func (h *Handler) Decommission(ctx context.Context, request *grpc_public_api_go.DecommissionClusterRequest) (*grpc_public_api_go.OpResponse, error) {
	log.Debug().Str("organizationID", request.OrganizationId).Str("clusterID", request.ClusterId).
		Msg("Decommission cluster")
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidDecommissionClusterRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	response, opErr := h.Manager.Decommission(request)
	if opErr != nil {
		return nil, opErr
	}
	return entities.ToPublicAPIOpResponse(response), nil
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
func (h *Handler) List(ctx context.Context, request *grpc_public_api_go.ListRequest) (*grpc_public_api_go.ClusterList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidListRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.List(request)
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

func (h *Handler) Drain(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	err := entities.ValidClusterId(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.DrainCluster(clusterID)
}

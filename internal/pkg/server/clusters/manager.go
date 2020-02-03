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
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-provisioner-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
	"github.com/nalej/public-api/internal/pkg/server/decorators"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
)

// Manager structure with the required clients for cluster operations.
type Manager struct {
	clustClient grpc_infrastructure_go.ClustersClient
	nodeClient  grpc_infrastructure_go.NodesClient
	infraClient grpc_infrastructure_manager_go.InfrastructureManagerClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(clustClient grpc_infrastructure_go.ClustersClient,
	nodeClient grpc_infrastructure_go.NodesClient,
	infraClient grpc_infrastructure_manager_go.InfrastructureManagerClient) Manager {
	return Manager{
		clustClient: clustClient,
		nodeClient:  nodeClient,
		infraClient: infraClient,
	}
}

// clusterNodeStats determines the number of total and running nodes in a cluster.
func (m *Manager) clusterNodesStats(organizationID string, clusterID string) (int64, int64, error) {
	runningNodes := 0

	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
	}
	ctx, cancel := common.GetContext()
	defer cancel()
	clusterNodes, err := m.nodeClient.ListNodes(ctx, cID)
	if err != nil {
		return 0, 0, err
	}
	for _, n := range clusterNodes.Nodes {
		if n.Status == grpc_infrastructure_go.InfraStatus_RUNNING {
			runningNodes++
		}
	}
	return int64(len(clusterNodes.Nodes)), int64(runningNodes), nil
}

// Install a new cluster adding it to the system.
func (m *Manager) Install(request *grpc_public_api_go.InstallRequest) (*grpc_common_go.OpResponse, error) {
	installRequest := &grpc_installer_go.InstallRequest{
		OrganizationId:    request.OrganizationId,
		ClusterId:         request.ClusterId,
		ClusterType:       request.ClusterType,
		InstallBaseSystem: request.InstallBaseSystem,
		KubeConfigRaw:     request.KubeConfigRaw,
		Hostname:          request.Hostname,
		Username:          request.Username,
		PrivateKey:        request.PrivateKey,
		Nodes:             request.Nodes,
		TargetPlatform:    grpc_installer_go.Platform(grpc_installer_go.Platform_value[request.TargetPlatform.String()]),
		StaticIpAddresses: request.StaticIpAddresses,
	}
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.infraClient.InstallCluster(ctx, installRequest)
}

// Provision and install a new cluster adding it to the system.
func (m *Manager) ProvisionAndInstall(request *grpc_provisioner_go.ProvisionClusterRequest) (*grpc_infrastructure_manager_go.ProvisionerResponse, error) {
	request.RequestId = uuid.NewV4().String()
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.infraClient.ProvisionAndInstallCluster(ctx, request)
}

// Scale the number of nodes in the cluster.
func (m *Manager) Scale(request *grpc_provisioner_go.ScaleClusterRequest) (*grpc_infrastructure_manager_go.ProvisionerResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.infraClient.Scale(ctx, request)
}

// Uninstall a existing cluster. This process will uninstall the nalej platform and
// remove the cluster from the list.
func (m *Manager) Uninstall(request *grpc_public_api_go.UninstallClusterRequest) (*grpc_common_go.OpResponse, error) {
	imPlatform, err := entities.ToInstallerTargetPlatform(request.TargetPlatform)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	imRequest := &grpc_installer_go.UninstallClusterRequest{
		OrganizationId: request.OrganizationId,
		ClusterId:      request.ClusterId,
		ClusterType:    request.ClusterType,
		KubeConfigRaw:  request.KubeConfigRaw,
		TargetPlatform: *imPlatform,
	}
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.infraClient.Uninstall(ctx, imRequest)
}

// Decommission an application cluster. This process will uninstall the nalej platform,
// decommission the cluster from the infrastructure provider, and remove the cluster from the list.
func (m *Manager) Decommission(request *grpc_public_api_go.DecommissionClusterRequest) (*grpc_common_go.OpResponse, error) {
	imPlatform, err := entities.ToInstallerTargetPlatform(request.TargetPlatform)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	ctx, cancel := common.GetContext()
	defer cancel()
	dRequest := &grpc_provisioner_go.DecommissionClusterRequest{
		OrganizationId:      request.OrganizationId,
		ClusterId:           request.ClusterId,
		ClusterType:         request.ClusterType,
		IsManagementCluster: false,
		TargetPlatform:      *imPlatform,
		AzureCredentials:    request.AzureCredentials,
		AzureOptions:        request.AzureOptions,
	}
	return m.infraClient.DecommissionCluster(ctx, dRequest)
}

func (m *Manager) extendInfo(source *grpc_infrastructure_go.Cluster) (*grpc_public_api_go.Cluster, error) {
	totalNodes, runningNodes, err := m.clusterNodesStats(source.OrganizationId, source.ClusterId)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPICluster(source, totalNodes, runningNodes), nil
}

func (m *Manager) Info(clusterID *grpc_infrastructure_go.ClusterId) (*grpc_public_api_go.Cluster, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	retrieved, err := m.infraClient.GetCluster(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	return m.extendInfo(retrieved)
}

// List all the clusters in an organization.
func (m *Manager) List(request *grpc_public_api_go.ListRequest) (*grpc_public_api_go.ClusterList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	list, err := m.infraClient.ListClusters(ctx, &grpc_organization_go.OrganizationId{
		OrganizationId: request.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	clusters := make([]*grpc_public_api_go.Cluster, 0)
	for _, c := range list.Clusters {
		toAdd, err := m.extendInfo(c)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, toAdd)
	}

	if request.Order != nil {
		sortOptions := decorators.NewOrderOptions(*request.Order)
		sortedClusters := decorators.ApplyDecorator(clusters, decorators.NewOrderDecorator(sortOptions))
		if sortedClusters.Error != nil {
			return nil, conversions.ToGRPCError(sortedClusters.Error)
		}
		clusters = sortedClusters.ClusterList
	}

	return &grpc_public_api_go.ClusterList{
		Clusters: clusters,
	}, nil
}

// Update the cluster information.
func (m *Manager) Update(updateClusterRequest *grpc_public_api_go.UpdateClusterRequest) (*grpc_public_api_go.Cluster, error) {
	log.Debug().Interface("request", updateClusterRequest).Msg("update cluster request")
	toSend := entities.ToInfraClusterUpdate(*updateClusterRequest)
	ctx, cancel := common.GetContext()
	defer cancel()
	updated, err := m.infraClient.UpdateCluster(ctx, toSend)
	if err != nil {
		return nil, err
	}
	result, err := m.extendInfo(updated)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (m *Manager) Cordon(clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.infraClient.CordonCluster(ctx, clusterID)
}

func (m *Manager) Uncordon(clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.infraClient.UncordonCluster(ctx, clusterID)
}

func (m *Manager) DrainCluster(clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.infraClient.DrainCluster(ctx, clusterID)

}

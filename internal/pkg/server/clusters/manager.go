/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package clusters

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
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
		clustClient: clustClient, nodeClient: nodeClient, infraClient: infraClient,
	}
}

// clusterNodeStats determines the number of total and running nodes in a cluster.
func (m *Manager) clusterNodesStats(organizationID string, clusterID string) (int64, int64, error) {
	runningNodes := 0

	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
	}
	clusterNodes, err := m.nodeClient.ListNodes(context.Background(), cID)
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
func (m *Manager) Install(request *grpc_public_api_go.InstallRequest) (*grpc_infrastructure_manager_go.InstallResponse, error) {
	installRequest := &grpc_installer_go.InstallRequest{
		OrganizationId:    request.OrganizationId,
		ClusterId:         request.ClusterId,
		ClusterType:       request.ClusterType,
		InstallBaseSystem: request.InstallBaseSystem,
		KubeConfigRaw:     request.KubeConfigRaw,
		Username:          request.Username,
		PrivateKey:        request.PrivateKey,
		Nodes:             request.Nodes,
	}
	return m.infraClient.InstallCluster(context.Background(), installRequest)
}

func (m * Manager) extendInfo(source *grpc_infrastructure_go.Cluster) (*grpc_public_api_go.Cluster, error) {
	totalNodes, runningNodes, err := m.clusterNodesStats(source.OrganizationId, source.ClusterId)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPICluster(source, totalNodes, runningNodes), nil
}

func (m *Manager) Info(clusterID *grpc_infrastructure_go.ClusterId) (*grpc_public_api_go.Cluster, error) {
	retrieved, err := m.clustClient.GetCluster(context.Background(), clusterID)
	if err != nil {
		return nil, err
	}
	return m.extendInfo(retrieved)
}

// List all the clusters in an organization.
func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ClusterList, error) {
	list, err := m.clustClient.ListClusters(context.Background(), organizationID)
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
	return &grpc_public_api_go.ClusterList{
		Clusters: clusters,
	}, nil
}

// Update the cluster information.
func (m *Manager) Update(updateClusterRequest *grpc_public_api_go.UpdateClusterRequest) (*grpc_common_go.Success, error) {
	toSend := entities.ToInfraClusterUpdate(*updateClusterRequest)
	_, err := m.clustClient.UpdateCluster(context.Background(), toSend)
	if err != nil {
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

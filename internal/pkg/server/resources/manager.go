/*
 * Copyright 2019 Nalej
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

package resources

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for resources operations.
type Manager struct {
	clustClient grpc_infrastructure_go.ClustersClient
	nodeClient  grpc_infrastructure_go.NodesClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(clustClient grpc_infrastructure_go.ClustersClient,
	nodeClient grpc_infrastructure_go.NodesClient) Manager {
	return Manager{
		clustClient: clustClient, nodeClient: nodeClient,
	}
}

func (m *Manager) getNumNodes(organizationID string, clusterID string) (int, derrors.Error) {
	// Return number of nodes in a cluster
	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
	}
	ctx, cancel := common.GetContext()
	defer cancel()
	clusterNodes, err := m.nodeClient.ListNodes(ctx, cID)
	if err != nil {
		return 0, conversions.ToDerror(err)
	}
	return len(clusterNodes.Nodes), nil
}

func (m *Manager) getSummary(organizationID *grpc_organization_go.OrganizationId) (int, int, derrors.Error) {
	// Obtain list of clusters
	totalNodes := 0
	ctx, cancel := common.GetContext()
	defer cancel()
	list, err := m.clustClient.ListClusters(ctx, organizationID)
	if err != nil {
		return 0, 0, conversions.ToDerror(err)
	}
	for _, c := range list.Clusters {
		n, err := m.getNumNodes(c.OrganizationId, c.ClusterId)
		if err != nil {
			return 0, 0, err
		}
		totalNodes += n
	}
	return len(list.Clusters), totalNodes, nil
}

func (m *Manager) Summary(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ResourceSummary, error) {
	totalClusters, totalNodes, err := m.getSummary(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_public_api_go.ResourceSummary{
		OrganizationId: organizationID.OrganizationId,
		TotalClusters:  int64(totalClusters),
		TotalNodes:     int64(totalNodes),
	}, nil
}

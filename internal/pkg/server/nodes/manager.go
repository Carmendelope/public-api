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

package nodes

import (
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for node operations.
type Manager struct {
	nodeClient grpc_infrastructure_go.NodesClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(nodeClient grpc_infrastructure_go.NodesClient) Manager {
	return Manager{
		nodeClient: nodeClient,
	}
}

// List retrieves information about the nodes of a cluster.
func (m *Manager) List(clusterId *grpc_infrastructure_go.ClusterId) (*grpc_public_api_go.NodeList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	nodes, err := m.nodeClient.ListNodes(ctx, clusterId)
	if err != nil {
		return nil, err
	}
	result := make([]*grpc_public_api_go.Node, 0)
	for _, n := range nodes.Nodes {
		result = append(result, entities.ToPublicAPINode(n))
	}
	return &grpc_public_api_go.NodeList{
		Nodes: result,
	}, nil
}

// UpdateNode allows the user to update the information of a node.
func (m *Manager) UpdateNode(request *grpc_public_api_go.UpdateNodeRequest) (*grpc_public_api_go.Node, error) {
	updateRequest := &grpc_infrastructure_go.UpdateNodeRequest{
		OrganizationId: request.OrganizationId,
		NodeId:         request.NodeId,
		AddLabels:      request.AddLabels,
		RemoveLabels:   request.RemoveLabels,
		Labels:         request.Labels,
	}
	ctx, cancel := common.GetContext()
	defer cancel()
	updated, err := m.nodeClient.UpdateNode(ctx, updateRequest)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPINode(updated), nil
}

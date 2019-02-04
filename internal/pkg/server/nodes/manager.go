/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package nodes

import (
	"context"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
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
	nodes, err := m.nodeClient.ListNodes(context.Background(), clusterId)
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

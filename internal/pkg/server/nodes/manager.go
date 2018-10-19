/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package nodes

import (
	"context"
	"github.com/nalej/grpc-infrastructure-go"
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

func (m *Manager) ClusterNodes(clusterId *grpc_infrastructure_go.ClusterId) (*grpc_infrastructure_go.NodeList, error) {
	return m.nodeClient.ListNodes(context.Background(), clusterId)
}
/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package agent

import (
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for agent operations.
type Manager struct {
	agentClient grpc_inventory_manager_go.AgentClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(agentClient grpc_inventory_manager_go.AgentClient) Manager {
	return Manager{
		agentClient: agentClient,
	}
}

func (m *Manager) CreateAgentJoinToken(edgeController *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_manager_go.AgentJoinToken, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.agentClient.CreateAgentJoinToken(ctx, edgeController)
}

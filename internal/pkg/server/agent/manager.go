/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package agent

import (
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
	"github.com/satori/go.uuid"
)

// Manager structure with the required clients for agent operations.
type Manager struct {
	agentClient grpc_inventory_manager_go.AgentClient
}

const (
	activateOp = "start"
	deactivateOp = "stop"
	monitoringPlugin = "metrics"
)

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

func (m *Manager) ActivateMonitoring(assetRequest *grpc_public_api_go.AssetMonitoringRequest) (*grpc_public_api_go.AgentOpResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	op := ""
	if assetRequest.Activate {
		op = activateOp
	}else{
		op = deactivateOp
	}

	reponse, err := m.agentClient.TriggerAgentOperation(ctx, &grpc_inventory_manager_go.AgentOpRequest{
		OrganizationId: assetRequest.OrganizationId,
		EdgeControllerId: assetRequest.EdgeControllerId,
		AssetId: assetRequest.AssetId,
		OperationId: uuid.NewV4().String(),
		Operation: op,
		Plugin: monitoringPlugin,
	})
	if err != nil {
		return nil, err
	}

	return entities.ToPublicAPIAgentOpRequest(reponse), nil

}

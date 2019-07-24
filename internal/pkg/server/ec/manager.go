/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package ec

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for node operations.
type Manager struct {
	eicClient grpc_inventory_manager_go.EICClient
	agentClient grpc_inventory_manager_go.AgentClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(eicClient grpc_inventory_manager_go.EICClient, agentClient grpc_inventory_manager_go.AgentClient) Manager {
	return Manager{
		eicClient: eicClient,
		agentClient:agentClient,
	}
}

func (m *Manager) CreateEICToken(organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.EICJoinToken, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.eicClient.CreateEICToken(ctx, organizationID)
}

func (m *Manager) UnlinkEIC(edgeControllerID *grpc_inventory_manager_go.UnlinkECRequest) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.eicClient.UnlinkEIC(ctx, edgeControllerID)
}

func (m*Manager) InstallAgent(request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_public_api_go.ECOpResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	response, err :=  m.agentClient.InstallAgent(ctx, request)

	if err != nil {
		return nil, nil
	}
	return entities.ToPublicAPIECOPResponse(response), nil

}

func (m *Manager) UpdateGeolocation(updateRequest *grpc_inventory_manager_go.UpdateGeolocationRequest) (*grpc_inventory_go.EdgeController, error){
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.eicClient.UpdateECGeolocation(ctx, updateRequest)
}
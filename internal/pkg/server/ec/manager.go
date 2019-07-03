/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package ec

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
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

func (m*Manager) InstallAgent(request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_inventory_manager_go.InstallAgentResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.agentClient.InstallAgent(ctx, request)
}

func (m *Manager) UpdateGeolocation(updateRequest *grpc_inventory_manager_go.UpdateGeolocationRequest) (*grpc_inventory_go.EdgeController, error){
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.eicClient.UpdateECGeolocation(ctx, updateRequest)
}
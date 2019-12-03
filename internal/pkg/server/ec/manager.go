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
	eicClient   grpc_inventory_manager_go.EICClient
	agentClient grpc_inventory_manager_go.AgentClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(eicClient grpc_inventory_manager_go.EICClient, agentClient grpc_inventory_manager_go.AgentClient) Manager {
	return Manager{
		eicClient:   eicClient,
		agentClient: agentClient,
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

func (m *Manager) InstallAgent(request *grpc_inventory_manager_go.InstallAgentRequest) (*grpc_public_api_go.ECOpResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	response, err := m.agentClient.InstallAgent(ctx, request)

	if err != nil {
		return nil, nil
	}
	return entities.ToPublicAPIECOPResponse(response), nil

}

func (m *Manager) UpdateGeolocation(updateRequest *grpc_inventory_manager_go.UpdateGeolocationRequest) (*grpc_inventory_go.EdgeController, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.eicClient.UpdateECGeolocation(ctx, updateRequest)
}

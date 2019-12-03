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
	activateOp       = "start"
	deactivateOp     = "stop"
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
	} else {
		op = deactivateOp
	}

	reponse, err := m.agentClient.TriggerAgentOperation(ctx, &grpc_inventory_manager_go.AgentOpRequest{
		OrganizationId:   assetRequest.OrganizationId,
		EdgeControllerId: assetRequest.EdgeControllerId,
		AssetId:          assetRequest.AssetId,
		OperationId:      uuid.NewV4().String(),
		Operation:        op,
		Plugin:           monitoringPlugin,
	})
	if err != nil {
		return nil, err
	}

	return entities.ToPublicAPIAgentOpRequest(reponse), nil

}

func (m *Manager) UninstallAgent(request *grpc_inventory_manager_go.UninstallAgentRequest) (*grpc_public_api_go.ECOpResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	response, err := m.agentClient.UninstallAgent(ctx, request)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIECOPResponse(response), nil
}

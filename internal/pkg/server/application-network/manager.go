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

package application_network

import (
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

type Manager struct {
	appNetClient grpc_application_manager_go.ApplicationNetworkClient
	appClient    grpc_application_manager_go.ApplicationManagerClient
}

func NewManager(
	applicationNetworkClient grpc_application_manager_go.ApplicationNetworkClient,
	applicationClient grpc_application_manager_go.ApplicationManagerClient,
) Manager {
	return Manager{appNetClient: applicationNetworkClient, appClient: applicationClient}
}

// AddConnection adds a new connection between one outbound and one inbound
func (m *Manager) AddConnection(connRequest *grpc_application_network_go.AddConnectionRequest) (*grpc_public_api_go.OpResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	appNetResponse, err := m.appNetClient.AddConnection(ctx, connRequest)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIOpResponse(appNetResponse), nil
}

// RemoveConnection removes a connection
func (m *Manager) RemoveConnection(connRequest *grpc_application_network_go.RemoveConnectionRequest) (*grpc_public_api_go.OpResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	appNetResponse, err := m.appNetClient.RemoveConnection(ctx, connRequest)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIOpResponse(appNetResponse), nil
}

// ListConnections retrieves a list all the established connections of an organization
func (m *Manager) ListConnections(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_network_go.ConnectionInstanceList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appNetClient.ListConnections(ctx, organizationID)
}

func (m *Manager) ListAvailableInstanceInbounds(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_manager_go.AvailableInstanceInboundList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appClient.ListAvailableInstanceInbounds(ctx, organizationID)
}

func (m *Manager) ListAvailableInstanceOutbounds(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_manager_go.AvailableInstanceOutboundList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appClient.ListAvailableInstanceOutbounds(ctx, organizationID)
}

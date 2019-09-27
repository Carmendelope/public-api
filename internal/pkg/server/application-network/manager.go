/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
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
}

func NewManager(client grpc_application_manager_go.ApplicationNetworkClient) Manager {
	return Manager{appNetClient: client}
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

func (m *Manager) ListAvailableInboundConnections(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_manager_go.AvailableInstanceInboundList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appNetClient.ListAvailableInstanceInbounds(ctx, organizationID)
}

func (m *Manager) ListAvailableOutboundConnections(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_manager_go.AvailableInstanceOutboundList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appNetClient.ListAvailableInstanceOutbounds(ctx, organizationID)
}

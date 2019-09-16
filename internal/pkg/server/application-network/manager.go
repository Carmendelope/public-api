/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

type Manager struct {
	appNetClient grpc_application_network_go.ApplicationNetworkClient
}

func NewManager(client grpc_application_network_go.ApplicationNetworkClient) Manager {
	return Manager{appNetClient: client}
}

// AddConnection adds a new connection between one outbound and one inbound
func (m *Manager) AddConnection(connRequest *grpc_application_network_go.AddConnectionRequest) (*grpc_application_network_go.ConnectionInstance, error){
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appNetClient.AddConnection(ctx, connRequest)
}

// RemoveConnection removes a connection
func (m *Manager) RemoveConnection(connRequest *grpc_application_network_go.RemoveConnectionRequest) (*grpc_common_go.Success, error){
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appNetClient.RemoveConnection(ctx, connRequest)
}

// ListConnections retrieves a list all the established connections of an organization
func (m *Manager) ListConnections(organizationID *grpc_organization_go.OrganizationId) (*grpc_application_network_go.ConnectionInstanceList, error){
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.appNetClient.ListConnections(ctx, organizationID)
}
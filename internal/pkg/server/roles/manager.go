/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package roles

import (
	"context"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-manager-go"
)

// Manager structure with the required clients for roles operations.
type Manager struct {
	client grpc_user_manager_go.UserManagerClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(client grpc_user_manager_go.UserManagerClient) Manager {
	return Manager{client}
}

func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_authx_go.RoleList, error) {
	return m.client.ListRoles(context.Background(), organizationID)
}
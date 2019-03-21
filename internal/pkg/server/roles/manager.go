/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package roles

import (
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/public-api/internal/pkg/server/common"
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
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.client.ListRoles(ctx, organizationID)
}

func (m *Manager) AssignRole(request *grpc_user_manager_go.AssignRoleRequest) (*grpc_user_manager_go.User, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.client.AssignRole(ctx, request)
}

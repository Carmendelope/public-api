/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package roles

import (
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
)

// Manager structure with the required clients for roles operations.
type Manager struct {

}

// NewManager creates a Manager using a set of clients.
func NewManager() Manager {
	return Manager{}
}

func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_authx_go.RoleList, error) {
	panic("implement me")
}
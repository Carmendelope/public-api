/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package roles

import (
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
)

type Manager struct {

}

func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_authx_go.RoleList, error) {
	panic("implement me")
}
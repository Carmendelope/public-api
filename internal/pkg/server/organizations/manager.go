/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organizations

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for organization operations.
type Manager struct {
	orgClient grpc_organization_go.OrganizationsClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(orgClient grpc_organization_go.OrganizationsClient) Manager {
	return Manager{orgClient: orgClient}
}

func (m *Manager) ToOrganizationInfo(organization *grpc_organization_go.Organization) *grpc_public_api_go.OrganizationInfo {
	return &grpc_public_api_go.OrganizationInfo{
		OrganizationId: organization.OrganizationId,
		Name:           organization.Name,
	}
}

func (m *Manager) Info(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.OrganizationInfo, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	retrieved, err := m.orgClient.GetOrganization(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	return m.ToOrganizationInfo(retrieved), nil
}

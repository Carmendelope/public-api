/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organizations

import (
	"context"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
)

// Manager structure with the required clients for organization operations.
type Manager struct {
	orgClient grpc_organization_go.OrganizationsClient
}

func NewManager(orgClient grpc_organization_go.OrganizationsClient) Manager {
	return Manager{orgClient:orgClient}
}

func (m * Manager) ToOrganizationInfo(organization * grpc_organization_go.Organization) *grpc_public_api_go.OrganizationInfo {
	return &grpc_public_api_go.OrganizationInfo{
		OrganizationId:       organization.OrganizationId,
		Name:                 organization.Name,
	}
}

func (m *Manager) Info(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.OrganizationInfo, error) {
	retrieved, err := m.orgClient.GetOrganization(context.Background(), organizationID)
	if err != nil {
		return nil, err
	}
	return m.ToOrganizationInfo(retrieved), nil
}
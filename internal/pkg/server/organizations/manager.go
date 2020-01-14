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

package organizations

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for organization operations.
type Manager struct {
	orgClient grpc_organization_manager_go.OrganizationsClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(orgClient grpc_organization_manager_go.OrganizationsClient) Manager {
	return Manager{orgClient: orgClient}
}

func (m *Manager) ToOrganizationInfo(organization *grpc_organization_manager_go.Organization) *grpc_organization_manager_go.Organization{
	return &grpc_organization_manager_go.Organization{
		OrganizationId: organization.OrganizationId,
		Name:           organization.Name,
	}
}

func (m *Manager) Info(organizationID *grpc_organization_go.OrganizationId) (*grpc_organization_manager_go.Organization, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.orgClient.GetOrganization(ctx, organizationID)
	//retrieved, err := m.orgClient.GetOrganization(ctx, organizationID)
	//if err != nil {
	//	log.Debug().Interface("error", err).Msg("error retrieving info")
	//	return nil, err
	//}
	//return m.ToOrganizationInfo(retrieved), nil
}

func (m *Manager) Update(updateRequest *grpc_organization_go.UpdateOrganizationRequest)  (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.orgClient.UpdateOrganization(ctx, updateRequest)
}

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
 *
 *
 */

package organization_settings

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
	"github.com/nalej/public-api/internal/pkg/server/decorators"
)

type Manager struct {
	settingClient grpc_organization_manager_go.OrganizationsClient
}

func NewManager(settingClient grpc_organization_manager_go.OrganizationsClient) Manager {
	return Manager{settingClient: settingClient}
}

func (m *Manager) Update(updateRequest *grpc_public_api_go.UpdateSettingRequest) (*grpc_common_go.Success, error) {

	ctx, cancel := common.GetContext()
	defer cancel()

	return m.settingClient.UpdateSettings(ctx, entities.ToUpdateSettingRequest(updateRequest))

}

func (m *Manager) List(organizationID *grpc_public_api_go.ListRequest) (*grpc_organization_manager_go.SettingList, error) {

	ctx, cancel := common.GetContext()
	defer cancel()

	list, err := m.settingClient.ListSettings(ctx, &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID.OrganizationId,
	})
	if err != nil {
		return nil, err
	}
	// if sorting requested -> apply the decorator
	if organizationID.Order != nil {
		sortOptions := decorators.OrderOptions{Field: organizationID.Order.Field, Asc: organizationID.Order.Order == grpc_common_go.Order_ASC}
		sortingResponse := decorators.ApplyDecorator(list.Settings, decorators.NewOrderDecorator(sortOptions))
		if sortingResponse.Error != nil {
			return nil, conversions.ToGRPCError(sortingResponse.Error)
		} else {
			list.Settings = sortingResponse.SettingList
		}
	}
	return list, nil
}

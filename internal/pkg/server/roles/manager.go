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

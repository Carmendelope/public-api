/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package users

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for users operations.
type Manager struct {
	umClient grpc_user_manager_go.UserManagerClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(client grpc_user_manager_go.UserManagerClient) Manager {
	return Manager{client}
}

func (m *Manager) Add(addUserRequest *grpc_public_api_go.AddUserRequest) (*grpc_public_api_go.User, error) {
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: addUserRequest.OrganizationId,
	}
	ctx, cancel := common.GetContext()
	defer cancel()
	role, err := m.umClient.ListRoles(ctx, orgID)
	if err != nil {
		return nil, err
	}
	var roleId string
	for _, r := range role.Roles {
		if r.Name == addUserRequest.RoleName {
			roleId = r.RoleId
		}
	}
	if roleId == "" {
		return nil, conversions.ToGRPCError(derrors.NewInvalidArgumentError("role not found"))
	}
	toAdd := &grpc_user_manager_go.AddUserRequest{
		OrganizationId: addUserRequest.OrganizationId,
		Email:          addUserRequest.Email,
		Password:       addUserRequest.Password,
		Name:           addUserRequest.Name,
		PhotoUrl:       "",
		RoleId:         roleId,
	}
	ctx2, cancel2 := common.GetContext()
	defer cancel2()
	added, err := m.umClient.AddUser(ctx2, toAdd)
	if err != nil {
		return nil, err
	}
	return &grpc_public_api_go.User{
		OrganizationId: added.OrganizationId,
		Email:          added.Email,
		Name:           added.Name,
		RoleName:       added.RoleName,
	}, nil
}

func (m *Manager) Info(userID *grpc_user_go.UserId) (*grpc_public_api_go.User, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	retrieved, err := m.umClient.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &grpc_public_api_go.User{
		OrganizationId: retrieved.OrganizationId,
		Email:          retrieved.Email,
		Name:           retrieved.Name,
		RoleName:       retrieved.RoleName,
	}, nil
}

func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.UserList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	list, err := m.umClient.ListUsers(ctx, organizationID)
	if err != nil {
		return nil, err
	}
	users := make([]*grpc_public_api_go.User, 0)
	for _, u := range list.Users {
		if !u.Internal {
			toAdd := &grpc_public_api_go.User{
				OrganizationId: u.OrganizationId,
				Email:          u.Email,
				Name:           u.Name,
				RoleName:       u.RoleName,
			}
			users = append(users, toAdd)
		}
	}
	return &grpc_public_api_go.UserList{
		Users: users,
	}, nil
}

func (m *Manager) Delete(userID *grpc_user_go.UserId) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.umClient.RemoveUser(ctx, userID)
}

func (m *Manager) Update(updateUserRequest *grpc_user_go.UpdateUserRequest) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.umClient.Update(ctx, updateUserRequest)
}

func (m *Manager) ResetPassword(changePasswordRequest *grpc_user_manager_go.ChangePasswordRequest) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.umClient.ChangePassword(ctx, changePasswordRequest)
}

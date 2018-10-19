/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package users

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
)

// Manager structure with the required clients for users operations.
type Manager struct {

}

// NewManager creates a Manager using a set of clients.
func NewManager() Manager {
	return Manager{}
}

func (m * Manager) Info(userID *grpc_user_go.UserId) (*grpc_public_api_go.User, error) {
	panic("implement me")
}

func (m * Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.UserList, error) {
	panic("implement me")
}

func (m * Manager) Delete(userID *grpc_user_go.UserId) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (m * Manager) ResetPassword(userID *grpc_user_go.UserId) (*grpc_public_api_go.PasswordResetResponse, error) {
	panic("implement me")
}

func (m * Manager) Update(updateUserRequest *grpc_user_go.UpdateUserRequest) (*grpc_common_go.Success, error) {
	panic("implement me")
}
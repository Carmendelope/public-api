/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
)

const emptyOrganizationId = "organization_id cannot be empty"
const emptyInstanceId = "app_instance_id cannot be empty"
const emptyDescriptorId = "app_descriptor_id cannot be empty"
const emptyClusterId = "cluster_id cannot be empty"
const emptyEmail = "email cannot be empty"
const emptyName = "name cannot be empty"
const emptyPassword = "password cannot be empty"
const emptyNewPassword = "new password cannot be empty"
const emptyRoleName = "role_name cannot be empty"
const emptyRoleID = "role_id cannot be empty"

func ValidOrganizationId(organizationID *grpc_organization_go.OrganizationId) derrors.Error {
	if organizationID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidClusterId(clusterID *grpc_infrastructure_go.ClusterId) derrors.Error {
	if clusterID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if clusterID.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	return nil
}

func ValidUserId(userID *grpc_user_go.UserId) derrors.Error {
	if userID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if userID.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAppInstanceID(appInstanceID *grpc_application_go.AppInstanceId) derrors.Error {
	if appInstanceID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if appInstanceID.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyInstanceId)
	}
	return nil
}

func ValidAppDescriptorID(appDescriptorID *grpc_application_go.AppDescriptorId) derrors.Error {
	if appDescriptorID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if appDescriptorID.AppDescriptorId == "" {
		return derrors.NewInvalidArgumentError(emptyDescriptorId)
	}
	return nil
}

func ValidUpdateClusterRequest(updateClusterRequest *grpc_public_api_go.UpdateClusterRequest) derrors.Error {
	if updateClusterRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if updateClusterRequest.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	return nil
}

func ValidAddUserRequest(request *grpc_public_api_go.AddUserRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	if request.Password == "" {
		return derrors.NewInvalidArgumentError(emptyPassword)
	}
	if request.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	if request.RoleName == "" {
		return derrors.NewInvalidArgumentError(emptyRoleName)
	}
	return nil
}

func ValidUpdateUserRequest(updateUserRequest *grpc_user_go.UpdateUserRequest) derrors.Error {
	if updateUserRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if updateUserRequest.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAddAppDescriptor(request *grpc_application_go.AddAppDescriptorRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if len(request.Services) == 0 {
		return derrors.NewInvalidArgumentError("expecting at least one service")
	}
	for _, s := range request.Services {
		if s.OrganizationId != request.OrganizationId {
			return derrors.NewInvalidArgumentError("organization_id mismatch")
		}
	}
	return nil
}

func ValidDeployRequest(request *grpc_application_manager_go.DeployRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AppDescriptorId == "" {
		return derrors.NewInvalidArgumentError(emptyDescriptorId)
	}
	return nil
}

func ValidInstallRequest(request *grpc_public_api_go.InstallRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidChangePasswordRequest(request *grpc_user_manager_go.ChangePasswordRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.Password == "" {
		return derrors.NewInvalidArgumentError(emptyPassword)
	}
	if request.NewPassword == "" {
		return derrors.NewInvalidArgumentError(emptyNewPassword)
	}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAssignRoleRequest(request *grpc_user_manager_go.AssignRoleRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.RoleId == "" {
		return derrors.NewInvalidArgumentError(emptyRoleID)
	}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

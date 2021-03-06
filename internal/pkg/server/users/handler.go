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

package users

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the user requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) Add(ctx context.Context, addUserRequest *grpc_public_api_go.AddUserRequest) (*grpc_public_api_go.User, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if addUserRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAddUserRequest(addUserRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Add(addUserRequest)
}

func (h *Handler) Info(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_public_api_go.User, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if userID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	if !rm.OrgPrimitive && userID.Email != rm.UserID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested user")
	}
	err = entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Info(userID)
}

func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.UserList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if organizationID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidOrganizationId(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.List(organizationID)
}

func (h *Handler) Delete(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_common_go.Success, error) {
	log.Debug().Str("organizationID", userID.OrganizationId).Str("email", userID.Email).Msg("delete user")
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if userID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Delete(userID)
}

func (h *Handler) ChangePassword(ctx context.Context, changePasswordRequest *grpc_user_manager_go.ChangePasswordRequest) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if changePasswordRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	if !rm.OrgPrimitive && changePasswordRequest.Email != rm.UserID {
		return nil, derrors.NewPermissionDeniedError("cannot reset password of selected user")
	}
	err = entities.ValidChangePasswordRequest(changePasswordRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.ResetPassword(changePasswordRequest)
}

func (h *Handler) Update(ctx context.Context, updateUserRequest *grpc_user_go.UpdateUserRequest) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if updateUserRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	log.Debug().Interface("rm", rm).Interface("updateUserRequest", updateUserRequest).Msg("Processing update request")
	if !rm.OrgPrimitive && updateUserRequest.Email != rm.UserID {
		return nil, derrors.NewPermissionDeniedError("cannot update the information of selected user")
	}
	err = entities.ValidUpdateUserRequest(updateUserRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Update(updateUserRequest)
}

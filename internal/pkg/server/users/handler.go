/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package users

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the user requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h *Handler) Info(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_public_api_go.User, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	if userID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewUnauthenticatedError("cannot access requested OrganizationID")
	}
	if userID.Email != rm.UserID {
		return nil, derrors.NewUnauthenticatedError("cannot access requested user")
	}
	err = entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Info(userID)
}

func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.UserList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	if organizationID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewUnauthenticatedError("cannot access requested OrganizationID")
	}
	err = entities.ValidOrganizationId(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.List(organizationID)
}

func (h *Handler) Delete(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	if userID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewUnauthenticatedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Delete(userID)
}

func (h *Handler) ResetPassword(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_public_api_go.PasswordResetResponse, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	if userID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewUnauthenticatedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.ResetPassword(userID)
}

func (h *Handler) Update(ctx context.Context, updateUserRequest *grpc_user_go.UpdateUserRequest) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	if updateUserRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewUnauthenticatedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUpdateUserRequest(updateUserRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Update(updateUserRequest)
}



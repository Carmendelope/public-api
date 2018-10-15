/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package users

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/entities"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h *Handler) Info(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_public_api_go.User, error) {
	err := entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Info(userID)
}

func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.UserList, error) {
	err := entities.ValidOrganizationId(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.List(organizationID)
}

func (h *Handler) Delete(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_common_go.Success, error) {
	err := entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Delete(userID)
}

func (h *Handler) ResetPassword(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_public_api_go.PasswordResetResponse, error) {
	err := entities.ValidUserId(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.ResetPassword(userID)
}

func (h *Handler) Update(ctx context.Context, updateUserRequest *grpc_user_go.UpdateUserRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidUpdateUserRequest(updateUserRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Update(updateUserRequest)
}



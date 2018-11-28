/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package roles

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the roles requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h*Handler) ToPublicRoleList(roles *grpc_authx_go.RoleList, internal bool) * grpc_public_api_go.RoleList{
	aux := make([]*grpc_authx_go.Role, 0)
	for _, r := range roles.Roles{
		if internal == r.Internal {
			aux = append(aux, r)
		}
	}
	result := make([]*grpc_public_api_go.Role, 0)
	for _, r := range roles.Roles{
		primitives := make([]string, 0)
		for _, p := range r.Primitives{
			primitives = append(primitives, p.String())
		}
		toAdd := &grpc_public_api_go.Role{
			OrganizationId:       r.OrganizationId,
			RoleId:               r.RoleId,
			Name:                 r.Name,
			Primitives:           primitives,
		}
		result = append(result, toAdd)
	}
	return &grpc_public_api_go.RoleList{
		Roles:                result,
	}
}

func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.RoleList, error) {
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
	roles, lErr := h.Manager.List(organizationID)
	if lErr != nil {
		return nil, lErr
	}
	return h.ToPublicRoleList(roles, false), nil
}

// ListInternal retrieves the list of internal roles inside an organization.
func (h*Handler) ListInternal(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.RoleList, error){
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
	roles, lErr := h.Manager.List(organizationID)
	if lErr != nil {
		return nil, lErr
	}
	return h.ToPublicRoleList(roles, true), nil
}

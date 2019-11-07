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
 *
 */

package roles

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the roles requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) ToPublicRoleList(roles *grpc_authx_go.RoleList, internal bool) *grpc_public_api_go.RoleList {
	log.Debug().Bool("internal", internal).Int("received", len(roles.Roles)).Msg("Transforming role list")
	aux := make([]*grpc_authx_go.Role, 0)
	for _, r := range roles.Roles {
		if internal == r.Internal {
			aux = append(aux, r)
		}
	}
	result := make([]*grpc_public_api_go.Role, 0)
	for _, r := range aux {
		primitives := make([]string, 0)
		for _, p := range r.Primitives {
			primitives = append(primitives, p.String())
		}
		toAdd := &grpc_public_api_go.Role{
			OrganizationId: r.OrganizationId,
			RoleId:         r.RoleId,
			Name:           r.Name,
			Primitives:     primitives,
		}
		result = append(result, toAdd)
	}
	return &grpc_public_api_go.RoleList{
		Roles: result,
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
func (h *Handler) ListInternal(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.RoleList, error) {
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

// AssignRole assigns a role to an existing user.
func (h *Handler) AssignRole(ctx context.Context, request *grpc_user_manager_go.AssignRoleRequest) (*grpc_user_manager_go.User, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAssignRoleRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	user, lErr := h.Manager.AssignRole(request)
	if lErr != nil {
		return nil, lErr
	}
	return user, nil
}

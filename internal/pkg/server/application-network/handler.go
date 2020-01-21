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

package application_network

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the applications requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// AddConnection adds a new connection between one outbound and one inbound
func (h *Handler) AddConnection(ctx context.Context, connRequest *grpc_application_network_go.AddConnectionRequest) (*grpc_public_api_go.OpResponse, error) {

	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if connRequest.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAddConnectionRequest(connRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	return h.Manager.AddConnection(connRequest)
}

// RemoveConnection removes a connection
func (h *Handler) RemoveConnection(ctx context.Context, request *grpc_application_network_go.RemoveConnectionRequest) (*grpc_public_api_go.OpResponse, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidRemoveConnectionRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	return h.Manager.RemoveConnection(request)
}

// ListConnections retrieves a list all the established connections of an organization
func (h *Handler) ListConnections(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ConnectionInstanceList, error) {
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
	connections, mErr := h.Manager.ListConnections(organizationID)
	if mErr != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return entities.ToPublicAPIConnectionList(connections), err
}

// ListAvailableInstanceInbounds retrieves a list of available inbounds of an organization
func (h *Handler) ListAvailableInstanceInbounds(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_application_manager_go.AvailableInstanceInboundList, error) {
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

	return h.Manager.ListAvailableInstanceInbounds(organizationID)
}

// ListAvailableInstanceOutbounds retrieves a list of available outbounds of an organization
func (h *Handler) ListAvailableInstanceOutbounds(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_application_manager_go.AvailableInstanceOutboundList, error) {
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

	return h.Manager.ListAvailableInstanceOutbounds(organizationID)
}

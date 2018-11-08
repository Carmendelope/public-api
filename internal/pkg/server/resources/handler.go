/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package resources

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the resources requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h * Handler) Summary(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ResourceSummary, error) {
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
	return h.Manager.Summary(organizationID)
}


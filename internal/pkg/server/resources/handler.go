/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package resources

import (
	"context"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/entities"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h * Handler) Summary(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ResourceSummary, error) {
	err := entities.ValidOrganizationId(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Summary(organizationID)
}


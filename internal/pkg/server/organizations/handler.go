/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organizations

import (
	"context"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the organizations requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h *Handler) Info(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.OrganizationInfo, error) {
	err := entities.ValidOrganizationId(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Info(organizationID)
}


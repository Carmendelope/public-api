/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package clusters

import (
	"context"
	"github.com/nalej/grpc-common-go"
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

func (h * Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ClusterList, error) {
	err := entities.ValidOrganizationId(organizationID)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.List(organizationID)
}

func (h * Handler) Update(ctx context.Context, updateClusterRequest *grpc_public_api_go.UpdateClusterRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidUpdateClusterRequest(updateClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.Update(updateClusterRequest)
}



/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package nodes

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the node requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) ClusterNodes(ctx context.Context, clusterId *grpc_infrastructure_go.ClusterId) (*grpc_infrastructure_go.NodeList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if clusterId.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidClusterId(clusterId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.ClusterNodes(clusterId)
}

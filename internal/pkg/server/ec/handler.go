/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package ec

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"golang.org/x/net/context"
)

// Handler structure for the ec requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (*Handler) CreateEICToken(context.Context, *grpc_organization_go.OrganizationId) (*grpc_inventory_manager_go.EICJoinToken, error) {
	panic("implement me")
}

func (*Handler) UnlinkEIC(context.Context, *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	panic("implement me")
}



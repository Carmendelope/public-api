/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package inventory

import (
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for node operations.
type Manager struct {
	invManagerClient grpc_inventory_manager_go.InventoryClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(invManagerClient grpc_inventory_manager_go.InventoryClient) Manager {
	return Manager{
		invManagerClient: invManagerClient,
	}
}

func (m * Manager) List(orgID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.InventoryList, error){
	ctx, cancel := common.GetContext()
	defer cancel()
	list, err := m.invManagerClient.List(ctx, orgID)
	if err != nil{
		return nil, err
	}

	devices := entities.ToPublicAPIDeviceArray(list.Devices)
	assets := entities.ToPublicAPIAssetArray(list.Assets)
	controllers := entities.ToPublicAPIControllerArray(list.Controllers)

	return &grpc_public_api_go.InventoryList{
		Devices:              devices,
		Assets:               assets,
		Controllers:          controllers,
	}, nil
}
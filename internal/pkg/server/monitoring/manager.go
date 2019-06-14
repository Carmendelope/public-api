/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"github.com/nalej/grpc-inventory-manager-go"
)

// Manager structure with the required clients for monitoring operations.
type Manager struct {
	client *grpc_inventory_manager_go.InventoryMonitoringClient
}

func NewManager(client grpc_inventory_manager_go.InventoryMonitoringClient) Manager {
	return Manager{
		client: &client,
	}
}

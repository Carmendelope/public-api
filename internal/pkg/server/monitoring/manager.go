/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"github.com/nalej/derrors"

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

func (m *Manager) ListMetrics(selector *grpc_inventory_manager_go.AssetSelector) (*grpc_inventory_manager_go.MetricsList, error) {
	return nil, derrors.NewUnimplementedError("ListMetrics not implemented")
}

func (m *Manager) QueryMetrics(request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	return nil, derrors.NewUnimplementedError("QueryMetrics not implemented")
}

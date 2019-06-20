/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"github.com/nalej/grpc-inventory-manager-go"

	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for monitoring operations.
type Manager struct {
	client grpc_inventory_manager_go.InventoryMonitoringClient
}

func NewManager(client grpc_inventory_manager_go.InventoryMonitoringClient) Manager {
	return Manager{
		client: client,
	}
}

func (m *Manager) ListMetrics(selector *grpc_inventory_manager_go.AssetSelector) (*grpc_inventory_manager_go.MetricsList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.client.ListMetrics(ctx, selector)
}

func (m *Manager) QueryMetrics(request *grpc_inventory_manager_go.QueryMetricsRequest) (*grpc_inventory_manager_go.QueryMetricsResult, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.client.QueryMetrics(ctx, request)
}

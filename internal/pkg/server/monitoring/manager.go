/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package monitoring

import (
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-monitoring-go"

	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for monitoring operations.
type Manager struct {
	client grpc_monitoring_go.AssetMonitoringClient
}

func NewManager(client grpc_monitoring_go.AssetMonitoringClient) Manager {
	return Manager{
		client: client,
	}
}

func (m *Manager) ListMetrics(selector *grpc_inventory_go.AssetSelector) (*grpc_monitoring_go.MetricsList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.client.ListMetrics(ctx, selector)
}

func (m *Manager) QueryMetrics(request *grpc_monitoring_go.QueryMetricsRequest) (*grpc_monitoring_go.QueryMetricsResult, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.client.QueryMetrics(ctx, request)
}

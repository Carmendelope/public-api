/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package monitoring

import (
	"github.com/nalej/grpc-monitoring-go"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

// Manager structure with the required clients for monitoring operations.
type Manager struct {
	client grpc_monitoring_go.MonitoringManagerClient
}

func NewManager(client grpc_monitoring_go.MonitoringManagerClient) Manager {
	return Manager{client: client}
}

func (m *Manager) GetClusterStats(request *grpc_monitoring_go.ClusterStatsRequest) (*grpc_monitoring_go.ClusterStats, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.client.GetClusterStats(ctx, request)
}

func (m *Manager) GetClusterSummary(request *grpc_monitoring_go.ClusterSummaryRequest) (*grpc_monitoring_go.ClusterSummary, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.client.GetClusterSummary(ctx, request)
}

func (m *Manager) GetOrganizationApplicationStats(request *grpc_monitoring_go.OrganizationApplicationStatsRequest) (*grpc_monitoring_go.OrganizationApplicationStatsResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()

	return m.client.GetOrganizationApplicationStats(ctx, request)
}

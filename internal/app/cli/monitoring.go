/*
 * Copyright 2020 Nalej
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

package cli

import (
	"github.com/nalej/grpc-monitoring-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"strings"
	"time"
)

type Monitoring struct {
	Connection
	Credentials
}

func NewMonitoring(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Monitoring {
	return &Monitoring{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (m *Monitoring) load() {
	err := m.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (m *Monitoring) getClient() (grpc_public_api_go.MonitoringClient, *grpc.ClientConn) {
	connection, err := m.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewMonitoringClient(connection)
	return client, connection
}

func (m *Monitoring) GetClusterStats(organizationId string, clusterId string, rangeMinutes int32, fields string) {
	if organizationId == "" {
		log.Fatal().Msg("organizationId cannot be empty")
	}
	if clusterId == "" {
		log.Fatal().Msg("clusterId cannot be empty")
	}

	m.load()
	context, cancel := m.GetContext()
	client, connection := m.getClient()
	defer connection.Close()
	defer cancel()

	request := &grpc_monitoring_go.ClusterStatsRequest{
		OrganizationId: organizationId,
		ClusterId:      clusterId,
		RangeMinutes:   rangeMinutes,
		Fields:         toPlatformStatusFields(strings.Split(fields, ",")),
	}

	clusterStats, err := client.GetClusterStats(context, request)
	m.PrintResultOrError(clusterStats, err, "cannot query cluster stats")
}

func toPlatformStatusFields(fields []string) []grpc_monitoring_go.PlatformStatsField {
	platformStatFields := make([]grpc_monitoring_go.PlatformStatsField, 0)
	for _, fieldName := range fields {
		if fieldName != "" {
			platformStatFieldValue, exists := grpc_monitoring_go.PlatformStatsField_value[strings.ToUpper(fieldName)]
			if exists {
				platformStatFields = append(platformStatFields, grpc_monitoring_go.PlatformStatsField(platformStatFieldValue))
			} else {
				log.Warn().Str("field", fieldName).Msg("Field name does not exist and will be ignored.")
			}
		}
	}
	return platformStatFields
}

func (m *Monitoring) GetClusterSummary(organizationId string, clusterId string, rangeMinutes int32) {
	if organizationId == "" {
		log.Fatal().Msg("organizationId cannot be empty")
	}
	if clusterId == "" {
		log.Fatal().Msg("clusterId cannot be empty")
	}

	m.load()
	context, cancel := m.GetContext()
	client, connection := m.getClient()
	defer connection.Close()
	defer cancel()

	request := &grpc_monitoring_go.ClusterSummaryRequest{
		OrganizationId: organizationId,
		ClusterId:      clusterId,
		RangeMinutes:   rangeMinutes,
	}

	clusterSummary, err := client.GetClusterSummary(context, request)
	m.PrintResultOrError(clusterSummary, err, "cannot query cluster summary")
}

func (m *Monitoring) GetOrganizationApplicationStats(organizationId string, watch bool) {
	if organizationId == "" {
		log.Fatal().Msg("organizationId cannot be empty")
	}

	m.load()
	context, cancel := m.GetContext()
	client, connection := m.getClient()
	defer connection.Close()
	defer cancel()

	request := &grpc_monitoring_go.OrganizationApplicationStatsRequest{
		OrganizationId: organizationId,
	}

	response, err := client.GetOrganizationApplicationStats(context, request)
	m.PrintResultOrError(response, err, "cannot query organization application stats")

	if watch {
		ticker := time.NewTicker(WatchSleep)
		previous := response
		for {
			_ = <-ticker.C
			context, cancel := m.GetContext()
			client, connection := m.getClient()

			request := &grpc_monitoring_go.OrganizationApplicationStatsRequest{
				OrganizationId: organizationId,
			}

			response, err := client.GetOrganizationApplicationStats(context, request)
			connection.Close()
			cancel()

			if err != nil || previous.Timestamp != response.Timestamp {
				m.PrintResultOrError(response, err, "cannot query organization application stats")
			}
			previous = response
		}
	}
}

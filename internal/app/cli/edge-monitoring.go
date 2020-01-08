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
	"strings"
	"time"

	"github.com/araddon/dateparse"

	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-monitoring-go"
	"github.com/nalej/grpc-public-api-go"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type AssetSelector struct {
	OrganizationId   string
	EdgeControllerId string
	AssetIds         []string
	GroupIds         []string
	Labels           map[string]string
}

func (a *AssetSelector) ToGRPC() *grpc_inventory_go.AssetSelector {
	selector := &grpc_inventory_go.AssetSelector{
		OrganizationId:   a.OrganizationId,
		EdgeControllerId: a.EdgeControllerId,
		AssetIds:         a.AssetIds,
		GroupIds:         a.GroupIds,
		Labels:           a.Labels,
	}

	return selector
}

type TimeRange struct {
	// Timestamps are strings to be parsed
	Timestamp  string
	Start      string
	End        string
	Resolution time.Duration
}

func dateParse(in string) int64 {
	if in == "" {
		return 0
	}

	t, err := dateparse.ParseLocal(in)
	if err != nil {
		log.Fatal().Str("timestamp", in).Err(err).Msg("invalid timestamp")
	}

	return t.UTC().Unix()
}

func (t *TimeRange) ToGRPC() *grpc_monitoring_go.QueryMetricsRequest_TimeRange {
	timeRange := &grpc_monitoring_go.QueryMetricsRequest_TimeRange{
		Timestamp:  dateParse(t.Timestamp),
		TimeStart:  dateParse(t.Start),
		TimeEnd:    dateParse(t.End),
		Resolution: int64(t.Resolution.Seconds()),
	}

	return timeRange
}

type InventoryMonitoring struct {
	Connection
	Credentials
}

func NewInventoryMonitoring(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *InventoryMonitoring {
	return &InventoryMonitoring{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (i *InventoryMonitoring) load() {
	err := i.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (i *InventoryMonitoring) getClient() (grpc_public_api_go.InventoryMonitoringClient, *grpc.ClientConn) {
	conn, err := i.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewInventoryMonitoringClient(conn)
	return client, conn
}

func (i *InventoryMonitoring) QueryMetrics(selector *AssetSelector, metrics []string, timeRange *TimeRange, aggr string) {
	i.load()
	ctx, cancel := i.GetContext()
	client, conn := i.getClient()
	defer conn.Close()
	defer cancel()

	aggrType, found := grpc_monitoring_go.AggregationType_value[aggr]
	if !found {
		methods := []string{}
		for method := range grpc_monitoring_go.AggregationType_value {
			methods = append(methods, method)
		}
		log.Fatal().Str("aggregation", aggr).Msg("Aggregation method not available. Available methods: " + strings.Join(methods, ", "))
	}

	query := &grpc_monitoring_go.QueryMetricsRequest{
		Assets:      selector.ToGRPC(),
		Metrics:     metrics,
		TimeRange:   timeRange.ToGRPC(),
		Aggregation: grpc_monitoring_go.AggregationType(aggrType),
	}

	result, err := client.QueryMetrics(ctx, query)
	i.PrintResultOrError(result, err, "cannot query inventory metrics")
}

func (i *InventoryMonitoring) ListMetrics(selector *AssetSelector) {
	i.load()
	ctx, cancel := i.GetContext()
	client, conn := i.getClient()
	defer conn.Close()
	defer cancel()

	metrics, err := client.ListMetrics(ctx, selector.ToGRPC())
	i.PrintResultOrError(metrics, err, "cannot list available inventory metrics")
}

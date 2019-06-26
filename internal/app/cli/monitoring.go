/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package cli

import (
	"strings"
	"time"

	"github.com/araddon/dateparse"

	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-public-api-go"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type AssetSelector struct {
	OrganizationId string
	EdgeControllerId string
	AssetIds []string
	GroupIds []string
	Labels map[string]string
}

func (a *AssetSelector) ToGRPC() *grpc_inventory_manager_go.AssetSelector {
	selector := &grpc_inventory_manager_go.AssetSelector{
		OrganizationId: a.OrganizationId,
		EdgeControllerId: a.EdgeControllerId,
		AssetIds: a.AssetIds,
		GroupIds: a.GroupIds,
		Labels: a.Labels,
	}

	return selector
}

type TimeRange struct {
	// Timestamps are strings to be parsed
	Timestamp string
	Start string
	End string
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

func (t *TimeRange) ToGRPC() *grpc_inventory_manager_go.QueryMetricsRequest_TimeRange {
	timeRange := &grpc_inventory_manager_go.QueryMetricsRequest_TimeRange{
		Timestamp: dateParse(t.Timestamp),
		TimeStart: dateParse(t.Start),
		TimeEnd: dateParse(t.End),
		Resolution: int64(t.Resolution.Seconds()),
	}

	return timeRange
}

type InventoryMonitoring struct {
	Connection
	Credentials
}

func NewInventoryMonitoring(address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *InventoryMonitoring {
	return &InventoryMonitoring{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output),
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

	aggrType, found := grpc_inventory_manager_go.AggregationType_value[aggr]
	if !found {
		methods := []string{}
		for method := range(grpc_inventory_manager_go.AggregationType_value) {
			methods = append(methods, method)
		}
		log.Fatal().Str("aggregation", aggr).Msg("Aggregation method not available. Available methods: " + strings.Join(methods, ", "))
	}

	query := &grpc_inventory_manager_go.QueryMetricsRequest{
		Assets: selector.ToGRPC(),
		Metrics: metrics,
		TimeRange: timeRange.ToGRPC(),
		Aggregation: grpc_inventory_manager_go.AggregationType(aggrType),
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

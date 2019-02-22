/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package unified_logging

import (
	"context"
	"github.com/nalej/grpc-unified-logging-go"
)

type Manager struct {
	unifiedLoggingClient grpc_unified_logging_go.CoordinatorClient
}

func NewManager(unifiedLoggingClient grpc_unified_logging_go.CoordinatorClient) Manager {
	return Manager{unifiedLoggingClient}
}

func (m *Manager) Search(request *grpc_unified_logging_go.SearchRequest) (*grpc_unified_logging_go.LogResponse, error) {
	return m.unifiedLoggingClient.Search(context.Background(), request)
}

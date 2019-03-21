/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package unified_logging

import (
	"github.com/nalej/grpc-unified-logging-go"
	"github.com/nalej/public-api/internal/pkg/server/common"
)

type Manager struct {
	unifiedLoggingClient grpc_unified_logging_go.CoordinatorClient
}

func NewManager(unifiedLoggingClient grpc_unified_logging_go.CoordinatorClient) Manager {
	return Manager{unifiedLoggingClient}
}

func (m *Manager) Search(request *grpc_unified_logging_go.SearchRequest) (*grpc_unified_logging_go.LogResponse, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.unifiedLoggingClient.Search(ctx, request)
}

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
 *
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

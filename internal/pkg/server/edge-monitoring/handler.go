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

package edge_monitoring

import (
	"context"

	"github.com/nalej/derrors"

	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-monitoring-go"
	"github.com/nalej/grpc-public-api-go"

	"github.com/nalej/grpc-utils/pkg/conversions"

	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the node requests.
type Handler struct {
	manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{
		manager: manager,
	}
}

func (h *Handler) ListMetrics(ctx context.Context, selector *grpc_inventory_go.AssetSelector) (*grpc_monitoring_go.MetricsList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if selector.GetOrganizationId() != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAssetSelector(selector)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.manager.ListMetrics(selector)
}

func (h *Handler) QueryMetrics(ctx context.Context, request *grpc_monitoring_go.QueryMetricsRequest) (*grpc_monitoring_go.QueryMetricsResult, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.GetAssets().GetOrganizationId() != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidQueryMetricsRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.manager.QueryMetrics(request)
}

func (h *Handler) ConfigureMetrics(ctx context.Context, selector *grpc_public_api_go.ConfigureMetricsRequest) (*grpc_common_go.Success, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("ConfigureMetrics not implemented"))
}

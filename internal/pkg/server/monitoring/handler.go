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
	"github.com/nalej/derrors"
	"golang.org/x/net/context"

	"github.com/nalej/grpc-monitoring-go"
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

func (h *Handler) GetClusterStats(context context.Context, request *grpc_monitoring_go.ClusterStatsRequest) (*grpc_monitoring_go.ClusterStats, error) {
	rm, err := authhelper.GetRequestMetadata(context)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.GetOrganizationId() != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}

	err = entities.ValidClusterStatsRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.manager.GetClusterStats(request)
}

func (h *Handler) GetClusterSummary(context context.Context, request *grpc_monitoring_go.ClusterSummaryRequest) (*grpc_monitoring_go.ClusterSummary, error) {
	rm, err := authhelper.GetRequestMetadata(context)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.GetOrganizationId() != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}

	err = entities.ValidClusterSummaryRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.manager.GetClusterSummary(request)
}

func (h *Handler) GetOrganizationApplicationStats(context context.Context, request *grpc_monitoring_go.OrganizationApplicationStatsRequest) (*grpc_monitoring_go.OrganizationApplicationStatsResponse, error) {
	rm, err := authhelper.GetRequestMetadata(context)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.GetOrganizationId() != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}

	err = entities.ValidOrganizationApplicationStatsRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.manager.GetOrganizationApplicationStats(request)
}

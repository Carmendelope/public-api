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

package provisioner

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-provisioner-go"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) ProvisionCluster(ctx context.Context, request *grpc_provisioner_go.ProvisionClusterRequest) (
	*grpc_provisioner_go.ProvisionClusterResponse, error) {
	return h.Manager.ProvisionCluster(request)
}

func (h *Handler) CheckProgress(ctx context.Context, request *grpc_common_go.RequestId) (
	*grpc_provisioner_go.ProvisionClusterResponse, error) {
	log.Debug().Msg("incoming check progress request")
	return h.Manager.CheckProgress(request)
}

func (h *Handler) RemoveProvision(ctx context.Context, request *grpc_common_go.RequestId) (*grpc_common_go.Success, error) {
	return h.Manager.RemoveProvision(request)
}

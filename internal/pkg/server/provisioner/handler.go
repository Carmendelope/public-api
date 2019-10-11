/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
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

func(h *Handler) ProvisionCluster(ctx context.Context, request *grpc_provisioner_go.ProvisionClusterRequest) (
    *grpc_provisioner_go.ProvisionClusterResponse, error) {
        return h.Manager.ProvisionCluster(request)
}

func(h *Handler) CheckProgress(ctx context.Context, request *grpc_common_go.RequestId) (
    *grpc_provisioner_go.ProvisionClusterResponse,error) {
        log.Debug().Msg("incoming check progress request")
        return h.Manager.CheckProgress(request)
}

func(h *Handler) RemoveProvision(ctx context.Context, request *grpc_common_go.RequestId)(*grpc_common_go.Success, error) {
    return h.Manager.RemoveProvision(request)
}
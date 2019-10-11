/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package provisioner

import (
    "context"
    "github.com/nalej/grpc-common-go"
    "github.com/nalej/grpc-provisioner-go"
)

type Manager struct {
    ProvisionerClient grpc_provisioner_go.ProvisionClient
}

func NewManager (provClient grpc_provisioner_go.ProvisionClient) Manager {
    return Manager{provClient}
}

func(m *Manager) ProvisionCluster(request *grpc_provisioner_go.ProvisionClusterRequest) (
    *grpc_provisioner_go.ProvisionClusterResponse, error) {
    return m.ProvisionerClient.ProvisionCluster(context.Background(), request)
}

func(m *Manager) CheckProgress(request *grpc_common_go.RequestId) (*grpc_provisioner_go.ProvisionClusterResponse,error) {
    return m.ProvisionerClient.CheckProgress(context.Background(), request)
}

func(m *Manager)RemoveProvision(request *grpc_common_go.RequestId)(*grpc_common_go.Success, error) {
    return m.ProvisionerClient.RemoveProvision(context.Background(), request)
}
/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package clusters

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
)

type Manager struct {

}

func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ClusterList, error) {
	panic("implement me")
}

func (m * Manager) Update(updateClusterRequest *grpc_public_api_go.UpdateClusterRequest) (*grpc_common_go.Success, error) {
	panic("implement me")
}
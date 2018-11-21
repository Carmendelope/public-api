/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
)

// ToInfraClusterUpdate transforms a public api update request into a infrastructure one.
func ToInfraClusterUpdate(update grpc_public_api_go.UpdateClusterRequest) *grpc_infrastructure_go.UpdateClusterRequest {

	result := &grpc_infrastructure_go.UpdateClusterRequest{
		OrganizationId:    update.OrganizationId,
		ClusterId:         update.ClusterId,
		UpdateName:        update.Name != "",
		Name:              update.Name,
		UpdateDescription: update.Description != "",
		Description:       update.Description,
		UpdateLabels:      update.Labels != nil,
		Labels:            update.Labels,
	}

	return result
}

func ToPublicAPICluster(source *grpc_infrastructure_go.Cluster, totalNodes int64, runningNodes int64) *grpc_public_api_go.Cluster {
	return &grpc_public_api_go.Cluster{
		OrganizationId:       source.OrganizationId,
		ClusterId:            source.ClusterId,
		Name:                 source.Name,
		Description:          source.Description,
		ClusterTypeName:      source.ClusterType.String(),
		MultitenantSupport:   source.Multitenant.String(),
		StatusName:           source.Status.String(),
		Labels:               source.Labels,
		TotalNodes:           totalNodes,
		RunningNodes:         runningNodes,
	}
}
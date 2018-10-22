/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
)

// ToInfraClusterUpdate transforms a public api update request into a infrastructure one.
func ToInfraClusterUpdate(update grpc_public_api_go.UpdateClusterRequest) * grpc_infrastructure_go.UpdateClusterRequest{

	result := &grpc_infrastructure_go.UpdateClusterRequest{
		OrganizationId:       update.OrganizationId,
		ClusterId:            update.ClusterId,
		UpdateName:           update.Name != "",
		Name:                 update.Name,
		UpdateDescription:    update.Description != "",
		Description:          update.Description,
		UpdateLabels:         update.Labels != nil,
		Labels:               update.Labels,
	}

	return result
}
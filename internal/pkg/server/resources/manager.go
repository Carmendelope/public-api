/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package resources

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
)

// Manager structure with the required clients for resources operations.
type Manager struct {

}

func (m * Manager) Summary(organizationID *grpc_organization_go.OrganizationId) (*grpc_public_api_go.ResourceSummary, error) {
	panic("implement me")
}
/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
)

type Organizations struct {
	Connection
	Credentials
}

func NewOrganizations(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Organizations {
	return &Organizations{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (o *Organizations) Info(organizationID string) *grpc_public_api_go.OrganizationInfo {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	err := o.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}

	c, err := o.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()
	ctx, cancel := o.GetContext()
	defer cancel()

	orgClient := grpc_public_api_go.NewOrganizationsClient(c)
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	info, iErr := orgClient.Info(ctx, orgID)

	o.PrintResultOrError(info, iErr, "cannot obtain organization info")
	return info
}

/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

type Organizations struct {
	Connection
	Credentials
}

func NewOrganizations(address string, port int) * Organizations {
	return &Organizations{
		Connection: *NewConnection(address, port),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (o * Organizations) Info(organizationID string) {
	err := o.LoadCredentials()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}

	c, err := o.GetConnection()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()
	ctx, cancel := o.GetContext()
	defer cancel()

	orgClient := grpc_public_api_go.NewOrganizationsClient(c)
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId:       organizationID,
	}
	info, iErr := orgClient.Info(ctx, orgID)
	if iErr != nil{
		log.Fatal().Str("trace", conversions.ToDerror(iErr).DebugReport()).Msg("cannot obtain organization info")
	}
	fmt.Println(info.String())

}



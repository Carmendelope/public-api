/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Inventory struct{
	Connection
	Credentials
}

func NewInventory(address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *Inventory {
	return &Inventory{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (i * Inventory) load() {
	err := i.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (i * Inventory) getClient() (grpc_public_api_go.InventoryClient, *grpc.ClientConn) {
	conn, err := i.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewInventoryClient(conn)
	return client, conn
}

func (i * Inventory) List(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	i.load()
	ctx, cancel := i.GetContext()
	client, conn := i.getClient()
	defer conn.Close()
	defer cancel()

	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}

	list, err := client.List(ctx, orgID)
	i.PrintResultOrError(list, err, "cannot retrieve asset list")

}


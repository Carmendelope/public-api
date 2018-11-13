/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Clusters struct {
	Connection
	Credentials
}

func NewClusters(address string, port int) * Clusters{
	return &Clusters{
		Connection: *NewConnection(address, port),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (c * Clusters) load() {
	err := c.LoadCredentials()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (c * Clusters) getClient() (grpc_public_api_go.ClustersClient, *grpc.ClientConn){
	conn, err := c.GetConnection()
	if err != nil{
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	clusterClient := grpc_public_api_go.NewClustersClient(conn)
	return clusterClient, conn
}

func (c * Clusters) List(organizationID string) {
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId:       organizationID,
	}
	clusters, err := client.List(ctx, orgID)
	if err != nil{
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot obtain cluster list")
	}
	c.PrintResult(clusters)
}

func (c* Clusters) Update(organizationID string, clusterID string, newName string, newDescription string) {
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	updateRequest := &grpc_public_api_go.UpdateClusterRequest{
		OrganizationId:       organizationID,
		ClusterId:            clusterID,
		Name:                 newName,
		Description:          newDescription,
	}
	success, err := client.Update(ctx, updateRequest)
	if err != nil{
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot obtain cluster list")
	}
	c.PrintResult(success)
}

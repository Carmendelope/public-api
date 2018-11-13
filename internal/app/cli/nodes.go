/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Nodes struct {
	Connection
	Credentials
}

func NewNodes(address string, port int) *Nodes {
	return &Nodes{
		Connection:  *NewConnection(address, port),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (n *Nodes) load() {
	err := n.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (n *Nodes) getClient() (grpc_public_api_go.NodesClient, *grpc.ClientConn) {
	conn, err := n.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	nodesClient := grpc_public_api_go.NewNodesClient(conn)
	return nodesClient, conn
}

func (n *Nodes) List(organizationID string, clusterID string) {
	n.load()
	ctx, cancel := n.GetContext()
	client, conn := n.getClient()
	defer conn.Close()
	defer cancel()

	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
	}
	list, err := client.ClusterNodes(ctx, cID)
	if err != nil {
		log.Fatal().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("cannot obtain node list")
	}
	n.PrintResult(list)
}

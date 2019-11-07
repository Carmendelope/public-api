/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package cli

import (
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Nodes struct {
	Connection
	Credentials
}

func NewNodes(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Nodes {
	return &Nodes{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
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

	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("clusterID cannot be empty")
	}

	n.load()
	ctx, cancel := n.GetContext()
	client, conn := n.getClient()
	defer conn.Close()
	defer cancel()

	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
	}
	list, err := client.List(ctx, cID)
	n.PrintResultOrError(list, err, "cannot list nodes")
}

func (n *Nodes) ModifyNodeLabels(organizationID string, nodeID string, add bool, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if nodeID == "" {
		log.Fatal().Msg("nodeID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}
	n.load()
	ctx, cancel := n.GetContext()
	client, conn := n.getClient()
	defer conn.Close()
	defer cancel()
	updateRequest := &grpc_public_api_go.UpdateNodeRequest{
		OrganizationId: organizationID,
		NodeId:         nodeID,
		AddLabels:      add,
		RemoveLabels:   !add,
		Labels:         GetLabels(rawLabels),
	}
	updated, err := client.UpdateNode(ctx, updateRequest)
	n.PrintResultOrError(updated, err, "cannot update node labels")
}

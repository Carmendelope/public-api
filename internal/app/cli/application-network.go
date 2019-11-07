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
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ApplicationNetwork struct {
	Connection
	Credentials
}

func NewApplicationNetwork(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *ApplicationNetwork {
	return &ApplicationNetwork{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (an *ApplicationNetwork) load() {
	err := an.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (an *ApplicationNetwork) getClient() (grpc_public_api_go.ApplicationNetworkClient, *grpc.ClientConn) {
	conn, err := an.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewApplicationNetworkClient(conn)
	return client, conn
}

func (an *ApplicationNetwork) AddConnection(organizationID string, sourceInstanceID string, outbound string, targetInstanceID string, inbound string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	an.load()
	ctx, cancel := an.GetContext()
	client, conn := an.getClient()
	defer conn.Close()
	defer cancel()

	connection := &grpc_application_network_go.AddConnectionRequest{
		OrganizationId:   organizationID,
		SourceInstanceId: sourceInstanceID,
		TargetInstanceId: targetInstanceID,
		InboundName:      inbound,
		OutboundName:     outbound,
	}

	added, err := client.AddConnection(ctx, connection)
	an.PrintResultOrError(added, err, "cannot add a new connection")
}

func (an *ApplicationNetwork) RemoveConnection(organizationID string, sourceInstanceID string, outbound string, targetInstanceID string, inbound string, force bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	an.load()
	ctx, cancel := an.GetContext()
	client, conn := an.getClient()
	defer conn.Close()
	defer cancel()

	connection := &grpc_application_network_go.RemoveConnectionRequest{
		OrganizationId:   organizationID,
		SourceInstanceId: sourceInstanceID,
		TargetInstanceId: targetInstanceID,
		InboundName:      inbound,
		OutboundName:     outbound,
		UserConfirmation: force,
	}

	added, err := client.RemoveConnection(ctx, connection)
	an.PrintResultOrError(added, err, "cannot remove the connection")
}

func (an *ApplicationNetwork) ListConnection(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	an.load()
	ctx, cancel := an.GetContext()
	client, conn := an.getClient()
	defer conn.Close()
	defer cancel()

	orgID := grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}

	added, err := client.ListConnections(ctx, &orgID)
	an.PrintResultOrError(added, err, "cannot list connections")
}

func (an *ApplicationNetwork) ListAvailableInbounds(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	an.load()
	ctx, cancel := an.GetContext()
	client, conn := an.getClient()
	defer conn.Close()
	defer cancel()

	orgID := grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}

	inboundList, err := client.ListAvailableInstanceInbounds(ctx, &orgID)
	an.PrintResultOrError(inboundList, err, "cannot list available inbound interfaces")
}

func (an *ApplicationNetwork) ListAvailableOutbounds(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	an.load()
	ctx, cancel := an.GetContext()
	client, conn := an.getClient()
	defer conn.Close()
	defer cancel()

	orgID := grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}

	outboundList, err := client.ListAvailableInstanceOutbounds(ctx, &orgID)
	an.PrintResultOrError(outboundList, err, "cannot list available outbound interfaces")
}

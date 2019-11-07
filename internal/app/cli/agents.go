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
	"encoding/json"
	"fmt"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"path/filepath"
)

const agetJoinTokenFile = "agentJoinToken.json"

type Agent struct {
	Connection
	Credentials
}

func NewAgent(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Agent {
	return &Agent{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (a *Agent) load() {
	err := a.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (a *Agent) getClient() (grpc_public_api_go.AgentClient, *grpc.ClientConn) {
	conn, err := a.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewAgentClient(conn)
	return client, conn
}

// CreateJoinToken request the creation of an agent join token. The result will be written into outputPath if set. If
// not the current working directory will be used.
func (a *Agent) CreateAgentJoinToken(organizationID string, edgeControllerID string, outputPath string) {

	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}

	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	eic := &grpc_inventory_go.EdgeControllerId{
		OrganizationId:   organizationID,
		EdgeControllerId: edgeControllerID,
	}
	token, err := client.CreateAgentJoinToken(ctx, eic)
	a.PrintResultOrError(token, err, "cannot create join token")
	if err == nil {
		a.writeAgentJoinToken(token, outputPath)
	}
}

// ActivateAgentMonitoring send a message to activate or deactivate the monitoring of an agent
func (a *Agent) ActivateAgentMonitoring(organizationID string, edgeControllerID string, assetID string, activate bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}
	if assetID == "" {
		log.Fatal().Msg("assetID cannot be empty")
	}

	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	request := &grpc_public_api_go.AssetMonitoringRequest{
		OrganizationId:   organizationID,
		EdgeControllerId: edgeControllerID,
		AssetId:          assetID,
		Activate:         activate,
	}

	token, err := client.ActivateMonitoring(ctx, request)
	a.PrintResultOrError(token, err, "cannot Activate Monitoring")

}

// writeJoinToken writes the EIC join token to a file so that it can be exported to the EIC.
func (a *Agent) writeAgentJoinToken(token *grpc_inventory_manager_go.AgentJoinToken, outputPath string) {
	if outputPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot determine current directory")
		}
		outputPath = cwd
	}
	outputFilePath := filepath.Join(outputPath, agetJoinTokenFile)
	marshaled, err := json.Marshal(token)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot marshal token information")
	}
	err = ioutil.WriteFile(outputFilePath, marshaled, 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("error writing agent token file")
	}
	fmt.Printf("\nAgent Token file: %s\n", outputFilePath)
}

func (a *Agent) UninstallAgent(organizationID string, assetID string, force bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	if assetID == "" {
		log.Fatal().Msg("assetID cannot be empty")
	}

	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	request := &grpc_inventory_manager_go.UninstallAgentRequest{
		OrganizationId: organizationID,
		AssetId:        assetID,
		Force:          force,
	}

	token, err := client.UninstallAgent(ctx, request)
	a.PrintResultOrError(token, err, "cannot uninstall agent")
}

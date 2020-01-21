/*
 * Copyright 2020 Nalej
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
 */

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"path/filepath"
)

type EdgeController struct {
	Connection
	Credentials
}

func NewEdgeController(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *EdgeController {
	return &EdgeController{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (ec *EdgeController) load() {
	err := ec.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (ec *EdgeController) getClient() (grpc_public_api_go.EdgeControllersClient, *grpc.ClientConn) {
	conn, err := ec.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewEdgeControllersClient(conn)
	return client, conn
}

func (ec *EdgeController) getInventoryClient() (grpc_public_api_go.InventoryClient, *grpc.ClientConn) {
	conn, err := ec.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewInventoryClient(conn)
	return client, conn
}

// CreateJoinToken request the creation of an EIC join token. The result will be written into outputPath if set. If
// not the current working directory will be used.
func (ec *EdgeController) CreateJoinToken(organizationID string, outputPath string) {

	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	ec.load()
	ctx, cancel := ec.GetContext()
	client, conn := ec.getClient()
	defer conn.Close()
	defer cancel()

	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	token, err := client.CreateEICToken(ctx, orgID)
	ec.PrintResultOrError(token, err, "cannot create join token")
	if err == nil {
		ec.writeJoinToken(token, outputPath)
	}
}

// writeJoinToken writes the EIC join token to a file so that it can be exported to the EIC.
func (ec *EdgeController) writeJoinToken(token *grpc_inventory_manager_go.EICJoinToken, outputPath string) {
	if outputPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot determine current directory")
		}
		outputPath = cwd
	}
	outputFilePath := filepath.Join(outputPath, "joinToken.json")
	marshaled, err := json.Marshal(token)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot marshal token information")
	}
	err = ioutil.WriteFile(outputFilePath, marshaled, 0600)
	if err != nil {
		log.Fatal().Err(err).Msg("error writing token file")
	}
	fmt.Printf("\nToken file: %s\n", outputFilePath)
	// TODO Add information about how to copy that in the EIC as embedded documentation.
}

func (ec *EdgeController) Unlink(organizationID string, edgeControllerID string, force bool) {

	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}

	ec.load()
	ctx, cancel := ec.GetContext()
	client, conn := ec.getClient()
	defer conn.Close()
	defer cancel()

	request := &grpc_inventory_manager_go.UnlinkECRequest{
		OrganizationId:   organizationID,
		EdgeControllerId: edgeControllerID,
		Force:            force,
	}
	success, err := client.UnlinkEIC(ctx, request)
	ec.PrintResultOrError(success, err, "cannot unlink edge controller")

}

func (ec *EdgeController) getInstallCredentials(username string, password string, publicKeyPath string, isSudoer bool) *grpc_inventory_manager_go.SSHCredentials {

	credentials := &grpc_inventory_manager_go.SSHCredentials{
		Username: username,
		IsSudoer: isSudoer,
	}

	if publicKeyPath != "" {

		path := GetPath(publicKeyPath)
		log.Debug().Str("publicKeyPath", path).Msg("loading public key from file")
		publicKey, err := ioutil.ReadFile(path)
		if err != nil {
			log.Fatal().Str("publicKeyPath", path).Msg("cannot load public key file")
		}

		credentials.Credentials = &grpc_inventory_manager_go.SSHCredentials_ClientCertificate{
			ClientCertificate: string(publicKey),
		}
	} else {
		credentials.Credentials = &grpc_inventory_manager_go.SSHCredentials_Password{
			Password: password,
		}
	}

	return credentials
}

func (ec *EdgeController) InstallAgent(organizationID string, edgeControllerID string, agentType grpc_inventory_manager_go.AgentType, targetHost string, username string, password string,
	publicKeyPath string, isSudoer bool) {

	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}
	if targetHost == "" {
		log.Fatal().Msg("targetHost cannot be empty")
	}
	if username == "" {
		log.Fatal().Msg("username cannot be empty")
	}
	if password == "" && publicKeyPath == "" {
		log.Fatal().Msg("either password or public key must be specified")
	}

	credentials := ec.getInstallCredentials(username, password, publicKeyPath, isSudoer)
	installRequest := &grpc_inventory_manager_go.InstallAgentRequest{
		OrganizationId:   organizationID,
		EdgeControllerId: edgeControllerID,
		AgentType:        agentType,
		Credentials:      credentials,
		TargetHost:       targetHost,
	}

	ec.load()
	ctx, cancel := ec.GetContext()
	client, conn := ec.getClient()
	defer conn.Close()
	defer cancel()

	result, err := client.InstallAgent(ctx, installRequest)
	ec.PrintResultOrError(result, err, "cannot trigger the install of an agent")

}

func (ec *EdgeController) UpdateGeolocation(organizationID string, edgeControllerID string, geolocation string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}
	ec.load()
	ctx, cancel := ec.GetContext()
	client, conn := ec.getClient()
	defer conn.Close()
	defer cancel()

	updateRequest := &grpc_inventory_manager_go.UpdateGeolocationRequest{
		OrganizationId:   organizationID,
		EdgeControllerId: edgeControllerID,
		Geolocation:      geolocation,
	}

	_, err := client.UpdateGeolocation(ctx, updateRequest)
	ec.PrintResultOrError(&grpc_common_go.Success{}, err, "cannot update geolocation")
}

func (ec *EdgeController) Update(organizationID string, edgeControllerID string, addLabel bool, removeLabel bool, labels map[string]string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}
	if addLabel == removeLabel {
		log.Fatal().Msg("cannot add and remove labels in the same operation")
	}

	ec.load()
	ctx, cancel := ec.GetContext()
	client, conn := ec.getInventoryClient()
	defer conn.Close()
	defer cancel()

	updateRequest := &grpc_inventory_go.UpdateEdgeControllerRequest{
		OrganizationId:   organizationID,
		EdgeControllerId: edgeControllerID,
		AddLabels:        addLabel,
		RemoveLabels:     removeLabel,
		Labels:           labels,
	}

	_, err := client.UpdateEdgeController(ctx, updateRequest)
	ec.PrintResultOrError(&grpc_common_go.Success{}, err, "cannot update ec")
}

func (ec *EdgeController) getECLabelRequest(organizationID string, edgeControllerID string, rawLabels string, addLabels bool) *grpc_inventory_go.UpdateEdgeControllerRequest {
	labels := GetLabels(rawLabels)
	return &grpc_inventory_go.UpdateEdgeControllerRequest{
		OrganizationId:      organizationID,
		EdgeControllerId:    edgeControllerID,
		AddLabels:           addLabels,
		RemoveLabels:        !addLabels,
		Labels:              labels,
		UpdateLastAlive:     false,
		UpdateGeolocation:   false,
		UpdateLastOpSummary: false,
	}
}

func (ec *EdgeController) AddLabelToEC(organizationID string, edgeControllerID string, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}

	ec.load()
	ctx, cancel := ec.GetContext()
	client, conn := ec.getInventoryClient()
	defer conn.Close()
	defer cancel()

	request := ec.getECLabelRequest(organizationID, edgeControllerID, rawLabels, true)

	success, err := client.UpdateEdgeController(ctx, request)
	ec.PrintResultOrError(success, err, "cannot add labels to EC")

}

func (ec *EdgeController) RemoveLabelFromEC(organizationID string, edgeControllerID string, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}

	ec.load()
	ctx, cancel := ec.GetContext()
	client, conn := ec.getInventoryClient()
	defer conn.Close()
	defer cancel()

	request := ec.getECLabelRequest(organizationID, edgeControllerID, rawLabels, false)

	success, err := client.UpdateEdgeController(ctx, request)
	ec.PrintResultOrError(success, err, "cannot remove labels from EC")
}

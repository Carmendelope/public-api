/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package cli

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"path/filepath"
)

type EdgeController struct{
	Connection
	Credentials
}

func NewEdgeController(address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *EdgeController {
	return &EdgeController{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (ec*EdgeController) load() {
	err := ec.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (ec*EdgeController) getClient() (grpc_public_api_go.EdgeControllersClient, *grpc.ClientConn) {
	conn, err := ec.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewEdgeControllersClient(conn)
	return client, conn
}

// CreateJoinToken request the creation of an EIC join token. The result will be written into outputPath if set. If
// not the current workding directory will be used.
func (ec * EdgeController) CreateJoinToken(organizationID string, outputPath string) {

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
	if err != nil{
		ec.writeJoinToken(token, outputPath)
	}
}

// writeJoinToken writes the EIC join token to a file so that it can be exported to the EIC.
func (ec * EdgeController) writeJoinToken(token * grpc_inventory_manager_go.EICJoinToken, outputPath string){
	if outputPath == "" {
		cwd, err  :=  os.Getwd()
		if err != nil{
			log.Fatal().Err(err).Msg("cannot determine current directory")
		}
		outputPath = cwd
	}
	outputFilePath := filepath.Join(outputPath, "joinToken.json")
	marshaled, err := json.Marshal(token)
	if err != nil{
		log.Fatal().Err(err).Msg("cannot marshal token information")
	}
	err = ioutil.WriteFile(outputFilePath, marshaled, 0600)
	if err != nil{
		log.Fatal().Err(err).Msg("error writing token file")
	}
	fmt.Printf("\nToken file: %s\n", outputFilePath)
	// TODO Add information about how to copy that in the EIC as embedded documentation.
}
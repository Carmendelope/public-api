/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Asset struct{
	Connection
	Credentials
}

func NewAsset (address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *Asset {
	return &Asset{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (a * Asset) load() {
	err := a.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (a * Asset) getClient() (grpc_public_api_go.InventoryClient, *grpc.ClientConn) {
	conn, err := a.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewInventoryClient(conn)
	return client, conn
}

func (a *Asset) UpdateLocation (organizationID string, assetID string, location string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if assetID == "" {
		log.Fatal().Msg("assetID cannot be empty")
	}
	if location == "" {
		log.Fatal().Msg("location cannot be empty")
	}
	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	updateRequest := &grpc_inventory_go.UpdateAssetRequest{
		OrganizationId: organizationID,
		AssetId: assetID,
		UpdateLocation: true,
		Location: &grpc_inventory_go.InventoryLocation{
			Geolocation: location,
		},
	}

	_, err := client.UpdateAsset (ctx, updateRequest)
	a.PrintResultOrError(&grpc_common_go.Success{}, err, "cannot update location")
}

func (a *Asset) Update (organizationID string, assetID string, addLabel bool, removeLabel bool, labels map[string]string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if assetID == "" {
		log.Fatal().Msg("assetID cannot be empty")
	}
	if addLabel == removeLabel {
		log.Fatal().Msg("cannot add and remove labels in the same operation")
	}

	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	updateRequest := &grpc_inventory_go.UpdateAssetRequest{
		OrganizationId:       organizationID,
		AssetId:              assetID,
		AddLabels:            addLabel,
		RemoveLabels:         removeLabel,
		Labels:               labels,
	}

	_, err := client.UpdateAsset (ctx, updateRequest)
	a.PrintResultOrError(&grpc_common_go.Success{}, err, "cannot update asset")
}

func (a*Asset) getAssetLabelRequest(organizationID string, assetID string, rawLabels string, addLabels bool) *grpc_inventory_go.UpdateAssetRequest{
	labels := GetLabels(rawLabels)
	return &grpc_inventory_go.UpdateAssetRequest{
		OrganizationId:       organizationID,
		AssetId:              assetID,
		AddLabels:            addLabels,
		RemoveLabels:         !addLabels,
		Labels:               labels,
		UpdateLastOpSummary:  false,
		UpdateLastAlive:      false,
		UpdateIp:             false,
		UpdateLocation:       false,
	}
}

func (a*Asset) AddLabelToAsset(organizationID string, assetID string, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if assetID == "" {
		log.Fatal().Msg("assetID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}

	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	request := a.getAssetLabelRequest(organizationID, assetID, rawLabels, true)

	success, err := client.UpdateAsset(ctx, request)
	a.PrintResultOrError(success, err, "cannot add labels to asset")

}

func (a*Asset) RemoveLabelFromAsset(organizationID string, assetID string, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if assetID == "" {
		log.Fatal().Msg("assetID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}

	a.load()
	ctx, cancel := a.GetContext()
	client, conn := a.getClient()
	defer conn.Close()
	defer cancel()

	request := a.getAssetLabelRequest(organizationID, assetID, rawLabels, false)

	success, err := client.UpdateAsset(ctx, request)
	a.PrintResultOrError(success, err, "cannot remove labels from asset")
}

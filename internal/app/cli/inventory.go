/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Inventory struct {
	Connection
	Credentials
}

func NewInventory(address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *Inventory {
	return &Inventory{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (i *Inventory) load() {
	err := i.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (i *Inventory) getClient() (grpc_public_api_go.InventoryClient, *grpc.ClientConn) {
	conn, err := i.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewInventoryClient(conn)
	return client, conn
}

func (i *Inventory) List(organizationID string) {
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

	inventory, err := client.List(ctx, orgID)
	i.PrintResultOrError(inventory, err, "cannot retrieve inventory list")

}

func (i *Inventory) GetControllerExtendedInfo(organizationID string, edgeControllerID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if edgeControllerID == "" {
		log.Fatal().Msg("edgeControllerID cannot be empty")
	}
	i.load()
	ctx, cancel := i.GetContext()
	client, conn := i.getClient()
	defer conn.Close()
	defer cancel()
	controllerID := &grpc_inventory_go.EdgeControllerId{
		OrganizationId:   organizationID,
		EdgeControllerId: edgeControllerID,
	}
	info, err := client.GetControllerExtendedInfo(ctx, controllerID)
	i.PrintResultOrError(info, err, "cannot extended controller information")
}

func (i *Inventory) GetAssetInfo(organizationID string, assetID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if assetID == "" {
		log.Fatal().Msg("assetID cannot be empty")
	}
	i.load()
	ctx, cancel := i.GetContext()
	client, conn := i.getClient()
	defer conn.Close()
	defer cancel()
	id := &grpc_inventory_go.AssetId{
		OrganizationId: organizationID,
		AssetId:        assetID,
	}
	info, err := client.GetAssetInfo(ctx, id)
	i.PrintResultOrError(info, err, "cannot asset information")
}

func (i *Inventory) GetDeviceInfo(organizationID string, deviceID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if deviceID == "" {
		log.Fatal().Msg("deviceID cannot be empty")
	}
	i.load()
	ctx, cancel := i.GetContext()
	client, conn := i.getClient()
	defer conn.Close()
	defer cancel()
	id := &grpc_inventory_manager_go.DeviceId{
		OrganizationId:	organizationID,
		AssetDeviceId:	deviceID,
	}
	info, err := client.GetDeviceInfo(ctx, id)
	i.PrintResultOrError(info, err, "cannot get device information")
}

func (i * Inventory) UpdateDeviceLocation (organizationID string, assetDeviceID string, location string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if assetDeviceID == "" {
		log.Fatal().Msg("deviceID cannot be empty")
	}
	if location == "" {
		log.Fatal().Msg("location cannot be empty")
	}
	i.load()
	ctx, cancel := i.GetContext()
	client, conn := i.getClient()
	defer conn.Close()
	defer cancel()

	request := &grpc_inventory_manager_go.UpdateDeviceLocationRequest{
		OrganizationId:	organizationID,
		AssetDeviceId:	assetDeviceID,
		UpdateLocation: true,
		Location: &grpc_inventory_go.InventoryLocation{
			Geolocation: location,
		},
	}

	info, err := client.UpdateDeviceLocation(ctx, request)
	i.PrintResultOrError(info, err, "cannot update device location")
}
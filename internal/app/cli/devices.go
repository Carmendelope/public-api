/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"strings"
)

type Devices struct{
	Connection
	Credentials
}

func NewDevices(address string, port int, insecure bool, caCertPath string) *Devices {
	return &Devices{
		Connection:  *NewConnection(address, port, insecure, caCertPath),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (d*Devices) load() {
	err := d.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (d*Devices) getClient() (grpc_public_api_go.DevicesClient, *grpc.ClientConn) {
	conn, err := d.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	client := grpc_public_api_go.NewDevicesClient(conn)
	return client, conn
}

func (d*Devices) AddDeviceGroup(organizationID string, name string, enabled bool, defaultConnectivity bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if name == "" {
		log.Fatal().Msg("name cannot be empty")
	}

	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()

	addRequest := &grpc_device_manager_go.AddDeviceGroupRequest{
		OrganizationId:            organizationID,
		Name:                      name,
		Enabled:                   enabled,
		DeviceDefaultConnectivity: defaultConnectivity,
	}
	added, err := client.AddDeviceGroup(ctx, addRequest)
	d.PrintResultOrError(added, err, "cannot add device group")
}

func (d*Devices) UpdateDeviceGroup(organizationID string, deviceGroupID string, updateEnabled bool, enabled bool, updateDefaultConnectivity bool, defaultConnectivy bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if deviceGroupID == "" {
		log.Fatal().Msg("deviceGroupID cannot be empty")
	}
	if !updateEnabled && !updateDefaultConnectivity{
		log.Fatal().Msg("Either updateEnabled or updateDefaultConnectivity must be set")
	}
	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()
	updateRequest := &grpc_device_manager_go.UpdateDeviceGroupRequest{
		OrganizationId:            organizationID,
		DeviceGroupId:             deviceGroupID,
		UpdateEnabled:             updateEnabled,
		Enabled:                   enabled,
		UpdateDeviceConnectivity:  updateDefaultConnectivity,
		DefaultDeviceConnectivity: defaultConnectivy,
	}
	updated, err := client.UpdateDeviceGroup(ctx, updateRequest)
	d.PrintResultOrError(updated, err, "cannot update device group")
}

func (d*Devices) RemoveDeviceGroup(organizationID string, deviceGroupID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if deviceGroupID == "" {
		log.Fatal().Msg("deviceGroupID cannot be empty")
	}

	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()

	dgID := &grpc_device_go.DeviceGroupId{
		OrganizationId:       organizationID,
		DeviceGroupId:        deviceGroupID,
	}
	_, err := client.RemoveDeviceGroup(ctx, dgID)
	d.PrintSuccessOrError(err, "cannot remove device group", "device group has been removed")
}

func (d*Devices) ListDeviceGroups(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()

	oID := &grpc_organization_go.OrganizationId{
		OrganizationId:       organizationID,
	}
	dgs, err := client.ListDeviceGroups(ctx, oID)
	d.PrintResultOrError(dgs, err, "cannot list device groups")
}

func (d*Devices) ListDevices(organizationID string, deviceGroupID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if deviceGroupID == "" {
		log.Fatal().Msg("deviceGroupID cannot be empty")
	}

	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()

	dgID := &grpc_device_go.DeviceGroupId{
		OrganizationId:       organizationID,
		DeviceGroupId:        deviceGroupID,
	}
	devices, err := client.ListDevices(ctx, dgID)
	d.PrintResultOrError(devices, err, "cannot list devices")
}

func (d*Devices) getDeviceLabelRequest(organizationID string, deviceGroupID string, deviceID string, rawLabels string) *grpc_device_manager_go.DeviceLabelRequest{
	labels := make(map[string]string, 0)

	split := strings.Split(rawLabels, ";")
	for _, l := range split{
		ls := strings.Split(l, ":")
		if len(ls) != 2{
			log.Fatal().Str("label", l).Msg("malformed label, expecting key:value")
		}
		labels[ls[0]] = ls[1]
	}

	return &grpc_device_manager_go.DeviceLabelRequest{
		OrganizationId:       organizationID,
		DeviceGroupId:        deviceGroupID,
		DeviceId:             deviceID,
		Labels:               labels,
	}
}

func (d*Devices) AddLabelToDevice(organizationID string, deviceGroupID string, deviceID string, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if deviceGroupID == "" {
		log.Fatal().Msg("deviceGroupID cannot be empty")
	}
	if deviceID == "" {
		log.Fatal().Msg("deviceID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}

	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()

	request := d.getDeviceLabelRequest(organizationID, deviceGroupID, deviceID, rawLabels)

	_, err := client.AddLabelToDevice(ctx, request)
	d.PrintSuccessOrError(err, "cannot add labels to device", "labels have been added")

}

func (d*Devices) RemoveLabelFromDevice(organizationID string, deviceGroupID string, deviceID string, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if deviceGroupID == "" {
		log.Fatal().Msg("deviceGroupID cannot be empty")
	}
	if deviceID == "" {
		log.Fatal().Msg("deviceID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}

	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()

	request := d.getDeviceLabelRequest(organizationID, deviceGroupID, deviceID, rawLabels)

	_, err := client.RemoveLabelFromDevice(ctx, request)
	d.PrintSuccessOrError(err, "cannot remove labels from device", "labels have been removed")
}

func (d*Devices) UpdateDevice(organizationID string, deviceGroupID string, deviceID string, enabled bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if deviceGroupID == "" {
		log.Fatal().Msg("deviceGroupID cannot be empty")
	}
	if deviceID == "" {
		log.Fatal().Msg("deviceID cannot be empty")
	}

	d.load()
	ctx, cancel := d.GetContext()
	client, conn := d.getClient()
	defer conn.Close()
	defer cancel()

	request := &grpc_device_manager_go.UpdateDeviceRequest{
		OrganizationId:       organizationID,
		DeviceGroupId:        deviceGroupID,
		DeviceId:             deviceID,
		Enabled:              enabled,
	}

	updated, err := client.UpdateDevice(ctx, request)
	d.PrintResultOrError(updated, err, "cannot update device")

}

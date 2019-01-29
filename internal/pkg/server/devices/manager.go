/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package devices

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
)

// Manager structure with the required clients for node operations.
type Manager struct {
	deviceClient grpc_device_manager_go.DevicesClient
}

// NewManager creates a Manager using a set of clients.
func NewManager(deviceClient grpc_device_manager_go.DevicesClient) Manager {
	return Manager{
		deviceClient: deviceClient,
	}
}

func (m *Manager) AddDeviceGroup(request *grpc_device_manager_go.AddDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error) {
	return m.deviceClient.AddDeviceGroup(context.Background(), request)
}

func (m *Manager) UpdateDeviceGroup(request *grpc_device_manager_go.UpdateDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error) {
	return m.deviceClient.UpdateDeviceGroup(context.Background(), request)
}

func (m *Manager) RemoveDeviceGroup(request *grpc_device_go.DeviceGroupId) (*grpc_common_go.Success, error) {
	return m.deviceClient.RemoveDeviceGroup(context.Background(), request)
}

func (m *Manager) ListDeviceGroups(request *grpc_organization_go.OrganizationId) (*grpc_device_manager_go.DeviceGroupList, error) {
	return m.deviceClient.ListDeviceGroups(context.Background(), request)
}

func (m *Manager) ListDevices(request *grpc_device_go.DeviceGroupId) (*grpc_device_manager_go.DeviceList, error) {
	return m.deviceClient.ListDevices(context.Background(), request)
}

func (m *Manager) AddLabelToDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	return m.deviceClient.AddLabelToDevice(context.Background(), request)
}

func (m *Manager) RemoveLabelFromDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	return m.deviceClient.RemoveLabelFromDevice(context.Background(), request)
}

func (m *Manager) UpdateDevice(request *grpc_device_manager_go.UpdateDeviceRequest) (*grpc_device_manager_go.Device, error) {
	return m.deviceClient.UpdateDevice(context.Background(), request)
}

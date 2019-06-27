/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package devices

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/public-api/internal/pkg/entities"
	"github.com/nalej/public-api/internal/pkg/server/common"
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
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.deviceClient.AddDeviceGroup(ctx, request)
}

func (m *Manager) UpdateDeviceGroup(request *grpc_device_manager_go.UpdateDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.deviceClient.UpdateDeviceGroup(ctx, request)
}

func (m *Manager) RemoveDeviceGroup(request *grpc_device_go.DeviceGroupId) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.deviceClient.RemoveDeviceGroup(ctx, request)
}

func (m *Manager) ListDeviceGroups(request *grpc_organization_go.OrganizationId) (*grpc_device_manager_go.DeviceGroupList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.deviceClient.ListDeviceGroups(ctx, request)
}

func (m *Manager) ListDevices(request *grpc_device_go.DeviceGroupId) (*grpc_public_api_go.DeviceList, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	list, err := m.deviceClient.ListDevices(ctx, request)
	if err != nil {
		return nil, err
	}

	return entities.ToPublicAPIDeviceList(list), nil

}

func (m *Manager) AddLabelToDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.deviceClient.AddLabelToDevice(ctx, request)
}

func (m *Manager) RemoveLabelFromDevice(request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.deviceClient.RemoveLabelFromDevice(ctx, request)
}

func (m *Manager) UpdateDevice(request *grpc_device_manager_go.UpdateDeviceRequest) (*grpc_public_api_go.Device, error) {
	ctx, cancel := common.GetContext()
	defer cancel()
	device, err :=  m.deviceClient.UpdateDevice(ctx, request)
	if err != nil {
		return nil, err
	}
	return entities.ToPublicAPIDevice(device), nil
}

func (m*Manager) RemoveDevice(deviceID *grpc_device_go.DeviceId) (*grpc_common_go.Success, error){
	ctx, cancel := common.GetContext()
	defer cancel()
	return m.deviceClient.RemoveDevice(ctx, deviceID)
}

func (m*Manager) GetDevice (deviceID *grpc_device_go.DeviceId) (*grpc_public_api_go.Device, error){
	ctx, cancel := common.GetContext()
	defer cancel()
	dmDevice, err := m.deviceClient.GetDevice(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	return entities.ToPublicAPIDevice(dmDevice), nil
}
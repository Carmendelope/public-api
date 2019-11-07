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

package devices

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/public-api/internal/pkg/authhelper"
	"github.com/nalej/public-api/internal/pkg/entities"
)

// Handler structure for the node requests.
type Handler struct {
	Manager Manager
}

func (h *Handler) GetDevice(ctx context.Context, deviceID *grpc_device_go.DeviceId) (*grpc_public_api_go.Device, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if deviceID.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidDeviceID(deviceID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.GetDevice(deviceID)
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) AddDeviceGroup(ctx context.Context, request *grpc_device_manager_go.AddDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidAddDeviceGroupRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.AddDeviceGroup(request)
}

func (h *Handler) UpdateDeviceGroup(ctx context.Context, request *grpc_device_manager_go.UpdateDeviceGroupRequest) (*grpc_device_manager_go.DeviceGroup, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUpdateDeviceGroupRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.UpdateDeviceGroup(request)
}

func (h *Handler) RemoveDeviceGroup(ctx context.Context, request *grpc_device_go.DeviceGroupId) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidDeviceGroupID(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.RemoveDeviceGroup(request)
}

func (h *Handler) ListDeviceGroups(ctx context.Context, request *grpc_organization_go.OrganizationId) (*grpc_device_manager_go.DeviceGroupList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidOrganizationId(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.ListDeviceGroups(request)
}

func (h *Handler) ListDevices(ctx context.Context, request *grpc_device_go.DeviceGroupId) (*grpc_public_api_go.DeviceList, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidDeviceGroupID(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.ListDevices(request)
}

func (h *Handler) AddLabelToDevice(ctx context.Context, request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidDeviceLabelRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.AddLabelToDevice(request)
}

func (h *Handler) RemoveLabelFromDevice(ctx context.Context, request *grpc_device_manager_go.DeviceLabelRequest) (*grpc_common_go.Success, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidDeviceLabelRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.RemoveLabelFromDevice(request)
}

func (h *Handler) UpdateDevice(ctx context.Context, request *grpc_device_manager_go.UpdateDeviceRequest) (*grpc_public_api_go.Device, error) {
	rm, err := authhelper.GetRequestMetadata(ctx)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	if request.OrganizationId != rm.OrganizationID {
		return nil, derrors.NewPermissionDeniedError("cannot access requested OrganizationID")
	}
	err = entities.ValidUpdateDeviceRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return h.Manager.UpdateDevice(request)
}

func (h *Handler) RemoveDevice(ctx context.Context, deviceID *grpc_device_go.DeviceId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidDeviceID(deviceID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.RemoveDevice(deviceID)
}

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
 */

package entities

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-log-download-manager-go"
	"github.com/nalej/grpc-monitoring-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-provisioner-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/rs/zerolog/log"
	"github.com/santhosh-tekuri/jsonschema"
	"strings"
	"sync"
)

const emptyOrganizationId = "organization_id cannot be empty"
const emptyInstanceId = "app_instance_id cannot be empty"
const emptyDescriptorId = "app_descriptor_id cannot be empty"
const emptyServiceGroupId = "service_group_id cannot be empty"
const emptyServiceGroupInstanceId = "service_group_instance_id cannot be empty"
const emptyServiceId = "service_id cannot be empty"
const emptyClusterId = "cluster_id cannot be empty"
const emptyNodeId = "node_id cannot be empty"
const emptyEmail = "email cannot be empty"
const emptyName = "name cannot be empty"
const emptyPassword = "password cannot be empty"
const emptyNewPassword = "new password cannot be empty"
const emptyRoleName = "role_name cannot be empty"
const emptyRoleID = "role_id cannot be empty"
const emptyDeviceGroupId = "device_group_id cannot be empty"
const emptyDeviceId = "device_id cannot be empty"
const emptyDeviceGroupApiKey = "device_group_api_key cannot be empty"
const emptyLabels = "labels cannot be empty"
const invalidSortOrder = "sort order can only be ascending or descending"
const emptyEdgeControllerId = "edge_controller_id cannot be empty"
const emptyAssetId = "asset_id cannot be empty"
const emptyAssetDeviceId = "asset_device_id cannot be empty"
const emptyGeolocation = "geolocation cannot be empty"

const emptySourceInstanceId = "source_instance_id cannot be empty"
const emptyTargetInstanceId = "target_instance_id cannot be empty"
const emptyInboundName = "inbound_name cannot be empty"
const emptyOutboundName = "outbound_name cannot be empty"

const emptyRequestId = "request_id cannot be empty"

const emptyKey = "key cannot be empty"


// --------- Application descriptor JSON Schema
type AppJSONSchema struct {
	// Singleton object used to validate application descriptors
	appDescriptorSchema *jsonschema.Schema
	// Singleton value
	singletonValidator sync.Once
}

// -------------------------------------------

// Local instance for the application descriptor validator
var AppDescValidator = AppJSONSchema{}

// Initialize the local AppDescValidator reading the schema from the filePath. This is a single run operation.
func InitializeJSON() derrors.Error {
	var err error
	AppDescValidator.singletonValidator.Do(func() {
		log.Debug().Msg("loading application descriptor validator schema...")
		compiler := jsonschema.NewCompiler()
		schemaURL := "http://nalej.com/app_descriptor.json"
		if derr := compiler.AddResource(schemaURL, strings.NewReader(APP_DESC_SCHEMA)); err != nil {
			log.Error().Err(err).Msg("impossible to add JSON schema definition")
			err = derr
			return
		}

		schema, schemaErr := compiler.Compile(schemaURL)
		if schemaErr != nil {
			log.Error().Err(err).Msg("impossible to load json schema for application descriptors")
			err = schemaErr
			return
		}
		AppDescValidator.appDescriptorSchema = schema
		log.Debug().Msg("schema validator loaded")
	})
	if err != nil {
		return derrors.NewInvalidArgumentError("impossible to load json schema for application descriptors", err)
	}
	return nil
}

func ValidOrganizationId(organizationID *grpc_organization_go.OrganizationId) derrors.Error {
	if organizationID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidUpdateOrganizationRequest(request *grpc_organization_go.UpdateOrganizationRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if !request.UpdateName && !request.UpdateFullAddress && !request.UpdateCity && !request.UpdateState &&
		!request.UpdateCountry && !request.UpdateZipCode && !request.UpdatePhoto {
		return derrors.NewInvalidArgumentError("some field must be modified")
	}
	return nil
}

func ValidClusterId(clusterID *grpc_infrastructure_go.ClusterId) derrors.Error {
	if clusterID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if clusterID.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	return nil
}

func ValidUserId(userID *grpc_user_go.UserId) derrors.Error {
	if userID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if userID.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAppInstanceID(appInstanceID *grpc_application_go.AppInstanceId) derrors.Error {
	if appInstanceID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if appInstanceID.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyInstanceId)
	}
	return nil
}

// Validate that the JSON descriptor for the application follows the current JSONSchema
func ValidAppDescriptorFormat(jsonContent []byte) derrors.Error {

	// Initialize JSON in case it is not working
	InitializeJSON()

	// TODO: commented until inbound-outbounds are consolidated∫
	//err := AppDescValidator.appDescriptorSchema.Validate(bytes.NewReader(jsonContent))

	//if err != nil {
	//	return derrors.NewInvalidArgumentError(err.Error())
	//}
	return nil
}

func ValidAppDescriptorID(appDescriptorID *grpc_application_go.AppDescriptorId) derrors.Error {
	if appDescriptorID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if appDescriptorID.AppDescriptorId == "" {
		return derrors.NewInvalidArgumentError(emptyDescriptorId)
	}
	return nil
}

func ValidUpdateAppDescriptor(request *grpc_application_go.UpdateAppDescriptorRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AppDescriptorId == "" {
		return derrors.NewInvalidArgumentError(emptyDescriptorId)
	}
	if request.AddLabels && request.RemoveLabels {
		return derrors.NewInvalidArgumentError("add_labels and remove_labels cannot be set at the same time")
	}
	if (request.AddLabels || request.RemoveLabels) && (len(request.Labels) == 0) {
		return derrors.NewInvalidArgumentError(emptyLabels)
	}
	return nil
}

func ValidUpdateClusterRequest(request *grpc_public_api_go.UpdateClusterRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	if request.AddLabels && request.RemoveLabels {
		return derrors.NewInvalidArgumentError("add_labels and remove_labels cannot be set at the same time")
	}
	if (request.AddLabels || request.RemoveLabels) && (len(request.Labels) == 0) {
		return derrors.NewInvalidArgumentError(emptyLabels)
	}
	return nil
}

func ValidUpdateNodeRequest(request *grpc_public_api_go.UpdateNodeRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.NodeId == "" {
		return derrors.NewInvalidArgumentError(emptyNodeId)
	}
	if request.AddLabels && request.RemoveLabels {
		return derrors.NewInvalidArgumentError("add_labels and remove_labels cannot be set at the same time")
	}
	if (request.AddLabels || request.RemoveLabels) && (len(request.Labels) == 0) {
		return derrors.NewInvalidArgumentError(emptyLabels)
	}
	return nil
}

func ValidAddUserRequest(request *grpc_public_api_go.AddUserRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	if request.Password == "" {
		return derrors.NewInvalidArgumentError(emptyPassword)
	}
	if request.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	if request.RoleName == "" {
		return derrors.NewInvalidArgumentError(emptyRoleName)
	}
	return nil
}

func ValidUpdateUserRequest(updateUserRequest *grpc_user_go.UpdateUserRequest) derrors.Error {
	if updateUserRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if updateUserRequest.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAddAppDescriptor(request *grpc_application_go.AddAppDescriptorRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if len(request.Groups) == 0 {
		return derrors.NewInvalidArgumentError("expecting at least one service group")
	}
	for _, g := range request.Groups {
		if len(g.Services) == 0 {
			return derrors.NewInvalidArgumentError(fmt.Sprintf("group %s has no services", g.Name))
		}

	}

	// NP-872. Check the device_ids is empty
	for _, rule := range request.Rules {
		if len(rule.DeviceGroupIds) > 0 {
			return derrors.NewInvalidArgumentError(fmt.Sprintf("rule %s cannot have device_group_ids", rule.Name))
		}
	}

	return nil
}

func ValidDeployRequest(request *grpc_application_manager_go.DeployRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AppDescriptorId == "" {
		return derrors.NewInvalidArgumentError(emptyDescriptorId)
	}
	return nil
}

func ValidUndeployRequest(request *grpc_application_manager_go.UndeployRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyInstanceId)
	}
	return nil
}

func ValidInstallRequest(request *grpc_public_api_go.InstallRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidScaleClusterRequest(request *grpc_provisioner_go.ScaleClusterRequest) derrors.Error {
	if request.RequestId != "" {
		return derrors.NewInvalidArgumentError("request_id is set by infrastructure-manager")
	}
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	if request.AzureCredentials == nil {
		return derrors.NewInvalidArgumentError("azure_credentials cannot be empty")
	}
	if request.AzureOptions == nil || request.AzureOptions.ResourceGroup == "" {
		return derrors.NewInvalidArgumentError("azure_options.resource_group cannot be empty")
	}

	return nil
}

func ValidUninstallClusterRequest(request *grpc_public_api_go.UninstallClusterRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	if request.KubeConfigRaw == "" {
		return derrors.NewInvalidArgumentError("kube_config_raw cannot be empty")
	}
	return nil
}

func ValidDecomissionClusterRequest(request *grpc_public_api_go.DecomissionClusterRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	return nil
}

func ValidChangePasswordRequest(request *grpc_user_manager_go.ChangePasswordRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.Password == "" {
		return derrors.NewInvalidArgumentError(emptyPassword)
	}
	if request.NewPassword == "" {
		return derrors.NewInvalidArgumentError(emptyNewPassword)
	}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidAssignRoleRequest(request *grpc_user_manager_go.AssignRoleRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.RoleId == "" {
		return derrors.NewInvalidArgumentError(emptyRoleID)
	}
	if request.Email == "" {
		return derrors.NewInvalidArgumentError(emptyEmail)
	}
	return nil
}

func ValidDeviceGroupID(deviceGroupID *grpc_device_go.DeviceGroupId) derrors.Error {
	if deviceGroupID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if deviceGroupID.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	return nil
}

func ValidDeviceID(deviceId *grpc_device_go.DeviceId) derrors.Error {
	if deviceId.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if deviceId.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if deviceId.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}
	return nil
}

func ValidAddDeviceGroupRequest(request *grpc_device_manager_go.AddDeviceGroupRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	return nil
}

func ValidUpdateDeviceGroupRequest(request *grpc_device_manager_go.UpdateDeviceGroupRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if !request.UpdateEnabled && !request.UpdateDeviceConnectivity {
		return derrors.NewInvalidArgumentError("either update_enabled or update_device_connectivity must be set")
	}
	return nil
}

func ValidRegisterDeviceRequest(request *grpc_device_manager_go.RegisterDeviceRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if request.DeviceGroupApiKey == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupApiKey)
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}

	return nil
}

func ValidDeviceLabelRequest(request *grpc_device_manager_go.DeviceLabelRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}
	if len(request.Labels) == 0 {
		return derrors.NewInvalidArgumentError(emptyLabels)
	}

	return nil
}

func ValidUpdateDeviceRequest(request *grpc_device_manager_go.UpdateDeviceRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceGroupId)
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyDeviceId)
	}

	return nil
}

func ValidSearchRequest(request *grpc_public_api_go.SearchRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}

	// validate the field dependencies,
	if request.ServiceGroupId != "" && request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyInstanceId)
	}
	if request.ServiceGroupInstanceId != "" {
		if request.AppInstanceId == "" {
			return derrors.NewInvalidArgumentError(emptyInstanceId)
		} else if request.ServiceGroupId == "" {
			return derrors.NewInvalidArgumentError(emptyServiceGroupId)
		}
	}
	if request.ServiceId != "" {
		if request.AppInstanceId == "" {
			return derrors.NewInvalidArgumentError(emptyInstanceId)
		} else if request.ServiceGroupId == "" {
			return derrors.NewInvalidArgumentError(emptyServiceGroupId)
		} else if request.ServiceGroupInstanceId == "" {
			return derrors.NewInvalidArgumentError(emptyServiceGroupInstanceId)
		}
	}
	if request.ServiceInstanceId != "" {
		if request.AppInstanceId == "" {
			return derrors.NewInvalidArgumentError(emptyInstanceId)
		} else if request.ServiceGroupId == "" {
			return derrors.NewInvalidArgumentError(emptyServiceGroupId)
		} else if request.ServiceGroupInstanceId == "" {
			return derrors.NewInvalidArgumentError(emptyServiceGroupInstanceId)
		} else if request.ServiceId == "" {
			return derrors.NewInvalidArgumentError(emptyServiceId)
		}
	}

	if request.Order != nil {
		if request.Order.Order != grpc_public_api_go.Order_ASC &&
			request.Order.Order != grpc_public_api_go.Order_DESC {
			return derrors.NewInvalidArgumentError(invalidSortOrder)
		}
	}
	return nil
}

func ValidDownloadRequestId(request *grpc_log_download_manager_go.DownloadRequestId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.RequestId == "" {
		return derrors.NewInvalidArgumentError(emptyRequestId)
	}

	return nil
}

func ValidDownloadLogRequest(request *grpc_log_download_manager_go.DownloadLogRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}

	return nil
}

func ValidClusterStatsRequest(request *grpc_monitoring_go.ClusterStatsRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyInstanceId)
	}

	return nil
}

func ValidClusterSummaryRequest(request *grpc_monitoring_go.ClusterSummaryRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyInstanceId)
	}

	return nil
}

func ValidOrganizationApplicationStatsRequest(request *grpc_monitoring_go.OrganizationApplicationStatsRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}

	return nil
}

func ValidEdgeControllerID(edgeControllerID *grpc_inventory_go.EdgeControllerId) derrors.Error {
	if edgeControllerID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if edgeControllerID.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}
	return nil
}

func ValidUnlinkECRequest(request *grpc_inventory_manager_go.UnlinkECRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}
	return nil
}

func ValidInstallAgentRequest(request *grpc_inventory_manager_go.InstallAgentRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError("edge_controller_id cannot be empty")
	}
	if request.TargetHost == "" {
		return derrors.NewInvalidArgumentError("target_host cannot be empty")
	}
	return nil
}

func ValidAssetID(assetID *grpc_inventory_go.AssetId) derrors.Error {
	if assetID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if assetID.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}
	return nil
}

func ValidUninstallAgentRequest(request *grpc_inventory_manager_go.UninstallAgentRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}
	return nil
}

func ValidAssetMonitoringRequest(request *grpc_public_api_go.AssetMonitoringRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}
	return nil
}

func ValidUpdateGeolocationRequest(request *grpc_inventory_manager_go.UpdateGeolocationRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}

	return nil
}

func ValidUpdateDeviceLocationRequest(request *grpc_inventory_manager_go.UpdateDeviceLocationRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AssetDeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetDeviceId)
	}
	if request.Location != nil && request.Location.Geolocation == "" {
		return derrors.NewInvalidArgumentError(emptyGeolocation)
	}
	return nil
}

func ValidAssetSelector(selector *grpc_inventory_go.AssetSelector) derrors.Error {
	if selector == nil {
		return derrors.NewInvalidArgumentError("empty asset selector")
	}
	if selector.GetOrganizationId() == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	return nil
}

func ValidTimeRange(timeRange *grpc_monitoring_go.QueryMetricsRequest_TimeRange) derrors.Error {
	if !(timeRange.GetTimestamp() == 0) {
		if timeRange.GetTimeStart() != 0 || timeRange.GetTimeEnd() != 0 || timeRange.GetResolution() != 0 {
			return derrors.NewInvalidArgumentError("timestamp is set; start, end and resolution should be 0").
				WithParams(timeRange.GetTimestamp(), timeRange.GetTimeStart(),
					timeRange.GetTimeEnd(), timeRange.GetResolution())
		}
	} else {
		if timeRange.GetTimeStart() == 0 && timeRange.GetTimeEnd() == 0 {
			return derrors.NewInvalidArgumentError("timestamp is not set; either start, end or both should be set").
				WithParams(timeRange.GetTimestamp(), timeRange.GetTimeStart(),
					timeRange.GetTimeEnd(), timeRange.GetResolution())
		}
	}

	return nil
}

func ValidQueryMetricsRequest(request *grpc_monitoring_go.QueryMetricsRequest) derrors.Error {
	// We check the asset selector so we know we have an organization ID.
	derr := ValidAssetSelector(request.GetAssets())
	if derr != nil {
		return derr
	}

	// Check the time range to either be a point in time or a range
	derr = ValidTimeRange(request.GetTimeRange())
	if derr != nil {
		return derr
	}

	if len(request.GetAssets().GetAssetIds()) != 1 && request.GetAggregation() == grpc_monitoring_go.AggregationType_NONE {
		return derrors.NewInvalidArgumentError("metrics for more than one asset requested without aggregation method")
	}

	return nil
}

func ValidDeviceId(request *grpc_inventory_manager_go.DeviceId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AssetDeviceId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetDeviceId)
	}
	return nil
}

func ValidUpdateAssetRequest(request *grpc_inventory_go.UpdateAssetRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}

	return nil
}

func ValidUpdateEdgeControllerRequest(request *grpc_inventory_go.UpdateEdgeControllerRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}

	return nil
}

func ValidAddConnectionRequest(conn *grpc_application_network_go.AddConnectionRequest) derrors.Error {
	if conn.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if conn.TargetInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyTargetInstanceId)
	}
	if conn.SourceInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptySourceInstanceId)
	}
	if conn.InboundName == "" {
		return derrors.NewInvalidArgumentError(emptyInboundName)
	}
	if conn.OutboundName == "" {
		return derrors.NewInvalidArgumentError(emptyOutboundName)
	}

	return nil
}

func ValidRemoveConnectionRequest(conn *grpc_application_network_go.RemoveConnectionRequest) derrors.Error {
	if conn.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if conn.TargetInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyTargetInstanceId)
	}
	if conn.SourceInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptySourceInstanceId)
	}
	if conn.InboundName == "" {
		return derrors.NewInvalidArgumentError(emptyInboundName)
	}
	if conn.OutboundName == "" {
		return derrors.NewInvalidArgumentError(emptyOutboundName)
	}

	return nil
}


func ValidUpdateSettingRequest(updateRequest *grpc_public_api_go.UpdateSettingRequest) derrors.Error {
	if updateRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if updateRequest.Key == "" {
		return derrors.NewInvalidArgumentError(emptyKey)
	}
	return nil
}

func ValidListRequest(request *grpc_public_api_go.ListRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}

	return nil
}
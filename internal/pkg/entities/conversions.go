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

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
)

// ToInfraClusterUpdate transforms a public api update request into a infrastructure one.
func ToInfraClusterUpdate(update grpc_public_api_go.UpdateClusterRequest) *grpc_infrastructure_go.UpdateClusterRequest {

	result := &grpc_infrastructure_go.UpdateClusterRequest{
		OrganizationId:                   update.OrganizationId,
		ClusterId:                        update.ClusterId,
		UpdateName:                       update.UpdateName,
		Name:                             update.Name,
		AddLabels:                        update.AddLabels,
		RemoveLabels:                     update.RemoveLabels,
		Labels:                           update.Labels,
		UpdateMillicoresConversionFactor: update.UpdateMillicoresConversionFactor,
		MillicoresConversionFactor:       update.MillicoresConversionFactor,
	}

	return result
}

func ToPublicAPICluster(source *grpc_infrastructure_go.Cluster, totalNodes int64, runningNodes int64) *grpc_public_api_go.Cluster {
	return &grpc_public_api_go.Cluster{
		OrganizationId:             source.OrganizationId,
		ClusterId:                  source.ClusterId,
		Name:                       source.Name,
		ClusterTypeName:            source.ClusterType.String(),
		MultitenantSupport:         source.Multitenant.String(),
		StatusName:                 source.ClusterStatus.String(),
		Status:                     source.ClusterStatus,
		Labels:                     source.Labels,
		TotalNodes:                 totalNodes,
		RunningNodes:               runningNodes,
		LastAliveTimestamp:         source.LastAliveTimestamp,
		MillicoresConversionFactor: source.MillicoresConversionFactor,
		State:                      source.State,
		StateName:                  source.State.String(),
	}
}

func ToPublicAPINode(source *grpc_infrastructure_go.Node) *grpc_public_api_go.Node {
	return &grpc_public_api_go.Node{
		OrganizationId: source.OrganizationId,
		ClusterId:      source.ClusterId,
		NodeId:         source.NodeId,
		Ip:             source.Ip,
		Labels:         source.Labels,
		StatusName:     source.Status.String(),
		StateName:      source.State.String(),
	}
}

func ToPublicAPIEndpoints(source []*grpc_application_go.Endpoint) []*grpc_public_api_go.Endpoint {
	result := make([]*grpc_public_api_go.Endpoint, 0)
	for _, e := range source {
		toAdd := &grpc_public_api_go.Endpoint{
			TypeName: e.Type.String(),
			Path:     e.Path,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIPorts(source []*grpc_application_go.Port) []*grpc_public_api_go.Port {
	result := make([]*grpc_public_api_go.Port, 0)
	for _, p := range source {
		toAdd := &grpc_public_api_go.Port{
			Name:         p.Name,
			InternalPort: p.InternalPort,
			ExposedPort:  p.ExposedPort,
			Endpoints:    ToPublicAPIEndpoints(p.Endpoints),
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIStorage(source []*grpc_application_go.Storage) []*grpc_public_api_go.Storage {
	result := make([]*grpc_public_api_go.Storage, 0)
	for _, s := range source {
		toAdd := &grpc_public_api_go.Storage{
			Size:      s.Size,
			MountPath: s.MountPath,
			TypeName:  s.Type.String(),
		}
		result = append(result, toAdd)
	}
	return result
}

func hideCredentials(credentials *grpc_application_go.ImageCredentials) *grpc_application_go.ImageCredentials {

	return &grpc_application_go.ImageCredentials{
		Username:         credentials.Username,
		Password:         "redacted",
		Email:            credentials.Email,
		DockerRepository: credentials.DockerRepository,
	}
}

func ToPublicAPIServiceInstances(source []*grpc_application_go.ServiceInstance) []*grpc_public_api_go.ServiceInstance {
	result := make([]*grpc_public_api_go.ServiceInstance, 0)

	for _, si := range source {
		endpoints := make([]string, len(si.Endpoints))
		for i, e := range si.Endpoints {
			endpoints[i] = e.Fqdn
		}
		credentials := si.Credentials
		if credentials != nil {
			credentials = hideCredentials(credentials)
		}

		toAdd := &grpc_public_api_go.ServiceInstance{
			OrganizationId:         si.OrganizationId,
			AppDescriptorId:        si.AppDescriptorId,
			AppInstanceId:          si.AppInstanceId,
			ServiceGroupId:         si.ServiceGroupId,
			ServiceGroupInstanceId: si.ServiceGroupInstanceId,
			ServiceId:              si.ServiceId,
			ServiceInstanceId:      si.ServiceInstanceId,
			Name:                   si.Name,
			TypeName:               si.Type.String(),
			Image:                  si.Image,
			Credentials:            credentials,
			Specs:                  si.Specs,
			Storage:                ToPublicAPIStorage(si.Storage),
			ExposedPorts:           ToPublicAPIPorts(si.ExposedPorts),
			EnvironmentVariables:   si.EnvironmentVariables,
			Configs:                si.Configs,
			Labels:                 si.Labels,
			DeployAfter:            si.DeployAfter,
			StatusName:             si.Status.String(),
			Endpoints:              endpoints,
			DeployedOnClusterId:    si.DeployedOnClusterId,
			RunArguments:           si.RunArguments,
			Info:                   si.Info,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIInstanceMetadata(metadata *grpc_application_go.InstanceMetadata) *grpc_public_api_go.InstanceMetadata {

	if metadata == nil {
		return nil
	}

	status := make(map[string]string, 0)
	for key, value := range metadata.Status {
		status[key] = value.String()
	}
	instance := &grpc_public_api_go.InstanceMetadata{
		OrganizationId:      metadata.OrganizationId,
		AppDescriptorId:     metadata.AppDescriptorId,
		AppInstanceId:       metadata.AppInstanceId,
		MonitoredInstanceId: metadata.MonitoredInstanceId,
		TypeName:            metadata.Type.String(),
		InstancesId:         metadata.InstancesId,
		DesiredReplicas:     metadata.DesiredReplicas,
		AvailableReplicas:   metadata.AvailableReplicas,
		UnavailableReplicas: metadata.UnavailableReplicas,
		StatusName:          status,
		Info:                metadata.Info,
	}
	return instance
}

func ToPublicAPIGroupInstances(source []*grpc_application_go.ServiceGroupInstance) []*grpc_public_api_go.ServiceGroupInstance {
	result := make([]*grpc_public_api_go.ServiceGroupInstance, 0)
	for _, sgi := range source {
		serviceInstance := make([]*grpc_public_api_go.ServiceInstance, len(sgi.ServiceInstances))
		serviceInstance = ToPublicAPIServiceInstances(sgi.ServiceInstances)
		var spec *grpc_public_api_go.ServiceGroupDeploymentSpecs
		if sgi.Specs != nil {
			spec = &grpc_public_api_go.ServiceGroupDeploymentSpecs{
				Replicas:            sgi.Specs.Replicas,
				MultiClusterReplica: sgi.Specs.MultiClusterReplica,
				DeploymentSelectors: sgi.Specs.DeploymentSelectors,
			}
		}

		toAdd := &grpc_public_api_go.ServiceGroupInstance{
			OrganizationId:         sgi.OrganizationId,
			AppDescriptorId:        sgi.AppDescriptorId,
			AppInstanceId:          sgi.AppInstanceId,
			ServiceGroupId:         sgi.ServiceGroupId,
			ServiceGroupInstanceId: sgi.ServiceGroupInstanceId,
			Name:                   sgi.Name,
			ServiceInstances:       serviceInstance,
			PolicyName:             sgi.Policy.String(),
			StatusName:             sgi.Status.String(),
			Metadata:               ToPublicAPIInstanceMetadata(sgi.Metadata),
			Specs:                  spec,
			Labels:                 sgi.Labels,
			GlobalFqdn:             sgi.GlobalFqdn,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPISecurityRules(source []*grpc_application_go.SecurityRule) []*grpc_public_api_go.SecurityRule {

	result := make([]*grpc_public_api_go.SecurityRule, 0)
	for _, sr := range source {
		toAdd := &grpc_public_api_go.SecurityRule{
			OrganizationId:           sr.OrganizationId,
			AppDescriptorId:          sr.AppDescriptorId,
			RuleId:                   sr.RuleId,
			Name:                     sr.Name,
			TargetServiceGroupName:   sr.TargetServiceGroupName,
			TargetServiceName:        sr.TargetServiceName,
			TargetPort:               sr.TargetPort,
			AccessName:               sr.Access.String(),
			AuthServiceGroupName:     sr.AuthServiceGroupName,
			AuthServices:             sr.AuthServices,
			DeviceGroupIds:           sr.DeviceGroupIds,
			DeviceGroupNames:         sr.DeviceGroupNames,
			InboundNetInterfaceName:  sr.InboundNetInterface,
			OutboundNetInterfaceName: sr.OutboundNetInterface,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIAppInstance(source *grpc_application_manager_go.AppInstance) *grpc_public_api_go.AppInstance {

	metadata := make([]*grpc_public_api_go.InstanceMetadata, 0)
	for _, met := range source.Metadata {
		metadata = append(metadata, ToPublicAPIInstanceMetadata(met))
	}

	return &grpc_public_api_go.AppInstance{
		OrganizationId:        source.OrganizationId,
		AppDescriptorId:       source.AppDescriptorId,
		AppInstanceId:         source.AppInstanceId,
		Name:                  source.Name,
		ConfigurationOptions:  source.ConfigurationOptions,
		EnvironmentVariables:  source.EnvironmentVariables,
		Labels:                source.Labels,
		Rules:                 ToPublicAPISecurityRules(source.Rules),
		Groups:                ToPublicAPIGroupInstances(source.Groups),
		StatusName:            source.Status.String(),
		Metadata:              metadata,
		Info:                  source.Info,
		InboundNetInterfaces:  source.InboundNetInterfaces,
		OutboundNetInterfaces: source.OutboundNetInterfaces,
		InboundConnections: ToPublicAPIConnectionList(&grpc_application_network_go.ConnectionInstanceList{
			Connections: source.InboundConnections,
		}).List,
		OutboundConnections: ToPublicAPIConnectionList(&grpc_application_network_go.ConnectionInstanceList{
			Connections: source.OutboundConnections,
		}).List,
	}
}

func ToPublicAPIAssetInfo(assetInfo *grpc_inventory_go.AssetInfo) *grpc_public_api_go.AssetInfo {
	if assetInfo == nil {
		return nil
	}
	return &grpc_public_api_go.AssetInfo{
		Hardware: assetInfo.Hardware,
		Storage:  assetInfo.Storage,
		Os:       ToPublicAPIAssetOS(assetInfo.Os),
	}
}

func ToPublicAPIDevice(device *grpc_device_manager_go.Device) *grpc_public_api_go.Device {
	return &grpc_public_api_go.Device{
		OrganizationId:   device.OrganizationId,
		DeviceGroupId:    device.DeviceGroupId,
		DeviceId:         device.DeviceId,
		RegisterSince:    device.RegisterSince,
		Labels:           device.Labels,
		Enabled:          device.Enabled,
		DeviceStatusName: device.DeviceStatus.String(),
		Location:         device.Location,
		AssetInfo:        ToPublicAPIAssetInfo(device.AssetInfo),
	}
}

func InventoryDeviceToPublicAPIDevice(device *grpc_inventory_manager_go.Device) *grpc_public_api_go.Device {
	return &grpc_public_api_go.Device{
		OrganizationId:   device.OrganizationId,
		DeviceGroupId:    device.DeviceGroupId,
		DeviceId:         device.DeviceId,
		AssetDeviceId:    device.AssetDeviceId,
		RegisterSince:    device.RegisterSince,
		Labels:           device.Labels,
		Enabled:          device.Enabled,
		DeviceStatusName: device.DeviceStatus.String(),
		Location:         device.Location,
		AssetInfo:        ToPublicAPIAssetInfo(device.AssetInfo),
	}
}

func ToPublicAPIDeviceList(list *grpc_device_manager_go.DeviceList) *grpc_public_api_go.DeviceList {
	result := make([]*grpc_public_api_go.Device, 0)
	for _, device := range list.Devices {
		toAdd := ToPublicAPIDevice(device)
		result = append(result, toAdd)
	}
	return &grpc_public_api_go.DeviceList{
		Devices: result,
	}
}

func ToPublicAPIDeviceArray(devices []*grpc_inventory_manager_go.Device) []*grpc_public_api_go.Device {
	result := make([]*grpc_public_api_go.Device, 0)
	for _, device := range devices {
		toAdd := InventoryDeviceToPublicAPIDevice(device)
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIAppParam(param *grpc_application_go.AppParameter) *grpc_public_api_go.AppParameter {
	if param == nil {
		return nil
	}
	return &grpc_public_api_go.AppParameter{
		Name:         param.Name,
		Description:  param.Description,
		Path:         param.Path,
		Type:         param.Type.String(),
		DefaultValue: param.DefaultValue,
		EnumValues:   param.EnumValues,
		Category:     param.Category.String(),
	}
}

func ToPublicAPIAssetOS(os *grpc_inventory_go.OperatingSystemInfo) *grpc_public_api_go.OperatingSystemInfo {
	if os == nil {
		return nil
	}
	return &grpc_public_api_go.OperatingSystemInfo{
		Name:         os.Name,
		Version:      os.Version,
		Class:        os.Class,
		ClassName:    os.Class.String(),
		Architecture: os.Architecture,
	}
}
func ToPublicApiAgentOpSummary(opSummary *grpc_inventory_go.AgentOpSummary) *grpc_public_api_go.AgentOpSummary {
	if opSummary == nil {
		return nil
	}
	return &grpc_public_api_go.AgentOpSummary{
		OperationId:  opSummary.OperationId,
		Timestamp:    opSummary.Timestamp,
		Status:       opSummary.Status,
		OpStatusName: opSummary.Status.String(),
		Info:         opSummary.Info,
	}
}

func ToPublicAPIAsset(asset *grpc_inventory_manager_go.Asset) *grpc_public_api_go.Asset {
	if asset == nil {
		return nil
	}
	return &grpc_public_api_go.Asset{
		OrganizationId:     asset.OrganizationId,
		EdgeControllerId:   asset.EdgeControllerId,
		AssetId:            asset.AssetId,
		AgentId:            asset.AgentId,
		Created:            asset.Created,
		Labels:             asset.Labels,
		Os:                 ToPublicAPIAssetOS(asset.Os),
		Hardware:           asset.Hardware,
		Storage:            asset.Storage,
		EicNetIp:           asset.EicNetIp,
		LastOpSummary:      ToPublicApiAgentOpSummary(asset.LastOpSummary),
		LastAliveTimestamp: asset.LastAliveTimestamp,
		Status:             asset.Status,
		StatusName:         asset.Status.String(),
		Location:           asset.Location,
	}
}

func ToPublicAPIAssetArray(assets []*grpc_inventory_manager_go.Asset) []*grpc_public_api_go.Asset {
	result := make([]*grpc_public_api_go.Asset, 0)
	for _, asset := range assets {
		if asset.Show {
			toAdd := ToPublicAPIAsset(asset)
			result = append(result, toAdd)
		}
	}
	return result
}

func ToPublicApiECOpSummary(opSummary *grpc_inventory_go.ECOpSummary) *grpc_public_api_go.ECOpSummary {
	if opSummary == nil {
		return nil
	}
	return &grpc_public_api_go.ECOpSummary{
		OperationId:  opSummary.OperationId,
		Timestamp:    opSummary.Timestamp,
		Status:       opSummary.Status,
		OpStatusName: opSummary.Status.String(),
		Info:         opSummary.Info,
	}
}

func ToPublicAPIController(controller *grpc_inventory_manager_go.EdgeController) *grpc_public_api_go.EdgeController {
	return &grpc_public_api_go.EdgeController{
		OrganizationId:     controller.OrganizationId,
		EdgeControllerId:   controller.EdgeControllerId,
		Created:            controller.Created,
		Name:               controller.Name,
		Labels:             controller.Labels,
		LastAliveTimestamp: controller.LastAliveTimestamp,
		Status:             controller.Status,
		StatusName:         controller.Status.String(),
		Location:           controller.Location,
		LastOpResult:       ToPublicApiECOpSummary(controller.LastOpResult),
		AssetInfo:          ToPublicAPIAssetInfo(controller.AssetInfo),
	}
}

func ToPublicAPIECOPResponse(response *grpc_inventory_manager_go.EdgeControllerOpResponse) *grpc_public_api_go.ECOpResponse {
	return &grpc_public_api_go.ECOpResponse{
		OrganizationId:   response.OperationId,
		EdgeControllerId: response.EdgeControllerId,
		OperationId:      response.OperationId,
		Timestamp:        response.Timestamp,
		Status:           response.Status.String(),
		Info:             response.Info,
	}
}

func ToPublicAPIControllerArray(controllers []*grpc_inventory_manager_go.EdgeController) []*grpc_public_api_go.EdgeController {
	result := make([]*grpc_public_api_go.EdgeController, 0)
	for _, controller := range controllers {
		if controller.Show {
			toAdd := ToPublicAPIController(controller)
			result = append(result, toAdd)
		}
	}
	return result
}

func ToPublicAPIAgentOpRequest(response *grpc_inventory_manager_go.AgentOpResponse) *grpc_public_api_go.AgentOpResponse {
	if response == nil {
		return nil
	}

	return &grpc_public_api_go.AgentOpResponse{
		OrganizationId:   response.OperationId,
		EdgeControllerId: response.EdgeControllerId,
		AssetId:          response.AssetId,
		OperationId:      response.OperationId,
		Timestamp:        response.Timestamp,
		Status:           response.Status.String(),
		Info:             response.Info,
	}
}

func ToPublicAPIOpResponse(response *grpc_common_go.OpResponse) *grpc_public_api_go.OpResponse {
	return &grpc_public_api_go.OpResponse{
		OrganizationId: response.OrganizationId,
		RequestId:      response.RequestId,
		OperationName:  response.OperationName,
		ElapsedTime:    response.ElapsedTime,
		Timestamp:      response.Timestamp,
		Status:         response.Status,
		StatusName:     response.Status.String(),
		Info:           response.Info,
		Error:          response.Error,
	}
}

func ToPublicAPIConnectionList(connectionInstanceList *grpc_application_network_go.ConnectionInstanceList) *grpc_public_api_go.ConnectionInstanceList {
	publicConnections := make([]*grpc_public_api_go.ConnectionInstance, len(connectionInstanceList.Connections))
	for i, connection := range connectionInstanceList.Connections {
		publicConnections[i] = &grpc_public_api_go.ConnectionInstance{
			OrganizationId:     connection.OrganizationId,
			ConnectionId:       connection.ConnectionId,
			SourceInstanceId:   connection.SourceInstanceId,
			SourceInstanceName: connection.SourceInstanceName,
			TargetInstanceId:   connection.TargetInstanceId,
			TargetInstanceName: connection.TargetInstanceName,
			InboundName:        connection.InboundName,
			OutboundName:       connection.OutboundName,
			OutboundRequired:   connection.OutboundRequired,
			StatusName:         connection.Status.String(),
			IpRange:            connection.IpRange,
			ZtNetworkId:        connection.ZtNetworkId,
		}
	}
	return &grpc_public_api_go.ConnectionInstanceList{List: publicConnections}
}

func NewSearchRequest(request *grpc_public_api_go.SearchRequest) *grpc_application_manager_go.SearchRequest {

	return &grpc_application_manager_go.SearchRequest{
		OrganizationId:         request.OrganizationId,
		AppDescriptorId:        request.AppDescriptorId,
		AppInstanceId:          request.AppInstanceId,
		ServiceGroupId:         request.ServiceGroupId,
		ServiceGroupInstanceId: request.ServiceGroupInstanceId,
		ServiceId:              request.ServiceId,
		ServiceInstanceId:      request.ServiceInstanceId,
		MsgQueryFilter:         request.MsgQueryFilter,
		From:                   request.From,
		To:                     request.To,
		IncludeMetadata:        true,
		NFirst:                 request.NFirst,
	}
}

func ToInstallerTargetPlatform(pbPlatform grpc_public_api_go.Platform) (*grpc_installer_go.Platform, derrors.Error) {
	var installerPlatform grpc_installer_go.Platform
	switch pbPlatform {
	case grpc_public_api_go.Platform_AZURE:
		installerPlatform = grpc_installer_go.Platform_AZURE
	case grpc_public_api_go.Platform_MINIKUBE:
		installerPlatform = grpc_installer_go.Platform_MINIKUBE
	default:
		log.Warn().Str("platform", pbPlatform.String()).Msg("unknown platform")
		return nil, derrors.NewInvalidArgumentError("unknown platform").WithParams(pbPlatform.String())
	}
	return &installerPlatform, nil
}

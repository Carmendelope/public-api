/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-inventory-manager-go"
	"github.com/nalej/grpc-public-api-go"
)

// ToInfraClusterUpdate transforms a public api update request into a infrastructure one.
func ToInfraClusterUpdate(update grpc_public_api_go.UpdateClusterRequest) *grpc_infrastructure_go.UpdateClusterRequest {

	result := &grpc_infrastructure_go.UpdateClusterRequest{
		OrganizationId:    update.OrganizationId,
		ClusterId:         update.ClusterId,
		UpdateName:        update.Name != "",
		Name:              update.Name,
		AddLabels: update.AddLabels,
		RemoveLabels: update.RemoveLabels,
		Labels:            update.Labels,
	}

	return result
}

func ToPublicAPICluster(source *grpc_infrastructure_go.Cluster, totalNodes int64, runningNodes int64) *grpc_public_api_go.Cluster {
	return &grpc_public_api_go.Cluster{
		OrganizationId:     source.OrganizationId,
		ClusterId:          source.ClusterId,
		Name:               source.Name,
		ClusterTypeName:    source.ClusterType.String(),
		MultitenantSupport: source.Multitenant.String(),
		StatusName:         source.Status.String(),
		Labels:             source.Labels,
		TotalNodes:         totalNodes,
		RunningNodes:       runningNodes,
		Cordon:				source.Cordon,
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


func hideCredentials (credentials * grpc_application_go.ImageCredentials) *grpc_application_go.ImageCredentials {

	return &grpc_application_go.ImageCredentials{
		Username: credentials.Username,
		Password: "redacted",
		Email: credentials.Email,
		DockerRepository:credentials.DockerRepository,
	}
}

func ToPublicAPIServiceInstances(source []*grpc_application_go.ServiceInstance) []*grpc_public_api_go.ServiceInstance {
	result := make([]*grpc_public_api_go.ServiceInstance, 0)


	for _, si := range source {
		endpoints := make([]string,len(si.Endpoints))
		for i,e := range si.Endpoints {
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

func ToPublicAPIInstanceMetadata (metadata * grpc_application_go.InstanceMetadata) *grpc_public_api_go.InstanceMetadata {

	if metadata == nil{
		return nil
	}

	status := make(map[string]string, 0)
	for key, value := range metadata.Status {
		status[key] = value.String()
	}
	instance := &grpc_public_api_go.InstanceMetadata{
		OrganizationId:		metadata.OrganizationId,
		AppDescriptorId:	metadata.AppDescriptorId,
		AppInstanceId: 		metadata.AppInstanceId,
		MonitoredInstanceId:metadata.MonitoredInstanceId,
		TypeName:           metadata.Type.String(),
		InstancesId:        metadata.InstancesId,
		DesiredReplicas:    metadata.DesiredReplicas,
		AvailableReplicas:  metadata.AvailableReplicas,
		UnavailableReplicas:metadata.UnavailableReplicas,
		StatusName:         status,
		Info:               metadata.Info,
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
				Replicas:         sgi.Specs.Replicas,
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
			OrganizationId:  sr.OrganizationId,
			AppDescriptorId: sr.AppDescriptorId,
			RuleId:          sr.RuleId,
			Name:            sr.Name,
			TargetServiceGroupName: sr.TargetServiceGroupName,
			TargetServiceName: sr.TargetServiceName,
			TargetPort:      sr.TargetPort,
			AccessName:      sr.Access.String(),
			AuthServiceGroupName: sr.AuthServiceGroupName,
			AuthServices:    sr.AuthServices,
			DeviceGroupNames:sr.DeviceGroupNames,
			DeviceGroupIds:  sr.DeviceGroupIds,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIAppInstance(source *grpc_application_go.AppInstance) *grpc_public_api_go.AppInstance {

	metadata := make ([]*grpc_public_api_go.InstanceMetadata, 0)
	for _, met := range source.Metadata {
		metadata = append(metadata, ToPublicAPIInstanceMetadata(met))
	}

	return &grpc_public_api_go.AppInstance{
		OrganizationId:       source.OrganizationId,
		AppDescriptorId:      source.AppDescriptorId,
		AppInstanceId:        source.AppInstanceId,
		Name:                 source.Name,
		ConfigurationOptions: source.ConfigurationOptions,
		EnvironmentVariables: source.EnvironmentVariables,
		Labels:               source.Labels,
		Rules:                ToPublicAPISecurityRules(source.Rules),
		Groups:               ToPublicAPIGroupInstances(source.Groups),
		StatusName:           source.Status.String(),
		Metadata:			  metadata,
		Info: 				  source.Info,
	}
}

func ToPublicAPIDevice(device * grpc_device_manager_go.Device) * grpc_public_api_go.Device  {
	return &grpc_public_api_go.Device{
		OrganizationId: device.OrganizationId,
		DeviceGroupId: device.DeviceGroupId,
		DeviceId: device.DeviceId,
		RegisterSince: device.RegisterSince,
		Labels: device.Labels,
		Enabled: device.Enabled,
		DeviceStatusName: device.DeviceStatus.String(),
	}
}

func InventoryDeviceToPublicAPIDevice(device * grpc_inventory_manager_go.Device) * grpc_public_api_go.Device  {
	return &grpc_public_api_go.Device{
		OrganizationId: device.OrganizationId,
		DeviceGroupId: device.DeviceGroupId,
		DeviceId: device.DeviceId,
		AssetDeviceId: device.AssetDeviceId,
		RegisterSince: device.RegisterSince,
		Labels: device.Labels,
		Enabled: device.Enabled,
		DeviceStatusName: device.DeviceStatus.String(),
	}
}

func ToPublicAPIDeviceList(list * grpc_device_manager_go.DeviceList) * grpc_public_api_go.DeviceList  {
	result := make([]*grpc_public_api_go.Device, 0)
	for _, device := range list.Devices {
		toAdd := ToPublicAPIDevice(device)
		result = append(result, toAdd)
	}
	return & grpc_public_api_go.DeviceList {
		Devices: result,
	}
}

func ToPublicAPIDeviceArray(devices [] * grpc_inventory_manager_go.Device) [] * grpc_public_api_go.Device  {
	result := make([]*grpc_public_api_go.Device, 0)
	for _, device := range devices {
		toAdd := InventoryDeviceToPublicAPIDevice(device)
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIAppParam (param *grpc_application_go.AppParameter) * grpc_public_api_go.AppParameter {
	return &grpc_public_api_go.AppParameter{
		Name: 			param.Name,
		Description:	param.Description,
		Path: 			param.Path,
		Type: 			param.Type.String(),
		DefaultValue: 	param.DefaultValue,
		EnumValues: 	param.EnumValues,
		Category:    	param.Category.String(),
	}
}

func ToPublicAPIAssetOS(os * grpc_inventory_go.OperatingSystemInfo) * grpc_public_api_go.OperatingSystemInfo{
	return &grpc_public_api_go.OperatingSystemInfo{
		Name:                 os.Name,
		Version:              os.Version,
		Class:                os.Class,
		ClassName:            os.Class.String(),
		Architecture:         os.Architecture,
	}
}

func ToPublicAPIAsset(asset *grpc_inventory_manager_go.Asset) * grpc_public_api_go.Asset{
	return &grpc_public_api_go.Asset{
		OrganizationId:       asset.OrganizationId,
		EdgeControllerId:     asset.EdgeControllerId,
		AssetId:              asset.AssetId,
		AgentId:              asset.AgentId,
		Created:              asset.Created,
		Labels:               asset.Labels,
		Os:                   ToPublicAPIAssetOS(asset.Os),
		Hardware:             asset.Hardware,
		Storage:              asset.Storage,
		EicNetIp:             asset.EicNetIp,
		LastOpSummary:        asset.LastOpSummary,
		LastAliveTimestamp:   asset.LastAliveTimestamp,
		Status:               asset.Status,
		StatusName:           asset.Status.String(),
	}
}

func ToPublicAPIAssetArray(assets [] * grpc_inventory_manager_go.Asset) [] * grpc_public_api_go.Asset  {
	result := make([]*grpc_public_api_go.Asset, 0)
	for _, asset := range assets {
		if asset.Show{
			toAdd := ToPublicAPIAsset(asset)
			result = append(result, toAdd)
		}
	}
	return result
}

func ToPublicAPIController(controller *grpc_inventory_manager_go.EdgeController) *grpc_public_api_go.EdgeController{
	return &grpc_public_api_go.EdgeController{
		OrganizationId:       controller.OrganizationId,
		EdgeControllerId:     controller.EdgeControllerId,
		Created:              controller.Created,
		Name:                 controller.Name,
		Labels:               controller.Labels,
		LastAliveTimestamp:   controller.LastAliveTimestamp,
		Status:               controller.Status,
		StatusName:           controller.Status.String(),
		Location:             controller.Location,
	}
}

func ToPublicAPIControllerArray(controllers []*grpc_inventory_manager_go.EdgeController) [] *grpc_public_api_go.EdgeController{
	result := make([]*grpc_public_api_go.EdgeController, 0)
	for _, controller := range controllers {
		if controller.Show{
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
		OrganizationId: response.OperationId,
		EdgeControllerId: response.EdgeControllerId,
		AssetId: response.AssetId,
		OperationId:response.OperationId,
		Timestamp: response.Timestamp,
		Status: response.Status.String(),
		Info: response.Info,
	}
}
/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
)

// ToInfraClusterUpdate transforms a public api update request into a infrastructure one.
func ToInfraClusterUpdate(update grpc_public_api_go.UpdateClusterRequest) *grpc_infrastructure_go.UpdateClusterRequest {

	result := &grpc_infrastructure_go.UpdateClusterRequest{
		OrganizationId:    update.OrganizationId,
		ClusterId:         update.ClusterId,
		UpdateName:        update.Name != "",
		Name:              update.Name,
		UpdateDescription: update.Description != "",
		Description:       update.Description,
		UpdateLabels:      update.Labels != nil,
		Labels:            update.Labels,
	}

	return result
}

func ToPublicAPICluster(source *grpc_infrastructure_go.Cluster, totalNodes int64, runningNodes int64) *grpc_public_api_go.Cluster {
	return &grpc_public_api_go.Cluster{
		OrganizationId:       source.OrganizationId,
		ClusterId:            source.ClusterId,
		Name:                 source.Name,
		Description:          source.Description,
		ClusterTypeName:      source.ClusterType.String(),
		MultitenantSupport:   source.Multitenant.String(),
		StatusName:           source.Status.String(),
		Labels:               source.Labels,
		TotalNodes:           totalNodes,
		RunningNodes:         runningNodes,
	}
}

func ToPublicAPINode(source * grpc_infrastructure_go.Node) *grpc_public_api_go.Node {
	return &grpc_public_api_go.Node{
		OrganizationId:       source.OrganizationId,
		ClusterId:            source.ClusterId,
		NodeId:               source.NodeId,
		Ip:                   source.Ip,
		Labels:               source.Labels,
		StatusName:           source.Status.String(),
		StateName:            source.State.String(),
	}
}

func ToPublicAPIEndpoints(source []*grpc_application_go.Endpoint) []*grpc_public_api_go.Endpoint {
	result := make([]*grpc_public_api_go.Endpoint, 0)
	for _, e := range source{
		toAdd := &grpc_public_api_go.Endpoint{
			TypeName:             e.Type.String(),
			Path:                 e.Path,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIPorts(source []*grpc_application_go.Port) []*grpc_public_api_go.Port {
	result := make([]*grpc_public_api_go.Port, 0)
	for _, p := range source{
		toAdd := &grpc_public_api_go.Port{
			Name:                 p.Name,
			InternalPort:         p.InternalPort,
			ExposedPort:          p.ExposedPort,
			Endpoints:            ToPublicAPIEndpoints(p.Endpoints),
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIStorage(source []*grpc_application_go.Storage) []*grpc_public_api_go.Storage {
	result := make([]*grpc_public_api_go.Storage, 0)
	for _, s := range source {
		toAdd := &grpc_public_api_go.Storage{
			Size:                 s.Size,
			MountPath:            s.MountPath,
			TypeName:             s.Type.String(),
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIServiceInstances(source []*grpc_application_go.ServiceInstance) [] * grpc_public_api_go.ServiceInstance {
	result := make([]*grpc_public_api_go.ServiceInstance, 0)
	for _, si := range source {
		toAdd := &grpc_public_api_go.ServiceInstance{
			OrganizationId:       si.OrganizationId,
			AppDescriptorId:      si.AppDescriptorId,
			AppInstanceId:        si.AppInstanceId,
			ServiceId:            si.ServiceId,
			Name:                 si.Name,
			Description:          si.Description,
			TypeName:             si.Type.String(),
			Image:                si.Image,
			Credentials:          si.Credentials,
			Specs:                si.Specs,
			Storage:              ToPublicAPIStorage(si.Storage),
			ExposedPorts:         ToPublicAPIPorts(si.ExposedPorts),
			EnvironmentVariables: si.EnvironmentVariables,
			Configs:              si.Configs,
			Labels:               si.Labels,
			DeployAfter:          si.DeployAfter,
			StatusName:           si.Status.String(),
			Endpoints: si.Endpoints,
			DeployedOnClusterId: si.DeployedOnClusterId,
			RunArguments:        si.RunArguments,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIGroupInstances(source []*grpc_application_go.ServiceGroupInstance) [] *grpc_public_api_go.ServiceGroupInstance {
	result := make([]*grpc_public_api_go.ServiceGroupInstance, 0)
	for _, sgi := range source {
		toAdd := &grpc_public_api_go.ServiceGroupInstance{
			OrganizationId:       sgi.OrganizationId,
			AppDescriptorId:      sgi.AppDescriptorId,
			AppInstanceId:        sgi.AppInstanceId,
			ServiceGroupId:       sgi.ServiceGroupId,
			Name:                 sgi.Name,
			Description:          sgi.Description,
			ServiceInstances:     sgi.ServiceInstances,
			PolicyName:           sgi.Policy.String(),
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPISecurityRules(source []*grpc_application_go.SecurityRule) [] *grpc_public_api_go.SecurityRule {
	result := make([]*grpc_public_api_go.SecurityRule, 0)
	for _, sr := range source {
		toAdd := &grpc_public_api_go.SecurityRule{
			OrganizationId:       sr.OrganizationId,
			AppDescriptorId:      sr.AppDescriptorId,
			RuleId:               sr.RuleId,
			Name:                 sr.Name,
			SourceServiceId:      sr.SourceServiceId,
			SourcePort:           sr.SourcePort,
			AccessName:           sr.Access.String(),
			AuthServices:         sr.AuthServices,
			DeviceGroups:         sr.DeviceGroups,
		}
		result = append(result, toAdd)
	}
	return result
}

func ToPublicAPIAppInstance(source *grpc_application_go.AppInstance) * grpc_public_api_go.AppInstance {
	return &grpc_public_api_go.AppInstance{
		OrganizationId:       source.OrganizationId,
		AppDescriptorId:      source.AppDescriptorId,
		AppInstanceId:        source.AppInstanceId,
		Name:                 source.Name,
		Description:          source.Description,
		ConfigurationOptions: source.ConfigurationOptions,
		EnvironmentVariables: source.EnvironmentVariables,
		Labels:               source.Labels,
		Rules:                ToPublicAPISecurityRules(source.Rules),
		Groups:               ToPublicAPIGroupInstances(source.Groups),
		Services:             ToPublicAPIServiceInstances(source.Services),
		StatusName:           source.Status.String(),
	}
}
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

package cli

import (
	"os"
	"time"

	"github.com/golang/protobuf/jsonpb"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-provisioner-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
)

const CliProvisionCheckSleepTime = time.Second * 20

type Provision struct {
	Connection
	Credentials
	KubeConfigOutputPath string
}

func NewProvision(address string, port int, insecure bool, useTLS bool, caCertPath string, output string,
	labelLength int, kubeConfigOutputPath string) *Provision {
	return &Provision{
		Connection:           *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials:          *NewEmptyCredentials(DefaultPath),
		KubeConfigOutputPath: kubeConfigOutputPath,
	}
}

func (p *Provision) ProvisionAndInstall(organizationId string, clusterName string, azureCredentialsPath string,
	azureDnsZoneName string, azureResourceGroup string, clusterType grpc_infrastructure_go.ClusterType, isManagementCluster bool,
	isProduction bool, kubernetesVersion string, nodeType string, numNodes int64, targetPlatform grpc_public_api_go.Platform,
	zone string) {
	err := p.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}

	azureCredentials, err := p.loadAzureCredentials(azureCredentialsPath)
	if err != nil {
		p.PrintResultOrError("", err, "error loading azure credentials")
		return
	}

	c, err := p.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()

	provClient := grpc_public_api_go.NewClustersClient(c)

	installerPlatform := p.convertTargetPlatform(targetPlatform)

	request := grpc_provisioner_go.ProvisionClusterRequest{
		OrganizationId:   organizationId,
		ClusterName:      clusterName,
		AzureCredentials: azureCredentials,
		AzureOptions: &grpc_provisioner_go.AzureProvisioningOptions{
			DnsZoneName:   azureDnsZoneName,
			ResourceGroup: azureResourceGroup,
		},
		ClusterType:         clusterType,
		IsManagementCluster: isManagementCluster,
		IsProduction:        isProduction,
		KubernetesVersion:   kubernetesVersion,
		NodeType:            nodeType,
		NumNodes:            numNodes,
		TargetPlatform:      installerPlatform,
		Zone:                zone,
	}

	ctx, cancel := p.GetContext()
	defer cancel()

	resp, errReq := provClient.ProvisionAndInstall(ctx, &request)

	p.PrintResultOrError(resp, errReq, "cannot provision cluster")
}

// Scale sends the ScaleClusterRequest to the public-api.
func (p *Provision) Scale(organizationID string, clusterID string, clusterType grpc_infrastructure_go.ClusterType,
	numNodes int64, targetPlatform grpc_public_api_go.Platform, azureCredentialsPath string,
	 azureResourceGroup string ) {
	err := p.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
	azureCredentials, err := p.loadAzureCredentials(azureCredentialsPath)
	if err != nil {
		p.PrintResultOrError("", err, "error loading azure credentials")
		return
	}

	c, err := p.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	defer c.Close()

	provClient := grpc_public_api_go.NewClustersClient(c)
	installerPlatform := p.convertTargetPlatform(targetPlatform)

	request := grpc_provisioner_go.ScaleClusterRequest{
		OrganizationId:   organizationID,
		ClusterId: clusterID,
		ClusterType:         clusterType,
		NumNodes:            numNodes,
		// The user may only scale application clusters through public-api-cli
		IsManagementCluster: false,
		TargetPlatform:      installerPlatform,
		AzureCredentials: azureCredentials,
		AzureOptions: &grpc_provisioner_go.AzureProvisioningOptions{
			ResourceGroup: azureResourceGroup,
		},
	}

	ctx, cancel := p.GetContext()
	defer cancel()
	resp, errReq := provClient.Scale(ctx, &request)
	p.PrintResultOrError(resp, errReq, "cannot scale cluster")
}

func (p *Provision) convertTargetPlatform(pbPlatform grpc_public_api_go.Platform) grpc_installer_go.Platform {
	var installerPlatform grpc_installer_go.Platform
	switch pbPlatform {
	case grpc_public_api_go.Platform_AZURE:
		installerPlatform = grpc_installer_go.Platform_AZURE
	case grpc_public_api_go.Platform_MINIKUBE:
		installerPlatform = grpc_installer_go.Platform_MINIKUBE
	default:
		log.Fatal().Str("platform", pbPlatform.String()).Msg("unknown platform")
	}

	return installerPlatform
}

// LoadAzureCredentials loads the content of a file into the grpc structure.
func (p *Provision) loadAzureCredentials(credentialsPath string) (*grpc_provisioner_go.AzureCredentials, derrors.Error) {
	credentials := &grpc_provisioner_go.AzureCredentials{}
	file, err := os.Open(credentialsPath)
	if err != nil {
		return nil, derrors.AsError(err, "cannot open credentials path")
	}
	// The unmarshalling using jsonpb is available due to the fact that the naming of the JSON fields produced
	// by Azure matches the ones defined in the protobuf json mapping.
	err = jsonpb.Unmarshal(file, credentials)
	if err != nil {
		return nil, derrors.AsError(err, "cannot unmarshal content")
	}
	log.Debug().Interface("tenantId", credentials.TenantId).Msg("azure credentials have been loaded")
	return credentials, nil
}

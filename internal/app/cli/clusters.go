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

package cli

import (
	"fmt"
	grpc_common_go "github.com/nalej/grpc-common-go"
	"io/ioutil"
	"math"
	"reflect"
	"time"

	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Clusters struct {
	Connection
	Credentials
}

func NewClusters(address string, port int, insecure bool, useTLS bool, caCertPath string, output string, labelLength int) *Clusters {
	return &Clusters{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output, labelLength),
		Credentials: *NewEmptyCredentials(DefaultPath),
	}
}

func (c *Clusters) load() {
	err := c.LoadCredentials()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot load credentials, try login first")
	}
}

func (c *Clusters) getClient() (grpc_public_api_go.ClustersClient, *grpc.ClientConn) {
	conn, err := c.GetConnection()
	if err != nil {
		log.Fatal().Str("trace", err.DebugReport()).Msg("cannot create the connection with the Nalej platform")
	}
	clusterClient := grpc_public_api_go.NewClustersClient(conn)
	return clusterClient, conn
}

func (c *Clusters) Install(
	organizationID string,
	kubeConfigPath string, ingressHostname string, username string, privateKeyPath string, nodes []string,
	targetPlatform grpc_public_api_go.Platform, useStaticIPAddresses bool, ipAddressIngress string) {

	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	staticIPAddresses := grpc_installer_go.StaticIPAddresses{
		UseStaticIp: useStaticIPAddresses,
		Ingress:     ipAddressIngress,
	}

	installRequest := &grpc_public_api_go.InstallRequest{
		OrganizationId:    organizationID,
		ClusterType:       grpc_infrastructure_go.ClusterType_KUBERNETES,
		Hostname:          ingressHostname,
		InstallBaseSystem: false,
		TargetPlatform:    targetPlatform,
		StaticIpAddresses: &staticIPAddresses,
	}

	if username != "" && privateKeyPath != "" && len(nodes) > 0 {
		installRequest.InstallBaseSystem = true
		log.Info().Msg("Base system will be installed")

		pk, err := ioutil.ReadFile(privateKeyPath)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot read private key")
		}
		installRequest.PrivateKey = string(pk)
		installRequest.Nodes = nodes
	}

	if kubeConfigPath != "" {
		kc, err := ioutil.ReadFile(kubeConfigPath)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot read kube config file")
		}
		installRequest.KubeConfigRaw = string(kc)
	}

	c.load()
	ctx, cancel := c.GetContext(time.Minute * 3)
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()

	log.Debug().Interface("request", installRequest).Msg("Install request")
	response, err := client.Install(ctx, installRequest)

	c.PrintResultOrError(response, err, "cannot install new cluster")
}

func (c *Clusters) Info(organizationID string, clusterID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("clusterID cannot be empty")
	}
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
	}
	retrieved, err := client.Info(ctx, cID)
	c.PrintResultOrError(retrieved, err, "cannot obtain cluster information")
}

func (c *Clusters) List(organizationID string, watch bool, orderBy string, desc bool) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	// check the name of the filed to sort by is correct
	orderByColumn := ""
	if orderBy != "" {
		switch orderBy {
		case "name":
			orderByColumn = orderBy
		case "status", "state":
			orderByColumn = orderBy + "_name"
		default:
			log.Fatal().Msg("only is allowed to sort by name, status or state")
		}
	}

	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()

	var order *grpc_common_go.OrderOptions
	if orderBy != "" {
		descOrder := grpc_common_go.Order_ASC
		if desc == true {
			descOrder = grpc_common_go.Order_DESC
		}
		order = &grpc_common_go.OrderOptions{
			Field:                orderByColumn,
			Order:                descOrder,
		}
	}else{
		order = nil
	}

	orgID := &grpc_public_api_go.ListRequest{
		OrganizationId: organizationID,
		Order: order,
	}
	previous, err := client.List(ctx, orgID)
	c.PrintResultOrError(previous, err, "cannot obtain cluster list")
	toCompare := make(map[string]*grpc_public_api_go.Cluster, 0)
	if watch {
		for _, p := range previous.Clusters {
			toCompare[p.ClusterId] = p
		}
	}
	for watch {
		watchCtx, watchCancel := c.GetContext()
		clusters, err := client.List(watchCtx, orgID)
		if err != nil {
			c.PrintResultOrError(clusters, err, "cannot obtain cluster list information")
		}
		toShow := make([]*grpc_public_api_go.Cluster, 0)
		for _, retrieved := range clusters.Clusters {
			found, exists := toCompare[retrieved.ClusterId]
			if !exists {
				toShow = append(toShow, retrieved)
			} else if c.clusterDiff(found, retrieved) {
				toShow = append(toShow, retrieved)
			}
			toCompare[retrieved.ClusterId] = retrieved
		}

		if len(toShow) > 0 {
			clusters.Clusters = toShow
			fmt.Println("")
			c.PrintResultOrError(clusters, err, "cannot obtain cluster list information")
		}

		watchCancel()
		time.Sleep(WatchSleep)
	}
}

func (c *Clusters) clusterDiff(previous *grpc_public_api_go.Cluster, current *grpc_public_api_go.Cluster) bool {
	return (current.Status != previous.Status) || (current.State != previous.State) ||
		!reflect.DeepEqual(current.Labels, previous.Labels)
}

func (c *Clusters) ModifyClusterLabels(organizationID string, clusterID string, add bool, rawLabels string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("clusterID cannot be empty")
	}
	if rawLabels == "" {
		log.Fatal().Msg("labels cannot be empty")
	}
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	updateRequest := &grpc_public_api_go.UpdateClusterRequest{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
		AddLabels:      add,
		RemoveLabels:   !add,
		Labels:         GetLabels(rawLabels),
	}
	updated, err := client.Update(ctx, updateRequest)
	c.PrintResultOrError(updated, err, "cannot update cluster labels")
}

func (c *Clusters) Update(organizationID string, clusterID string, newName string, millicoresConversionFactor float64) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("clusterID cannot be empty")
	}

	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	updateRequest := &grpc_public_api_go.UpdateClusterRequest{
		OrganizationId:                   organizationID,
		ClusterId:                        clusterID,
		UpdateName:                       newName != "",
		Name:                             newName,
		UpdateMillicoresConversionFactor: !math.IsNaN(millicoresConversionFactor),
		MillicoresConversionFactor:       millicoresConversionFactor,
	}
	success, err := client.Update(ctx, updateRequest)
	c.PrintResultOrError(success, err, "cannot update cluster information")
}

func (c *Clusters) CordonCluster(organizationID string, clusterID string) {
	if organizationID == "" {
		log.Fatal().Msg("organization ID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("cluster ID cannot be empty")
	}
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	clusterIDReq := &grpc_infrastructure_go.ClusterId{
		ClusterId:      clusterID,
		OrganizationId: organizationID,
	}
	success, err := client.Cordon(ctx, clusterIDReq)
	c.PrintResultOrError(success, err, "cannot cordon cluster")
}

func (c *Clusters) UncordonCluster(organizationID string, clusterID string) {
	if organizationID == "" {
		log.Fatal().Msg("organization ID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("cluster ID cannot be empty")
	}
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	clusterIDReq := &grpc_infrastructure_go.ClusterId{
		ClusterId:      clusterID,
		OrganizationId: organizationID,
	}
	success, err := client.Uncordon(ctx, clusterIDReq)
	c.PrintResultOrError(success, err, "cannot uncordon cluster")
}

func (c *Clusters) DrainCluster(organizationID string, clusterID string) {
	if organizationID == "" {
		log.Fatal().Msg("organization ID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("cluster ID cannot be empty")
	}
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	clusterIDReq := &grpc_infrastructure_go.ClusterId{
		ClusterId:      clusterID,
		OrganizationId: organizationID,
	}
	success, err := client.Drain(ctx, clusterIDReq)
	c.PrintResultOrError(success, err, "cannot drain cluster")
}

// Uninstall a existing cluster. This process will uninstall the nalej platform and
// remove the cluster from the list.
func (c *Clusters) Uninstall(organizationID string, clusterID string, kubeConfigPath string, targetPlatform grpc_public_api_go.Platform) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("clusterID cannot be empty")
	}
	if kubeConfigPath == "" {
		log.Fatal().Msg("kubeConfigPath cannot be empty")
	}
	var kubeConfigRaw = ""

	kc, err := ioutil.ReadFile(kubeConfigPath)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read kubeConfig file")
	}
	kubeConfigRaw = string(kc)

	c.load()

	uninstallRequest := &grpc_public_api_go.UninstallClusterRequest{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
		ClusterType:    grpc_infrastructure_go.ClusterType_KUBERNETES,
		KubeConfigRaw:  kubeConfigRaw,
		TargetPlatform: targetPlatform,
	}
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	success, err := client.Uninstall(ctx, uninstallRequest)
	c.PrintResultOrError(success, err, "cannot uninstall cluster")
}

/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"io/ioutil"
	"time"

	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-monitor-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Clusters struct {
	Connection
	Credentials
}

func NewClusters(address string, port int, insecure bool, useTLS bool, caCertPath string, output string) *Clusters {
	return &Clusters{
		Connection:  *NewConnection(address, port, insecure, useTLS, caCertPath, output),
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

func (c *Clusters) Monitor(organizationID string, clusterID string, rangeMinutes int32) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}
	if clusterID == "" {
		log.Fatal().Msg("clusterID cannot be empty")
	}
	if rangeMinutes < 0 {
		log.Fatal().Int32("range", rangeMinutes).Msg("range should be positive")
	}
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	request := &grpc_infrastructure_monitor_go.ClusterSummaryRequest{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
		RangeMinutes:   rangeMinutes,
	}
	retrieved, err := client.Monitor(ctx, request)
	c.PrintResultOrError(retrieved, err, "cannot obtain cluster monitoring information")
}

func (c *Clusters) List(organizationID string) {
	if organizationID == "" {
		log.Fatal().Msg("organizationID cannot be empty")
	}

	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: organizationID,
	}
	clusters, err := client.List(ctx, orgID)
	c.PrintResultOrError(clusters, err, "cannot obtain cluster list")
}

func (c* Clusters) ModifyClusterLabels(organizationID string, clusterID string, add bool, rawLabels string){
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
		AddLabels: add,
		RemoveLabels: !add,
		Labels: GetLabels(rawLabels),
	}
	updated, err := client.Update(ctx, updateRequest)
	c.PrintResultOrError(updated, err, "cannot update cluster labels")
}

func (c *Clusters) Update(organizationID string, clusterID string, newName string) {
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
		OrganizationId: organizationID,
		ClusterId:      clusterID,
		UpdateName: true,
		Name:           newName,
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
		ClusterId: clusterID,
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
		ClusterId: clusterID,
		OrganizationId: organizationID,
	}
	success, err := client.Uncordon(ctx, clusterIDReq)
	c.PrintResultOrError(success, err, "cannot uncordon cluster")
}
/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"io/ioutil"
	"time"
)

type Clusters struct {
	Connection
	Credentials
}

func NewClusters(address string, port int, insecure bool, caCertPath string) *Clusters {
	return &Clusters{
		Connection:  *NewConnection(address, port, insecure, caCertPath),
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
	organizationID string, clusterID string,
	kubeConfigPath string, ingressHostname string, username string, privateKeyPath string, nodes []string,
	useCoreDNS bool, targetPlatform grpc_public_api_go.Platform) {
	installRequest := &grpc_public_api_go.InstallRequest{
		OrganizationId:    organizationID,
		ClusterId:         clusterID,
		ClusterType:       grpc_infrastructure_go.ClusterType_KUBERNETES,
		Hostname: ingressHostname,
		InstallBaseSystem: false,
		TargetPlatform: targetPlatform,
	}

	if useCoreDNS {
		installRequest.UseKubeDns = false
		installRequest.UseCoreDns = true
	} else {
		installRequest.UseKubeDns = true
		installRequest.UseCoreDns = false
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

	response, err := client.Install(ctx, installRequest)
	c.PrintResultOrError(response, err, "cannot install new cluster")
}

func (c * Clusters) Info(organizationID string, clusterID string){
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	cID := &grpc_infrastructure_go.ClusterId{
		OrganizationId:       organizationID,
		ClusterId:            clusterID,
	}
	retrieved, err := client.Info(ctx, cID)
	c.PrintResultOrError(retrieved, err, "cannot obtain cluster information")
}

func (c *Clusters) List(organizationID string) {
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

func (c *Clusters) Update(organizationID string, clusterID string, newName string, newDescription string) {
	c.load()
	ctx, cancel := c.GetContext()
	client, conn := c.getClient()
	defer conn.Close()
	defer cancel()
	updateRequest := &grpc_public_api_go.UpdateClusterRequest{
		OrganizationId: organizationID,
		ClusterId:      clusterID,
		Name:           newName,
		Description:    newDescription,
	}
	success, err := client.Update(ctx, updateRequest)
	c.PrintResultOrError(success, err, "cannot update cluster information")
}

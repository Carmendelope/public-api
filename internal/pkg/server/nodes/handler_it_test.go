/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

 /*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
*/

package nodes

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/public-api/internal/pkg/server/ithelpers"
	"github.com/nalej/public-api/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"os"
)

var _ = ginkgo.Describe("Nodes", func() {

	const NumNodes = 10

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress = os.Getenv("IT_SM_ADDRESS")
	)

	if systemModelAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server * grpc.Server
	// grpc test listener
	var listener * bufconn.Listener
	// client
	var orgClient grpc_organization_go.OrganizationsClient
	var clustClient grpc_infrastructure_go.ClustersClient
	var nodeClient grpc_infrastructure_go.NodesClient
	var smConn * grpc.ClientConn
	var client grpc_public_api_go.NodesClient

	// Target organization.
	var targetOrganization * grpc_organization_go.Organization
	var targetCluster * grpc_infrastructure_go.Cluster

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)
		clustClient = grpc_infrastructure_go.NewClustersClient(smConn)
		nodeClient = grpc_infrastructure_go.NewNodesClient(smConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(nodeClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterNodesServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewNodesClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
		targetCluster = ithelpers.CreateCluster(targetOrganization, fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), clustClient)
		ithelpers.CreateNodes(targetCluster, NumNodes, clustClient, nodeClient)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.It("should be able to list the nodes in a clusters", func(){

		clusterID := &grpc_infrastructure_go.ClusterId{
			OrganizationId:       targetCluster.OrganizationId,
			ClusterId:            targetCluster.ClusterId,
		}
		list, err := client.ClusterNodes(context.Background(), clusterID)

		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(list.Nodes)).To(gomega.Equal(NumNodes))
	})

})
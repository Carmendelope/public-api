/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
*/

package resources

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

func createEnvironment(
	targetOrganization *grpc_organization_go.Organization,
	numClusters int, numNodes int,
	clustClient grpc_infrastructure_go.ClustersClient,
	nodeClient grpc_infrastructure_go.NodesClient,
	){
	for clusterIndex := 0; clusterIndex < numClusters; clusterIndex++{
		cluster := ithelpers.CreateCluster(targetOrganization,
			fmt.Sprintf("cluster-%d", clusterIndex), clustClient)
		ithelpers.CreateNodes(cluster, numNodes, clustClient, nodeClient)
	}
}

var _ = ginkgo.Describe("Resources", func() {

	const NumNodes = 10
	const NumClusters = 5

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
	var client grpc_public_api_go.ResourcesClient

	// Target organization.
	var targetOrganization * grpc_organization_go.Organization

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)
		clustClient = grpc_infrastructure_go.NewClustersClient(smConn)
		nodeClient = grpc_infrastructure_go.NewNodesClient(smConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(clustClient, nodeClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterResourcesServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewResourcesClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
		createEnvironment(targetOrganization, NumClusters, NumNodes, clustClient, nodeClient)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.It("should be able to obtain the summary", func(){

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId:       targetOrganization.OrganizationId,
		}

		summary, err := client.Summary(context.Background(), organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(summary.TotalClusters).To(gomega.Equal(int64(NumClusters)))
		gomega.Expect(summary.TotalNodes).To(gomega.Equal(int64(NumClusters * NumNodes)))
	})

})
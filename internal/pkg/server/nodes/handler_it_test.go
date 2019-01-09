/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
*/

package nodes

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
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

	if !utils.RunIntegrationTests() {
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
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var orgClient grpc_organization_go.OrganizationsClient
	var clustClient grpc_infrastructure_go.ClustersClient
	var nodeClient grpc_infrastructure_go.NodesClient
	var smConn *grpc.ClientConn
	var client grpc_public_api_go.NodesClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	var targetCluster *grpc_infrastructure_go.Cluster
	var token string
	var devToken string
	var opeToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(interceptor.NewConfig(ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

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
		token = ithelpers.GenerateToken("email@nalej.com",
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
		devToken = ithelpers.GenerateToken("dev@nalej.com",
			targetOrganization.OrganizationId, "Developer", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_APPS})
		opeToken = ithelpers.GenerateToken("op@nalej.com",
			targetOrganization.OrganizationId, "Operator", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_RESOURCES})
	})

	ginkgo.AfterSuite(func() {

		ithelpers.NewTestCleaner(smConn).DeleteOrganizationClusters(targetOrganization.OrganizationId)

		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.It("should be able to list the nodes in a clusters", func() {

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner should be able to list the nodes in a clusters"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: false, Msg: "Developer should NOT be able to list the nodes in a clusters"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: true, Msg: "Operator should be able to list the nodes in a clusters"})

		clusterID := &grpc_infrastructure_go.ClusterId{
			OrganizationId: targetCluster.OrganizationId,
			ClusterId:      targetCluster.ClusterId,
		}
		for _, test := range tests {

			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()

			list, err := client.List(ctx, clusterID)

			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(len(list.Nodes)).To(gomega.Equal(NumNodes))
			}else{
				gomega.Expect(err).NotTo(gomega.Succeed())
			}

		}

	})

})

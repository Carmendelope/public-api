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

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
IT_ORGMNG_ADDRESS=localhost:8950
*/

package resources

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
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
	targetOrganization *grpc_organization_manager_go.Organization,
	numClusters int, numNodes int,
	clustClient grpc_infrastructure_go.ClustersClient,
	nodeClient grpc_infrastructure_go.NodesClient,
) {
	for clusterIndex := 0; clusterIndex < numClusters; clusterIndex++ {
		cluster := ithelpers.CreateCluster(targetOrganization,
			fmt.Sprintf("cluster-%d", clusterIndex), clustClient)
		ithelpers.CreateNodes(cluster, numNodes, clustClient, nodeClient)
	}
}

var _ = ginkgo.Describe("Resources", func() {

	const NumNodes = 10
	const NumClusters = 5

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress = os.Getenv("IT_SM_ADDRESS")
		orgManagerAddress = os.Getenv("IT_ORGMNG_ADDRESS")
	)

	if systemModelAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var orgClient grpc_organization_manager_go.OrganizationsClient
	var clustClient grpc_infrastructure_go.ClustersClient
	var nodeClient grpc_infrastructure_go.NodesClient
	var smConn *grpc.ClientConn
	var orgConn *grpc.ClientConn
	var client grpc_public_api_go.ResourcesClient

	// Target organization.
	var targetOrganization *grpc_organization_manager_go.Organization
	var token string
	var devToken string
	var opeToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(interceptor.NewConfig(
			ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		clustClient = grpc_infrastructure_go.NewClustersClient(smConn)
		nodeClient = grpc_infrastructure_go.NewNodesClient(smConn)

		orgConn	 = utils.GetConnection(orgManagerAddress)
		orgClient = grpc_organization_manager_go.NewOrganizationsClient(orgConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(clustClient, nodeClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterResourcesServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewResourcesClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
		createEnvironment(targetOrganization, NumClusters, NumNodes, clustClient, nodeClient)
		token = ithelpers.GenerateToken("email@nalej.com",
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
		devToken = ithelpers.GenerateToken("dev@nalej.com", targetOrganization.OrganizationId, "Developer",
			"secret", []grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_APPS})
		opeToken = ithelpers.GenerateToken("ope@nalej.com", targetOrganization.OrganizationId, "Operator",
			"secret", []grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_RESOURCES})
	})

	ginkgo.AfterSuite(func() {

		ithelpers.NewTestCleaner(smConn).DeleteOrganizationClusters(targetOrganization.OrganizationId)

		server.Stop()
		listener.Close()
		smConn.Close()
		orgConn.Close()
	})

	ginkgo.It("should be able to obtain the summary", func() {

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner should be able to obtain the summary"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: false, Msg: "Developer should NOT be able to obtain the summary"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: true, Msg: "Operator should be able to obtain the summary"})

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		for _, test := range tests {
			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			summary, err := client.Summary(ctx, organizationID)
			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(summary.TotalClusters).To(gomega.Equal(int64(NumClusters)))
				gomega.Expect(summary.TotalNodes).To(gomega.Equal(int64(NumClusters * NumNodes)))
			} else {
				gomega.Expect(err).NotTo(gomega.Succeed())
			}
		}
	})

})

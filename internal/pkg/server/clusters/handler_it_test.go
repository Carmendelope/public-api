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
IT_ORGMGR_ADDRESS=localhost:8950
IT_INFRAMGR_ADDRESS=localhost:8860
IT_IM_COORD_ADDRESS=localhost:8423
*/

package clusters

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	grpc_common_go "github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
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

var _ = ginkgo.Describe("Clusters", func() {

	const NumNodes = 10

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress    = os.Getenv("IT_SM_ADDRESS")
		orgManagerAddress     = os.Getenv("IT_ORGMGR_ADDRESS")
		infraManagerAddress   = os.Getenv("IT_INFRAMGR_ADDRESS")
		monitorManagerAddress = os.Getenv("IT_MONITORING_MANAGER_ADDRESS")
	)

	if systemModelAddress == "" || infraManagerAddress == "" || monitorManagerAddress == "" {
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
	var infraClient grpc_infrastructure_manager_go.InfrastructureManagerClient
	var smConn *grpc.ClientConn
	var orgConn *grpc.ClientConn
	var infraConn *grpc.ClientConn
	var client grpc_public_api_go.ClustersClient

	// Target organization.
	var targetOrganization *grpc_organization_manager_go.Organization
	var targetCluster *grpc_infrastructure_go.Cluster
	var token string
	var devToken string
	var opeToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(interceptor.NewConfig(ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		clustClient = grpc_infrastructure_go.NewClustersClient(smConn)
		nodeClient = grpc_infrastructure_go.NewNodesClient(smConn)
		infraConn = utils.GetConnection(infraManagerAddress)
		infraClient = grpc_infrastructure_manager_go.NewInfrastructureManagerClient(infraConn)
		orgConn = utils.GetConnection(orgManagerAddress)
		orgClient = grpc_organization_manager_go.NewOrganizationsClient(orgConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(clustClient, nodeClient, infraClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterClustersServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewClustersClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
		targetCluster = ithelpers.CreateCluster(targetOrganization, fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), clustClient)
		ithelpers.CreateNodes(targetCluster, NumNodes, clustClient, nodeClient)
		token = ithelpers.GenerateToken("email@nalej.com",
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})

		devToken = ithelpers.GenerateToken("dev@nalej.com",
			targetOrganization.OrganizationId, "Developer", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_APPS})

		opeToken = ithelpers.GenerateToken("oper@nalej.com",
			targetOrganization.OrganizationId, "Operator", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_RESOURCES})
	})

	ginkgo.AfterSuite(func() {

		ithelpers.NewTestCleaner(smConn).DeleteOrganizationClusters(targetOrganization.OrganizationId)

		server.Stop()
		listener.Close()
		smConn.Close()
		orgConn.Close()
	})

	ginkgo.It("should be able to retrieve the information of a cluster", func() {
		clusterID := &grpc_infrastructure_go.ClusterId{
			OrganizationId: targetCluster.OrganizationId,
			ClusterId:      targetCluster.ClusterId,
		}

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner should be able to retrieve the information of a cluster"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: false, Msg: "Developer should NOT be able to retrieve the information of a cluster"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: true, Msg: "Operator should be able  to retrieve the information of a cluster"})

		for _, test := range tests {
			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			retrieved, err := client.Info(ctx, clusterID)
			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(retrieved.ClusterId).Should(gomega.Equal(targetCluster.ClusterId))
				gomega.Expect(retrieved.MultitenantSupport).Should(gomega.Equal(targetCluster.Multitenant.String()))
				gomega.Expect(retrieved.ClusterTypeName).Should(gomega.Equal(targetCluster.ClusterType.String()))
			} else {
				gomega.Expect(err).NotTo(gomega.Succeed())
			}
		}
	})

	ginkgo.FIt("should be able to list the clusters", func() {

		organizationID := &grpc_public_api_go.ListRequest{
			OrganizationId: targetOrganization.OrganizationId,
			Order: &grpc_common_go.OrderOptions{
				Field: "name",
				Order: grpc_common_go.Order_ASC,
			},
		}

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner should be able to list the clusters"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: false, Msg: "Developer should NOT be able to list the clusters"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: true, Msg: "Operator should be able to list the clusters"})

		for _, test := range tests {
			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			clusters, err := client.List(ctx, organizationID)
			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(len(clusters.Clusters)).To(gomega.Equal(1))
				c0 := clusters.Clusters[0]
				gomega.Expect(c0.TotalNodes).Should(gomega.Equal(int64(NumNodes)))
				gomega.Expect(c0.RunningNodes).Should(gomega.Equal(int64(0)))
			} else {
				gomega.Expect(err).NotTo(gomega.Succeed())
			}
		}

	})

	ginkgo.It("should be able to update a cluster", func() {

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner should be able to update a cluster"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: false, Msg: "Developer should NOT be able to update a cluster"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: true, Msg: "Operator should be able to update a cluster"})

		newLabels := make(map[string]string, 0)
		newLabels["nk"] = "nv"

		for _, test := range tests {
			updateRequest := &grpc_public_api_go.UpdateClusterRequest{
				OrganizationId: targetCluster.OrganizationId,
				ClusterId:      targetCluster.ClusterId,
				Name:           "newName: " + ithelpers.GenerateUUID(),
				Labels:         newLabels,
			}

			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			done, err := client.Update(ctx, updateRequest)
			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(done).ToNot(gomega.BeNil())

				clusterID := &grpc_infrastructure_go.ClusterId{
					OrganizationId: targetCluster.OrganizationId,
					ClusterId:      targetCluster.ClusterId,
				}
				ctx2, cancel2 := ithelpers.GetContext(token)
				defer cancel2()
				retrieved, err := clustClient.GetCluster(ctx2, clusterID)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(retrieved.Name).Should(gomega.Equal(updateRequest.Name))
				gomega.Expect(retrieved.Labels).Should(gomega.Equal(updateRequest.Labels))
			} else {
				gomega.Expect(err).NotTo(gomega.Succeed())
			}
		}

	})

})

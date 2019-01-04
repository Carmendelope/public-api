/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
IT_INFRAMGR_ADDRESS=localhost:8860
*/

package clusters

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-infrastructure-manager-go"
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

var _ = ginkgo.Describe("Clusters", func() {

	const NumNodes = 10

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress  = os.Getenv("IT_SM_ADDRESS")
		infraManagerAddress = os.Getenv("IT_INFRAMGR_ADDRESS")
	)

	if systemModelAddress == "" || infraManagerAddress == "" {
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
	var infraClient grpc_infrastructure_manager_go.InfrastructureManagerClient
	var smConn *grpc.ClientConn
	var infraConn *grpc.ClientConn
	var client grpc_public_api_go.ClustersClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	var targetCluster *grpc_infrastructure_go.Cluster
	var token string
	var devToken string
	var opeToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		//authConfig := ithelpers.GetAuthConfig(
		//	"/public_api.Clusters/Info",
		//	"/public_api.Clusters/List",
		//	"/public_api.Clusters/Update")
		//server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(
		//	interceptor.NewConfig(authConfig, "secret", ithelpers.AuthHeader)))

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(interceptor.NewConfig(ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)
		clustClient = grpc_infrastructure_go.NewClustersClient(smConn)
		nodeClient = grpc_infrastructure_go.NewNodesClient(smConn)
		infraConn = utils.GetConnection(infraManagerAddress)
		infraClient = grpc_infrastructure_manager_go.NewInfrastructureManagerClient(infraConn)

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
		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.It("should be able to retrieve the information of a cluster", func(){
		clusterID := &grpc_infrastructure_go.ClusterId{
			OrganizationId:       targetCluster.OrganizationId,
			ClusterId:            targetCluster.ClusterId,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		retrieved, err := client.Info(ctx, clusterID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved.ClusterId).Should(gomega.Equal(targetCluster.ClusterId))
		gomega.Expect(retrieved.MultitenantSupport).Should(gomega.Equal(targetCluster.Multitenant.String()))
		gomega.Expect(retrieved.ClusterTypeName).Should(gomega.Equal(targetCluster.ClusterType.String()))
	})

	ginkgo.It("Developer should NOT be able to retrieve the information of a cluster", func(){
		clusterID := &grpc_infrastructure_go.ClusterId{
			OrganizationId:       targetCluster.OrganizationId,
			ClusterId:            targetCluster.ClusterId,
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()
		_, err := client.Info(ctx, clusterID)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	ginkgo.It("Operator should be able to retrieve the information of a cluster", func(){
		clusterID := &grpc_infrastructure_go.ClusterId{
			OrganizationId:       targetCluster.OrganizationId,
			ClusterId:            targetCluster.ClusterId,
		}
		ctx, cancel := ithelpers.GetContext(opeToken)
		defer cancel()
		retrieved, err := client.Info(ctx, clusterID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved.ClusterId).Should(gomega.Equal(targetCluster.ClusterId))
		gomega.Expect(retrieved.MultitenantSupport).Should(gomega.Equal(targetCluster.Multitenant.String()))
		gomega.Expect(retrieved.ClusterTypeName).Should(gomega.Equal(targetCluster.ClusterType.String()))
	})


	ginkgo.It("should be able to list the clusters", func() {

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		clusters, err := client.List(ctx, organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(clusters.Clusters)).To(gomega.Equal(1))
		c0 := clusters.Clusters[0]
		gomega.Expect(c0.TotalNodes).Should(gomega.Equal(int64(NumNodes)))
		gomega.Expect(c0.RunningNodes).Should(gomega.Equal(int64(0)))
	})
	ginkgo.It("Developer should NOT be able to list the clusters", func() {

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()
		_, err := client.List(ctx, organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Operator should be able to list the clusters", func() {

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(opeToken)
		defer cancel()
		clusters, err := client.List(ctx, organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(clusters.Clusters)).To(gomega.Equal(1))
		c0 := clusters.Clusters[0]
		gomega.Expect(c0.TotalNodes).Should(gomega.Equal(int64(NumNodes)))
		gomega.Expect(c0.RunningNodes).Should(gomega.Equal(int64(0)))
	})

	ginkgo.It("should be able to update a cluster", func() {
		newLabels := make(map[string]string, 0)
		newLabels["nk"] = "nv"
		updateRequest := &grpc_public_api_go.UpdateClusterRequest{
			OrganizationId: targetCluster.OrganizationId,
			ClusterId:      targetCluster.ClusterId,
			Name:           "newName",
			Description:    "newDescription",
			Labels:         newLabels,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		done, err := client.Update(ctx, updateRequest)
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
		gomega.Expect(retrieved.Description).Should(gomega.Equal(updateRequest.Description))
		gomega.Expect(retrieved.Labels).Should(gomega.Equal(updateRequest.Labels))
	})
	ginkgo.It("Developer should NOT be able to update a cluster", func() {
		newLabels := make(map[string]string, 0)
		newLabels["pp"] = "pp"
		updateRequest := &grpc_public_api_go.UpdateClusterRequest{
			OrganizationId: targetCluster.OrganizationId,
			ClusterId:      targetCluster.ClusterId,
			Name:           "newName",
			Description:    "newDescription",
			Labels:         newLabels,
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()
		_, err := client.Update(ctx, updateRequest)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	ginkgo.It("Operator should be able to update a cluster", func() {
		newLabels := make(map[string]string, 0)
		newLabels["op"] = "OP"
		updateRequest := &grpc_public_api_go.UpdateClusterRequest{
			OrganizationId: targetCluster.OrganizationId,
			ClusterId:      targetCluster.ClusterId,
			Name:           "newName",
			Description:    "newDescription",
			Labels:         newLabels,
		}
		ctx, cancel := ithelpers.GetContext(opeToken)
		defer cancel()
		done, err := client.Update(ctx, updateRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(done).ToNot(gomega.BeNil())

		clusterID := &grpc_infrastructure_go.ClusterId{
			OrganizationId: targetCluster.OrganizationId,
			ClusterId:      targetCluster.ClusterId,
		}
		ctx2, cancel2 := ithelpers.GetContext(opeToken)
		defer cancel2()
		retrieved, err := clustClient.GetCluster(ctx2, clusterID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved.Name).Should(gomega.Equal(updateRequest.Name))
		gomega.Expect(retrieved.Description).Should(gomega.Equal(updateRequest.Description))
		gomega.Expect(retrieved.Labels).Should(gomega.Equal(updateRequest.Labels))
	})

})

/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
IT_APPMGR_ADDRESS=localhost:8910
*/

package applications

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-manager-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
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

var _ = ginkgo.Describe("Applications", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress = os.Getenv("IT_SM_ADDRESS")
		appManagerAddress  = os.Getenv("IT_APPMGR_ADDRESS")
	)

	if systemModelAddress == "" || appManagerAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var orgClient grpc_organization_go.OrganizationsClient
	var appClient grpc_application_manager_go.ApplicationManagerClient
	var smConn *grpc.ClientConn
	var appConn *grpc.ClientConn
	var client grpc_public_api_go.ApplicationsClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	var token string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		authConfig := ithelpers.GetAuthConfig(
			"/public_api.Applications/AddAppDescriptor", "/public_api.Applications/ListAppDescriptors",
			"/public_api.Applications/GetAppDescriptor", "/public_api.Applications/Deploy", "/public_api.Applications/Undeploy",
			"/public_api.Applications/ListAppInstances", "/public_api.Applications/GetAppInstance")
		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(
			interceptor.NewConfig(authConfig, "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)
		appConn = utils.GetConnection(appManagerAddress)
		appClient = grpc_application_manager_go.NewApplicationManagerClient(appConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(appClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterApplicationsServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewApplicationsClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
		token = ithelpers.GenerateToken("email@nalej.com",
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		smConn.Close()
		appConn.Close()
	})

	ginkgo.Context("descriptors", func() {
		ginkgo.It("should be able to register a new descriptor", func() {
			toAdd := ithelpers.GetAddDescriptorRequest(targetOrganization.OrganizationId)
			ctx, cancel := ithelpers.GetContext(token)
			defer cancel()
			added, err := client.AddAppDescriptor(ctx, toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.AppDescriptorId).ShouldNot(gomega.BeEmpty())
		})
		ginkgo.It("should be able to get the information of a descriptor", func() {
			toAdd := ithelpers.GetAddDescriptorRequest(targetOrganization.OrganizationId)
			ctx, cancel := ithelpers.GetContext(token)
			defer cancel()
			added, err := client.AddAppDescriptor(ctx, toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added.AppDescriptorId).ShouldNot(gomega.BeEmpty())
			appDescriptorID := &grpc_application_go.AppDescriptorId{
				OrganizationId:  added.OrganizationId,
				AppDescriptorId: added.AppDescriptorId,
			}
			ctx2, cancel2 := ithelpers.GetContext(token)
			defer cancel2()
			retrieved, err := client.GetAppDescriptor(ctx2, appDescriptorID)
			gomega.Expect(retrieved.AppDescriptorId).Should(gomega.Equal(added.AppDescriptorId))
		})
		ginkgo.It("should be able to list the existing descriptors", func() {
			numDescriptors := 5
			org := ithelpers.CreateOrganization(fmt.Sprintf("list-desc-%d", ginkgo.GinkgoRandomSeed()), orgClient)
			orgToken := ithelpers.GenerateToken("email@nalej.com",
				org.OrganizationId, "Owner", "secret",
				[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
			for i := 0; i < numDescriptors; i++ {
				toAdd := ithelpers.GetAddDescriptorRequest(org.OrganizationId)
				ctx, cancel := ithelpers.GetContext(orgToken)
				defer cancel()
				added, err := client.AddAppDescriptor(ctx, toAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(added.AppDescriptorId).ShouldNot(gomega.BeEmpty())
			}

			organizationID := &grpc_organization_go.OrganizationId{
				OrganizationId: org.OrganizationId,
			}
			ctx, cancel := ithelpers.GetContext(orgToken)
			defer cancel()
			list, err := client.ListAppDescriptors(ctx, organizationID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list.Descriptors)).Should(gomega.Equal(numDescriptors))
		})
	})

	ginkgo.Context("instances", func() {

		var targetDescriptor *grpc_application_go.AppDescriptor

		ginkgo.BeforeEach(func() {
			targetDescriptor = ithelpers.CreateAppDescriptor(targetOrganization.OrganizationId, appClient)
		})

		ginkgo.FIt("should be able to deploy an application", func() {
			toDeploy := &grpc_application_manager_go.DeployRequest{
				OrganizationId:  targetDescriptor.OrganizationId,
				AppDescriptorId: targetDescriptor.AppDescriptorId,
				Name:            "deploy-test",
				Description:     "deploy-test",
			}
			ctx, cancel := ithelpers.GetContext(token)
			defer cancel()
			deployed, err := client.Deploy(ctx, toDeploy)
			if err != nil {
				log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error")
			}
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(deployed.AppInstanceId).ShouldNot(gomega.BeEmpty())
		})

		ginkgo.PIt("should be able to undeploy an application", func() {

		})

		ginkgo.PIt("should be able to list the running instances", func() {

		})

		ginkgo.PIt("should be able to retrieve the information of a running instance", func() {

		})
	})

})

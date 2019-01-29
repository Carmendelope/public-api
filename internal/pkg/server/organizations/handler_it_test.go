/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
*/

package organizations

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
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

var _ = ginkgo.Describe("Organizations", func() {

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
	var smConn *grpc.ClientConn
	var client grpc_public_api_go.OrganizationsClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	var token string
	var devToken string
	var opeToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(interceptor.NewConfig(
			ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))
		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(orgClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterOrganizationsServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewOrganizationsClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
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
		server.Stop()
		listener.Close()
		smConn.Close()
	})

	ginkgo.It("should be able to retrieve an existing organization", func() {

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner should be able to retrieve an existing organization"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: true, Msg: "Developer should be able to retrieve an existing organization"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: true, Msg: "Operator should be able to retrieve an existing organization"})

		orgID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		for _, test := range tests {
			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			info, err := client.Info(ctx, orgID)
			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(info.OrganizationId).Should(gomega.Equal(targetOrganization.OrganizationId))
				gomega.Expect(info.Name).Should(gomega.Equal(targetOrganization.Name))
			} else {
				gomega.Expect(err).NotTo(gomega.Succeed())
			}
		}

	})

	ginkgo.It("should fail on an organization that does not exists", func() {

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: false, Msg: "should fail on an organization that does not exists"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: false, Msg: "should fail on an organization that does not exists"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: false, Msg: "should fail on an organization that does not exists"})

		orgID := &grpc_organization_go.OrganizationId{
			OrganizationId: "does-not-exists",
		}
		for _, test := range tests {
			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			info, err := client.Info(ctx, orgID)
			if test.Success {
				gomega.Expect(info).NotTo(gomega.BeNil())
			} else {
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(info).To(gomega.BeNil())
			}

		}

	})

})

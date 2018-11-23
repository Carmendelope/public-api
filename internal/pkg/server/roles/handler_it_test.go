/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
IT_USER_MANAGER_ADDRESS=localhost:8920
*/

package roles

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-manager-go"
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

var _ = ginkgo.Describe("Roles", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress = os.Getenv("IT_SM_ADDRESS")
		userManagerAddress = os.Getenv("IT_USER_MANAGER_ADDRESS")
	)

	if systemModelAddress == "" || userManagerAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var orgClient grpc_organization_go.OrganizationsClient
	var umClient grpc_user_manager_go.UserManagerClient
	var smConn *grpc.ClientConn
	var umConn *grpc.ClientConn
	var client grpc_public_api_go.RolesClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	var targetRole *grpc_authx_go.Role
	var token string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		authConfig := ithelpers.GetAuthConfig("/public_api.Roles/List")
		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(
			interceptor.NewConfig(authConfig, "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)
		umConn = utils.GetConnection(userManagerAddress)
		umClient = grpc_user_manager_go.NewUserManagerClient(umConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(umClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterRolesServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewRolesClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
		targetRole = ithelpers.CreateRole(targetOrganization.OrganizationId, umClient)
		token = ithelpers.GenerateToken("email@nalej.com",
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		smConn.Close()
		umConn.Close()
	})

	ginkgo.It("should be able to list the roles in the system", func() {
		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		roleList, err := client.List(ctx, organizationID)

		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(roleList.Roles)).Should(gomega.Equal(1))
		gomega.Expect(roleList.Roles[0].RoleId).Should(gomega.Equal(targetRole.RoleId))
	})

})

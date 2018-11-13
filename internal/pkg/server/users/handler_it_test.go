/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
IT_USER_MANAGER_ADDRESS=localhost:8920
*/

package users

import (
	"context"
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-public-api-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/public-api/internal/pkg/server/ithelpers"
	"github.com/nalej/public-api/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
	"os"
)

var _ = ginkgo.Describe("Users", func() {

	const NumUsers = 5

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
	var client grpc_public_api_go.UsersClient

	// Target organization.
	var targetOrganization *grpc_organization_go.Organization
	var targetRole *grpc_authx_go.Role
	var targetUser *grpc_user_manager_go.User
	var token string

	ginkgo.BeforeSuite(func() {
		rand.Seed(ginkgo.GinkgoRandomSeed())
		listener = test.GetDefaultListener()
		authConfig := ithelpers.GetAuthConfig(
			"/public_api.Users/Info", "/public_api.Users/List",
			"/public_api.Users/Delete", "/public_api.Users/ResetPassword",
			"/public_api.Users/Update")
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
		grpc_public_api_go.RegisterUsersServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewUsersClient(conn)
		targetOrganization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient)
		targetRole = ithelpers.CreateRole(targetOrganization.OrganizationId, umClient)
		targetUser = ithelpers.CreateUser(targetOrganization.OrganizationId, targetRole.RoleId, umClient)
		token = ithelpers.GenerateToken(targetUser.Email,
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		smConn.Close()
		umConn.Close()
	})

	ginkgo.It("should be able to retrieve the user information", func() {
		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		info, err := client.Info(ctx, userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(info.OrganizationId).Should(gomega.Equal(targetUser.OrganizationId))
		gomega.Expect(info.Email).Should(gomega.Equal(targetUser.Email))
		gomega.Expect(info.Name).Should(gomega.Equal(targetUser.Name))
		gomega.Expect(info.RoleName).Should(gomega.Equal(targetRole.Name))
	})

	ginkgo.It("should be able list users in an organization", func() {
		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		list, err := client.List(ctx, organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(list.Users)).Should(gomega.Equal(1))
	})

	ginkgo.PIt("should be able to delete a user", func() {
		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}
		success, err := client.Delete(context.Background(), userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})

	ginkgo.PIt("should be able to reset the password of a user", func() {
		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}
		reset, err := client.ResetPassword(context.Background(), userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(reset.NewPassword).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.PIt("should be able to update an existing user", func() {
		updateUserRequest := &grpc_user_go.UpdateUserRequest{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
			Name:           "newName",
			Role:           "newRole",
		}
		success, err := client.Update(context.Background(), updateUserRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})
})

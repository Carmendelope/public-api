/*
 * Copyright 2019 Nalej
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
IT_USER_MANAGER_ADDRESS=localhost:8920
*/

package users

import (
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
	var devToken string
	var opToken string

	ginkgo.BeforeSuite(func() {
		rand.Seed(ginkgo.GinkgoRandomSeed())
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(
			interceptor.NewConfig(ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

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
		token = ithelpers.GenerateToken("owner@nalej.com",
			targetOrganization.OrganizationId, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})
		devToken = ithelpers.GenerateToken("dev@nalej.com",
			targetOrganization.OrganizationId, "Developer", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_APPS})
		opToken = ithelpers.GenerateToken("ope@nalej.com",
			targetOrganization.OrganizationId, "OPerator", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_RESOURCES})

	})

	ginkgo.AfterSuite(func() {
		ithelpers.NewTestCleaner(umConn)
		ithelpers.DeleteAllUsers(targetOrganization.OrganizationId, umClient)
		server.Stop()
		listener.Close()
		smConn.Close()
		umConn.Close()
	})

	ginkgo.BeforeEach(func() {
		//	ithelpers.DeleteAllUsers(targetOrganization.OrganizationId, umClient)
		targetUser = ithelpers.CreateUser(targetOrganization.OrganizationId, targetRole.RoleId, umClient)
	})

	ginkgo.It("should be able to add a new user", func() {
		addRequest := &grpc_public_api_go.AddUserRequest{
			OrganizationId: targetOrganization.OrganizationId,
			Email:          fmt.Sprintf("random%d@nalej.com", rand.Int()),
			Password:       "password",
			Name:           "Name",
			RoleName:       targetRole.Name,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		added, err := client.Add(ctx, addRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added.OrganizationId).Should(gomega.Equal(addRequest.OrganizationId))
		gomega.Expect(added.Email).Should(gomega.Equal(addRequest.Email))
		gomega.Expect(added.RoleName).Should(gomega.Equal(addRequest.RoleName))
	})
	ginkgo.It("Developer should NOT be able to add a new user", func() {
		addRequest := &grpc_public_api_go.AddUserRequest{
			OrganizationId: targetOrganization.OrganizationId,
			Email:          fmt.Sprintf("developer%d@nalej.com", rand.Int()),
			Password:       "password",
			Name:           "Name",
			RoleName:       targetRole.Name,
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()
		_, err := client.Add(ctx, addRequest)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})
	ginkgo.It("Operator should NOT be able to add a new user", func() {
		addRequest := &grpc_public_api_go.AddUserRequest{
			OrganizationId: targetOrganization.OrganizationId,
			Email:          fmt.Sprintf("operator%d@nalej.com", rand.Int()),
			Password:       "password",
			Name:           "Name",
			RoleName:       targetRole.Name,
		}
		ctx, cancel := ithelpers.GetContext(opToken)
		defer cancel()
		_, err := client.Add(ctx, addRequest)
		gomega.Expect(err).NotTo(gomega.Succeed())
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

	ginkgo.It("Developer should NOT be able to retrieve the user information", func() {

		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()
		_, err := client.Info(ctx, userID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	ginkgo.It("Operator should NOT be able to retrieve the user information", func() {

		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}
		ctx, cancel := ithelpers.GetContext(opToken)
		defer cancel()
		_, err := client.Info(ctx, userID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	ginkgo.It("should be able list users in an organization", func() {

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()
		list, err := client.List(ctx, organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())

	})
	ginkgo.It("Developer should NOT be able list users in an organization", func() {

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()
		_, err := client.List(ctx, organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Operator should NOT be able list users in an organization", func() {

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		ctx, cancel := ithelpers.GetContext(opToken)
		defer cancel()
		_, err := client.List(ctx, organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	ginkgo.It("should be able to delete a user", func() {

		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}

		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()

		success, err := client.Delete(ctx, userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})

	ginkgo.It("Developer should NOT be able to delete a user", func() {

		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}

		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()

		_, err := client.Delete(ctx, userID)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	ginkgo.It("Operator should NOT be able to delete a user", func() {

		userID := &grpc_user_go.UserId{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
		}

		ctx, cancel := ithelpers.GetContext(opToken)
		defer cancel()

		_, err := client.Delete(ctx, userID)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	ginkgo.It("should be able to update an existing user", func() {
		updateUserRequest := &grpc_user_go.UpdateUserRequest{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
			Name:           "newName",
			PhotoUrl:       "newURL",
		}
		ctx, cancel := ithelpers.GetContext(token)
		defer cancel()

		success, err := client.Update(ctx, updateUserRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})
	ginkgo.It("Developer should NOT be able to update an existing user", func() {

		updateUserRequest := &grpc_user_go.UpdateUserRequest{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
			Name:           "newName",
			PhotoUrl:       "newURL",
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()

		_, err := client.Update(ctx, updateUserRequest)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})
	ginkgo.It("Operator should NOT be able to update an existing user", func() {

		updateUserRequest := &grpc_user_go.UpdateUserRequest{
			OrganizationId: targetUser.OrganizationId,
			Email:          targetUser.Email,
			Name:           "newName",
			PhotoUrl:       "newURL",
		}
		ctx, cancel := ithelpers.GetContext(devToken)
		defer cancel()

		_, err := client.Update(ctx, updateUserRequest)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

})

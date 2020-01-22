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
IT_USER_MANAGER_ADDRESS=localhost:8920
IT_ORGMNG_ADDRESS=localhost:8950

*/

package roles

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
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
		orgManagerAddress  = os.Getenv("IT_ORGMNG_ADDRESS")
		userManagerAddress = os.Getenv("IT_USER_MANAGER_ADDRESS")
	)

	if userManagerAddress == "" || orgManagerAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var orgClient grpc_organization_manager_go.OrganizationsClient
	var umClient grpc_user_manager_go.UserManagerClient
	var orgConn *grpc.ClientConn

	var umConn *grpc.ClientConn
	var client grpc_public_api_go.RolesClient

	// Target organization.
	var targetOrganization *grpc_organization_manager_go.Organization
	var targetRole *grpc_authx_go.Role
	var token string
	var devToken string
	var opeToken string

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(interceptor.NewConfig(
			ithelpers.GetAllAuthConfig(), "secret", ithelpers.AuthHeader)))

		orgConn = utils.GetConnection(orgManagerAddress)
		orgClient = grpc_organization_manager_go.NewOrganizationsClient(orgConn)
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
		umConn.Close()
		orgConn.Close()
	})

	ginkgo.It("should be able to list the roles in the system", func() {

		tests := make([]utils.TestResult, 0)
		tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner should be able to list the roles in the system"})
		tests = append(tests, utils.TestResult{Token: devToken, Success: false, Msg: "Developer should NOT be able to list the roles in the system"})
		tests = append(tests, utils.TestResult{Token: opeToken, Success: false, Msg: "Operator should NOT be able to list the roles in the system"})

		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.OrganizationId,
		}
		for _, test := range tests {
			ctx, cancel := ithelpers.GetContext(test.Token)
			defer cancel()
			roleList, err := client.List(ctx, organizationID)

			if test.Success {
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(len(roleList.Roles)).Should(gomega.Equal(1))
				gomega.Expect(roleList.Roles[0].RoleId).Should(gomega.Equal(targetRole.RoleId))
			} else {
				gomega.Expect(err).NotTo(gomega.Succeed())
			}
		}

	})

})

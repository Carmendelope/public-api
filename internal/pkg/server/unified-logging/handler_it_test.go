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
 *
 */

/*
RUN_INTEGRATION_TEST=true
IT_SM_ADDRESS=localhost:8800
IT_UL_COORD_ADDRESS=localhost:8323
*/

package unified_logging

import (
	"fmt"
	"github.com/nalej/authx/pkg/interceptor"
	"github.com/nalej/grpc-application-manager-go"
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
	"time"
)

var _ = ginkgo.Describe("Unified Logging", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		systemModelAddress    = os.Getenv("IT_SM_ADDRESS")
		unifiedLoggingAddress = os.Getenv("IT_UL_COORD_ADDRESS")
	)

	if systemModelAddress == "" || unifiedLoggingAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var ulClient grpc_application_manager_go.UnifiedLoggingClient
	var ulConn *grpc.ClientConn
	var orgClient grpc_organization_go.OrganizationsClient
	var smConn *grpc.ClientConn
	var client grpc_public_api_go.UnifiedLoggingClient

	var organization, appInstance, sgInstance string
	var token string
	var devToken string
	var operToken string

	var from, to time.Time

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()

		authConfig := ithelpers.GetAllAuthConfig()
		server = grpc.NewServer(interceptor.WithServerAuthxInterceptor(interceptor.NewConfig(authConfig, "secret", ithelpers.AuthHeader)))

		smConn = utils.GetConnection(systemModelAddress)
		orgClient = grpc_organization_go.NewOrganizationsClient(smConn)

		ulConn = utils.GetConnection(unifiedLoggingAddress)
		ulClient = grpc_application_manager_go.NewUnifiedLoggingClient(ulConn)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())

		manager := NewManager(ulClient)
		handler := NewHandler(manager)
		grpc_public_api_go.RegisterUnifiedLoggingServer(server, handler)
		test.LaunchServer(server, listener)

		client = grpc_public_api_go.NewUnifiedLoggingClient(conn)

		// Need organization, application descriptor, application instance, service group instance
		organization = ithelpers.CreateOrganization(fmt.Sprintf("testOrg-%d", ginkgo.GinkgoRandomSeed()), orgClient).GetOrganizationId()
		// Instances don't have to exist, we search for them anyway and get empty result
		appInstance = fmt.Sprintf("testAppInstance-%d", ginkgo.GinkgoRandomSeed())
		sgInstance = fmt.Sprintf("testSGInstance-%d", ginkgo.GinkgoRandomSeed())

		from = time.Unix(0, 0)

		to = time.Now()

		token = ithelpers.GenerateToken("email@nalej.com",
			organization, "Owner", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG})

		devToken = ithelpers.GenerateToken("dev@nalej.com", organization, "Developer", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_APPS})

		operToken = ithelpers.GenerateToken("oper@nalej.com", organization, "Operator", "secret",
			[]grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_PROFILE, grpc_authx_go.AccessPrimitive_RESOURCES})
	})

	ginkgo.AfterSuite(func() {
		testCleaner := ithelpers.NewTestCleaner(smConn)
		testCleaner.DeleteOrganizationDescriptors(organization)

		server.Stop()
		listener.Close()
		smConn.Close()
		ulConn.Close()
	})

	ginkgo.Context("search", func() {
		ginkgo.It("should be able to search logs of an application instance", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner"})
			tests = append(tests, utils.TestResult{Token: devToken, Success: true, Msg: "Developer"})
			tests = append(tests, utils.TestResult{Token: operToken, Success: false, Msg: "Operator"})

			request := &grpc_public_api_go.SearchRequest{
				OrganizationId: organization,
				AppInstanceId:  appInstance,
			}
			for _, test := range tests {
				ginkgo.By(test.Msg)
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				result, err := client.Search(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(result.OrganizationId).Should(gomega.Equal(organization))
					gomega.Expect(result.From).Should(gomega.BeNil())
					gomega.Expect(result.To).Should(gomega.BeNil())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
				}
			}
		})

		ginkgo.It("should be able to search logs of a service group instance", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner"})
			tests = append(tests, utils.TestResult{Token: devToken, Success: true, Msg: "Developer"})
			tests = append(tests, utils.TestResult{Token: operToken, Success: false, Msg: "Operator"})

			request := &grpc_public_api_go.SearchRequest{
				OrganizationId:         organization,
				AppInstanceId:          appInstance,
				ServiceGroupInstanceId: sgInstance,
			}
			for _, test := range tests {
				ginkgo.By(test.Msg)
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				result, err := client.Search(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(result.OrganizationId).Should(gomega.Equal(organization))
					gomega.Expect(result.From).Should(gomega.BeNil())
					gomega.Expect(result.To).Should(gomega.BeNil())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
				}
			}
		})

		ginkgo.It("should be able to search logs with a message filter", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner"})
			tests = append(tests, utils.TestResult{Token: devToken, Success: true, Msg: "Developer"})
			tests = append(tests, utils.TestResult{Token: operToken, Success: false, Msg: "Operator"})

			request := &grpc_public_api_go.SearchRequest{
				OrganizationId:         organization,
				AppInstanceId:          appInstance,
				ServiceGroupInstanceId: sgInstance,
				MsgQueryFilter:         "message filter",
			}
			for _, test := range tests {
				ginkgo.By(test.Msg)
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				result, err := client.Search(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(result.OrganizationId).Should(gomega.Equal(organization))
					gomega.Expect(result.From).Should(gomega.BeNil())
					gomega.Expect(result.To).Should(gomega.BeNil())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
				}
			}
		})

		ginkgo.It("should be able to search logs with a time constraint", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner"})
			tests = append(tests, utils.TestResult{Token: devToken, Success: true, Msg: "Developer"})
			tests = append(tests, utils.TestResult{Token: operToken, Success: false, Msg: "Operator"})

			request := &grpc_public_api_go.SearchRequest{
				OrganizationId: organization,
				AppInstanceId:  appInstance,
				From:           from.Unix(),
				To:             to.Unix(),
			}
			for _, test := range tests {
				ginkgo.By(test.Msg)
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				result, err := client.Search(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(result.OrganizationId).Should(gomega.Equal(organization))
					// We don't check from/to, as we're dealing with empty data in this
					// test. This means there are no real minimum and maximum timestamps
					// and from/to are nil.
					gomega.Expect(result.GetFrom()).Should(gomega.BeNil())
					gomega.Expect(result.GetTo()).Should(gomega.BeNil())
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
				}
			}
		})

		ginkgo.It("should be able to retrieve logs in descending order", func() {
			tests := make([]utils.TestResult, 0)
			tests = append(tests, utils.TestResult{Token: token, Success: true, Msg: "Owner"})
			tests = append(tests, utils.TestResult{Token: devToken, Success: true, Msg: "Developer"})
			tests = append(tests, utils.TestResult{Token: operToken, Success: false, Msg: "Operator"})

			request := &grpc_public_api_go.SearchRequest{
				OrganizationId: organization,
				AppInstanceId:  appInstance,
				Order:          &grpc_public_api_go.OrderOptions{Order: grpc_public_api_go.Order_ASC, Field: "timestamp"},
			}
			for _, test := range tests {
				ginkgo.By(test.Msg)
				ctx, cancel := ithelpers.GetContext(test.Token)
				defer cancel()
				result, err := client.Search(ctx, request)
				if test.Success {
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(result.OrganizationId).Should(gomega.Equal(organization))
				} else {
					gomega.Expect(err).NotTo(gomega.Succeed())
				}
			}
		})

		ginkgo.It("should not accept an empty search request", func() {
			request := &grpc_public_api_go.SearchRequest{}

			ctx, cancel := ithelpers.GetContext(token)
			defer cancel()
			_, err := client.Search(ctx, request)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})

		ginkgo.It("should not accept a search request without application instance", func() {
			request := &grpc_public_api_go.SearchRequest{
				OrganizationId: organization,
			}

			ctx, cancel := ithelpers.GetContext(token)
			defer cancel()
			_, err := client.Search(ctx, request)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
})
